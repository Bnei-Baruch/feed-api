package recommendations

import (
	"database/sql"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type Recommender struct {
	Suggesters []core.Suggester
}

func MakeRecommender(db *sql.DB) *Recommender {
	return &Recommender{Suggesters: []core.Suggester{
		core.MakeRoundRobinSuggester([]core.Suggester{
			MakeLastClipsSameTagSuggester(db),
			MakeLastContentUnitsSuggester(db),
			MakePrevContentUnitsSuggester(db),
			MakeLastLessonsSameTagSuggester(db),
			MakeLastProgramsSameTagSuggester(db),
		}),
	}}
}

func (f *Recommender) Recommend(r core.MoreRequest) ([]core.ContentItem, error) {
	suggestions := [][]core.ContentItem(nil)
	for _, suggester := range f.Suggesters {
		if s, err := suggester.More(r); err != nil {
			return nil, err
		} else {
			suggestions = append(suggestions, s)
		}
	}
	return Merge(r, suggestions)
}

func Merge(r core.MoreRequest, suggestions [][]core.ContentItem) ([]core.ContentItem, error) {
	uids := make(map[string]bool)
	mergedFeed := []core.ContentItem(nil)
	for _, s := range suggestions {
		for _, contentItem := range s {
			if _, ok := uids[contentItem.UID]; !ok {
				uids[contentItem.UID] = true
				mergedFeed = append(mergedFeed, contentItem)
			}
		}
	}
	return mergedFeed[0:utils.MinInt(len(r.CurrentFeed)+r.MoreItems, len(mergedFeed))], nil
}
