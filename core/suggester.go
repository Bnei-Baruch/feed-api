package core

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/feed-api/data_models"
)

const (
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
	DEFAULT     = "default"
)

// Keys are content types or collection mids. Values are one of the constants:
// SUBSCRIBE, UNSUBSCRIBE, DEFAULT
type Subscriptions map[string]string

type Recommend struct {
	Uid string `json:"uid"`
}

// Define filter constrain for unit being watched.
type FilterSelectorEnum int

const (
	UnitContentTypes       FilterSelectorEnum = 0
	CollectionContentTypes                    = 1
	Tags                                      = 2
	Sources                                   = 3
	Collections                               = 4
	SameTag                                   = 5
	SameSource                                = 6
	SameCollection                            = 7
	WatchingNowFilter                         = 8
	PopularFilter                             = 9
	AgeFilter                                 = 10
)

var FILTER_STRING_TO_VALUE = map[string]FilterSelectorEnum{
	"FUnitContentTypes":       UnitContentTypes,
	"FCollectionContentTypes": CollectionContentTypes,
	"FTags":                   Tags,
	"FSources":                Sources,
	"FCollections":            Collections,
	"FSameTag":                SameTag,
	"FSameSource":             SameSource,
	"FSameCollection":         SameCollection,
	"FWatchingNowFilter":      WatchingNowFilter,
	"FPopularFilter":          PopularFilter,
	"FAgeFilter":              AgeFilter,
}

type SuggesterFilter struct {
	FilterSelector FilterSelectorEnum `json:"filter_selector,omitempty" form:"filter_selector,omitempty"`
	Args           []string           `json:"args,omitempty" form:"args,omitempty"`
}

// Defines order constrain for selecting content units / collection.
type OrderSelectorEnum int

const (
	Last        OrderSelectorEnum = 0
	Next                          = 1
	Prev                          = 2
	Rand                          = 3
	Popular                       = 4
	WatchingNow                   = 5
)

var ORDER_STRING_TO_VALUE = map[string]OrderSelectorEnum{
	"OLast":        Last,
	"ONext":        Next,
	"OPrev":        Prev,
	"ORand":        Rand,
	"OPopular":     Popular,
	"OWatchingNow": WatchingNow,
}

type SuggesterSpec struct {
	Name          string            `json:"name,omitempty" form:"name,omitempty"`
	Filters       []SuggesterFilter `json:"filters,omitempty" form:"filters,omitempty"`
	OrderSelector OrderSelectorEnum `json:"order_selector,omitempty" form:"order_selector,omitempty"`
	Specs         []SuggesterSpec   `json:"specs,omitempty" form:"specs,omitempty"`

	// Deprecated
	Args       []string `json:"args,omitempty" form:"args,omitempty"`
	SecondArgs []string `json:"second_args,omitempty" form:"second_args,omitempty"`
}

type MoreOptions struct {
	ContentTypes Subscriptions `json:"content_types" form:"content_type"`
	// Map from collection content type to Subscriptions.
	Collections map[string]Subscriptions `json:"collections" form:"collections"`

	Recommend Recommend        `json:"recommend" form:"recommend"`
	Spec      *SuggesterSpec   `json:"spec,omitempty" form:"spec,omitempty"`
	Specs     []*SuggesterSpec `json:"specs,omitempty" form:"specs,omitempty"`

	Languages []string `json:"languages,omitempty" form:"languages,omitempty"`
	SkipUids  []string `json:"skip_uids" form:"skip_uids"`

	DebugTimestamp *int64 `json:"debug_timestamp,omitempty" form:"debug_timestamp,omitempty"`

	WatchingNowMin *int64 `json:"watching_now_min,omitempty" form:"watching_now_min,omitempty"`
	PopularMin     *int64 `json:"popular_min,omitempty" form:"popular_min,omitempty"`

	WithPosts bool `json:"with_posts,omitempty" form:"with_posts,omitempty"`
}

type MoreRequest struct {
	Namespace   string        `json:"namespace,omitempty" form:"namespace,omitempty"`
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
	Suggester     string    `json:"suggester"`
	FeedOrder     int64     `json:"feed_order"`
}

type Suggester interface {
	More(request MoreRequest) ([]ContentItem, error)

	MarshalSpec() (SuggesterSpec, error)
	UnmarshalSpec(suggesterContext SuggesterContext, spec SuggesterSpec) error
}

type SuggesterContext struct {
	DB         *sql.DB
	DataModels *data_models.DataModels
	Cache      map[string]interface{}
}

type MakeSuggesterFunc func(suggesterContext SuggesterContext) Suggester

var Suggesters = map[string]MakeSuggesterFunc{}

func RegisterSuggester(name string, makeFunc MakeSuggesterFunc) {
	Suggesters[name] = makeFunc
}

func MakeSuggesterFromName(suggesterContext SuggesterContext, name string) (Suggester, error) {
	if makeSuggesterFunc, ok := Suggesters[name]; !ok {
		return nil, errors.New(fmt.Sprintf("Did not find suggester %s in registry.", name))
	} else {
		return makeSuggesterFunc(suggesterContext), nil
	}
}

func CurrentFeedsToUidsString(currentFeed []ContentItem) string {
	parts := []string(nil)
	for i := range currentFeed {
		parts = append(parts, currentFeed[i].UID)
	}
	return strings.Join(parts, ",")
}
