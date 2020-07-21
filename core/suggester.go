package core

import (
	"time"
)

const (
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
	DEFAULT     = "default"
)

// Keys are content types or collection mids. Values are one of the constants:
// SUBSCRIBE, UNSUBSCRIBE, DEFAULT
type Subscriptions map[string]string

type MoreOptions struct {
	ContentTypes Subscriptions `json:"content_types" form:"content_type"`
	// Map from collection content type to Subscriptions.
	Collections map[string]Subscriptions `json:"collections" form:"collections"`
}

type MoreRequest struct {
	CurrentFeed []ContentItem `json:"current_feed" form:"current_feed"`
	MoreItems   int           `json:"more_items" form:"more_items"`
	Options     MoreOptions   `json:"options" form:"options"`
}

type ContentItem struct {
	UID           string    `json:"uid"`
	ContentType   string    `json:"content_type"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	OriginalOrder []int64   `json:"original_order"`
}

type Suggester interface {
	More(request MoreRequest) ([]ContentItem, error)
}
