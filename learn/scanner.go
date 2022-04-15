package learn

import (
	"github.com/Bnei-Baruch/feed-api/databases/chronicles/models"
)

type EntriesScanner interface {
	ScanAll() ([]*models.Entry, error)
}
