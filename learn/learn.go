package learn

import (
	"encoding/json"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/core"
	cModels "github.com/Bnei-Baruch/feed-api/databases/chronicles/models"
	"github.com/Bnei-Baruch/feed-api/databases/mdb"
	mdbModels "github.com/Bnei-Baruch/feed-api/databases/mdb/models"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type SelectedData struct {
	Uid string `json:"uid,omitempty"`
}

type RecommendData struct {
	RequestData core.MoreRequest `json:"request_data,ompitempty"`
}

func LoadUids(uidsToLoad map[string]bool) (unitsMap map[string]interface{}, err error) {
	uids := make(map[string]interface{})
	// Content units
	var units mdbModels.ContentUnitSlice
	if units, err = mdbModels.ContentUnits(qm.WhereIn("uid in ?", utils.ToInterfaceSlice(utils.StringKeys(uidsToLoad))...)).All(common.LocalMdb); err != nil {
		return nil, err
	}
	for _, unit := range units {
		uids[unit.UID] = unit
	}
	// Collections
	var collections mdbModels.CollectionSlice
	if collections, err = mdbModels.Collections(qm.WhereIn("uid in ?", utils.ToInterfaceSlice(utils.StringKeys(uidsToLoad))...)).All(common.LocalMdb); err != nil {
		return nil, err
	}
	for _, collection := range collections {
		uids[collection.UID] = collection
	}
	// Tags
	var tags mdbModels.TagSlice
	if tags, err = mdbModels.Tags(qm.WhereIn("uid in ?", utils.ToInterfaceSlice(utils.StringKeys(uidsToLoad))...)).All(common.LocalMdb); err != nil {
		return nil, err
	}
	for _, tag := range tags {
		uids[tag.UID] = tag
	}
	// Sources
	var sources mdbModels.SourceSlice
	if sources, err = mdbModels.Sources(qm.WhereIn("uid in ?", utils.ToInterfaceSlice(utils.StringKeys(uidsToLoad))...)).All(common.LocalMdb); err != nil {
		return nil, err
	}
	for _, source := range sources {
		uids[source.UID] = source
	}
	return uids, nil
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
	return recommendData, err
}

type RecommendPair struct {
	RecommendedClientEventID string
	SelectedClientEventID    string
}

func Learn(prodChronicles bool, chroniclesUrl string) error {
	log.Infof("Reading recommendations...")

	var err error
	selected := cModels.EntrySlice(nil)
	if prodChronicles {
		if selected, err = MakeScanner("recommend-selected", "22" /* 2022-01-03 lastReadId*/, chroniclesUrl).ScanAll(); err != nil {
			return err
		}
	} else {
		if selected, err = LoadEntries("recommend-selected"); err != nil {
			return err
		}
	}
	selectedMap := EntriesMap(selected)

	recommended := cModels.EntrySlice(nil)
	if prodChronicles {
		if recommended, err = MakeScanner("recommend", "22" /* 2022-01-03 lastReadId*/, chroniclesUrl).ScanAll(); err != nil {
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
	recommendedToSelected := make(map[string][]RecommendPair)
	selectedToRecommended := make(map[string][]RecommendPair)
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
					log.Infof("Could not find recommend for recommend selected. FlowClientID: %s", selectEntry.ClientFlowID.String)
				} else {
					if recommendData, err := GetRecommendData(recommendEntry); err != nil {
						log.Warnf("Failed getting recommend data, skipping.", err)
						continue
					} else {
						uidsToLoad[selectedData.Uid] = true
						rUid := recommendData.RequestData.Options.Recommend.Uid
						uidsToLoad[rUid] = true
						recommendedToSelected[rUid] = append(recommendedToSelected[rUid], RecommendPair{recommendEntry.ClientEventID.String, selectEntry.ClientEventID.String})
						selectedToRecommended[selectedData.Uid] = append(selectedToRecommended[selectedData.Uid], RecommendPair{recommendEntry.ClientEventID.String, selectEntry.ClientEventID.String})
						// log.Infof("%+v %s => %s", selectEntry.CreatedAt, recommendData.RequestData.Options.Recommend.Uid, selectedData.Uid)
					}
				}
			}
		}
	}
	log.Infof("Recommend selected: %d, without flow: %d, flow not found: %d", len(selected), noClientFlowId, flowNotFound)

	rToSTypeMap := make(map[string]map[string]int)
	sToRTypeMap := make(map[string]map[string]int)

	log.Infof("Loading %d uids", len(uidsToLoad))
	log.Infof("%+v", uidsToLoad)
	recommendNotLoaded := 0
	selectedNotLoaded := 0
	if uids, err := LoadUids(uidsToLoad); err != nil {
		return err
	} else {
		log.Infof("Loaded %d uids.", len(uids))

		for rUid, rToSPairs := range recommendedToSelected {
			if _, ok := uids[rUid]; !ok {
				log.Infof("Recommend uid %s could not be loaded, skipping.", rUid)
				recommendNotLoaded++
				continue
			}
			log.Infof("%s %25s ==>", rUid, ContentItemType(uids[rUid]))
			for _, rToSPair := range rToSPairs {
				var sUid string
				selectedEntry := selectedMap[rToSPair.SelectedClientEventID]
				if selectData, err := GetSelectedData(selectedEntry); err != nil {
					return err
				} else {
					sUid = selectData.Uid
				}
				if _, ok := uids[sUid]; !ok {
					log.Infof("Selected uid %s could not be loaded, skipping.", sUid)
					selectedNotLoaded++
					continue
				}
				log.Infof("\t%s %25s %15s", sUid, ContentItemType(uids[sUid]), selectedEntry.CreatedAt.Format("2006-02-01 15:04"))
				rContentItemType := ContentItemType(uids[rUid])
				sContentItemType := ContentItemType(uids[sUid])
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
