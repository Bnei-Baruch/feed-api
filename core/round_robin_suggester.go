package core

import (
	//"fmt"

	"github.com/Bnei-Baruch/feed-api/utils"
)

type RoundRobinSuggester struct {
	suggesters []Suggester
}

func MakeRoundRobinSuggester(suggesters []Suggester) *RoundRobinSuggester {
	return &RoundRobinSuggester{suggesters: suggesters}
}

func SplitContentItems(contentItems []ContentItem, numLists int) [][]ContentItem {
	currentFeeds := make([][]ContentItem, numLists)
	feedIndex := 0
	for i := range contentItems {
		currentFeeds[feedIndex] = append(currentFeeds[feedIndex], contentItems[i])
		feedIndex = (feedIndex + 1) % numLists
	}
	return currentFeeds
}

func (suggester *RoundRobinSuggester) More(request MoreRequest) ([]ContentItem, error) {
	allItems := [][]ContentItem(nil)
	currentFeeds := SplitContentItems(request.CurrentFeed, len(suggester.suggesters))
	maxLength := 0
	for i, s := range suggester.suggesters {
		suggesterRequest := request
		suggesterRequest.CurrentFeed = currentFeeds[i]
		if items, err := s.More(suggesterRequest); err != nil {
			return nil, err
		} else {
			allItems = append(allItems, items)
			maxLength = utils.MaxInt(maxLength, len(items))
		}
	}
	roundRobin := []ContentItem(nil)
	// Shift all items to continue from the right place from previous items.
	offset := len(request.CurrentFeed) % len(suggester.suggesters)
	allItems = append(allItems[offset:len(allItems)], allItems[0:offset]...)
	for i := 0; i < maxLength; i++ {
		for _, items := range allItems {
			if i < len(items) {
				roundRobin = append(roundRobin, items[i])
			}
		}
	}
	ret := roundRobin[0:utils.MinInt(request.MoreItems, len(roundRobin))]
	//fmt.Printf("RoundRobin:\n")
	//for i, ci := range ret {
	//	fmt.Printf("%d: %+v\n", i+1, ci)
	//}
	return ret, nil
}
