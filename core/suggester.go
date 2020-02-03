package core

import (
	"time"
)

type MoreRequest struct {
	CurrentFeed []ContentItem `json:"current_feed" form:"current_feed"`
	MoreItems   int           `json:"more_items" form:"more_items"`
}

type ContentItem struct {
	UID         string    `json:"uid"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type Suggester interface {
	More(request MoreRequest) ([]ContentItem, error)
}
