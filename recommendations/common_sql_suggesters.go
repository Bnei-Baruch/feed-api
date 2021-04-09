package recommendations

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

const DATE_FIELD = `coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text)::date`
const COLLECTION_DATE_FIELD = `coalesce(c.properties->>'start_date', c.created_at::text)::date`

// Filter out all content units which are lesson preps, e.g., part 0 which is shorter then 20 minutes (1200 seconds).
const FILTER_LESSON_PREP = `and not(cu.type_id = %d and (cu.properties->>'part')::int = 0 and (cu.properties->>'duration')::int < 1200)`

type LastContentUnitsSuggester struct {
	SqlSuggester
}

func MakeLastContentUnitsSuggester(db *sql.DB) *LastContentUnitsSuggester {
	return &LastContentUnitsSuggester{SqlSuggester: SqlSuggester{db, LastContentUnitsContentTypesGenSql([]string(nil)), "LastContentUnitsSuggester"}}
}

func (suggester *LastContentUnitsSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name}, nil
}

func (suggester *LastContentUnitsSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "LastContentUnitsSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'LastContentUnitsSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("LastContentUnitsSuggester expected to have some arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("LastContentUnitsSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	return nil
}

type LastClipsSuggester struct {
	SqlSuggester
}

func MakeLastClipsSuggester(db *sql.DB) *LastClipsSuggester {
	return &LastClipsSuggester{SqlSuggester: SqlSuggester{db, LastContentUnitsContentTypesGenSql([]string{consts.CT_CLIP}), "LastClipsSuggester"}}
}

func (suggester *LastClipsSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name}, nil
}

func (suggester *LastClipsSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "LastClipsSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'LastClipsSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("LastClipsSuggester expected to have some arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("LastClipsSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	return nil
}

type LastLessonsSuggester struct {
	SqlSuggester
}

func MakeLastLessonsSuggester(db *sql.DB) *LastLessonsSuggester {
	return &LastLessonsSuggester{SqlSuggester: SqlSuggester{
		db,
		LastContentUnitsContentTypesGenSql([]string{consts.CT_LESSON_PART, consts.CT_VIRTUAL_LESSON, consts.CT_WOMEN_LESSON}),
		"LastLessonsSuggester",
	}}
}

func (suggester *LastLessonsSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name}, nil
}

func (suggester *LastLessonsSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "LastLessonsSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'LastLessonsSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("LastLessonsSuggester expected to have some arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("LastLessonsSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	return nil
}

type LastProgramsSuggester struct {
	SqlSuggester
}

func MakeLastProgramsSuggester(db *sql.DB) *LastProgramsSuggester {
	return &LastProgramsSuggester{SqlSuggester: SqlSuggester{
		db,
		LastContentUnitsContentTypesGenSql([]string{consts.CT_VIDEO_PROGRAM_CHAPTER}),
		"LastProgramsSuggester",
	}}
}

func (suggester *LastProgramsSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name}, nil
}

