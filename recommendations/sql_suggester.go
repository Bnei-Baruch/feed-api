package recommendations

import (
	"database/sql"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/queries"
	log "github.com/sirupsen/logrus"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/mdb"
)

type GenerateSqlFunc func(request core.MoreRequest) string

type SqlSuggester struct {
	db     *sql.DB
	genSql GenerateSqlFunc
	name   string
}

func MakeSqlSuggester(db *sql.DB, genSql GenerateSqlFunc, name string) *SqlSuggester {
	return &SqlSuggester{db: db, genSql: genSql}
}

func (s *SqlSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	query := s.genSql(request)
	log.Infof("Suggester: %s Query: [%s]", s.name, query)
	if query == "" {
		return []core.ContentItem(nil), nil
	}
	rows, err := queries.Raw(s.db, query).Query()
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
		ret = append(ret, core.ContentItem{UID: uid, Date: date, CreatedAt: createdAt, ContentType: contentType, Suggester: s.name})
	}
	log.Infof("ret: %+v", ret)
	return ret, nil
}
