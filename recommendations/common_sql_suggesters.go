package recommendations

import (
	"database/sql"
	"fmt"

	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/utils"
)

const DATE_FIELD = `coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text)::date`

type LastContentUnitsSuggester struct {
	SqlSuggester
}

func MakeLastContentUnitsSuggester(db *sql.DB) *LastContentUnitsSuggester {
	return &LastContentUnitsSuggester{SqlSuggester: SqlSuggester{db, LastContentUnitsGenSql, "LastContentUnitsSuggester"}}
}

func LastContentUnitsGenSql(request core.MoreRequest) string {
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
			order by date desc, created_at desc
			limit %d;
		`,
		DATE_FIELD,
		request.Options.Recommend.Uid,
		request.Options.Recommend.Uid,
		request.MoreItems,
	)
}

type PrevContentUnitsSuggester struct {
	SqlSuggester
}

func MakePrevContentUnitsSuggester(db *sql.DB) *PrevContentUnitsSuggester {
	return &PrevContentUnitsSuggester{SqlSuggester: SqlSuggester{db, PrevContentUnitsGenSql, "PrevContentUnitsSuggester"}}
}

func PrevContentUnitsGenSql(request core.MoreRequest) string {
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
			order by date desc, created_at desc
			limit %d;
		`,
		DATE_FIELD,
		DATE_FIELD,
		request.Options.Recommend.Uid,
		request.Options.Recommend.Uid,
		request.MoreItems,
	)
}

type LastContentTypesSameTagSuggester struct {
	SqlSuggester
}

func MakeLastClipsSameTagSuggester(db *sql.DB) *LastContentTypesSameTagSuggester {
	return &LastContentTypesSameTagSuggester{SqlSuggester: SqlSuggester{db, LastContentTypesSameTag([]string{consts.CT_CLIP}), "LastClipsSameTagSuggester"}}
}

func MakeLastLessonsSameTagSuggester(db *sql.DB) *LastContentTypesSameTagSuggester {
	return &LastContentTypesSameTagSuggester{
		SqlSuggester: SqlSuggester{
			db,
			LastContentTypesSameTag([]string{consts.CT_LESSON_PART, consts.CT_VIRTUAL_LESSON, consts.CT_WOMEN_LESSON}),
			"LastLessonsSameTagSuggester",
		},
	}
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
				order by date desc, created_at desc
				limit %d;
			`,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes)),
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
				order by date desc, created_at desc
				limit %d;
			`,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			request.Options.Recommend.Uid,
			utils.InClause("and c.type_id in", core.ContentTypesToContentIds(contentTypes)),
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
				order by random()
				limit %d;
			`,
			DATE_FIELD,
			request.Options.Recommend.Uid,
			utils.InClause("and cu.type_id in", core.ContentTypesToContentIds(contentTypes)),
			request.MoreItems,
		)
	}
}
