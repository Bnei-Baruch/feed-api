package learn

import (
	"encoding/json"
	"math/rand"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/core"
	cModels "github.com/Bnei-Baruch/feed-api/databases/chronicles/models"
	"github.com/Bnei-Baruch/feed-api/databases/mdb"
	mdbModels "github.com/Bnei-Baruch/feed-api/databases/mdb/models"
	"github.com/Bnei-Baruch/feed-api/utils"
)

// Desired scale between positive and negative examples.
const NEGATIVE_MULTIPLIER = 5

type SelectedData struct {
	Uid string `json:"uid,omitempty"`
}

type TypedItem struct {
	Uid         string `json:"uid,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}

type RecommendData struct {
	RequestData     core.MoreRequest       `json:"request_data,ompitempty"`
	Recommendations map[string][]TypedItem `json:"recommendations,omitempty"`
}

func uidsMapToUnits(uidsMap map[string]int, unitsMap map[string]interface{}) []interface{} {
	ret := []interface{}(nil)
	for uid, count := range uidsMap {
		if unit, ok := unitsMap[uid]; ok {
			for i := 0; i < count; i++ {
				ret = append(ret, unit)
			}
		}
	}
	return ret
}

func uidsToUnits(uids []string, unitsMap map[string]interface{}) []interface{} {
	ret := []interface{}(nil)
	for _, uid := range uids {
		if unit, ok := unitsMap[uid]; ok {
			ret = append(ret, unit)
		}
	}
	return ret
}

func sampleRandomUids(skipUidsMap map[string]int, count int) (uids []string, err error) {
	// Content units
	var units mdbModels.ContentUnitSlice
	if units, err = mdbModels.ContentUnits(qm.Select("uid"), qm.OrderBy("random()"), qm.Limit(count)).All(common.LocalMdb); err != nil {
		return nil, err
	}
	for _, unit := range units {
		uids = append(uids, unit.UID)
	}
	return uids, nil
}

func LoadUids(uidsToLoad map[string]bool) (unitsMap map[string]interface{}, err error) {
	unitsMap = make(map[string]interface{})
	// Content units
	var units mdbModels.ContentUnitSlice
	if units, err = mdbModels.ContentUnits(qm.WhereIn("uid in ?", utils.ToInterfaceSlice(utils.StringKeys(uidsToLoad))...)).All(common.LocalMdb); err != nil {
		return nil, err
	}
	for _, unit := range units {
		unitsMap[unit.UID] = unit
	}
	// Collections
	var collections mdbModels.CollectionSlice
	if collections, err = mdbModels.Collections(qm.WhereIn("uid in ?", utils.ToInterfaceSlice(utils.StringKeys(uidsToLoad))...)).All(common.LocalMdb); err != nil {
		return nil, err
	}
	for _, collection := range collections {
		unitsMap[collection.UID] = collection
	}
	// Tags
	var tags mdbModels.TagSlice
	if tags, err = mdbModels.Tags(qm.WhereIn("uid in ?", utils.ToInterfaceSlice(utils.StringKeys(uidsToLoad))...)).All(common.LocalMdb); err != nil {
		return nil, err
	}
	for _, tag := range tags {
		unitsMap[tag.UID] = tag
	}
	// Sources
	var sources mdbModels.SourceSlice
	if sources, err = mdbModels.Sources(qm.WhereIn("uid in ?", utils.ToInterfaceSlice(utils.StringKeys(uidsToLoad))...)).All(common.LocalMdb); err != nil {
		return nil, err
	}
	for _, source := range sources {
		unitsMap[source.UID] = source
	}
	return unitsMap, nil
}

func LoadEntries(clientEventType string) (cModels.EntrySlice, error) {
	entries, err := cModels.Entries(qm.Where("client_event_type = ?", clientEventType)).All(common.LocalChroniclesDb)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt.Before(entries[j].CreatedAt)
	})
	if len(entries) == 0 {
		log.Infof("No selected")
	} else {
		log.Infof("Loaded %d %s. First at: %+v. Last at: %+v", len(entries), clientEventType, entries[0].CreatedAt, entries[len(entries)-1].CreatedAt)
	}
	return entries, err
}

func EntriesMap(entries cModels.EntrySlice) map[string]*cModels.Entry {
	m := make(map[string]*cModels.Entry)
	for _, entry := range entries {
		if clientEventId := entry.ClientEventID.String; clientEventId == "" {
			log.Warnf("Empty client event id for entry: %s", entry.ID)
		} else {
			m[clientEventId] = entry
		}
	}
	return m
}

func ContentItemType(item interface{}) string {
	if cu, ok := item.(*mdbModels.ContentUnit); ok {
		return mdb.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name
	} else if c, ok := item.(*mdbModels.Collection); ok {
		return mdb.CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name
	} else if _, ok := item.(*mdbModels.Tag); ok {
		return "TAG"
	} else if _, ok := item.(*mdbModels.Source); ok {
		return "SOURCE"
	}
	return ""
}

func GetSelectedData(entry *cModels.Entry) (SelectedData, error) {
	var selectedData SelectedData
	err := json.Unmarshal(entry.Data.JSON, &selectedData)
	return selectedData, err
}

func GetRecommendData(entry *cModels.Entry) (RecommendData, error) {
	var recommendData RecommendData
	err := json.Unmarshal(entry.Data.JSON, &recommendData)
	if err != nil {
		log.Warnf("Failed unmarshling %+v", string(entry.Data.JSON))
	}
	return recommendData, err
}

type RecommendClientEventIDPair struct {
	RecommendedClientEventID string
	SelectedClientEventID    string
}

type RecommendUidsPair struct {
	Recommended string
	Selected    string
}

type RecommendItemsPair struct {
	Recommended interface{}
	Selected    interface{}
}

func RecommendsToNegativeExamples(recommends cModels.EntrySlice, rToSMap map[string]map[string]bool) ([]RecommendUidsPair, error) {
	ret := []RecommendUidsPair(nil)
	for _, entry := range recommends {
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
	return ret, nil
}

func PrintTypesMap(unitsMap map[string]interface{}, recommendedToSelected map[string][]RecommendClientEventIDPair, selectedMap map[string]*cModels.Entry) error {
	rToSTypeMap := make(map[string]map[string]int)
	sToRTypeMap := make(map[string]map[string]int)

	recommendNotLoaded := 0
	selectedNotLoaded := 0

	for rUid, rToSPairs := range recommendedToSelected {
		if _, ok := unitsMap[rUid]; !ok {
			log.Infof("Recommend uid %s could not be loaded, skipping.", rUid)
			recommendNotLoaded++
			continue
		}
		log.Infof("%s %25s ==>", rUid, ContentItemType(unitsMap[rUid]))
		for _, rToSPair := range rToSPairs {
			var sUid string
			selectedEntry := selectedMap[rToSPair.SelectedClientEventID]
			if selectData, err := GetSelectedData(selectedEntry); err != nil {
				return err
			} else {
				sUid = selectData.Uid
			}
			if _, ok := unitsMap[sUid]; !ok {
				log.Infof("Selected uid %s could not be loaded, skipping.", sUid)
				selectedNotLoaded++
				continue
			}
			log.Infof("\t%s %25s %15s", sUid, ContentItemType(unitsMap[sUid]), selectedEntry.CreatedAt.Format("2006-02-01 15:04"))
			rContentItemType := ContentItemType(unitsMap[rUid])
			sContentItemType := ContentItemType(unitsMap[sUid])
			if _, ok := rToSTypeMap[rContentItemType]; !ok {
				rToSTypeMap[rContentItemType] = make(map[string]int)
			}
			rToSTypeMap[rContentItemType][sContentItemType] = rToSTypeMap[rContentItemType][sContentItemType] + 1
			if _, ok := sToRTypeMap[sContentItemType]; !ok {
				sToRTypeMap[sContentItemType] = make(map[string]int)
			}
			sToRTypeMap[sContentItemType][rContentItemType] = sToRTypeMap[sContentItemType][rContentItemType] + 1
		}
	}

	log.Infof("Could not load recommended %d, selected %d", recommendNotLoaded, selectedNotLoaded)
	log.Info("Selected to Recommended types:")
	for rContentItemType, sTypeMap := range rToSTypeMap {
		log.Infof("%-25s => ", rContentItemType)
		keys := utils.StringKeys(sTypeMap)
		sort.Slice(keys, func(i, j int) bool {
			return sTypeMap[keys[i]] > sTypeMap[keys[j]]
		})
		for _, sType := range keys {
			log.Infof("\t\t%5d %25s", sTypeMap[sType], sType)
		}
	}
	log.Info("Recommended to Selected types:")
	for sContentItemType, rTypeMap := range sToRTypeMap {
		log.Infof("%-25s => ", sContentItemType)
		keys := utils.StringKeys(rTypeMap)
		sort.Slice(keys, func(i, j int) bool {
			return rTypeMap[keys[i]] > rTypeMap[keys[j]]
		})
		for _, rType := range keys {
			log.Infof("\t\t%5d %25s", rTypeMap[rType], rType)
		}
	}
	return nil
}

func uidsToItems(uidsPairs []RecommendUidsPair, unitsMap map[string]interface{}) []RecommendItemsPair {
	ret := []RecommendItemsPair(nil)
	for _, pair := range uidsPairs {
		rCu, rOk := unitsMap[pair.Recommended]
		if sCu, sOk := unitsMap[pair.Selected]; rOk && sOk {
			ret = append(ret, RecommendItemsPair{Recommended: rCu, Selected: sCu})
		} else {
			log.Warnf("Skipping uids pair due to one of the uids not loaded: %s, %s", pair.Recommended, pair.Selected)
		}
	}
	return ret
}

func Learn(prodChronicles bool, chroniclesUrl string) error {
	log.Infof("Reading recommendations...")

	var err error
	selected := cModels.EntrySlice(nil)
	// lastReadId := "22" // 2021-12-08
	// lastReadId := "23" // 2021-12-30
	lastReadId := "24" // 2022-01-21
	// lastReadId := "25" // 2022-02-12
	if prodChronicles {
		if selected, err = MakeScanner("recommend-selected", lastReadId, chroniclesUrl).ScanAll(); err != nil {
			return err
		}
	} else {
		if selected, err = LoadEntries("recommend-selected"); err != nil {
			return err
		}
	}
	// selectedMap := EntriesMap(selected)

	recommended := cModels.EntrySlice(nil)
	if prodChronicles {
		if recommended, err = MakeScanner("recommend", lastReadId, chroniclesUrl).ScanAll(); err != nil {
			return err
		}
	} else {
		if recommended, err = LoadEntries("recommend"); err != nil {
			return err
		}
	}
	recommendedMap := EntriesMap(recommended)

	noClientFlowId := 0
	flowNotFound := 0
	uidsToLoad := make(map[string]bool)
	sUids := make(map[string]int)
	sUidsCount := 0
	recommendedToSelected := make(map[string][]RecommendClientEventIDPair)
	rToSUidsMap := make(map[string]map[string]bool)
	selectedToRecommended := make(map[string][]RecommendClientEventIDPair)
	positiveUidsPairs := []RecommendUidsPair(nil)

	for _, selectEntry := range selected {
		// log.Infof("Entry: %+v", selectEntry)
		if !selectEntry.ClientFlowID.Valid || selectEntry.ClientFlowID.String == "" {
			noClientFlowId++
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
			if selectEntry.ClientFlowID.Valid && selectEntry.ClientFlowID.String != "" {
				if recommendEntry, ok := recommendedMap[selectEntry.ClientFlowID.String]; !ok {
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
						sUids[selectedData.Uid]++
						sUidsCount += 1
						uidsToLoad[selectedData.Uid] = true
						rUid := recommendData.RequestData.Options.Recommend.Uid
						positiveUidsPairs = append(positiveUidsPairs, RecommendUidsPair{Recommended: rUid, Selected: selectedData.Uid})
						uidsToLoad[rUid] = true
						recommendedToSelected[rUid] = append(recommendedToSelected[rUid], RecommendClientEventIDPair{recommendEntry.ClientEventID.String, selectEntry.ClientEventID.String})
						selectedToRecommended[selectedData.Uid] = append(selectedToRecommended[selectedData.Uid], RecommendClientEventIDPair{recommendEntry.ClientEventID.String, selectEntry.ClientEventID.String})
						// log.Infof("%+v %s => %s", selectEntry.CreatedAt, recommendData.RequestData.Options.Recommend.Uid, selectedData.Uid)
						if _, ok := rToSUidsMap[rUid]; !ok {
							rToSUidsMap[rUid] = make(map[string]bool)
						}
						rToSUidsMap[rUid][selectedData.Uid] = true
					}
				}
			}
		}
	}
	log.Infof("Recommend selected: %d, without flow: %d, flow not found: %d", len(selected), noClientFlowId, flowNotFound)

	// Random negative uids.
	/*randomSize := sUidsCount * NEGATIVE_MULTIPLIER
	randomUids, err := sampleRandomUids(sUids, randomSize)
	if err != nil {
		return err
	}
	for _, uid := range randomUids {
		uidsToLoad[uid] = true
	}*/

	// Negative pairs from logs.
	negativeUidsPairs, err := RecommendsToNegativeExamples(recommended, rToSUidsMap)
	if err != nil {
		return err
	}
	for _, pair := range negativeUidsPairs {
		uidsToLoad[pair.Recommended] = true
		uidsToLoad[pair.Selected] = true
	}

	log.Infof("Loading %d uids", len(uidsToLoad))
	unitsMap, err := LoadUids(uidsToLoad)
	if err != nil {
		return err
	}
	log.Infof("Loaded %d uids.", len(unitsMap))

	// Prints some statistical info
	//if err := PrintTypesMap(unitsMap, recommendedToSelected, selectedMap); err != nil {
	//	return err
	//}

	positiveItemPairs := uidsToItems(positiveUidsPairs, unitsMap)
	negativeItemPairs := uidsToItems(negativeUidsPairs, unitsMap)
	// positiveUnits := uidsMapToUnits(sUids, unitsMap)
	// negativeUnits := uidsToUnits(randomUids, unitsMap)
	log.Infof("Learning classifier %d positives and %d negatives.", len(positiveItemPairs), len(negativeItemPairs))
	if len(negativeItemPairs) > len(positiveItemPairs)*NEGATIVE_MULTIPLIER {
		rand.Shuffle(len(negativeItemPairs), func(i, j int) {
			negativeItemPairs[i], negativeItemPairs[j] = negativeItemPairs[j], negativeItemPairs[i]
		})
		negativeItemPairs = negativeItemPairs[:len(positiveItemPairs)*NEGATIVE_MULTIPLIER]
		log.Infof("Using only %d negatives.", len(negativeItemPairs))
	}
	return CrossValidate(positiveItemPairs, negativeItemPairs)
}