func (suggester *LastProgramsSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "LastProgramsSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'LastProgramsSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("LastProgramsSuggester expected to have some arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("LastProgramsSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	return nil
}

func LastContentUnitsContentTypesGenSql(contentTypes []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		contentTypesClause := ""
		if len(contentTypes) != 0 {
			contentTypesClause = utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes))
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu
				where
					cu.secure = 0 AND cu.published IS TRUE
					%s %s %s %s %s
				order by date desc, created_at desc
				limit %d;
			`,
			DATE_FIELD,
			utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			contentTypesClause,
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			core.FilterByLanguageSql(request.Options.Languages),
			core.CONTENT_UNIT_PERSON_RAV,
			request.MoreItems,
		)
	}
}

type LastContentUnitsSameCollectionSuggester struct {
	SqlSuggester
}

func MakeLastContentUnitsSameCollectionSuggester(db *sql.DB) *LastContentUnitsSameCollectionSuggester {
	return &LastContentUnitsSameCollectionSuggester{SqlSuggester: SqlSuggester{db, LastContentUnitsSameCollectionGenSql, "LastContentUnitsSameCollectionSuggester"}}
}

func (suggester *LastContentUnitsSameCollectionSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name}, nil
}

func (suggester *LastContentUnitsSameCollectionSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "LastContentUnitsSameCollectionSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'LastContentUnitsSameCollectionSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("LastContentUnitsSameCollectionSuggester expected to have some arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("LastContentUnitsSameCollectionSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	return nil
}

func LastContentUnitsSameCollectionGenSql(request core.MoreRequest) string {
	if request.Options.Recommend.Uid == "" {
		return ""
	}
	return fmt.Sprintf(`
      select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
      from
        collections as c,
        content_units as cu,
        collections_content_units as ccu,
        (select ccu.collection_id as collection_id
         from content_units as cu, collections_content_units as ccu
         where cu.id = ccu.content_unit_id and cu.uid = '%s') as d
      where
        cu.secure = 0 AND cu.published IS TRUE and
        c.id = d.collection_id and
        c.id = ccu.collection_id and 
        cu.id = ccu.content_unit_id
        %s %s %s %s
      order by date desc, created_at desc
      limit %d;
    `,
		DATE_FIELD,
		request.Options.Recommend.Uid,
		utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
		fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
		core.FilterByLanguageSql(request.Options.Languages),
		core.CONTENT_UNIT_PERSON_RAV,
		request.MoreItems,
	)
}

type NextContentUnitsSameSourceSuggester struct {
	SqlSuggester
	contentTypes []string
}

func MakeNextContentUnitsSameSourceSuggester(db *sql.DB, contentTypes []string) *NextContentUnitsSameSourceSuggester {
	return &NextContentUnitsSameSourceSuggester{
		SqlSuggester: SqlSuggester{db, NextContentUnitsSameSourceGenSql(contentTypes), "NextContentUnitsSameSourceSuggester"},
		contentTypes: contentTypes,
	}
}

func (suggester *NextContentUnitsSameSourceSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name, Args: suggester.contentTypes}, nil
}

func (suggester *NextContentUnitsSameSourceSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "NextContentUnitsSameSourceSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'NextContentUnitsSameSourceSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) == 0 {
		return errors.New("NextContentUnitsSameSourceSuggester expected to have some arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("NextContentUnitsSameSourceSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	suggester.contentTypes = spec.Args
	suggester.SqlSuggester.genSql = NextContentUnitsSameSourceGenSql(suggester.contentTypes)
	return nil
}

func NextContentUnitsSameSourceGenSql(contentTypes []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu,
					content_units_sources as cus,
					(select cus.source_id as source_id, cu.created_at, %s as date
					 from content_units as cu, content_units_sources as cus
					 where cu.uid = '%s' and cus.content_unit_id = cu.id) as d
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.id = cus.content_unit_id and
					cus.source_id = d.source_id and
					(date > d.date or (date = d.date and cu.created_at > d.created_at))
					%s %s %s %s %s
				order by date asc, created_at asc
				limit %d;
			`,
			DATE_FIELD,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes)),
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			core.FilterByLanguageSql(request.Options.Languages),
			core.CONTENT_UNIT_PERSON_RAV,
			request.MoreItems,
		)
	}
}

type PrevContentUnitsSameSourceSuggester struct {
	SqlSuggester
	contentTypes []string
}

func MakePrevContentUnitsSameSourceSuggester(db *sql.DB, contentTypes []string) *PrevContentUnitsSameSourceSuggester {
	return &PrevContentUnitsSameSourceSuggester{
		SqlSuggester: SqlSuggester{db, PrevContentUnitsSameSourceGenSql(contentTypes), "PrevContentUnitsSameSourceSuggester"},
		contentTypes: contentTypes,
	}
}

func (suggester *PrevContentUnitsSameSourceSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name, Args: suggester.contentTypes}, nil
}

func (suggester *PrevContentUnitsSameSourceSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "PrevContentUnitsSameSourceSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'PrevContentUnitsSameSourceSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) == 0 {
		return errors.New("PrevContentUnitsSameSourceSuggester expected to have some arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("PrevContentUnitsSameSourceSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	suggester.contentTypes = spec.Args
	suggester.SqlSuggester.genSql = PrevContentUnitsSameSourceGenSql(suggester.contentTypes)
	return nil
}

func PrevContentUnitsSameSourceGenSql(contentTypes []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu,
					content_units_sources as cus,
					(select cus.source_id as source_id, cu.created_at, %s as date
					 from content_units as cu, content_units_sources as cus
					 where cu.uid = '%s' and cus.content_unit_id = cu.id) as d
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.id = cus.content_unit_id and
					cus.source_id = d.source_id and
					(date < d.date or (date = d.date and cu.created_at < d.created_at))
					%s %s %s %s %s
				order by date desc, created_at desc
				limit %d;
			`,
			DATE_FIELD,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes)),
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			core.FilterByLanguageSql(request.Options.Languages),
			core.CONTENT_UNIT_PERSON_RAV,
			request.MoreItems,
		)
	}
}

type PrevContentUnitsSameCollectionSuggester struct {
	SqlSuggester
}

func MakePrevContentUnitsSameCollectionSuggester(db *sql.DB) *PrevContentUnitsSameCollectionSuggester {
	return &PrevContentUnitsSameCollectionSuggester{SqlSuggester: SqlSuggester{db, PrevContentUnitsSameCollectionGenSql, "PrevContentUnitsSameCollectionSuggester"}}
}

func (suggester *PrevContentUnitsSameCollectionSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name}, nil
}

func (suggester *PrevContentUnitsSameCollectionSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "PrevContentUnitsSameCollectionSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'PrevContentUnitsSameCollectionSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("PrevContentUnitsSameCollectionSuggester expected to have some arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("PrevContentUnitsSameCollectionSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	return nil
}

func PrevContentUnitsSameCollectionGenSql(request core.MoreRequest) string {
	if request.Options.Recommend.Uid == "" {
		return ""
	}
	return fmt.Sprintf(`
			select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
			from
				collections as c,
				content_units as cu,
				collections_content_units as ccu,
				(select ccu.collection_id as collection_id, cu.created_at, %s as date
				 from content_units as cu, collections_content_units as ccu
				 where cu.id = ccu.content_unit_id and cu.uid = '%s') as d
			where
				cu.secure = 0 AND cu.published IS TRUE and
				c.id = d.collection_id and
				c.id = ccu.collection_id and 
				cu.id = ccu.content_unit_id and
				(date < d.date or (date = d.date and cu.created_at < d.created_at))
				%s %s %s %s
			order by date desc, created_at desc
			limit %d;
		`,
		DATE_FIELD,
		DATE_FIELD,
		request.Options.Recommend.Uid,
		utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
		fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
		core.FilterByLanguageSql(request.Options.Languages),
		core.CONTENT_UNIT_PERSON_RAV,
		request.MoreItems,
	)
}

type LastContentTypesSameTagSuggester struct {
	SqlSuggester
	contentTypes []string
}

func MakeLastContentTypesSameTagSuggester(db *sql.DB, contentTypes []string) *LastContentTypesSameTagSuggester {
	return &LastContentTypesSameTagSuggester{
		SqlSuggester: SqlSuggester{db, ContentTypesSameTagGen(contentTypes, core.Last), "LastContentTypesSameTagSuggester"},
		contentTypes: contentTypes,
	}
}

func (suggester *LastContentTypesSameTagSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name, Args: suggester.contentTypes}, nil
}

func (suggester *LastContentTypesSameTagSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != suggester.SqlSuggester.Name {
		return errors.New(fmt.Sprintf("Expected suggester name to be: '%s', got: '%s'.", suggester.Name, spec.Name))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("%s expected to have no suggesters, got %d.", suggester.Name, len(spec.Specs)))
	}
	if suggester.Name == "LastContentTypesSameTagSuggester" {
		if len(spec.Args) == 0 {
			return errors.New("LastContentTypesSameTagSuggester expected to have some arguments.")
		}
		suggester.contentTypes = spec.Args
		suggester.SqlSuggester.genSql = ContentTypesSameTagGen(suggester.contentTypes, core.Last)
	} else {
		if len(spec.Args) != 0 {
			return errors.New(fmt.Sprintf("%s expected to have some arguments.", suggester.Name))
		}
	}
	return nil
}

func MakeLastClipsSameTagSuggester(db *sql.DB) *LastContentTypesSameTagSuggester {
	return &LastContentTypesSameTagSuggester{
		SqlSuggester: SqlSuggester{db, ContentTypesSameTagGen([]string{consts.CT_CLIP}, core.Last), "LastClipsSameTagSuggester"},
	}
}

type LastLessonsSameTagSuggester struct {
	core.RoundRobinSuggester
}

func MakeLastLessonsSameTagSuggester(db *sql.DB) *LastLessonsSameTagSuggester {
	return &LastLessonsSameTagSuggester{*core.MakeRoundRobinSuggester([]core.Suggester{
		&LastContentTypesSameTagSuggester{
			SqlSuggester: SqlSuggester{
				db,
				ContentTypesSameTagGen([]string{consts.CT_LESSON_PART}, core.Last),
				"LastLessonPartSameTagSuggester",
			},
		},
		&LastContentTypesSameTagSuggester{
			SqlSuggester: SqlSuggester{
				db,
				ContentTypesSameTagGen([]string{consts.CT_VIRTUAL_LESSON}, core.Last),
				"LastVirtualLessonSameTagSuggester",
			},
		},
		&LastContentTypesSameTagSuggester{
			SqlSuggester: SqlSuggester{
				db,
				ContentTypesSameTagGen([]string{consts.CT_WOMEN_LESSON}, core.Last),
				"LastWomenLessonsSameTagSuggester",
			},
		},
		&LastContentTypesSameTagSuggester{
			SqlSuggester: SqlSuggester{
				db,
				ContentTypesSameTagGen([]string{consts.CT_LECTURE}, core.Last),
				"LastLectureSameTagSuggester",
			},
		},
	})}
}

func (suggester *LastLessonsSameTagSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: "LastLessonsSameTagSuggester"}, nil
}

func (suggester *LastLessonsSameTagSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "LastLessonsSameTagSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'LastLessonsSameTagSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("LastLessonsSameTagSuggester expected to have 0 arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("LastLessonsSameTagSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	return nil
}

func MakeLastProgramsSameTagSuggester(db *sql.DB) *LastContentTypesSameTagSuggester {
	return &LastContentTypesSameTagSuggester{
		SqlSuggester: SqlSuggester{
			db,
			ContentTypesSameTagGen([]string{consts.CT_VIDEO_PROGRAM_CHAPTER}, core.Last),
			"LastProgramsSameTagSuggester",
		},
	}
}

func ContentTypesSameTagGen(contentTypes []string, orderSelector core.OrderSelectorEnum) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		contentTypesSql := ""
		if len(contentTypes) > 0 {
			contentTypesSql = utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes))
		}
		orderSelectorWhere := ""
		if orderSelector == core.Next {
			orderSelectorWhere = "and (date > d.date or (date = d.date and cu.created_at > d.created_at))"
		} else if orderSelector == core.Prev {
			orderSelectorWhere = "and (date < d.date or (date = d.date and cu.created_at < d.created_at))"
		}
		orderSelectorOrderBy := "order by date desc, created_at desc"
		if orderSelector == core.Next {
			orderSelectorOrderBy = "order by date asc, created_at asc"
		} else if orderSelector == core.Rand {
			orderSelectorOrderBy = "order by random()"
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu,
					content_units_tags as cut,
					(select t.id as tag_id, cu.created_at, %s as date
					 from content_units as cu, content_units_tags as cut, tags as t
					 where t.id = cut.tag_id and cut.content_unit_id = cu.id and cu.uid = '%s') as d
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.id = cut.content_unit_id and
					cut.tag_id = d.tag_id
					%s %s %s %s %s %s
				%s
				limit %d;
			`,
			DATE_FIELD,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			orderSelectorWhere,
			utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			contentTypesSql,
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			core.FilterByLanguageSql(request.Options.Languages),
			core.CONTENT_UNIT_PERSON_RAV,
			orderSelectorOrderBy,
			request.MoreItems,
		)
	}
}

