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
			MakeNextContentUnitsSameSourceSuggester([]string{consts.CT_LESSON_PART}, db),
			MakePrevContentUnitsSameSourceSuggester([]string{consts.CT_LESSON_PART}, db),
			MakeLastCollectionSameSourceSuggester([]string{consts.CT_LESSONS_SERIES}, db),
			core.MakeCompletionSuggester([]core.Suggester{MakeLastLessonsSameTagSuggester(db), MakeLastLessonsSuggester(db)}),
			core.MakeCompletionSuggester([]core.Suggester{MakeLastProgramsSameTagSuggester(db), MakeLastProgramsSuggester(db)}),
			MakeLastCongressSameTagSuggester(db),
		}),
		MakeRandomContentTypesSuggester([]string{consts.CT_CLIP, consts.CT_LESSON_PART, consts.CT_VIDEO_PROGRAM_CHAPTER}, db),
	})}
}

func (recommender *Recommender) Recommend(r core.MoreRequest) ([]core.ContentItem, error) {
	return recommender.Suggester.More(r)
}
