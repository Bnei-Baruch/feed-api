package api

import (
	"time"

	"github.com/Bnei-Baruch/feed-api/utils"
)

type ItemsRequest struct {
	Offset int `json:"offset,omitempty" form:"offset" binding:"omitempty,min=0"`
}

type Item struct {
	ID          string      `json:"id"`
	ContentType string      `json:"content_type"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	FilmDate    *utils.Date `json:"film_date,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
}

type ContentUnit struct {
	Item
	Collections      map[string]*Collection  `json:"collections,omitempty"`
	DerivedUnits     map[string]*ContentUnit `json:"derived_units,omitempty"`
	Duration         float64                 `json:"duration,omitempty"`
	Files            []*File                 `json:"files,omitempty"`
	NameInCollection string                  `json:"name_in_collection,omitempty"`
	OriginalLanguage string                  `json:"original_language,omitempty"`
	Publishers       []string                `json:"publishers,omitempty"`
	SourceUnits      map[string]*ContentUnit `json:"source_units,omitempty"`
	Sources          []string                `json:"sources,omitempty"`
	Tags             []string                `json:"tags,omitempty"`
	mdbID            int64
}

type Collection struct {
	Item
	City            string         `json:"city,omitempty"`
	ContentUnits    []*ContentUnit `json:"content_units,omitempty"`
	Country         string         `json:"country,omitempty"`
	DefaultLanguage string         `json:"default_language,omitempty"`
	EndDate         *utils.Date    `json:"end_date,omitempty"`
	FullAddress     string         `json:"full_address,omitempty"`
	Genres          []string       `json:"genres,omitempty"`
	HolidayID       string         `json:"holiday_id,omitempty"`
	Number          int            `json:"number,omitempty"`
	SourceID        string         `json:"source_id,omitempty"`
	StartDate       *utils.Date    `json:"start_date,omitempty"`
}

type File struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Duration  float64   `json:"duration,omitempty"`
	Language  string    `json:"language,omitempty"`
	MimeType  string    `json:"mimetype,omitempty"`
	Type      string    `json:"type,omitempty"`
	SubType   string    `json:"subtype,omitempty"`
	VideoSize string    `json:"video_size,omitempty"`
	CreatedAt time.Time `json:"-"`
}

type ItemInterface interface {
}

type ItemsResponse struct {
	Total int64         `json:"total"`
	Items []interface{} `json:"items"`
}
