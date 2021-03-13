package recommendations

import (
	"database/sql"

	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/core"
)

type Recommender struct {
	Suggester core.Suggester
}

func MakeRecommender(db *sql.DB) *Recommender {
	return &Recommender{Suggester: core.MakeCompletionSuggester([]core.Suggester{
		core.MakeRoundRobinSuggester([]core.Suggester{
			core.MakeCompletionSuggester([]core.Suggester{MakeLastClipsSameTagSuggester(db), MakeLastClipsSuggester(db)}),
			MakeLastContentUnitsSameCollectionSuggester(db),
			MakePrevContentUnitsSameCollectionSuggester(db),
			MakeNextContentUnitsSameSourceSuggester(db, []string{consts.CT_LESSON_PART}),
			MakePrevContentUnitsSameSourceSuggester(db, []string{consts.CT_LESSON_PART}),
			MakeLastCollectionSameSourceSuggester(db, []string{consts.CT_LESSONS_SERIES}),
			core.MakeCompletionSuggester([]core.Suggester{MakeLastLessonsSameTagSuggester(db), MakeLastLessonsSuggester(db)}),
			core.MakeCompletionSuggester([]core.Suggester{MakeLastProgramsSameTagSuggester(db), MakeLastProgramsSuggester(db)}),
			MakeLastCongressSameTagSuggester(db),
		}),
		MakeRandomContentTypesSuggester(db, []string{consts.CT_CLIP, consts.CT_LESSON_PART, consts.CT_VIDEO_PROGRAM_CHAPTER}, []string(nil) /* tagUids */),
	})}
}

func (recommender *Recommender) Recommend(r core.MoreRequest) ([]core.ContentItem, error) {
	return recommender.Suggester.More(r)
}
