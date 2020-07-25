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
		core.MakeSortSuggester(core.MakeRoundRobinSuggester([]core.Suggester{
			MakeLastChaptersSuggester(db),
			// Not implemented yet.
			//MakePrevChapterSuggester(db),
			// MakeSameTopicSuggester(db),
		})),
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
	mergedFeed := []core.ContentItem(nil)
	for _, s := range suggestions {
		mergedFeed = append(mergedFeed, s...)
	}
	return mergedFeed[0:utils.MinInt(len(r.CurrentFeed)+r.MoreItems, len(mergedFeed))], nil
}
