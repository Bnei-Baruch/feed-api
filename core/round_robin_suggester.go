package core

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/feed-api/utils"
)

type RoundRobinSuggester struct {
	Name       string
	Suggesters []Suggester
}

func MakeRoundRobinSuggester(suggesters []Suggester) *RoundRobinSuggester {
	return &RoundRobinSuggester{
		Name:       "RoundRobinSuggester",
		Suggesters: suggesters,
	}
}

func init() {
	RegisterSuggester("RoundRobinSuggester", func(suggesterContext SuggesterContext) Suggester {
		return MakeRoundRobinSuggester([]Suggester(nil))
	})
}

func (suggester *RoundRobinSuggester) MarshalSpec() (SuggesterSpec, error) {
	var specs []SuggesterSpec
	for i := range suggester.Suggesters {
		if spec, err := suggester.Suggesters[i].MarshalSpec(); err != nil {
			return SuggesterSpec{}, err
		} else {
			specs = append(specs, spec)
		}
	}
	return SuggesterSpec{
		Name:  suggester.Name,
		Specs: specs,
	}, nil
}

func (suggester *RoundRobinSuggester) UnmarshalSpec(suggesterContext SuggesterContext, spec SuggesterSpec) error {
	if spec.Name != suggester.Name {
		return errors.New(fmt.Sprintf("Expected suggester name to be: '%s', got: '%s'.", suggester.Name, spec.Name))
	}
	if len(spec.Args) != 0 {
		return errors.New(fmt.Sprintf("%s expected to have no arguments.", suggester.Name))
	}
	if len(spec.Specs) == 0 {
		return errors.New(fmt.Sprintf("%s expected to have some suggesters, got 0.", suggester.Name))
	}
	for i := range spec.Specs {
		if newSuggester, err := MakeSuggesterFromName(suggesterContext, spec.Specs[i].Name); err != nil {
			return err
		} else {
			if err := newSuggester.UnmarshalSpec(suggesterContext, spec.Specs[i]); err != nil {
				return err
			}
			suggester.Suggesters = append(suggester.Suggesters, newSuggester)
		}
	}
	return nil
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
	currentFeeds := SplitContentItems(request.CurrentFeed, len(suggester.Suggesters))
	maxLength := 0
	for i, s := range suggester.Suggesters {
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
	// Shift all items to continue from the right suggester from previous items.
	// This might work wrong if:
	// a) Data changed from pervious call (which is ok).
	// b) There were duplicate UIDs in suggesters which make modulo the wrong action here.
	offset := len(request.CurrentFeed) % len(suggester.Suggesters)
	allItems = append(allItems[offset:len(allItems)], allItems[0:offset]...)
	uids := make(map[string]bool)
	for _, contentItem := range request.CurrentFeed {
		uids[contentItem.UID] = true
	}
	allItemsIndexes := []int(nil)
	for range allItems {
		allItemsIndexes = append(allItemsIndexes, 0)
	}
	// Eventually we need no more than request.MoreItems
	appendLength := utils.MaxInt(maxLength, request.MoreItems)
	for i := 0; i < appendLength; i++ {
		for j, items := range allItems {
			// Find in this suggester the next uid which we don't see in the feed.
			for ; allItemsIndexes[j] < len(items); allItemsIndexes[j]++ {
				if _, ok := uids[items[allItemsIndexes[j]].UID]; !ok {
					uids[items[allItemsIndexes[j]].UID] = true
					roundRobin = append(roundRobin, items[allItemsIndexes[j]])
					break
				}
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
