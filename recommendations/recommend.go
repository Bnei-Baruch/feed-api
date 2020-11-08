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
			MakeLastClipsSameTagSuggester(db),
			MakeLastContentUnitsSuggester(db),
			MakePrevContentUnitsSuggester(db),
			MakeLastLessonsSameTagSuggester(db),
			MakeLastProgramsSameTagSuggester(db),
			MakeLastCongressSameTagSuggester(db),
		}),
		MakeRandomContentTypesSuggester([]string{consts.CT_CLIP, consts.CT_LESSON_PART, consts.CT_VIDEO_PROGRAM_CHAPTER}, db),
	})}
}

func (recommender *Recommender) Recommend(r core.MoreRequest) ([]core.ContentItem, error) {
	return recommender.Suggester.More(r)
}