type LastCollectionSameSourceSuggester struct {
	SqlSuggester
	contentTypes []string
}

func MakeLastCollectionSameSourceSuggester(db *sql.DB, contentTypes []string) *LastCollectionSameSourceSuggester {
	return &LastCollectionSameSourceSuggester{
		SqlSuggester: SqlSuggester{
			db,
			LastCollectionSameSource(contentTypes),
			"LastCollectionSameSourceSuggester",
		},
		contentTypes: contentTypes,
	}
}

func (suggester *LastCollectionSameSourceSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name, Args: suggester.contentTypes}, nil
}

func (suggester *LastCollectionSameSourceSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "LastCollectionSameSourceSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'LastCollectionSameSourceSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) == 0 {
		return errors.New("LastCollectionSameSourceSuggester expected to have some arguments.")
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("LastCollectionSameSourceSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	suggester.contentTypes = spec.Args
	suggester.SqlSuggester.genSql = LastCollectionSameSource(suggester.contentTypes)
	return nil
}

func LastCollectionSameSource(contentTypes []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		return fmt.Sprintf(`
				select c.type_id, c.uid as uid, %s as date, c.created_at as created_at
				from
					content_units as cu,
					content_units_sources as cus,
					sources as s,
					collections as c,
					%s
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.uid = '%s' and
					cus.content_unit_id = cu.id and
					cus.source_id = s.id and
					s.uid = c.properties->>'source'
					%s %s %s
				order by date desc, created_at desc
				limit %d;
			`,
			COLLECTION_DATE_FIELD,
			core.CollectionsByFirstUnitLanguagesTableSql(contentTypes, request.Options.Languages),
			request.Options.Recommend.Uid,
			utils.InClause("and c.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			utils.InClause("and c.type_id in", core.ContentTypesToContentIds(contentTypes)),
			core.COLLECTIONS_BY_FIRST_CONTENT_UNIT_CLAUSE,
			request.MoreItems,
		)
	}
}

func MakeLastCongressSameTagSuggester(db *sql.DB) *LastContentTypesSameTagSuggester {
	return &LastContentTypesSameTagSuggester{
		SqlSuggester: SqlSuggester{
			db,
			LastCollectionContentTypesSameTag([]string{consts.CT_CONGRESS}),
			"LastCongressSameTagSuggester",
		},
	}
}

func LastCollectionContentTypesSameTag(contentTypes []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu,
					content_units_tags as cut,
					collections as c,
					collections_content_units ccu,
					(select t.id as tag_id
					 from content_units as cu, content_units_tags as cut, tags as t
					 where t.id = cut.tag_id and cut.content_unit_id = cu.id and cu.uid = '%s') as d
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.id = cut.content_unit_id and
					cut.tag_id = d.tag_id and
					c.id = ccu.collection_id and
					cu.id = ccu.content_unit_id
					%s %s %s %s %s
				order by date desc, created_at desc
				limit %d;
			`,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			utils.InClause("and c.type_id in", core.ContentTypesToContentIds(contentTypes)),
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			core.FilterByLanguageSql(request.Options.Languages),
			core.CONTENT_UNIT_PERSON_RAV,
			request.MoreItems,
		)
	}
}

type RandomContentTypesSuggester struct {
	SqlSuggester
	contentTypes []string
	tagUids      []string
}

func MakeRandomContentTypesSuggester(db *sql.DB, contentTypes []string, tagUids []string) *RandomContentTypesSuggester {
	return &RandomContentTypesSuggester{
		SqlSuggester: SqlSuggester{
			db,
			RandomContentTypes(contentTypes, tagUids),
			"RandomContentTypesSuggester",
		},
		contentTypes: contentTypes,
		tagUids:      tagUids,
	}
}

func (suggester *RandomContentTypesSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name, Args: suggester.contentTypes, SecondArgs: suggester.tagUids}, nil
}

func (suggester *RandomContentTypesSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "RandomContentTypesSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'RandomContentTypesSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("RandomContentTypesSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	suggester.contentTypes = spec.Args
	suggester.tagUids = spec.SecondArgs
	suggester.SqlSuggester.genSql = RandomContentTypes(suggester.contentTypes, suggester.tagUids)
	return nil
}

func RandomContentTypes(contentTypes []string, tagUids []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		contentTypesSql := ""
		if len(contentTypes) > 0 {
			contentTypesSql = utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes))
		}
		tagsFromSql := ""
		tagsWhereSql := ""
		if len(tagUids) > 0 {
			tagsFromSql = ", content_units_tags as cut, tags as t"
			tagsWhereSql = utils.InClause("and cu.id = cut.content_unit_id and cut.tag_id = t.id and t.uid in", tagUids)
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu
					%s
				where
					cu.secure = 0 AND cu.published IS TRUE
					%s %s %s %s %s %s
				order by random()
				limit %d;
			`,
			DATE_FIELD,
			tagsFromSql,
			utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			contentTypesSql,
			tagsWhereSql,
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			core.FilterByLanguageSql(request.Options.Languages),
			core.CONTENT_UNIT_PERSON_RAV,
			request.MoreItems,
		)
	}
}

