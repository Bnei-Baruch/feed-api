package core

import (
	"fmt"
	"sort"
)

type SortSuggester struct {
	suggester Suggester
}

func MakeSortSuggester(suggester Suggester) *SortSuggester {
	return &SortSuggester{suggester: suggester}
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
		currentFeed[i].OriginalOrder = currentFeed[i].OriginalOrder[1:len(currentFeed[i].OriginalOrder)]
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
