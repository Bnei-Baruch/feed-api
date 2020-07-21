package core

import (
	//"sort"

	"github.com/Bnei-Baruch/feed-api/utils"
)

func Merge(r MoreRequest, suggestions [][]ContentItem) ([]ContentItem, error) {
	mergedFeed := append([]ContentItem(nil), r.CurrentFeed...)
	for _, s := range suggestions {
		mergedFeed = append(mergedFeed, s...)
	}
	//sort.SliceStable(mergedFeed, func(i, j int) bool {
	//	return mergedFeed[i].CreatedAt.After(mergedFeed[j].CreatedAt)
	//})
	return mergedFeed[0:utils.MinInt(len(r.CurrentFeed)+r.MoreItems, len(mergedFeed))], nil
}