type RandomContentUnitsSameSourceSuggester struct {
	SqlSuggester
	contentTypes []string
	tagUids      []string
}

func MakeRandomContentUnitsSameSourceSuggester(db *sql.DB, contentTypes []string, tagUids []string) *RandomContentUnitsSameSourceSuggester {
	return &RandomContentUnitsSameSourceSuggester{
		SqlSuggester: SqlSuggester{db, RandomContentUnitsSameSourceGenSql(contentTypes, tagUids), "RandomContentUnitsSameSourceSuggester"},
		contentTypes: contentTypes,
		tagUids:      tagUids,
	}
}

func (suggester *RandomContentUnitsSameSourceSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name, Args: suggester.contentTypes, SecondArgs: suggester.tagUids}, nil
}

func (suggester *RandomContentUnitsSameSourceSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "RandomContentUnitsSameSourceSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'RandomContentUnitsSameSourceSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("RandomContentUnitsSameSourceSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	suggester.contentTypes = spec.Args
	suggester.tagUids = spec.SecondArgs
	suggester.SqlSuggester.genSql = RandomContentUnitsSameSourceGenSql(suggester.contentTypes, suggester.tagUids)
	return nil
}

func RandomContentUnitsSameSourceGenSql(contentTypes []string, tagUids []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		contentTypesSql := ""
		if len(contentTypes) > 0 {
			contentTypesSql = utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes))
		}
		tagsFromSql := ""
		tagsWhereSql := ""
		if len(tagUids) > 0 {
			tagsFromSql = ", content_units_tags as cut, tags as t"
			tagsWhereSql = utils.InClause("and cu.id = cut.content_unit_id and cut.tag_id = t.id and t.uid in", tagUids)
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu,
					content_units_sources as cus
					%s
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.id = cus.content_unit_id and
					cus.source_id in (
						select cus.source_id
						from content_units as cu, content_units_sources as cus
						where cu.uid = '%s' and cus.content_unit_id = cu.id
					)
					%s %s %s %s %s %s
				order by random()
				limit %d;
			`,
			DATE_FIELD,
			tagsFromSql,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			contentTypesSql,
			tagsWhereSql,
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			core.FilterByLanguageSql(request.Options.Languages),
			core.CONTENT_UNIT_PERSON_RAV,
			request.MoreItems,
		)
	}
}

type ContentTypesSameTagSuggester struct {
	SqlSuggester
	contentTypes  []string
	orderSelector core.OrderSelectorEnum
}

func MakeContentTypesSameTagSuggester(db *sql.DB, contentTypes []string, orderSelector core.OrderSelectorEnum) *ContentTypesSameTagSuggester {
	return &ContentTypesSameTagSuggester{
		SqlSuggester:  SqlSuggester{db, ContentTypesSameTagGen(contentTypes, orderSelector), "ContentTypesSameTagSuggester"},
		contentTypes:  contentTypes,
		orderSelector: orderSelector,
	}
}

func (suggester *ContentTypesSameTagSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name, Args: suggester.contentTypes, OrderSelector: suggester.orderSelector}, nil
}

func (suggester *ContentTypesSameTagSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "ContentTypesSameTagSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'ContentTypesSameTagSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("ContentTypesSameTagSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	suggester.contentTypes = spec.Args
	suggester.orderSelector = spec.OrderSelector
	suggester.SqlSuggester.genSql = ContentTypesSameTagGen(suggester.contentTypes, suggester.orderSelector)
	return nil
}

type ContentUnitCollectionSuggester struct {
	SqlSuggester
	contentTypes []string
}

func MakeContentUnitCollectionSuggester(db *sql.DB, contentTypes []string) *ContentUnitCollectionSuggester {
	return &ContentUnitCollectionSuggester{
		SqlSuggester: SqlSuggester{db, ContentUnitCollectionGen(contentTypes), "ContentUnitCollectionSuggester"},
		contentTypes: contentTypes,
	}
}

func (suggester *ContentUnitCollectionSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name, Args: suggester.contentTypes}, nil
}

func (suggester *ContentUnitCollectionSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "ContentUnitCollectionSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'ContentUnitCollectionSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("ContentUnitCollectionSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	suggester.contentTypes = spec.Args
	suggester.SqlSuggester.genSql = ContentUnitCollectionGen(suggester.contentTypes)
	return nil
}

func ContentUnitCollectionGen(contentTypes []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		return fmt.Sprintf(`
				select c.type_id, c.uid as uid, %s as date, c.created_at as created_at
				from
					content_units as cu,
					collections_content_units as ccu,
					collections as c
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.uid = '%s' and
					ccu.content_unit_id = cu.id and
					c.id = ccu.collection_id
					%s %s
				order by date desc, created_at desc
				limit %d;
			`,
			COLLECTION_DATE_FIELD,
			request.Options.Recommend.Uid,
			utils.InClause("and c.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			utils.InClause("and c.type_id in", core.ContentTypesToContentIds(contentTypes)),
			request.MoreItems,
		)
	}
}

type ContentUnitsSuggester struct {
	SqlSuggester
	filters       []core.SuggesterFilter
	orderSelector core.OrderSelectorEnum
}

func MakeContentUnitsSuggester(db *sql.DB, filters []core.SuggesterFilter, orderSelector core.OrderSelectorEnum) *ContentUnitsSuggester {
	return &ContentUnitsSuggester{
		SqlSuggester:  SqlSuggester{db, ContentUnitsSqlGen(filters, orderSelector), "NewContentUnitsSuggester"},
		filters:       filters,
		orderSelector: orderSelector,
	}
}

func (suggester *ContentUnitsSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: suggester.Name, Filters: suggester.filters, OrderSelector: suggester.orderSelector}, nil
}

func (suggester *ContentUnitsSuggester) UnmarshalSpec(db *sql.DB, spec core.SuggesterSpec) error {
	if spec.Name != "NewContentUnitsSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'NewContentUnitsSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("NewContentUnitsSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	suggester.filters = spec.Filters
	suggester.orderSelector = spec.OrderSelector
	suggester.SqlSuggester.genSql = ContentUnitsSqlGen(suggester.filters, suggester.orderSelector)
	return nil
}

func ContentUnitsSqlGen(filters []core.SuggesterFilter, orderSelector core.OrderSelectorEnum) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		filtersFrom := []string(nil)
		filtersWhere := []string(nil)
		collectionRelationAdded := false
		addCollectionRelation := func() {
			if !collectionRelationAdded {
				filtersFrom = append(filtersFrom, `,
						collections as c,
						collections_content_units as ccu
					`)
				filtersWhere = append(filtersWhere, `and
						c.id = ccu.collection_id and
						cu.id = ccu.content_unit_id
					`)
				collectionRelationAdded = true
			}
		}
		sourcesRelationAdded := false
		addSourcesRelation := func() {
			if !sourcesRelationAdded {
				filtersFrom = append(filtersFrom, ", content_units_sources as cus")
				sourcesRelationAdded = true
			}
		}
		tagsRelationAdded := false
		addTagsRelation := func() {
			if !tagsRelationAdded {
				filtersFrom = append(filtersFrom, ", content_units_tags as cut")
				tagsRelationAdded = true
			}
		}
		for _, filter := range filters {
			switch filter.FilterSelector {
			case core.UnitContentTypes:
				filtersWhere = append(filtersWhere, utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(filter.Args)))
			case core.CollectionContentTypes:
				addCollectionRelation()
				filtersWhere = append(filtersWhere, utils.InClause("and c.type_id in", core.ContentTypesToContentIds(filter.Args)))
			case core.Tags:
				addTagsRelation()
				filtersFrom = append(filtersFrom, ", tags as t")
				filtersWhere = append(filtersWhere, utils.InClause(`and
					cut.content_unit_id = cu.id and
					cut.tag_id = t.id and
					t.uid in`, filter.Args))
			case core.Sources:
				addSourcesRelation()
				filtersFrom = append(filtersFrom, ", sources as s")
				filtersWhere = append(filtersWhere, utils.InClause(`and 
					cus.content_unit_id = cu.id and
					cus.source_id = s.id and
					s.uid in`, filter.Args))
			case core.Collections:
				addCollectionRelation()
				filtersWhere = append(filtersWhere, utils.InClause("and c.uid in", filter.Args))
			case core.SameTag:
				addTagsRelation()
				filtersFrom = append(filtersFrom, fmt.Sprintf(`,
					(select t.id as tag_id, cu.created_at
					 from content_units as cu, content_units_tags as cut, tags as t
					 where t.id = cut.tag_id and cut.content_unit_id = cu.id and cu.uid = '%s') as same_tag
				`, request.Options.Recommend.Uid))
				filtersWhere = append(filtersWhere, `and
					cu.id = cut.content_unit_id and
					cut.tag_id = same_tag.tag_id
				`)
			case core.SameCollection:
				addCollectionRelation()
				filtersFrom = append(filtersFrom, fmt.Sprintf(`,
					(select ccu.collection_id as collection_id
					 from content_units as cu, collections_content_units as ccu
					 where cu.id = ccu.content_unit_id and cu.uid = '%s') as same_collection
				`, request.Options.Recommend.Uid))
				filtersWhere = append(filtersWhere, "and c.id = same_collection.collection_id")
			case core.SameSource:
				addSourcesRelation()
				filtersFrom = append(filtersFrom, fmt.Sprintf(`,
					(select cus.source_id as source_id, cu.created_at
					 from content_units as cu, content_units_sources as cus
					 where cu.uid = '%s' and cus.content_unit_id = cu.id) as same_source
				`, request.Options.Recommend.Uid))
				filtersWhere = append(filtersWhere, `and
					cu.id = cus.content_unit_id and
					cus.source_id = same_source.source_id
				`)
			default:
				log.Errorf("Did not expect filter selector enum %d", filter.FilterSelector)
			}
		}
		if orderSelector == core.Next || orderSelector == core.Prev {
			filtersFrom = append(filtersFrom, fmt.Sprintf(`,
				(select cu.created_at, %s as date
				 from content_units as cu
				 where cu.uid = '%s') as d
			`, DATE_FIELD, request.Options.Recommend.Uid))
			if orderSelector == core.Next {
				filtersWhere = append(filtersWhere, "and (date > d.date or (date = d.date and cu.created_at > d.created_at))")
			} else {
				filtersWhere = append(filtersWhere, "and (date < d.date or (date = d.date and cu.created_at < d.created_at))")
			}
		}
		orderSelectorOrderBy := "order by date desc, created_at desc"
		if orderSelector == core.Next {
			orderSelectorOrderBy = "order by date asc, created_at asc"
		} else if orderSelector == core.Rand {
			orderSelectorOrderBy = "order by random()"
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu
					%s
				where
					cu.secure = 0 AND cu.published IS TRUE
					%s %s %s %s %s
				%s
				limit %d;
			`,
			DATE_FIELD,
			strings.Join(filtersFrom, " "),
			utils.InClause("and cu.uid not in", append(request.Options.SkipUids, request.Options.Recommend.Uid)),
			strings.Join(filtersWhere, " "),
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			core.FilterByLanguageSql(request.Options.Languages),
			core.CONTENT_UNIT_PERSON_RAV,
			orderSelectorOrderBy,
			request.MoreItems,
		)
	}
}

