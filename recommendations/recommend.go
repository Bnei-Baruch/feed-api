package recommendations

import (
	"github.com/Bnei-Baruch/feed-api/core"
)

type Recommender struct {
	Suggester core.Suggester
}

func MakeRecommender(suggesterContext core.SuggesterContext) (*Recommender, error) {
	if s, err := core.MakeDefaultSuggester(suggesterContext); err != nil {
		return nil, err
	} else {
		return &Recommender{s}, err
	}
}

func (recommender *Recommender) Recommend(r core.MoreRequest) ([]core.ContentItem, error) {
	return recommender.Suggester.More(r)
}
