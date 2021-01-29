package core

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/queries"
	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/feed-api/mdb"
)

type CollectionSuggester struct {
	db          *sql.DB
	contentType string
}

func MakeCollectionSuggester(db *sql.DB, contentType string) *CollectionSuggester {
	return &CollectionSuggester{
		db:          db,
		contentType: contentType,
	}
}

func init() {
	RegisterSuggester("CollectionSuggester", func(db *sql.DB) Suggester { return MakeCollectionSuggester(db, "") })
}

func (suggester *CollectionSuggester) MarshalSpec() (SuggesterSpec, error) {
	return SuggesterSpec{
		Name: "CollectionSuggester",
		Args: []string{suggester.contentType},
	}, nil
}

func (suggester *CollectionSuggester) UnmarshalSpec(db *sql.DB, spec SuggesterSpec) error {
	if spec.Name != "CollectionSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'CollectionSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 1 {
		return errors.New("CollectionSuggester expected to have only one argument.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("CollectionSuggester expected to have 0 specs got %d.", len(spec.Specs)))
	}
	suggester.contentType = spec.Args[0]
	return nil
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
		ret = append(ret, ContentItem{UID: uid, Date: date, CreatedAt: createdAt, ContentType: suggester.contentType, Suggester: "CollectionSuggester"})
	}
	return ret, nil
}
