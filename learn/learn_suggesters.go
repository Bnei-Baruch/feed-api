package learn

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tma15/gonline"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/databases/chronicles/models"
	"github.com/pkg/errors"
)

const MIN_RETRAIN_POSITIVES = 100
const THREE_HOURS_AGO = time.Hour * -3

type Examples struct {
	Features *[]map[string]float64
	Labels   *[]string
}

func (e *Examples) Len() int {
	return len(*e.Features)
}

func (e *Examples) Append(examples Examples) {
	f := append(*e.Features, *examples.Features...)
	l := append(*e.Labels, *examples.Labels...)
	e.Features = &f
	e.Labels = &l
}

func (e *Examples) Trunc(s int) {
	f := (*e.Features)[:s]
	l := (*e.Labels)[:s]
	e.Features = &f
	e.Labels = &l
}

func MakeExamples() Examples {
	f := make([]map[string]float64, 0)
	l := make([]string, 0)
	return Examples{&f, &l}
}

type LearnSuggester struct {
	localChroniclesDb *sql.DB
	lastReadId        string
	selected          map[string]*models.Entry
	recommended       map[string]*models.Entry

	learner      gonline.LearnerInterface
	posExamples  Examples
	negExamples  Examples
	evalExamples Examples
}

func MakeLearnSuggester(localChroniclesDb *sql.DB) *LearnSuggester {
	return &LearnSuggester{
		localChroniclesDb,
		"",
		make(map[string]*models.Entry),
		make(map[string]*models.Entry),
		gonline.NewArow(10.),
		MakeExamples(),
		MakeExamples(),
		MakeExamples(),
	}
}

func (s *LearnSuggester) Name() string {
	return "LearnSuggester"
}

func (s *LearnSuggester) Interval() time.Duration {
	return time.Duration(time.Minute)
}

func UpdateEntriesMap(entries models.EntrySlice, m map[string]*models.Entry) {
	for _, entry := range entries {
		if clientEventId := entry.ClientEventID.String; clientEventId == "" {
			log.Warnf("Empty client event id for entry: %s", entry.ID)
		} else {
			m[clientEventId] = entry
		}
	}
}

func RemoveOldEntries(lastCreatedAt time.Time, m map[string]*models.Entry) {
	removed := 0
	t := lastCreatedAt.Add(THREE_HOURS_AGO)
	for id, entry := range m {
		if entry.CreatedAt.Before(t) {
			removed += 1
			delete(m, id)
		}
	}
	log.Infof("Removed %d old entries", removed)
}

func RecommendsMapToNegativeExamples(recommends map[string]*models.Entry, rToSMap map[string]map[string]bool, lastCreatedAt time.Time) ([]RecommendUidsPair, error) {
	ret := []RecommendUidsPair(nil)
	skippedFresh := 0
	for _, entry := range recommends {
		if entry.CreatedAt.After(lastCreatedAt.Add(time.Hour * -3)) {
			skippedFresh += 1
			// Skip fresh recommend entries as they might still have selections.
			// Take only recommend entries older than 3 hours as negatives.
			continue
		}
		// Converts one recommendation to number of pairs of
		// source uid => recommend uid (filter actually selected uids)
		recommendData, err := GetRecommendData(entry)
		if err != nil {
			log.Warnf("Failed unmarshling recommend data: %+v", err)
			continue
		}
		rUid := recommendData.RequestData.Options.Recommend.Uid
		for _, typedItems := range recommendData.Recommendations {
			for _, typedItem := range typedItems {
				sUid := typedItem.Uid
				if sMap, rOk := rToSMap[rUid]; rOk {
					if _, sOk := sMap[sUid]; sOk {
						continue // Don't take positive example as negative.
					}
				}
				ret = append(ret, RecommendUidsPair{Recommended: rUid, Selected: sUid})
			}
		}
	}
	log.Infof("Skipped %d recommend entries, cannot take them for negatives.", skippedFresh)
	return ret, nil
}

