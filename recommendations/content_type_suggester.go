package recommendations

import (
	"database/sql"
	"fmt"

	"github.com/Bnei-Baruch/sqlboiler/queries"
	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/mdb"
)

type ContentTypeSuggester struct {
	db           *sql.DB
	contentTypes []string
	suggesters   []core.Suggester
}

func MakeContentTypeSuggester(db *sql.DB, contentTypes []string, suggesters []core.Suggester) core.Suggester {
	return &ContentTypeSuggester{db, contentTypes, suggesters}
}

func (suggester *ContentTypeSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	var specs []core.SuggesterSpec
	for i := range suggester.suggesters {
		if spec, err := suggester.suggesters[i].MarshalSpec(); err != nil {
			return core.SuggesterSpec{}, err
		} else {
			specs = append(specs, spec)
		}
	}
	return core.SuggesterSpec{
		Name:  "ContentTypeSuggester",
		Args:  suggester.contentTypes,
		Specs: specs,
	}, nil
}

func (suggester *ContentTypeSuggester) UnmarshalSpec(suggesterContext core.SuggesterContext, spec core.SuggesterSpec) error {
	if spec.Name != "ContentTypeSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'ContentTypeSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) == 0 {
		return errors.New("ContentTypeSuggester expected to have some arguments.")
	}
	if len(spec.Specs) == 0 {
		return errors.New("ContentTypeSuggester expected to have some suggesters, got 0.")
	}
	suggester.contentTypes = spec.Args
	for i := range spec.Specs {
		if newSuggester, err := core.MakeSuggesterFromName(suggesterContext, spec.Specs[i].Name); err != nil {
			return err
		} else {
			if err := newSuggester.UnmarshalSpec(suggesterContext, spec.Specs[i]); err != nil {
				return err
			}
			suggester.suggesters = append(suggester.suggesters, newSuggester)
		}
	}
	return nil
}

func (s *ContentTypeSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	if request.Options.Recommend.Uid == "" {
		return []core.ContentItem(nil), nil
	}
	rows, err := queries.Raw(s.db, fmt.Sprintf(`
		select cu.type_id
		from
		  content_units as cu
		where
		  cu.uid = '%s'
	`, request.Options.Recommend.Uid)).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contentType := ""
	for rows.Next() {
		var typeId int64
		err := rows.Scan(&typeId)
		contentType = mdb.CONTENT_TYPE_REGISTRY.ByID[typeId].Name
		if err != nil {
			return nil, err
		}
	}
	for i := range s.contentTypes {
		if s.contentTypes[i] == contentType || s.contentTypes[i] == "*" {
			if i < len(s.suggesters) {
				return s.suggesters[i].More(request)
			}
		}
	}
	return []core.ContentItem(nil), nil
}

func init() {
	core.RegisterSuggester("ContentTypeSuggester", func(suggesterContext core.SuggesterContext) core.Suggester {
		return MakeContentTypeSuggester(suggesterContext.DB, []string(nil), []core.Suggester(nil))
	})
}
