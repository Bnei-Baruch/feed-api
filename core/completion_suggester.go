package core

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Given slice of Suggesters will take the first one, if not enough
// suggestions will take the second one, ect...
// The OriginalOrder here is the index of the suggester, e.g.,
// 0, 0, 0, 1, 1, 2, 2, 2, 2 (3 items from suggeter #0, 2 from #1, 4 from #4)
type CompletionSuggester struct {
	suggesters []Suggester
}

func MakeCompletionSuggester(suggesters []Suggester) *CompletionSuggester {
	return &CompletionSuggester{suggesters}
}

func init() {
	RegisterSuggester("CompletionSuggester", func(db *sql.DB) Suggester { return MakeCompletionSuggester([]Suggester(nil)) })
}

func (suggester *CompletionSuggester) MarshalSpec() (SuggesterSpec, error) {
	log.Infof("Marshal!")
	var specs []SuggesterSpec
	for i := range suggester.suggesters {
		if spec, err := suggester.suggesters[i].MarshalSpec(); err != nil {
			return SuggesterSpec{}, err
		} else {
			specs = append(specs, spec)
		}
	}
	return SuggesterSpec{
		Name:  "CompletionSuggester",
		Specs: specs,
	}, nil
}

func (suggester *CompletionSuggester) UnmarshalSpec(db *sql.DB, spec SuggesterSpec) error {
	if spec.Name != "CompletionSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'CompletionSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("CompletionSuggester expected to have no arguments.")
	}
	if len(spec.Specs) == 0 {
		return errors.New("CompletionSuggester expected to have some suggesters, got 0.")
	}
	for i := range spec.Specs {
		if newSuggester, err := MakeSuggesterFromName(db, spec.Specs[i].Name); err != nil {
			return err
		} else {
			if err := newSuggester.UnmarshalSpec(db, spec.Specs[i]); err != nil {
				return err
			}
			suggester.suggesters = append(suggester.suggesters, newSuggester)
		}
	}
	return nil
}

func (suggester *CompletionSuggester) More(request MoreRequest) ([]ContentItem, error) {
	if len(request.CurrentFeed) >= request.MoreItems {
		return request.CurrentFeed, nil
	}
	allItems := make([][]ContentItem, len(suggester.suggesters))
	suggestedSize := 0
	for i := range request.CurrentFeed {
		order := request.CurrentFeed[i].OriginalOrder[0]
		request.CurrentFeed[i].OriginalOrder = request.CurrentFeed[i].OriginalOrder[1:]
		allItems[order] = append(allItems[order], request.CurrentFeed[i])
		suggestedSize++
	}
	for i, s := range suggester.suggesters {
		suggesterRequest := request
		suggesterRequest.CurrentFeed = allItems[i]
		suggestedSize -= len(allItems[i])
		err := error(nil)
		if allItems[i], err = s.More(suggesterRequest); err != nil {
			return nil, err
		} else {
			for j := range allItems[i] {
				allItems[i][j].OriginalOrder = append([]int64{int64(i)}, allItems[i][j].OriginalOrder...)
			}
			suggestedSize += len(allItems[i])
		}
		if suggestedSize >= request.MoreItems {
			break
		}
	}
	completion := []ContentItem(nil)
	for i := range allItems {
		j := 0
		for len(completion) < request.MoreItems && j < len(allItems[i]) {
			completion = append(completion, allItems[i][j])
			j++
		}
	}
	//fmt.Printf("Completion:\n")
	//for i, ci := range completion {
	//	fmt.Printf("%d: %+v\n", i+1, ci)
	//}
	return completion, nil
}
