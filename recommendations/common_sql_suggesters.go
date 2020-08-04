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
				cu.id = ccu.content_unit_id
			order by date desc, created_at desc
			limit %d;
		`,
		DATE_FIELD,
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
				(date < d.date or (date = d.date and cu.created_at < d.created_at))
			order by date desc, created_at desc
			limit %d;
		`,
		DATE_FIELD,
		DATE_FIELD,
		request.Options.Recommend.Uid,
		request.MoreItems,
	)
}

type LastClipsSameTagSuggester struct {
	SqlSuggester
}

func MakeLastClipsSameTagSuggester(db *sql.DB) *LastClipsSameTagSuggester {
	return &LastClipsSameTagSuggester{SqlSuggester: SqlSuggester{db, LastClipsSameTag, "LastClipsSameTagSuggester"}}
}

func LastClipsSameTag(request core.MoreRequest) string {
	if request.Options.Recommend.Uid == "" {
		return ""
	}
	uids := []string{request.Options.Recommend.Uid}
	return fmt.Sprintf(`
		select t.type_id, t.uid, t.date, t.created_at from (
			select cu.type_id as type_id, cu.uid as uid, %s as date, cu.created_at as created_at, ROW_NUMBER() OVER(PARTITION BY cut.tag_id order by %s desc) as r
			from content_units as cu, content_units_tags as cut
			where cu.id = cut.content_unit_id and cut.tag_id in (
				select t.id
				from content_units as cu, content_units_tags as cut, tags as t
				where t.id = cut.tag_id and cut.content_unit_id = cu.id %s and cu.secure = 0 AND cu.published IS TRUE
			) %s and cu.secure = 0 AND cu.published IS TRUE
		) as t where t.r <= %d %s
		order by t.date desc, t.created_at desc
		`,
		DATE_FIELD,
		DATE_FIELD,
		utils.InClause("and cu.uid in", uids),
		utils.InClause("and cu.type_id in", core.ContentTypesToContentIds([]string{consts.CT_CLIP})),
		request.MoreItems,
		utils.InClause("and t.uid not in", uids),
	)
}
