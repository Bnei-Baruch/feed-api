package data_models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/feed-api/databases/chronicles/models"
	"github.com/Bnei-Baruch/feed-api/instrumentation"
	"github.com/Bnei-Baruch/feed-api/utils"
)

const (
	MAX_WINDOW_SIZE     = 2000000
	SCAN_SIZE           = 50000
	MAX_INTERVAL        = time.Duration(10 * 1000 * time.Millisecond)
	MIN_INTERVAL        = time.Duration(100 * time.Millisecond)
	HTTP_RETRIES        = 3
	DELETE_INSERT_RATIO = 10
)

type ScanHttpErrorRetry struct {
	err error
}

func (e *ScanHttpErrorRetry) Error() string { return fmt.Sprintf("%+v", e.err) }

func (e *ScanHttpErrorRetry) Is(target error) bool {
	_, ok := target.(*ScanHttpErrorRetry)
	return ok
}

type ChroniclesWindowModel struct {
	localChroniclesDb *sql.DB
	name              string
	interval          time.Duration

	chroniclesUrl string
	httpClient    *http.Client
	httpRetries   int64
	lastReadId    string
	prevReadId    string

	refreshCount int64
	totalCount   int64
}

func (m *ChroniclesWindowModel) LastReadId() string {
	return m.lastReadId
}

func (m *ChroniclesWindowModel) PrevReadId() string {
	return m.prevReadId
}

func MakeChroniclesWindowModel(localChroniclesDb *sql.DB, chroniclesUrl string) *ChroniclesWindowModel {
	lastReadId := ""
	if entry, err := models.Entries(qm.OrderBy("id desc"), qm.Limit(1)).One(localChroniclesDb); err == nil && entry != nil {
		lastReadId = entry.ID
	}

	log.Debugf("Chronicles Window last read id: %+s", lastReadId)

	return &ChroniclesWindowModel{
		localChroniclesDb,
		"ChroniclesWindowModel",
		MAX_INTERVAL,
		chroniclesUrl,
		&http.Client{Timeout: 5 * time.Second},
		HTTP_RETRIES,
		lastReadId,
		"",
		0,
		0,
	}
}

func (m *ChroniclesWindowModel) Name() string {
	return m.name
}

func (m *ChroniclesWindowModel) Interval() time.Duration {
	return m.interval
}

type ScanResponse struct {
	Entries []*models.Entry `json:"entries"`
}

