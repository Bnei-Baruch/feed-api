package core

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

// Sort items by date and then by create at time.
type SortSuggester struct {
	suggester Suggester
}

func MakeSortSuggester(suggester Suggester) *SortSuggester {
	return &SortSuggester{
		suggester: suggester,
	}
}

func init() {
	RegisterSuggester("SortSuggester", func(suggesterContext SuggesterContext) Suggester { return MakeSortSuggester(nil) })
}

func (suggester *SortSuggester) MarshalSpec() (SuggesterSpec, error) {
	if spec, err := suggester.suggester.MarshalSpec(); err != nil {
		return SuggesterSpec{}, err
	} else {
		return SuggesterSpec{
			Name:  "SortSuggester",
			Specs: []SuggesterSpec{spec},
		}, nil
	}
}

func (suggester *SortSuggester) UnmarshalSpec(suggesterContext SuggesterContext, spec SuggesterSpec) error {
	if spec.Name != "SortSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'SortSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("SortSuggester expected to have no arguments.")
	}
	if len(spec.Specs) != 1 {
		return errors.New(fmt.Sprintf("SortSuggester expected to have 1 suggesters, got %d.", len(spec.Specs)))
	}
	if newSuggester, err := MakeSuggesterFromName(suggesterContext, spec.Specs[0].Name); err != nil {
		return err
	} else {
		suggester.suggester = newSuggester
		if err := suggester.suggester.UnmarshalSpec(suggesterContext, spec.Specs[0]); err != nil {
			return err
		}
	}
	return nil
}

func printFeed(contentItems []ContentItem) {
	for i := range contentItems {
		fmt.Printf("%d   %+v\n", i, contentItems[i])
	}
}

func Unsort(contentItems []ContentItem) []ContentItem {
	currentFeed := []ContentItem(nil)
	for i := range contentItems {
		currentFeed = append(currentFeed, contentItems[i])
	}
	fmt.Printf("\ncurrentFeed:\n")
	printFeed(currentFeed)
	sort.SliceStable(currentFeed, func(i, j int) bool {
		return currentFeed[i].OriginalOrder[0] < currentFeed[j].OriginalOrder[0]
	})
	fmt.Printf("\nUnsorted currentFeed:\n")
	printFeed(currentFeed)
	// Remove the original order.
	for i := range currentFeed {
		currentFeed[i].OriginalOrder = currentFeed[i].OriginalOrder[1:]
	}
	fmt.Printf("\nOrder removed currentFeed:\n")
	printFeed(currentFeed)
	return currentFeed
}

func (suggester *SortSuggester) More(request MoreRequest) ([]ContentItem, error) {
	currentFeed := Unsort(request.CurrentFeed)
	suggesterRequest := request
	suggesterRequest.CurrentFeed = currentFeed
	items, err := suggester.suggester.More(suggesterRequest)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].OriginalOrder = append([]int64{int64(i + len(currentFeed))}, items[i].OriginalOrder...)
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Date.Equal(items[j].Date) {
			return items[i].CreatedAt.After(items[j].CreatedAt)
		} else {
			return items[i].Date.After(items[j].Date)
		}
	})
	//fmt.Printf("Sort:\n")
	//for i, ci := range items {
	//	fmt.Printf("%d: %+v\n", i+1, ci)
	//}
	return items, nil
}
