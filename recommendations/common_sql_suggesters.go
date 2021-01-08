package recommendations

import (
	"database/sql"
	"fmt"

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

type LastClipsSuggester struct {
	SqlSuggester
}

func MakeLastClipsSuggester(db *sql.DB) *LastClipsSuggester {
	return &LastClipsSuggester{SqlSuggester: SqlSuggester{db, LastContentUnitsContentTypesGenSql([]string{consts.CT_CLIP}), "LastClipsSuggester"}}
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
					cu.secure = 0 AND cu.published IS TRUE and
					cu.uid != '%s'
					%s
					%s
				order by date desc, created_at desc
				limit %d;
			`,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			contentTypesClause,
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
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
        cu.id = ccu.content_unit_id and
        cu.uid != '%s'
        %s
      order by date desc, created_at desc
      limit %d;
    `,
		DATE_FIELD,
		request.Options.Recommend.Uid,
		request.Options.Recommend.Uid,
		fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
		request.MoreItems,
	)
}

type NextContentUnitsSameSourceSuggester struct {
	SqlSuggester
}

func MakeNextContentUnitsSameSourceSuggester(contentTypes []string, db *sql.DB) *NextContentUnitsSameSourceSuggester {
	return &NextContentUnitsSameSourceSuggester{SqlSuggester: SqlSuggester{db, NextContentUnitsSameSourceGenSql(contentTypes), "NextContentUnitsSameSourceSuggester"}}
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
					(date > d.date or (date = d.date and cu.created_at > d.created_at)) and
					cu.uid != '%s'
					%s
					%s
				order by date asc, created_at asc
				limit %d;
			`,
			DATE_FIELD,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes)),
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			request.MoreItems,
		)
	}
}

type PrevContentUnitsSameSourceSuggester struct {
	SqlSuggester
}

func MakePrevContentUnitsSameSourceSuggester(contentTypes []string, db *sql.DB) *PrevContentUnitsSameSourceSuggester {
	return &PrevContentUnitsSameSourceSuggester{SqlSuggester: SqlSuggester{db, PrevContentUnitsSameSourceGenSql(contentTypes), "PrevContentUnitsSameSourceSuggester"}}
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
					(date < d.date or (date = d.date and cu.created_at < d.created_at)) and
					cu.uid != '%s'
					%s
					%s
				order by date desc, created_at desc
				limit %d;
			`,
			DATE_FIELD,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes)),
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
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
				(date < d.date or (date = d.date and cu.created_at < d.created_at)) and
				cu.uid != '%s'
				%s
			order by date desc, created_at desc
			limit %d;
		`,
		DATE_FIELD,
		DATE_FIELD,
		request.Options.Recommend.Uid,
		request.Options.Recommend.Uid,
		fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
		request.MoreItems,
	)
}

type LastContentTypesSameTagSuggester struct {
	SqlSuggester
}

func MakeLastClipsSameTagSuggester(db *sql.DB) *LastContentTypesSameTagSuggester {
	return &LastContentTypesSameTagSuggester{SqlSuggester: SqlSuggester{db, LastContentTypesSameTag([]string{consts.CT_CLIP}), "LastClipsSameTagSuggester"}}
}

func MakeLastLessonsSameTagSuggester(db *sql.DB) *core.RoundRobinSuggester {
	return core.MakeRoundRobinSuggester([]core.Suggester{
		&LastContentTypesSameTagSuggester{
			SqlSuggester: SqlSuggester{
				db,
				LastContentTypesSameTag([]string{consts.CT_LESSON_PART}),
				"LastLessonPartSameTagSuggester",
			},
		},
		&LastContentTypesSameTagSuggester{
			SqlSuggester: SqlSuggester{
				db,
				LastContentTypesSameTag([]string{consts.CT_VIRTUAL_LESSON}),
				"LastVirtualLessonSameTagSuggester",
			},
		},
		&LastContentTypesSameTagSuggester{
			SqlSuggester: SqlSuggester{
				db,
				LastContentTypesSameTag([]string{consts.CT_WOMEN_LESSON}),
				"LastWomenLessonsSameTagSuggester",
			},
		},
		&LastContentTypesSameTagSuggester{
			SqlSuggester: SqlSuggester{
				db,
				LastContentTypesSameTag([]string{consts.CT_LECTURE}),
				"LastLectureSameTagSuggester",
			},
		},
	})
}

func MakeLastProgramsSameTagSuggester(db *sql.DB) *LastContentTypesSameTagSuggester {
	return &LastContentTypesSameTagSuggester{
		SqlSuggester: SqlSuggester{
			db,
			LastContentTypesSameTag([]string{consts.CT_VIDEO_PROGRAM_CHAPTER}),
			"LastProgramsSameTagSuggester",
		},
	}
}

func LastContentTypesSameTag(contentTypes []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu,
					content_units_tags as cut,
					(select t.id as tag_id
					 from content_units as cu, content_units_tags as cut, tags as t
					 where t.id = cut.tag_id and cut.content_unit_id = cu.id and cu.uid = '%s') as d
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.id = cut.content_unit_id and
					cut.tag_id = d.tag_id and
					cu.uid != '%s' %s
					%s
				order by date desc, created_at desc
				limit %d;
			`,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes)),
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			request.MoreItems,
		)
	}
}

type LastCollectionSameSourceSuggester struct {
	SqlSuggester
}

func MakeLastCollectionSameSourceSuggester(contentTypes []string, db *sql.DB) *LastCollectionSameSourceSuggester {
	return &LastCollectionSameSourceSuggester{
		SqlSuggester: SqlSuggester{
			db,
			LastCollectionSameSource(contentTypes),
			"LastCollectionSameSourceSuggester",
		},
	}
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
					collections as c
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.uid = '%s' and
					cus.content_unit_id = cu.id and
					cus.source_id = s.id and
					s.uid = c.properties->>'source'
					%s
				order by date desc, created_at desc
				limit %d;
			`,
			COLLECTION_DATE_FIELD,
			request.Options.Recommend.Uid,
			utils.InClause("and c.type_id in", core.ContentTypesToContentIds(contentTypes)),
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
					cu.id = ccu.content_unit_id and
					cu.uid != '%s' %s
					%s
				order by date desc, created_at desc
				limit %d;
			`,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			request.Options.Recommend.Uid,
			utils.InClause("and c.type_id in", core.ContentTypesToContentIds(contentTypes)),
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			request.MoreItems,
		)
	}
}

type RandomContentTypesSuggester struct {
	SqlSuggester
}

func MakeRandomContentTypesSuggester(contentTypes []string, db *sql.DB) *RandomContentTypesSuggester {
	return &RandomContentTypesSuggester{
		SqlSuggester: SqlSuggester{
			db,
			RandomContentTypes(contentTypes),
			"RandomContentTypes",
		},
	}
}

func RandomContentTypes(contentTypes []string) GenerateSqlFunc {
	return func(request core.MoreRequest) string {
		if request.Options.Recommend.Uid == "" {
			return ""
		}
		return fmt.Sprintf(`
				select cu.type_id, cu.uid as uid, %s as date, cu.created_at as created_at
				from
					content_units as cu
				where
					cu.secure = 0 AND cu.published IS TRUE and
					cu.uid != '%s' %s
					%s
				order by random()
				limit %d;
			`,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes)),
			fmt.Sprintf(FILTER_LESSON_PREP, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID),
			request.MoreItems,
		)
	}
}
