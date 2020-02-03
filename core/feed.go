package core

import (
	"database/sql"
)

type Feed struct {
	Suggesters []Suggester
}

func MakeFeed(db *sql.DB) *Feed {
	return &Feed{Suggesters: []Suggester{MakeMorningLessonSuggester(db)}}
}

func (f *Feed) More(r MoreRequest) ([]ContentItem, error) {
	suggestions := [][]ContentItem(nil)
	for _, suggester := range f.Suggesters {
		if s, err := suggester.More(r); err != nil {
			return nil, err
		} else {
			suggestions = append(suggestions, s)
		}
	}
	return Merge(r.CurrentFeed, suggestions)
}