func (m *ChroniclesWindowModel) ScanChroniclesEntries() ([]*models.Entry, error) {
	log.Debugf("Scanning chronicles entries, last successfull [%s]", m.lastReadId)
	resp, err := m.httpClient.Post(
		fmt.Sprintf("%s/scan", m.chroniclesUrl),
		"application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"id":"%s","limit":%d}`, m.lastReadId, SCAN_SIZE))))
	if err != nil {
		log.Infof("Non http error %d %+v", m.httpRetries, err)
		if m.httpRetries > 0 {
			m.httpRetries -= 1
			return nil, &ScanHttpErrorRetry{err}
		} else {
			return nil, err
		}
	}
	m.httpRetries = HTTP_RETRIES
	if resp.StatusCode != http.StatusOK { // OK
		// TODO: Consider adding the body as error message.
		return nil, errors.New(fmt.Sprintf("Response code %d for scan: %s.", resp.StatusCode, resp.Status))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var scanResponse ScanResponse
	if err = json.Unmarshal(body, &scanResponse); err != nil {
		return nil, err
	}

	if len(scanResponse.Entries) > 0 {
		m.prevReadId = m.lastReadId
		m.lastReadId = scanResponse.Entries[len(scanResponse.Entries)-1].ID
	}

	return scanResponse.Entries, nil
}

type ABTestingData struct {
	AB map[string]string `json: "ab,omitempty"`
}

type SearchSelectedData struct {
	Rank *int64 `json: "rank,omitempty"`
}

func (m *ChroniclesWindowModel) Refresh() error {
	log.Debugf("Scanning entries...")
	if entries, err := m.ScanChroniclesEntries(); err != nil {
		log.Infof("Scan error: %+v.", err)
		retryError := &ScanHttpErrorRetry{}
		if errors.Is(err, retryError) {
			log.Infof("Scan http error: %+v. Skipping and retrying.", err)
			return nil
		}
		return err
	} else {
		// Insert entries to local table.
		if len(entries) == SCAN_SIZE {
			m.interval = utils.MaxDuration(m.interval/2, MIN_INTERVAL)
		} else {
			m.interval = utils.MinDuration(m.interval*2, MAX_INTERVAL)
		}
		log.Debugf("Updated interval to %s", m.interval)
		log.Debugf("Inserting %d entries.", len(entries))
		start := time.Now()
		allValues := []string(nil)
		for _, entry := range entries {
			instrumentation.Stats.EntriesCounterVec.WithLabelValues(entry.ClientEventType).Inc()
			switch entry.ClientEventType {
			case "recommend":
				ab := ""
				if entry.Data.Valid {
					var abd ABTestingData
					if err := json.Unmarshal(entry.Data.JSON, &abd); err != nil {
						return err
					}
					if version, ok := abd.AB["recommend"]; ok {
						ab = version
					}
				}
				instrumentation.Stats.RecommendCounter.WithLabelValues(ab).Inc()
			case "recommend-selected":
				ab := ""
				if entry.Data.Valid {
					var abd ABTestingData
					if err := json.Unmarshal(entry.Data.JSON, &abd); err != nil {
						return err
					}
					if version, ok := abd.AB["recommend"]; ok {
						ab = version
					}
				}
				instrumentation.Stats.RecommendSelectedCounter.WithLabelValues(ab).Inc()
			case "search":
				instrumentation.Stats.SearchCounter.Inc()
			case "search-selected":
				instrumentation.Stats.SearchSelectedCounter.Inc()
				if entry.Data.Valid {
					var ssd SearchSelectedData
					if err := json.Unmarshal(entry.Data.JSON, &ssd); err != nil {
						return err
					}
					if ssd.Rank != nil {
						instrumentation.Stats.SearchSelectedRankHistogram.Observe(float64(*ssd.Rank))
					}
				} else {
					log.Warnf("Unexpected null json for search-selected entry: %s", entry.ID)
				}
			case "autocomplete":
				instrumentation.Stats.AutocompleteCounter.Inc()
			case "autocomplete-selected":
				instrumentation.Stats.AutocompleteSelectedCounter.Inc()
			}

			entryValues := []string{
				fmt.Sprintf("'%s'", entry.ID),
				fmt.Sprintf("timestamp with time zone 'epoch' + %d * interval '1 microseconds'", entry.CreatedAt.UnixNano()/1000),
				fmt.Sprintf("'%s'", entry.UserID),
				fmt.Sprintf("'%s'", entry.IPAddr),
				fmt.Sprintf("'%s'", strings.ReplaceAll(entry.UserAgent, "'", "''")),
				fmt.Sprintf("'%s'", entry.Namespace),
				utils.NullStringToValue(entry.ClientEventID),
				fmt.Sprintf("'%s'", entry.ClientEventType),
				utils.NullStringToValue(entry.ClientFlowID),
				utils.NullStringToValue(entry.ClientFlowType),
				utils.NullStringToValue(entry.ClientSessionID),
				utils.NullJSONToValue(entry.Data),
			}
			allValues = append(allValues, fmt.Sprintf("(%s)", strings.Join(entryValues, ",")))
		}

		entryAllColumns := []string{"id", "created_at", "user_id", "ip_addr", "user_agent", "namespace", "client_event_id", "client_event_type", "client_flow_id", "client_flow_type", "client_session_id", "data"}
		insertQuery := fmt.Sprintf("INSERT INTO entries (%s) VALUES %s", strings.Join(entryAllColumns, ","), strings.Join(allValues, ","))
		if result, err := queries.Raw(insertQuery).Exec(m.localChroniclesDb); err != nil {
			log.Warnf("SQL: %s", insertQuery)
			log.Warnf("Failed inserting into local chronicles %+v", err)
			return err
		} else {
			if rowsInserted, err := result.RowsAffected(); err != nil {
				log.Warnf("Failed getting inserted count %+v", err)
				return err
			} else {
				m.totalCount += rowsInserted
				log.Debugf("Inserted %d entries to local chronicles by offset. Total of %d.", rowsInserted, m.totalCount)
			}
		}

		log.Debugf("Insert done in %s", time.Now().Sub(start))

		m.refreshCount += 1
		if m.interval == MAX_INTERVAL || m.refreshCount == DELETE_INSERT_RATIO {
			m.refreshCount = 0
			log.Debugf("Checking delete window.")
			start = time.Now()
			if entry, err := models.Entries(qm.OrderBy("id desc"), qm.Offset(MAX_WINDOW_SIZE)).One(m.localChroniclesDb); err == nil && entry != nil {
				log.Debugf("Deleting from id: %s (%s)", entry.ID, time.Now().Sub(start))
				start = time.Now()
				if result, err := queries.Raw(fmt.Sprintf("delete from entries where id <= '%s'", entry.ID)).Exec(m.localChroniclesDb); err != nil {
					log.Warnf("Failed deleting from local chronicles %+v", err)
					return err
				} else {
					if rowsDeleted, err := result.RowsAffected(); err != nil {
						log.Warnf("Failed getting deleted count %+v", err)
						return err
					} else {
						log.Debugf("Deleted %d entries from local chronicles by offset.", rowsDeleted)
					}
				}
			} else {
				log.Debugf("No delete required.")
			}
			log.Debugf("Delete done in %s", time.Now().Sub(start))
		} else {
			log.Debugf("Skipping delete")
		}
	}
	return nil
}
