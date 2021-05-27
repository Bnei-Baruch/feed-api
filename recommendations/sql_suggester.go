package recommendations

import (
	"database/sql"
	"time"

	"github.com/volatiletech/sqlboiler/queries"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/databases/mdb"
)

type GenerateSqlFunc func(request core.MoreRequest) string

type SqlSuggester struct {
	db     *sql.DB
	genSql GenerateSqlFunc
	Name   string
}

func (s *SqlSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	query := s.genSql(request)
	if query == "" {
		return []core.ContentItem(nil), nil
	}
	rows, err := queries.Raw(query).Query(s.db)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := []core.ContentItem(nil)
	for rows.Next() {
		var typeId int64
		var uid string
		var date time.Time
		var createdAt time.Time
		err := rows.Scan(&typeId, &uid, &date, &createdAt)
		contentType := mdb.CONTENT_TYPE_REGISTRY.ByID[typeId].Name
		if err != nil {
			return nil, err
		}
		ret = append(ret, core.ContentItem{UID: uid, Date: date, CreatedAt: createdAt, ContentType: contentType, Suggester: s.Name})
	}
	return ret, nil
}