func (s *LearnSuggester) refreshExamples() (Examples, Examples, string, error) {
	// Load more selected and recommended entries.
	var err error
	selected := models.EntrySlice(nil)
	if selected, err = MakeDbScanner("recommend-selected", s.lastReadId, s.localChroniclesDb).ScanAll(); err != nil {
		return Examples{}, Examples{}, s.lastReadId, err
	}
	recommended := models.EntrySlice(nil)
	if recommended, err = MakeDbScanner("recommend", s.lastReadId, s.localChroniclesDb).ScanAll(); err != nil {
		return Examples{}, Examples{}, s.lastReadId, err
	}

	var lastCreatedAt *time.Time = nil
	if len(selected) > 0 {
		s.lastReadId = selected[len(selected)-1].ID
		lastCreatedAt = &selected[len(selected)-1].CreatedAt
	}
	newSelected := make(map[string]*models.Entry)
	UpdateEntriesMap(selected, s.selected)
	UpdateEntriesMap(selected, newSelected)
	if len(recommended) > 0 {
		recommendedReadId := recommended[len(recommended)-1].ID
		if strings.Compare(s.lastReadId, recommendedReadId) == -1 {
			s.lastReadId = recommendedReadId
			lastCreatedAt = &recommended[len(recommended)-1].CreatedAt
		}
	}
	UpdateEntriesMap(recommended, s.recommended)

	// Convert possible recommend / selected entries to examples.
	positiveUidsPairs := []RecommendUidsPair(nil)
	rToSUidsMap := make(map[string]map[string]bool)
	uidsToLoad := make(map[string]bool)
	noClientFlowId := 0
	flowNotFound := 0

	for _, selectEntry := range s.selected {
		// log.Infof("Entry: %+v", selectEntry)
		if !selectEntry.ClientFlowID.Valid || selectEntry.ClientFlowID.String == "" {
			noClientFlowId++
			if noClientFlowId < 10 {
				log.Infof("No client flow id. Select id: %s", selectEntry.ID)
			} else if noClientFlowId < 11 {
				log.Infof("Skipping other no client flow id selected.")
			}
			continue
		}
		if !selectEntry.Data.Valid {
			log.Warnf("Non valid select data: %s", selectEntry.ID)
			continue
		}
		if selectedData, err := GetSelectedData(selectEntry); err != nil {
			log.Warnf("Failed getting select data, skipping.", err)
			continue
		} else {
			if recommendEntry, ok := s.recommended[selectEntry.ClientFlowID.String]; !ok {
				flowNotFound++
				if flowNotFound < 10 {
					log.Infof("Could not find recommend for recommend selected. Select id: %s FlowClientID: %s", selectEntry.ID, selectEntry.ClientFlowID.String)
				} else if flowNotFound < 11 {
					log.Infof("Skipping other not found recommend for recommend selected.")
				}
			} else {
				if recommendData, err := GetRecommendData(recommendEntry); err != nil {
					log.Warnf("Failed getting recommend data, skipping.", err)
					continue
				} else {
					// sUids[selectedData.Uid]++
					// sUidsCount += 1
					uidsToLoad[selectedData.Uid] = true
					rUid := recommendData.RequestData.Options.Recommend.Uid
					// Add positive examples only from newSelected entries. Other selected entries are stored to calculate properly negative examples.
					if _, ok := newSelected[selectedData.Uid]; ok {
						positiveUidsPairs = append(positiveUidsPairs, RecommendUidsPair{Recommended: rUid, Selected: selectedData.Uid})
					}
					uidsToLoad[rUid] = true
					// recommendedToSelected[rUid] = append(recommendedToSelected[rUid], RecommendClientEventIDPair{recommendEntry.ClientEventID.String, selectEntry.ClientEventID.String})
					// selectedToRecommended[selectedData.Uid] = append(selectedToRecommended[selectedData.Uid], RecommendClientEventIDPair{recommendEntry.ClientEventID.String, selectEntry.ClientEventID.String})
					// log.Infof("%+v %s => %s", selectEntry.CreatedAt, recommendData.RequestData.Options.Recommend.Uid, selectedData.Uid)
					if _, ok := rToSUidsMap[rUid]; !ok {
						rToSUidsMap[rUid] = make(map[string]bool)
					}
					rToSUidsMap[rUid][selectedData.Uid] = true
				}
			}
		}
	}
	log.Infof("Recommend selected: %d, recommended: %d, without flow id: %d, flow id not found: %d", len(s.selected), len(s.recommended), noClientFlowId, flowNotFound)

	// Negative pairs from logs.
	negativeUidsPairs, err := RecommendsMapToNegativeExamples(s.recommended, rToSUidsMap, *lastCreatedAt)
	if err != nil {
		return Examples{}, Examples{}, s.lastReadId, err
	}
	for _, pair := range negativeUidsPairs {
		uidsToLoad[pair.Recommended] = true
		uidsToLoad[pair.Selected] = true
	}

	log.Infof("Loading %d uids", len(uidsToLoad))
	unitsMap, err := LoadUids(uidsToLoad)
	if err != nil {
		return Examples{}, Examples{}, s.lastReadId, err
	}
	log.Infof("Loaded %d uids.", len(unitsMap))

	positiveItemPairs := uidsToItems(positiveUidsPairs, unitsMap)
	negativeItemPairs := uidsToItems(negativeUidsPairs, unitsMap)
	// positiveUnits := uidsMapToUnits(sUids, unitsMap)
	// negativeUnits := uidsToUnits(randomUids, unitsMap)
	log.Infof("Learning classifier %d positives and %d negatives.", len(positiveItemPairs), len(negativeItemPairs))
	/*if len(negativeItemPairs) > len(positiveItemPairs)*NEGATIVE_MULTIPLIER {
		rand.Shuffle(len(negativeItemPairs), func(i, j int) {
			negativeItemPairs[i], negativeItemPairs[j] = negativeItemPairs[j], negativeItemPairs[i]
		})
		negativeItemPairs = negativeItemPairs[:len(positiveItemPairs)*NEGATIVE_MULTIPLIER]
		log.Infof("Using only %d negatives.", len(negativeItemPairs))
	}*/

	x_pos, y_pos, x_neg, y_neg, err := PrepareExamples(positiveItemPairs, negativeItemPairs)

	// Clear old selected/recommended entries.
	RemoveOldEntries(*lastCreatedAt, s.selected)
	RemoveOldEntries(*lastCreatedAt, s.recommended)

	return Examples{x_pos, y_pos}, Examples{x_neg, y_neg}, s.lastReadId, err
}

