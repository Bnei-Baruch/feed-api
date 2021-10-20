package core

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Bnei-Baruch/feed-api/utils"
)

// Given a Suggester will take limited number of suggestions from that suggester.
// The |limit| will include already suggested (current feed) items as well so additional
// calls to More will not bring items more than |limit|.
type LimitSuggester struct {
	limit     int
	suggester Suggester
}

func MakeLimitSuggester(limit int, suggester Suggester) *LimitSuggester {
	return &LimitSuggester{limit, suggester}
}

func init() {
	RegisterSuggester("LimitSuggester", func(suggesterContext SuggesterContext) Suggester {
		return MakeLimitSuggester(0, nil)
	})
}

func (suggester *LimitSuggester) MarshalSpec() (SuggesterSpec, error) {
	if spec, err := suggester.suggester.MarshalSpec(); err != nil {
		return SuggesterSpec{}, err
	} else {
		return SuggesterSpec{
			Args:  []string{strconv.Itoa(suggester.limit)},
			Name:  "LimitSuggester",
			Specs: []SuggesterSpec{spec},
		}, nil
	}
}

func (suggester *LimitSuggester) UnmarshalSpec(suggesterContext SuggesterContext, spec SuggesterSpec) error {
	if spec.Name != "LimitSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'LimitSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Args) != 1 {
		return errors.New(fmt.Sprintf("LimitSuggester expected to have one argument, got %d.", len(spec.Args)))
	}
	if limit, err := strconv.Atoi(spec.Args[0]); err == nil {
		suggester.limit = limit
	} else {
		return errors.New(fmt.Sprintf("Failed converting arg (%s) to decimal.", spec.Args[0]))
	}
	if len(spec.Specs) != 1 {
		return errors.New(fmt.Sprintf("LimitSuggester expected to have one suggester, got %d.", len(spec.Specs)))
	}
	if newSuggester, err := MakeSuggesterFromName(suggesterContext, spec.Specs[0].Name); err != nil {
		return err
	} else {
		if err := newSuggester.UnmarshalSpec(suggesterContext, spec.Specs[0]); err != nil {
			return err
		}
		suggester.suggester = newSuggester
	}
	return nil
}

func (suggester *LimitSuggester) More(request MoreRequest) ([]ContentItem, error) {
	// We will fetch total number of limit items (including previous More calls, e.g., len(request.CurrentFeed)).
	// No more than |request.MoreItems| ofcourse.
	request.MoreItems = utils.MaxInt(utils.MinInt(suggester.limit-len(request.CurrentFeed), request.MoreItems), 0)
	log.Debugf("LimitSuggester MoreRequest:%d  %+v  %+v", request.MoreItems, CurrentFeedsToUidsString(request.CurrentFeed), request.Options.SkipUids)
	if request.MoreItems == 0 {
		return []ContentItem(nil), nil
	}
	if items, err := suggester.suggester.More(request); err != nil {
		return nil, err
	} else {
		return items, nil
	}
}
