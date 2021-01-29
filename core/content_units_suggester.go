package core

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/queries"
	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/feed-api/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

const (
	CONTENT_UNIT_SUGGESTER_NAME = "ContentUnitsSuggester"
)

type ContentUnitsSuggester struct {
	db           *sql.DB
	contentTypes []string
	name         string
}

func MakeContentUnitsSuggester(db *sql.DB, contentTypes []string) *ContentUnitsSuggester {
	return &ContentUnitsSuggester{
		db:           db,
		contentTypes: contentTypes,
		name:         CONTENT_UNIT_SUGGESTER_NAME,
	}
}

func init() {
	RegisterSuggester(CONTENT_UNIT_SUGGESTER_NAME, func(db *sql.DB) Suggester { return MakeContentUnitsSuggester(db, []string(nil)) })
}

func (suggester *ContentUnitsSuggester) MarshalSpec() (SuggesterSpec, error) {
	return SuggesterSpec{
		Name: CONTENT_UNIT_SUGGESTER_NAME,
		Args: suggester.contentTypes,
	}, nil
}

func (suggester *ContentUnitsSuggester) UnmarshalSpec(db *sql.DB, spec SuggesterSpec) error {
	if spec.Name != CONTENT_UNIT_SUGGESTER_NAME {
		return errors.New(fmt.Sprintf("Expected suggester name to be: '%s', got: '%s'.", CONTENT_UNIT_SUGGESTER_NAME, spec.Name))
	}
	if len(spec.Args) == 0 {
		return errors.New(fmt.Sprintf("%s expected to have some arguments.", CONTENT_UNIT_SUGGESTER_NAME))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("%s expected to have 0 specs, got %d.", CONTENT_UNIT_SUGGESTER_NAME, len(spec.Specs)))
	}
	suggester.contentTypes = spec.Args
	return nil
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
