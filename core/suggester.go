package core

import (
	"database/sql"
	"fmt"
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
	UnitContentTypes FilterSelectorEnum = iota
	CollectionContentTypes
	Tags
	Sources
	Collections
	SameTag
	SameSource
	SameCollection
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
}

type SuggesterFilter struct {
	FilterSelector FilterSelectorEnum `json:"filter_selector,omitempty" form:"filter_selector,omitempty"`
	Args           []string           `json:"args,omitempty" form:"args,omitempty"`
}

// Defines order constrain for selecting content units / collection.
type OrderSelectorEnum int

const (
	Last OrderSelectorEnum = iota
	Next
	Prev
	Rand
	// (TBD) Popular
)

var ORDER_STRING_TO_VALUE = map[string]OrderSelectorEnum{
	"OLast": Last,
	"ONext": Next,
	"OPrev": Prev,
	"ORand": Rand,
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

	Recommend Recommend      `json:"recommend" form:"recommend"`
	Spec      *SuggesterSpec `json:"spec,omitempty" form:"spec,omitempty"`

	Languages []string `json:"languages,omitempty" form:"languages,omitempty"`
	SkipUids  []string `json:"skip_uids"`
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
	Suggester     string    `json:"suggester"`
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