func init() {
	core.RegisterSuggester("LastClipsSameTagSuggester", func(db *sql.DB) core.Suggester { return MakeLastClipsSameTagSuggester(db) })
	core.RegisterSuggester("LastClipsSuggester", func(db *sql.DB) core.Suggester { return MakeLastClipsSuggester(db) })
	core.RegisterSuggester("LastCollectionSameSourceSuggester", func(db *sql.DB) core.Suggester { return MakeLastCollectionSameSourceSuggester(db, []string(nil)) })
	core.RegisterSuggester("LastCongressSameTagSuggester", func(db *sql.DB) core.Suggester { return MakeLastCongressSameTagSuggester(db) })
	core.RegisterSuggester("LastContentTypesSameTagSuggester", func(db *sql.DB) core.Suggester { return MakeLastContentTypesSameTagSuggester(db, []string(nil)) })
	core.RegisterSuggester("LastContentUnitsSameCollectionSuggester", func(db *sql.DB) core.Suggester { return MakeLastContentUnitsSameCollectionSuggester(db) })
	core.RegisterSuggester("LastContentUnitsSuggester", func(db *sql.DB) core.Suggester { return MakeLastContentUnitsSuggester(db) })
	core.RegisterSuggester("LastLessonsSameTagSuggester", func(db *sql.DB) core.Suggester { return MakeLastLessonsSameTagSuggester(db) })
	core.RegisterSuggester("LastLessonsSuggester", func(db *sql.DB) core.Suggester { return MakeLastLessonsSuggester(db) })
	core.RegisterSuggester("LastProgramsSameTagSuggester", func(db *sql.DB) core.Suggester { return MakeLastProgramsSameTagSuggester(db) })
	core.RegisterSuggester("LastProgramsSuggester", func(db *sql.DB) core.Suggester { return MakeLastProgramsSuggester(db) })
	core.RegisterSuggester("NextContentUnitsSameSourceSuggester", func(db *sql.DB) core.Suggester { return MakeNextContentUnitsSameSourceSuggester(db, []string(nil)) })
	core.RegisterSuggester("PrevContentUnitsSameCollectionSuggester", func(db *sql.DB) core.Suggester { return MakePrevContentUnitsSameCollectionSuggester(db) })
	core.RegisterSuggester("PrevContentUnitsSameSourceSuggester", func(db *sql.DB) core.Suggester { return MakePrevContentUnitsSameSourceSuggester(db, []string(nil)) })
	core.RegisterSuggester("RandomContentTypesSuggester", func(db *sql.DB) core.Suggester {
		return MakeRandomContentTypesSuggester(db, []string(nil), []string(nil))
	})
	core.RegisterSuggester("RandomContentUnitsSameSourceSuggester", func(db *sql.DB) core.Suggester {
		return MakeRandomContentUnitsSameSourceSuggester(db, []string(nil), []string(nil))
	})
	core.RegisterSuggester("ContentTypesSameTagSuggester", func(db *sql.DB) core.Suggester { return MakeContentTypesSameTagSuggester(db, []string(nil), core.Last) })
	core.RegisterSuggester("ContentUnitCollectionSuggester", func(db *sql.DB) core.Suggester { return MakeContentUnitCollectionSuggester(db, []string(nil)) })
	core.RegisterSuggester("NewContentUnitsSuggester", func(db *sql.DB) core.Suggester {
		return MakeContentUnitsSuggester(db, []core.SuggesterFilter(nil), core.Last)
	})
}
