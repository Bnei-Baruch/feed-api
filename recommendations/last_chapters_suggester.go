package recommendations

import (
	"database/sql"
	"fmt"

	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/utils"
)

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
	uids := []string{request.Options.Recommend.Uid}
	dateField := `coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text)::date`
	return fmt.Sprintf(`
		select t.type_id, t.uid, t.date, t.created_at from (
			select cu.type_id, ROW_NUMBER() OVER(PARTITION BY c.uid order by %s desc) as r, cu.uid as uid, %s as date, cu.created_at as created_at
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
		utils.InClause("cu.uid not in", uids),
		utils.InClause("cu.uid in", uids),
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
	dateField := `coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text)::date`
	return fmt.Sprintf(`
		select t.type_id, t.uid, t.date, t.created_at from (
			select cu.type_id as type_id, cu.uid as uid, %s as date, cu.created_at as created_at, ROW_NUMBER() OVER(PARTITION BY cut.tag_id order by %s desc) as r
			from content_units as cu, content_units_tags as cut
			where cu.id = cut.content_unit_id and cut.tag_id in (
				select t.id
				from content_units as cu, content_units_tags as cut, tags as t
				where t.id = cut.tag_id and cut.content_unit_id = cu.id %s
			) %s
		) as t where t.r <= %d %s
		order by t.date desc
		`,
		dateField,
		dateField,
		utils.InClause("and cu.uid in", uids),
		utils.InClause("and cu.type_id in", core.ContentTypesToContentIds([]string{consts.CT_CLIP})),
		request.MoreItems,
		utils.InClause("and t.uid not in", uids),
	)
}
