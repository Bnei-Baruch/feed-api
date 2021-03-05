package core

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/queries"
	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/feed-api/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

// Filter collection to have language.
const COLLECTIONS_BY_FIRST_CONTENT_UNIT_CLAUSE = `
  and c.id = collection_first_content_unit_languages.collection_id
`

func CollectionsByFirstUnitLanguagesTableSql(contentTypes []string, languages []string) string {
	orClauses := []string{}
	for i := range languages {
		if languages[i] == "he" {
			orClauses = append(orClauses, "cfcu.original_language = 'he'")
			languages = append(languages[:i], languages[i+1:]...)
			break
		}
	}
	if len(languages) > 0 {
		orClauses = append(orClauses, fmt.Sprintf(`
			0 < (
				select
					count(f.language)
				from
					files as f
				where
					f.content_unit_id = cfcu.content_unit_id
					and
					f.mime_type in ('video/mp4', 'audio/mpeg')
					%s
				)
		`, utils.InClause("and f.language in ", languages)))
	}

	return fmt.Sprintf(`
		(
			select
				cfcu.collection_id as collection_id,
				cfcu.original_language as original_language
			from 
				(
					select
						distinct on (c.id)
						c.id as collection_id,
						coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text)::date as date,
						cu.id as content_unit_id,
						cu.properties->>'original_language' as original_language
					from
						content_units as cu,
						collections as c,
						collections_content_units as ccu
					where
						ccu.content_unit_id = cu.id and
						ccu.collection_id = c.id
						%s
					order by c.id, date desc
				) as cfcu
			where
			%s
		) as collection_first_content_unit_languages
	`,
		utils.InClause("and c.type_id in", ContentTypesToContentIds(contentTypes)),
		strings.Join(orClauses, " or "))
}

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
	for _, uid := range request.Options.SkipUids {
		currentLessonUIDs = append(currentLessonUIDs, uid)
	}
	return suggester.fetchCollection(currentLessonUIDs, request.MoreItems, request.Options.Languages)
}

func (suggester *CollectionSuggester) fetchCollection(currentLessonUIDs []string, moreItems int, languages []string) ([]ContentItem, error) {
	uidsQuery := ""
	if len(currentLessonUIDs) > 0 {
		quoted := []string(nil)
		for _, uid := range currentLessonUIDs {
			quoted = append(quoted, fmt.Sprintf("'%s'", uid))
		}
		uidsQuery = fmt.Sprintf("and c.uid not in (%s)", strings.Join(quoted, ","))
	}
	query := fmt.Sprintf(`
		select
			c.uid, MAX((coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text))::date) as date, MAX(cu.created_at) as last_created_at
		from
			collections as c,
			content_units as cu,
			collections_content_units as ccu,
			%s
		where 
			c.id = ccu.collection_id
			and cu.id = ccu.content_unit_id
			and cu.secure = 0 AND cu.published IS TRUE
			%s
			and c.type_id = %d
			%s
		group by c.uid
		order by date desc, last_created_at desc
		limit %d;
		`,
		CollectionsByFirstUnitLanguagesTableSql([]string{suggester.contentType}, languages),
		uidsQuery,
		mdb.CONTENT_TYPE_REGISTRY.ByName[suggester.contentType].ID,
		COLLECTIONS_BY_FIRST_CONTENT_UNIT_CLAUSE,
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
		err := rows.Scan(&uid, &date, &createdAt)
		if err != nil {
			return nil, err
		}
		ret = append(ret, ContentItem{UID: uid, Date: date, CreatedAt: createdAt, ContentType: suggester.contentType, Suggester: "CollectionSuggester"})
	}
	return ret, nil
}
