package data_models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/feed-api/databases/chronicles/models"
	"github.com/Bnei-Baruch/feed-api/instrumentation"
	"github.com/Bnei-Baruch/feed-api/utils"
)

const (
	MAX_WINDOW_SIZE = 2000000
	SCAN_SIZE       = 1000
	MAX_INTERVAL    = time.Duration(time.Minute)
	MIN_INTERVAL    = time.Duration(100 * time.Millisecond)
)

type ChroniclesWindowModel struct {
	localChroniclesDb *sql.DB
	name              string
	interval          time.Duration

	chroniclesUrl string
	httpClient    *http.Client
	lastReadId    string
	prevReadId    string
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

	return &ChroniclesWindowModel{
		localChroniclesDb,
		"ChroniclesWindowModel",
		time.Duration(time.Minute),
		chroniclesUrl,
		&http.Client{Timeout: 5 * time.Second},
		lastReadId,
		"",
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
	log.Infof("Scanning chronicles entries, last successfull [%s]", m.lastReadId)
	resp, err := m.httpClient.Post(
		fmt.Sprintf("%s/scan", m.chroniclesUrl),
		"application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"id":"%s","limit":%d}`, m.lastReadId, SCAN_SIZE))))
	if err != nil {
		return nil, err
	}
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

func (m *ChroniclesWindowModel) Refresh() error {
	if entries, err := m.ScanChroniclesEntries(); err != nil {
		return err
	} else {
		// Insert entries to local table.
		log.Infof("Inserting %d entries.", len(entries))
		if len(entries) == SCAN_SIZE {
			m.interval = utils.MaxDuration(m.interval/2, MIN_INTERVAL)
		} else {
			m.interval = utils.MinDuration(m.interval*2, MAX_INTERVAL)
		}
		log.Infof("Updated interval to %s", m.interval)
		start := time.Now()
		for _, entry := range entries {
			//log.Infof("Entry: %+v", entry)
			switch entry.ClientEventType {
			case "recommend":
				instrumentation.Stats.RecommendCounter.Inc()
			case "recommend-selected":
				instrumentation.Stats.RecommendSelectedCounter.Inc()
			}
			if err := entry.Insert(m.localChroniclesDb, boil.Infer()); err != nil {
				return err
			}
		}
		log.Infof("Insert done in %s", time.Now().Sub(start))
		if entry, err := models.Entries(qm.OrderBy("id desc"), qm.Offset(MAX_WINDOW_SIZE)).One(m.localChroniclesDb); err == nil && entry != nil {
			log.Infof("Deleting from id: %s", entry.ID)
			if result, err := queries.Raw(fmt.Sprintf("delete from entries where id <= '%s'", entry.ID)).Exec(m.localChroniclesDb); err != nil {
				log.Infof("Failed deleting %+v", err)
				return err
			} else {
				if rowsDeleted, err := result.RowsAffected(); err != nil {
					log.Infof("Failed getting deleted count %+v", err)
					return err
				} else {
					log.Infof("Deleted %d entries from local chronicles by offset.", rowsDeleted)
				}
			}
		} else {
			log.Infof("No delete required.")
		}
		log.Infof("Delete done in %s", time.Now().Sub(start))
	}
	return nil
}
