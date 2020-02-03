package core

import (
	"sort"
)

func Merge(currentFeed []ContentItem, suggestions [][]ContentItem) ([]ContentItem, error) {
	mergedFeed := append([]ContentItem(nil), currentFeed...)
	for _, s := range suggestions {
		mergedFeed = append(mergedFeed, s...)
	}
	sort.SliceStable(mergedFeed, func(i, j int) bool {
		return mergedFeed[i].CreatedAt.Before(mergedFeed[j].CreatedAt)
	})
	return mergedFeed, nil
}
