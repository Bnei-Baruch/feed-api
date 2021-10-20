package core

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Sort items by date and then by create at time.
type SortSuggester struct {
	suggesters []Suggester
}

func MakeSortSuggester(suggesters []Suggester) *SortSuggester {
	return &SortSuggester{
		suggesters: suggesters,
	}
}

func init() {
	RegisterSuggester("SortSuggester", func(suggesterContext SuggesterContext) Suggester { return MakeSortSuggester(nil) })
}

func (suggester *SortSuggester) MarshalSpec() (SuggesterSpec, error) {
	var specs []SuggesterSpec
	for i := range suggester.suggesters {
		if spec, err := suggester.suggesters[i].MarshalSpec(); err != nil {
			return SuggesterSpec{}, err
		} else {
			specs = append(specs, spec)
		}
	}
	return SuggesterSpec{
		Name:  "SortSuggester",
		Specs: specs,
	}, nil
}

func (suggester *SortSuggester) UnmarshalSpec(suggesterContext SuggesterContext, spec SuggesterSpec) error {
	if spec.Name != "SortSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'SortSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New("SortSuggester expected to have no arguments.")
	}
	if len(spec.Specs) == 0 {
		return errors.New(fmt.Sprintf("SortSuggester expected to have some suggesters got 0."))
	}
	for i := range spec.Specs {
		if newSuggester, err := MakeSuggesterFromName(suggesterContext, spec.Specs[i].Name); err != nil {
			return err
		} else {
			if err := newSuggester.UnmarshalSpec(suggesterContext, spec.Specs[i]); err != nil {
				return err
			}
			suggester.suggesters = append(suggester.suggesters, newSuggester)
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
	// fmt.Printf("\ncurrentFeed:\n")
	// printFeed(currentFeed)
	sort.SliceStable(currentFeed, func(i, j int) bool {
		return currentFeed[i].OriginalOrder[0] < currentFeed[j].OriginalOrder[0]
	})
	// fmt.Printf("\nUnsorted currentFeed:\n")
	// printFeed(currentFeed)
	// Remove the original order.
	for i := range currentFeed {
		currentFeed[i].OriginalOrder = currentFeed[i].OriginalOrder[1:]
	}
	// fmt.Printf("\nOrder removed currentFeed:\n")
	// printFeed(currentFeed)
	return currentFeed
}

func (suggester *SortSuggester) Split(contentItems []ContentItem) [][]ContentItem {
	allItems := make([][]ContentItem, len(suggester.suggesters))
	for i := range contentItems {
		order := contentItems[i].OriginalOrder[0]
		contentItems[i].OriginalOrder = contentItems[i].OriginalOrder[1:]
		allItems[order] = append(allItems[order], contentItems[i])
	}
	//fmt.Printf("\nAfter split currentFeed:\n")
	//for i := range allItems {
	//	fmt.Printf("%d:\n", i)
	//	printFeed(allItems[i])
	//}
	return allItems
}

func (suggester *SortSuggester) More(request MoreRequest) ([]ContentItem, error) {
	maxOriginalOrder := int64(0)
	for i := range request.CurrentFeed {
		sortOrder := request.CurrentFeed[i].OriginalOrder[0]
		if maxOriginalOrder < sortOrder {
			maxOriginalOrder = sortOrder
		}
		// Skip uids from other suggesters from previous requests (e.g., from current feed)
		request.Options.SkipUids = append(request.Options.SkipUids, request.CurrentFeed[i].UID)
	}
	allCurrentItems := suggester.Split(Unsort(request.CurrentFeed))
	moreItems := []ContentItem(nil)
	for order := range allCurrentItems {
		currentFeed := allCurrentItems[order]
		suggesterRequest := request
		suggesterRequest.CurrentFeed = currentFeed
		// fmt.Printf("More[%d]: %+v\n", order, suggesterRequest)
		log.Debugf("SortSuggester %d MoreRequest:%d  %+v  %+v", order, suggesterRequest.MoreItems, CurrentFeedsToUidsString(suggesterRequest.CurrentFeed), suggesterRequest.Options.SkipUids)
		items, err := suggester.suggesters[order].More(suggesterRequest)
		if err != nil {
			return nil, err
		}
		for i := range items {
			items[i].OriginalOrder = append([]int64{int64(len(moreItems)) + maxOriginalOrder, int64(order)}, items[i].OriginalOrder...)
			moreItems = append(moreItems, items[i])
			// Skip uids which previous suggester suggested.
			request.Options.SkipUids = append(request.Options.SkipUids, items[i].UID)
		}
	}
	// fmt.Printf("After More:\n")
	//for i, ci := range moreItems {
	//	fmt.Printf("%d: %+v\n", i+1, ci)
	//}
	sort.SliceStable(moreItems, func(i, j int) bool {
		if moreItems[i].Date.Equal(moreItems[j].Date) {
			return moreItems[i].CreatedAt.After(moreItems[j].CreatedAt)
		} else {
			return moreItems[i].Date.After(moreItems[j].Date)
		}
	})
	//fmt.Printf("Sort:\n")
	//for i, ci := range moreItems {
	//	fmt.Printf("%d: %+v\n", i+1, ci)
	//}
	if len(moreItems) > request.MoreItems {
		moreItems = moreItems[:request.MoreItems]
	}
	//fmt.Printf("Cut:\n")
	//for i, ci := range moreItems {
	//	fmt.Printf("%d: %+v\n", i+1, ci)
	//}
	return moreItems, nil
}