func ShuffleSplitExamplesTenth(examples Examples) (Examples, Examples) {
	gonline.ShuffleData(examples.Features, examples.Labels)
	l := len(*examples.Features) / 10
	headFeatures := (*examples.Features)[:l]
	headLabels := (*examples.Labels)[:l]
	tailFeatures := (*examples.Features)[l:]
	tailLabels := (*examples.Labels)[l:]
	return Examples{&headFeatures, &headLabels}, Examples{&tailFeatures, &tailLabels}
}

func (s *LearnSuggester) Refresh() error {
	posExamples, negExamples, lastReadId, err := s.refreshExamples()
	if err != nil {
		return err
	}
	log.Infof("Positive examples: %d, negative examples: %d, prevReadId: %s lastReadId: %s", posExamples.Len(), negExamples.Len(), s.lastReadId, lastReadId)
	s.lastReadId = lastReadId
	s.posExamples.Append(posExamples)
	s.negExamples.Append(negExamples)

	if len(*s.posExamples.Features) < MIN_RETRAIN_POSITIVES {
		log.Info("Too few positive examples, skipping for now.")
		return nil
	}

	// If not enought negative examples normalize nagatives
	// between NEGATIVE_MULTIPLIER-1 and NEGATIVE_MULTIPLIER of positives.
	negExamples = Examples{}
	gonline.ShuffleData(s.negExamples.Features, s.negExamples.Labels)
	for len(*negExamples.Features) < (NEGATIVE_MULTIPLIER-1)*len(*s.posExamples.Features) {
		negExamples.Append(s.negExamples)
	}
	// Remove neg examples to max of NEGATIVE_MULTIPLIER if we duplicated them.
	if len(*negExamples.Features) > NEGATIVE_MULTIPLIER*len(*s.posExamples.Features) && len(*negExamples.Features) > len(*s.negExamples.Features) {
		negExamples.Trunc(NEGATIVE_MULTIPLIER * len(*s.posExamples.Features))
	}

	// If there is too much negatives (they were taken only once), duplicate positives.
	posExamples = Examples{}
	gonline.ShuffleData(s.posExamples.Features, s.posExamples.Labels)
	for NEGATIVE_MULTIPLIER*len(*posExamples.Features) < len(*negExamples.Features) {
		posExamples.Append(s.posExamples)
	}

	log.Infof("Training with new examples positives: %d, negatives: %d.", posExamples.Len(), negExamples.Len())
	posExamples.Append(negExamples)
	evalExamples, trainExamples := ShuffleSplitExamplesTenth(posExamples)
	s.evalExamples.Append(evalExamples)

	// Train on the new examples.
	times := 5
	for j := 0; j < times; j++ {
		gonline.ShuffleData(trainExamples.Features, trainExamples.Labels)
		s.learner.Fit(trainExamples.Features, trainExamples.Labels)
	}

	// Remove the examples used.
	s.posExamples = Examples{}
	s.negExamples = Examples{}

	// Evaluate new learned classifier.
	log.Infof("Iterative validation over %d examples.", evalExamples.Len())
	EvalLearner(0, &s.learner, s.evalExamples.Features, s.evalExamples.Labels)

	return nil
}

func (s *LearnSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	return nil, nil
}

func (s *LearnSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: "LearnSuggester"}, nil
}

func (s *LearnSuggester) UnmarshalSpec(suggesterContext core.SuggesterContext, spec core.SuggesterSpec) error {
	if spec.Name != "LearnSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'LearnSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("LearnSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	if len(spec.Filters) != 0 {
		return errors.New(fmt.Sprintf("LearnSuggester expected to have no filters, got %d.", len(spec.Filters)))
	}
	return nil
}
