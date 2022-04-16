package learn

import (
	"database/sql"

	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/databases/chronicles/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ChroniclesDbScanner struct {
	localChroniclesDb *sql.DB
	lastReadId        string
	eventType         string
}

func MakeDbScanner(eventType string, lastReadId string, localChronicledDb *sql.DB) EntriesScanner {
	return &ChroniclesDbScanner{
		localChronicledDb,
		lastReadId,
		eventType,
	}
}

func (s *ChroniclesDbScanner) ScanAll() ([]*models.Entry, error) {
	entries, err := models.Entries(qm.Where("id > ?", s.lastReadId), qm.And("client_event_type = ?", s.eventType), qm.OrderBy("id asc")).All(common.LocalChroniclesDb)
	if len(entries) > 0 {
		s.lastReadId = entries[len(entries)-1].ID
	}
	return entries, err
}
