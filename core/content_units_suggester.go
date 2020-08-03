package core

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/queries"

	"github.com/Bnei-Baruch/feed-api/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type ContentUnitsSuggester struct {
	db           *sql.DB
	contentTypes []string
	name         string
}

func MakeContentUnitsSuggester(db *sql.DB, contentTypes []string) *ContentUnitsSuggester {
	return &ContentUnitsSuggester{db: db, contentTypes: contentTypes, name: "ContentUnitsSuggester"}
}

func (suggester *ContentUnitsSuggester) More(request MoreRequest) ([]ContentItem, error) {
	currentUIDs := []string(nil)
	for _, ci := range request.CurrentFeed {
		if utils.StringInSlice(ci.ContentType, suggester.contentTypes) {
			currentUIDs = append(currentUIDs, ci.UID)
		}
	}
	return suggester.fetchContentUnits(currentUIDs, request.MoreItems)
}

func ContentTypesToContentIds(contentTypes []string) []string {
	contentTypesIds := []string(nil)
	for _, contentType := range contentTypes {
		contentTypesIds = append(contentTypesIds, fmt.Sprintf("%d", mdb.CONTENT_TYPE_REGISTRY.ByName[contentType].ID))
	}
	return contentTypesIds
}

func (suggester *ContentUnitsSuggester) fetchContentUnits(currentUIDs []string, moreItems int) ([]ContentItem, error) {
	query := fmt.Sprintf(`
		select cu.uid, (coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text))::date as date, cu.created_at, cu.type_id
		from content_units as cu
		where cu.secure = 0 AND cu.published IS TRUE
		%s
		%s
		order by date desc, cu.created_at desc
		limit %d;
		`,
		utils.InClause("and cu.uid not in", currentUIDs),
		utils.InClause("and cu.type_id in", ContentTypesToContentIds(suggester.contentTypes)),
		moreItems)
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
		var type_id int64
		err := rows.Scan(&uid, &date, &createdAt, &type_id)
		if err != nil {
			return nil, err
		}
		contentType := mdb.CONTENT_TYPE_REGISTRY.ByID[type_id].Name
		ret = append(ret, ContentItem{UID: uid, Date: date, CreatedAt: createdAt, ContentType: contentType, Suggester: suggester.name})
	}
	return ret, nil
}
