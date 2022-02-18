package learn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Bnei-Baruch/feed-api/databases/chronicles/models"
	"github.com/Bnei-Baruch/feed-api/utils"
)

const (
	INITIAL_SCAN_SIZE = 500
	MAX_SCAN_SIZE     = 20000
	RETRY_INTERVAL    = time.Duration(100 * time.Millisecond)
	HTTP_RETRIES      = 3
	HTTP_TIMEOUT      = 5 * time.Second
)

type ScanResponse struct {
	Entries []*models.Entry `json:"entries"`
}

type ScanHttpErrorRetry struct {
	err error
}

func (e *ScanHttpErrorRetry) Error() string { return fmt.Sprintf("%+v", e.err) }

func (e *ScanHttpErrorRetry) Is(target error) bool {
	_, ok := target.(*ScanHttpErrorRetry)
	return ok
}

type ChroniclesScanner struct {
	chroniclesUrl   string
	httpClient      *http.Client
	httpRetries     int64
	httpLastLatency time.Duration
	scanSize        int64
	lastReadId      string
	prevReadId      string

	eventType string
}

func MakeScanner(eventType string, lastReadId string, chroniclesUrl string) *ChroniclesScanner {
	return &ChroniclesScanner{
		chroniclesUrl,
		&http.Client{Timeout: HTTP_TIMEOUT},
		HTTP_RETRIES,
		0,
		INITIAL_SCAN_SIZE,
		lastReadId,
		"",
		eventType,
	}
}

func (s *ChroniclesScanner) Scan() ([]*models.Entry, error) {
	log.Debugf("Scanning chronicles entries, last successfull [%s]", s.lastReadId)
	start := time.Now()
	resp, err := s.httpClient.Post(
		fmt.Sprintf("%s/scan", s.chroniclesUrl),
		"application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"id":"%s","limit":%d,"event_types":["%s"]}`, s.lastReadId, s.scanSize, s.eventType))))
	diff := time.Now().Sub(start)
	if err != nil {
		log.Infof("Non http error %d %+v", s.httpRetries, err)
		if s.httpRetries > 0 {
			s.httpRetries -= 1
			s.scanSize = utils.MaxInt64(INITIAL_SCAN_SIZE, 2*s.scanSize/3)
			return nil, &ScanHttpErrorRetry{err}
		} else {
			return nil, err
		}
	}
	if diff < HTTP_TIMEOUT/2 {
		s.scanSize = utils.MinInt64(MAX_SCAN_SIZE, 12*s.scanSize/10)
	} else {
		s.scanSize = utils.MaxInt64(INITIAL_SCAN_SIZE, 10*s.scanSize/15)
	}
	s.httpRetries = HTTP_RETRIES
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK { // OK
		// TODO: Consider adding the body as error message.
		bodyString := "[Empty body or error reading]"
		if err == nil {
			bodyString = string(body)
		}
		return nil, errors.New(fmt.Sprintf("Response code %d for scan: %s. Body: %s", resp.StatusCode, resp.Status, bodyString))
	}
	// Body read error.
	if err != nil {
		return nil, err
	}

	var scanResponse ScanResponse
	if err = json.Unmarshal(body, &scanResponse); err != nil {
		return nil, err
	}

	if len(scanResponse.Entries) > 0 {
		s.prevReadId = s.lastReadId
		s.lastReadId = scanResponse.Entries[len(scanResponse.Entries)-1].ID
	}

	return scanResponse.Entries, nil
}

// Will close input channel when done.
func (s *ChroniclesScanner) StartScanning(errChan chan error, c chan []*models.Entry) {
	log.Debugf("Scanning entries to channel...")
	go func() {
		for {
			entries, err := s.Scan()
			if err != nil {
				log.Infof("Scan error: %+v.", err)
				retryError := &ScanHttpErrorRetry{}
				if errors.Is(err, retryError) {
					log.Infof("Scan http error: %+v. Skipping and retrying.", err)
					time.Sleep(RETRY_INTERVAL)
					continue // Retry again.
				}
				errChan <- err
				break // Exit on unexpected error.
			}
			if len(entries) == 0 {
				errChan <- nil // We are done.
				break
			}
			log.Infof("Adding %d entries", len(entries))
			c <- entries
		}
		log.Infof("Closing err channel")
		close(errChan)
		log.Infof("Closing entries channel")
		close(c)
	}()
}

func (s *ChroniclesScanner) ScanAll() ([]*models.Entry, error) {
	errChan := make(chan error, 1)
	entriesChan := make(chan []*models.Entry, 0)
	s.StartScanning(errChan, entriesChan)

	log.Infof("Reading entries...")
	entries := []*models.Entry(nil)
	for {
		log.Infof("Fetching entries...")
		e, opened := <-entriesChan
		log.Infof("Got %d entries adding to %d read items.", len(e), len(entries))
		if opened {
			entries = append(entries, e...)
		} else {
			break
		}
	}
	return entries, <-errChan
}
