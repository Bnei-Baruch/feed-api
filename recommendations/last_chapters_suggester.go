package recommendations

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/queries"
	log "github.com/sirupsen/logrus"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type LastChaptersSuggester struct {
	db *sql.DB
}

func MakeLastChaptersSuggester(db *sql.DB) *LastChaptersSuggester {
	return &LastChaptersSuggester{db: db}
}

func (suggester *LastChaptersSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	currentUIDs := []string(nil)
	for _, ci := range request.CurrentFeed {
		currentUIDs = append(currentUIDs, ci.UID)
	}
	return suggester.fetchLastChapters(currentUIDs, request.MoreItems)
}

func (suggester *LastChaptersSuggester) fetchLastChapters(currentUIDs []string, moreItems int) ([]core.ContentItem, error) {
	if len(currentUIDs) == 0 {
		return []core.ContentItem(nil), nil
	}
	log.Infof("uids: %+v", currentUIDs)
	dateField := `coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text)::date`
	query := fmt.Sprintf(`
		select t.type_id, t.uid, t.date, t.created_as from (
			select cu.type_id, ROW_NUMBER() OVER(PARTITION BY c.uid order by %s desc) as r, cu.uid as uid, %s as date, cu.created_at as created_as
			from collections as c, content_units as cu, collections_content_units as ccu
			where c.id = ccu.collection_id and cu.id = ccu.content_unit_id and
			%s
			and c.uid in (
				select c.uid
				from collections as c, content_units as cu, collections_content_units as ccu
				where %s and c.id = ccu.collection_id and cu.id = ccu.content_unit_id
			) 
			order by date desc
		) as t where t.r <= %d;
		`,
		dateField,
		dateField,
		utils.InClause("cu.uid not in", currentUIDs),
		utils.InClause("cu.uid in", currentUIDs),
		moreItems,
	)
	log.Infof("Query: %s", query)
	rows, err := queries.Raw(suggester.db, query).Query()
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
		ret = append(ret, core.ContentItem{UID: uid, Date: date, CreatedAt: createdAt, ContentType: contentType})
	}
	log.Infof("ret: %+v", ret)
	return ret, nil
}
