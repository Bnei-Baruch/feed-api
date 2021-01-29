package core

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
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

type SuggesterSpec struct {
	Name  string          `json:"name,omitempty" form:"name,omitempty"`
	Args  []string        `json:"args,omitempty" form:"args,omitempty"`
	Specs []SuggesterSpec `json:"specs,omitempty" form:"specs,omitempty"`
}

type MoreOptions struct {
	ContentTypes Subscriptions `json:"content_types" form:"content_type"`
	// Map from collection content type to Subscriptions.
	Collections map[string]Subscriptions `json:"collections" form:"collections"`
	Recommend   Recommend                `json:"recommend" form:"recommend"`
	Spec        *SuggesterSpec           `json:"spec,omitempty" form:"spec,omitempty"`
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
	UnmarshalSpec(db *sql.DB, spec SuggesterSpec) error
}

type MakeSuggesterFunc func(db *sql.DB) Suggester

var Suggesters = map[string]MakeSuggesterFunc{}

func RegisterSuggester(name string, makeFunc MakeSuggesterFunc) {
	Suggesters[name] = makeFunc
}

func MakeSuggesterFromName(db *sql.DB, name string) (Suggester, error) {
	if makeSuggesterFunc, ok := Suggesters[name]; !ok {
		return nil, errors.New(fmt.Sprintf("Did not find suggester %s in registry.", name))
	} else {
		return makeSuggesterFunc(db), nil
	}
}
