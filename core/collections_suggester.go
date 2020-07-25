package core

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/queries"

	"github.com/Bnei-Baruch/feed-api/mdb"
)

type CollectionSuggester struct {
	db          *sql.DB
	contentType string
}

func MakeCollectionSuggester(db *sql.DB, contentType string) *CollectionSuggester {
	return &CollectionSuggester{db: db, contentType: contentType}
}

func (suggester *CollectionSuggester) More(request MoreRequest) ([]ContentItem, error) {
	currentLessonUIDs := []string(nil)
	for _, ci := range request.CurrentFeed {
		if ci.ContentType == suggester.contentType {
			currentLessonUIDs = append(currentLessonUIDs, ci.UID)
		}
	}
	return suggester.fetchCollection(currentLessonUIDs, request.MoreItems)
}

func (suggester *CollectionSuggester) fetchCollection(currentLessonUIDs []string, moreItems int) ([]ContentItem, error) {
	uidsQuery := ""
	if len(currentLessonUIDs) > 0 {
		quoted := []string(nil)
		for _, uid := range currentLessonUIDs {
			quoted = append(quoted, fmt.Sprintf("'%s'", uid))
		}
		uidsQuery = fmt.Sprintf("and c.uid not in (%s)", strings.Join(quoted, ","))
	}
	query := fmt.Sprintf(`
		select c.uid, MAX((coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text))::date) as date, MAX(cu.created_at) as last_created_at
		from collections as c, content_units as cu, collections_content_units as ccu
		where c.id = ccu.collection_id
		and cu.id = ccu.content_unit_id
		and cu.secure = 0 AND cu.published IS TRUE
		%s
		and c.type_id = %d
		group by c.uid
		order by date desc, last_created_at desc
		limit %d;
		`, uidsQuery, mdb.CONTENT_TYPE_REGISTRY.ByName[suggester.contentType].ID, moreItems)
	rows, err := queries.Raw(suggester.db, query).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := []ContentItem(nil)
	for rows.Next() {
		var uid string
		var date time.Time
		var createdAt time.Time
		err := rows.Scan(&uid, &date, &createdAt)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ContentItem{UID: uid, Date: date, CreatedAt: createdAt, ContentType: suggester.contentType})
	}
	return ret, nil
}
