package core

import (
	"database/sql"
	"sort"

	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type Feed struct {
	Suggester Suggester
}

func MakeFeed(db *sql.DB) *Feed {
	return &Feed{Suggester: MakeSortSuggester(MakeRoundRobinSuggester([]Suggester{
		// 1. Morning lesson.
		MakeCollectionSuggester(db, consts.CT_DAILY_LESSON),
		// 2. Additional lessons.
		MakeRoundRobinSuggester([]Suggester{
			MakeContentUnitsSuggester(db, []string{consts.CT_LECTURE}),
			MakeContentUnitsSuggester(db, []string{consts.CT_VIRTUAL_LESSON}),
			MakeContentUnitsSuggester(db, []string{consts.CT_WOMEN_LESSON}),
			MakeContentUnitsSuggester(db, []string{consts.CT_EVENT_PART}),
		}),
		// 3. TODO: Twitter.
		// 4. Programs.
		MakeContentUnitsSuggester(db, []string{
			consts.CT_VIDEO_PROGRAM_CHAPTER,
		}),
		// 5. Blog (Article?, Publication?).
		MakeRoundRobinSuggester([]Suggester{
			MakeContentUnitsSuggester(db, []string{consts.CT_BLOG_POST}),
			MakeContentUnitsSuggester(db, []string{consts.CT_ARTICLE}),
			MakeContentUnitsSuggester(db, []string{consts.CT_PUBLICATION}),
		}),
		// 6. Yeshivat + Mean.
		MakeRoundRobinSuggester([]Suggester{
			MakeContentUnitsSuggester(db, []string{consts.CT_FRIENDS_GATHERING}),
			MakeContentUnitsSuggester(db, []string{consts.CT_MEAL}),
		}),
		// 7. Clip.
		MakeContentUnitsSuggester(db, []string{
			consts.CT_CLIP,
		}),
	}))}
}

func (f *Feed) More(r MoreRequest) ([]ContentItem, error) {
	return f.Suggester.More(r)
}

func Merge(r MoreRequest, suggestions [][]ContentItem) ([]ContentItem, error) {
	mergedFeed := append([]ContentItem(nil), r.CurrentFeed...)
	for _, s := range suggestions {
		mergedFeed = append(mergedFeed, s...)
	}
	sort.SliceStable(mergedFeed, func(i, j int) bool {
		return mergedFeed[i].CreatedAt.After(mergedFeed[j].CreatedAt)
	})
	return mergedFeed[0:utils.MinInt(len(r.CurrentFeed)+r.MoreItems, len(mergedFeed))], nil
}
