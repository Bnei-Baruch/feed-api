package recommendations

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/queries"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/data_models"
	"github.com/Bnei-Baruch/feed-api/databases/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type DataContentUnitsSuggester struct {
	suggesterContext core.SuggesterContext
	filters          []core.SuggesterFilter
	orderSelector    core.OrderSelectorEnum
}

func MakeDataContentUnitsSuggester(suggesterContext core.SuggesterContext, filters []core.SuggesterFilter, orderSelector core.OrderSelectorEnum) *DataContentUnitsSuggester {
	return &DataContentUnitsSuggester{suggesterContext, filters, orderSelector}
}

func (s *DataContentUnitsSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: "DataContentUnitsSuggester", Filters: s.filters, OrderSelector: s.orderSelector}, nil
}

func (s *DataContentUnitsSuggester) UnmarshalSpec(suggesterContext core.SuggesterContext, spec core.SuggesterSpec) error {
	if spec.Name != "DataContentUnitsSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'DataContentUnitsSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("DataContentUnitsSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	s.filters = spec.Filters
	s.orderSelector = spec.OrderSelector
	return nil
}

type ContentUnitRecommendInfo struct {
	data_models.ContentUnitInfo
	Tags    []string
	Sources []string
}

func LoadContentUnitRecommendInfo(uid string, suggesterContext core.SuggesterContext) (*ContentUnitRecommendInfo, error) {
	if _, ok := suggesterContext.Cache["recommendInfo"]; !ok {
		sqlStr := fmt.Sprintf(`
			select
				cu.type_id,
				(coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text))::date as date,
				cu.created_at,
				array_agg(distinct t.uid) as tags,
				array_agg(distinct s.uid) as sources
			from
				content_units as cu
			left outer join
				content_units_tags as cut
			on
				cu.id = cut.content_unit_id
			left outer join
				tags as t
			on
				t.id = cut.tag_id
			left outer join
				content_units_sources as cus
			on
				cu.id = cus.content_unit_id
			left outer join
				sources as s
			on
				s.id = cus.source_id
			where
				cu.uid = '%s'
			group by
				cu.type_id, date, cu.created_at;
		`, uid)
		rows, err := queries.Raw(sqlStr).Query(suggesterContext.DB)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		first := true
		var typeId int64
		var date time.Time
		var createdAt time.Time
		tags := []sql.NullString(nil)
		sources := []sql.NullString(nil)
		for rows.Next() {
			if !first {
				return nil, errors.New(fmt.Sprintf("Expected only one row for content unit uid: %s", uid))
			}
			err := rows.Scan(&typeId, &date, &createdAt, pq.Array(&tags), pq.Array(&sources))
			if err != nil {
				return nil, err
			}
			first = false
		}
		suggesterContext.Cache["recommendInfo"] = &ContentUnitRecommendInfo{
			data_models.ContentUnitInfo{typeId, uid, date, createdAt, true, false, 0},
			utils.NullStringSliceToStringSlice(tags),
			utils.NullStringSliceToStringSlice(sources),
		}
	}
	return suggesterContext.Cache["recommendInfo"].(*ContentUnitRecommendInfo), nil
}

func (s *DataContentUnitsSuggester) initialUids(request core.MoreRequest) map[string]bool {
	init_start := time.Now()
	defer func() {
		utils.Profile("DataContentUnitsSuggester.More.Init", time.Now().Sub(init_start))
	}()
	var ok bool
	if _, ok = s.suggesterContext.Cache["uids"]; !ok {
		prefiltered_start := time.Now()
		dm := s.suggesterContext.DataModels
		uids := utils.CopyStringMap(dm.ContentUnitsInfo.Prefiltered())
		log.Debugf("prefilter uids: %d", len(uids))
		utils.Profile("DataContentUnitsSuggester.More.Prefiltered", time.Now().Sub(prefiltered_start))

		language_start := time.Now()
		utils.IntersectMaps(uids, dm.LanguagesContentUnitsFilter.FilterValues(request.Options.Languages))
		utils.Profile("DataContentUnitsSuggester.More.Language", time.Now().Sub(language_start))

		persons_start := time.Now()
		utils.IntersectMaps(uids, dm.PersonsContentUnitsFilter.FilterValues([]string{data_models.RAV_PERSON_UID}))
		utils.Profile("DataContentUnitsSuggester.More.Persons", time.Now().Sub(persons_start))
		s.suggesterContext.Cache["uids"] = uids
	}
	return utils.CopyStringMap(s.suggesterContext.Cache["uids"].(map[string]bool))
}

func (s *DataContentUnitsSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	start := time.Now()
	defer func() {
		utils.Profile("DataContentUnitsSuggester.More", time.Now().Sub(start))
	}()
	if recommendInfo, err := LoadContentUnitRecommendInfo(request.Options.Recommend.Uid, s.suggesterContext); err != nil {
		return nil, err
	} else {
		dm := s.suggesterContext.DataModels
		uids := s.initialUids(request)

		log.Debugf("before skip uids: %d", len(uids))
		skip_start := time.Now()
		for _, uid := range request.Options.SkipUids {
			delete(uids, uid)
		}
		delete(uids, recommendInfo.Uid)
		utils.Profile("DataContentUnitsSuggester.More.Skip", time.Now().Sub(skip_start))

		filters_start := time.Now()
		suggesterNameParts := []string{"DataContentUnitsSuggester"}
		log.Debugf("before filter uids: %d", len(uids))
		for _, filter := range s.filters {
			filter_start := time.Now()
			switch filter.FilterSelector {
			case core.UnitContentTypes:
				suggesterNameParts = append(suggesterNameParts, ";UnitContentTypes:[", strings.Join(filter.Args, ","), "]")
				int64Args := core.ContentTypesToInt64Ids(filter.Args)
				utils.FilterMap(uids, func(uid string) bool {
					return utils.Int64InSlice(dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo).TypeId, int64Args)
				})
			case core.CollectionContentTypes:
				suggesterNameParts = append(suggesterNameParts, ";CollectionContentTypes:[", strings.Join(filter.Args, ","), "]")
				int64Args := core.ContentTypesToInt64Ids(filter.Args)
				utils.FilterMap(uids, func(uid string) bool {
					collectionsUids := dm.ContentUnitsCollectionsFilter.FilterValues([]string{uid})
					for collectionUid, _ := range collectionsUids {
						if utils.Int64InSlice(dm.CollectionsInfo.Data(collectionUid).(*data_models.CollectionInfo).TypeId, int64Args) {
							return true
						}
					}
					return false
				})
			case core.Tags:
				suggesterNameParts = append(suggesterNameParts, ";Tags:[", strings.Join(filter.Args, ","), "]")
				utils.IntersectMaps(uids, dm.TagsContentUnitsFilter.FilterValues(filter.Args))
			case core.Sources:
				suggesterNameParts = append(suggesterNameParts, ";Sources:[", strings.Join(filter.Args, ","), "]")
				utils.IntersectMaps(uids, dm.SourcesContentUnitsFilter.FilterValues(filter.Args))
			case core.Collections:
				suggesterNameParts = append(suggesterNameParts, ";Collections:[", strings.Join(filter.Args, ","), "]")
				utils.IntersectMaps(uids, dm.CollectionsContentUnitsFilter.FilterValues(filter.Args))
			case core.SameTag:
				suggesterNameParts = append(suggesterNameParts, ";SameTag:[", strings.Join(filter.Args, ","), "]")
				utils.IntersectMaps(uids, dm.TagsContentUnitsFilter.FilterValues(recommendInfo.Tags))
			case core.SameCollection:
				suggesterNameParts = append(suggesterNameParts, ";SameCollection:[", strings.Join(filter.Args, ","), "]")
				utils.IntersectMaps(uids, dm.CollectionsContentUnitsFilter.FilterValues(utils.SliceFromMap(dm.ContentUnitsCollectionsFilter.FilterValues([]string{request.Options.Recommend.Uid}))))
			case core.SameSource:
				suggesterNameParts = append(suggesterNameParts, ";SameSource:[", strings.Join(filter.Args, ","), "]")
				utils.IntersectMaps(uids, dm.SourcesContentUnitsFilter.FilterValues(recommendInfo.Sources))
			case core.WatchingNowFilter:
				suggesterNameParts = append(suggesterNameParts, ";WatchingNowFilter:[", strings.Join(filter.Args, ","), "]")
				if watchingNow, err := dm.SqlDataModel.AllWatchingNow(); err != nil {
					return nil, err
				} else {
					watchingNowUids := make(map[string]bool, len(watchingNow))
					for uid, count := range watchingNow {
						if count > 0 {
							watchingNowUids[uid] = true
						}
					}
					utils.IntersectMaps(uids, watchingNowUids)
				}
			default:
				log.Errorf("Did not expect filter selector enum %d", filter.FilterSelector)
			}
			utils.Profile(fmt.Sprintf("DataContentUnitsSuggester.More.Filter[%d]", filter.FilterSelector), time.Now().Sub(filter_start))
		}
		log.Debugf("Debug timestamp: %+v", request.Options.DebugTimestamp)
		if request.Options.DebugTimestamp != nil {
			from := time.Unix(*request.Options.DebugTimestamp, 0)
			log.Debugf("From: %+v", from)
			l := len(uids)
			utils.FilterMap(uids, func(uid string) bool {
				return from.After(dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo).Date)
			})
			log.Debugf("Debug timestamp from %d => %d", l, len(uids))
		}
		utils.Profile("DataContentUnitsSuggester.More.Filters", time.Now().Sub(filters_start))
		sort_start := time.Now()
		uidsSlice := utils.SliceFromMap(uids)
		log.Debugf("uidsSlice: %d", len(uidsSlice))
		sortLastPrev := func(i, j int) bool {
			icu := dm.ContentUnitsInfo.Data(uidsSlice[i]).(*data_models.ContentUnitInfo)
			jcu := dm.ContentUnitsInfo.Data(uidsSlice[j]).(*data_models.ContentUnitInfo)
			if icu.Date.Equal(jcu.Date) {
				if icu.CreatedAt.Equal(jcu.CreatedAt) {
					return false
				}
				return icu.CreatedAt.After(jcu.CreatedAt)
			}
			return icu.Date.After(jcu.Date)
		}
		switch s.orderSelector {
		case core.Last:
			suggesterNameParts = append(suggesterNameParts, ";Last")
			sort.SliceStable(uidsSlice, sortLastPrev)
		case core.Next:
			suggesterNameParts = append(suggesterNameParts, ";Next")
			uidsSlice = utils.Filter(uidsSlice, func(uid string) bool {
				info := dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo)
				return recommendInfo.Date.Before(info.Date) || (recommendInfo.Date.Equal(info.Date) && recommendInfo.CreatedAt.Before(info.CreatedAt))
			})
			sort.SliceStable(uidsSlice, func(i, j int) bool {
				icu := dm.ContentUnitsInfo.Data(uidsSlice[i]).(*data_models.ContentUnitInfo)
				jcu := dm.ContentUnitsInfo.Data(uidsSlice[j]).(*data_models.ContentUnitInfo)
				if icu.Date.Equal(jcu.Date) {
					if icu.CreatedAt.Equal(jcu.CreatedAt) {
						return false
					}
					return icu.CreatedAt.Before(jcu.CreatedAt)
				}
				return icu.Date.Before(jcu.Date)
			})
		case core.Prev:
			suggesterNameParts = append(suggesterNameParts, ";Prev")
			uidsSlice = utils.Filter(uidsSlice, func(uid string) bool {
				info := dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo)
				return recommendInfo.Date.After(info.Date) || (recommendInfo.Date.Equal(info.Date) && recommendInfo.CreatedAt.After(info.CreatedAt))
			})
			sort.SliceStable(uidsSlice, sortLastPrev)
		case core.Rand:
			suggesterNameParts = append(suggesterNameParts, ";Rand")
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(uidsSlice), func(i, j int) { uidsSlice[i], uidsSlice[j] = uidsSlice[j], uidsSlice[i] })
		case core.Popular:
			suggesterNameParts = append(suggesterNameParts, ";Popular")
			if err := dm.SqlDataModel.SortPopular(uidsSlice); err != nil {
				return nil, err
			}
		case core.WatchingNow:
			suggesterNameParts = append(suggesterNameParts, ";WatchingNow")
			if err := dm.SqlDataModel.SortWatchingNow(uidsSlice); err != nil {
				return nil, err
			}
		}
		utils.Profile("DataContentUnitsSuggester.More.Sort", time.Now().Sub(sort_start))
		ret_start := time.Now()
		ret := []core.ContentItem(nil)
		if request.MoreItems <= 0 {
			return ret, nil
		}
		for _, uid := range uidsSlice {
			cuInfo := dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo)
			contentType := mdb.CONTENT_TYPE_REGISTRY.ByID[cuInfo.TypeId].Name
			ret = append(ret, core.ContentItem{UID: uid, Date: cuInfo.Date, CreatedAt: cuInfo.CreatedAt, ContentType: contentType, Suggester: strings.Join(suggesterNameParts, "")})
			if len(ret) >= request.MoreItems {
				break
			}
		}
		utils.Profile("DataContentUnitsSuggester.More.BuildRet", time.Now().Sub(ret_start))
		return ret, nil
	}
}

type DataCollectionsSuggester struct {
	suggesterContext core.SuggesterContext
	filters          []core.SuggesterFilter
	orderSelector    core.OrderSelectorEnum
}

func MakeDataCollectionsSuggester(suggesterContext core.SuggesterContext, filters []core.SuggesterFilter, orderSelector core.OrderSelectorEnum) *DataCollectionsSuggester {
	return &DataCollectionsSuggester{suggesterContext, filters, orderSelector}
}

func (s *DataCollectionsSuggester) MarshalSpec() (core.SuggesterSpec, error) {
	return core.SuggesterSpec{Name: "DataCollectionsSuggester", Filters: s.filters, OrderSelector: s.orderSelector}, nil
}

func (s *DataCollectionsSuggester) UnmarshalSpec(suggesterContext core.SuggesterContext, spec core.SuggesterSpec) error {
	if spec.Name != "DataCollectionsSuggester" {
		return errors.New(fmt.Sprintf("Expected suggester name to be: 'DataCollectionsSuggester', got: '%s'.", spec.Name))
	}
	if len(spec.Specs) != 0 {
		return errors.New(fmt.Sprintf("DataCollectionsSuggester expected to have no suggesters, got %d.", len(spec.Specs)))
	}
	s.filters = spec.Filters
	s.orderSelector = spec.OrderSelector
	return nil
}

func (s *DataCollectionsSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	if recommendInfo, err := LoadContentUnitRecommendInfo(request.Options.Recommend.Uid, s.suggesterContext); err != nil {
		return nil, err
	} else {
		dm := s.suggesterContext.DataModels
		uids := dm.CollectionsInfo.Keys()
		uidsMap := make(map[string]bool, len(uids))
		for _, uid := range uids {
			uidsMap[uid] = true
		}
		utils.FilterMap(uidsMap, func(uid string) bool {
			return uid != recommendInfo.Uid && !utils.StringInSlice(uid, request.Options.SkipUids)
		})
		suggesterNameParts := []string{"DataCollectionsSuggester"}
		for _, filter := range s.filters {
			switch filter.FilterSelector {
			case core.UnitContentTypes:
				log.Errorf("Did not expect UnitContentType filter for DataCollectionsSuggester")
			case core.CollectionContentTypes:
				suggesterNameParts = append(suggesterNameParts, ";CollectionContentTypes:[", strings.Join(filter.Args, ","), "]")
				int64Args := core.ContentTypesToInt64Ids(filter.Args)
				utils.FilterMap(uidsMap, func(uid string) bool {
					return utils.Int64InSlice(dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo).TypeId, int64Args)
				})
			case core.Tags:
				log.Errorf("Did not expect Tags filter for DataCollectionsSuggester")
			case core.Sources:
				log.Errorf("Did not expect Sources filter for DataCollectionsSuggester")
			case core.Collections:
				suggesterNameParts = append(suggesterNameParts, ";Collections:[", strings.Join(filter.Args, ","), "]")
				utils.FilterMap(uidsMap, func(uid string) bool {
					return utils.StringInSlice(uid, filter.Args)
				})
			case core.SameTag:
				log.Errorf("Did not expect SameTag filter for DataCollectionsSuggester")
			case core.SameCollection:
				suggesterNameParts = append(suggesterNameParts, ";SameCollection:[", strings.Join(filter.Args, ","), "]")
				utils.IntersectMaps(uidsMap, dm.ContentUnitsCollectionsFilter.FilterValues([]string{request.Options.Recommend.Uid}))
			case core.SameSource:
				suggesterNameParts = append(suggesterNameParts, ";SameSource:[", strings.Join(filter.Args, ","), "]")
				utils.FilterMap(uidsMap, func(uid string) bool {
					return utils.StringInSlice(dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo).SourceUid, recommendInfo.Sources)
				})
			default:
				log.Errorf("Did not expect filter selector enum %d", filter.FilterSelector)
			}
		}
		log.Debugf("Debug timestamp: %+v", request.Options.DebugTimestamp)
		if request.Options.DebugTimestamp != nil {
			from := time.Unix(*request.Options.DebugTimestamp, 0)
			log.Debugf("From: %+v", from)
			l := len(uids)
			utils.FilterMap(uidsMap, func(uid string) bool {
				return from.After(dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo).Date)
			})
			log.Debugf("Debug timestamp from %d => %d", l, len(uids))
		}
		uids = utils.SliceFromMap(uidsMap)
		sortLastPrev := func(i, j int) bool {
			icu := dm.CollectionsInfo.Data(uids[i]).(*data_models.CollectionInfo)
			jcu := dm.CollectionsInfo.Data(uids[j]).(*data_models.CollectionInfo)
			if icu.Date.Equal(jcu.Date) {
				if icu.CreatedAt.Equal(jcu.CreatedAt) {
					return false
				}
				return icu.CreatedAt.After(jcu.CreatedAt)
			}
			return icu.Date.After(jcu.Date)
		}
		switch s.orderSelector {
		case core.Last:
			suggesterNameParts = append(suggesterNameParts, ";Last")
			sort.SliceStable(uids, sortLastPrev)
		case core.Next:
			log.Errorf("Did not expect Next order for DataCollectionsSuggester")
		case core.Prev:
			log.Errorf("Did not expect Prev order for DataCollectionsSuggester")
		case core.Rand:
			suggesterNameParts = append(suggesterNameParts, ";Rand")
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(uids), func(i, j int) { uids[i], uids[j] = uids[j], uids[i] })
		}
		ret := []core.ContentItem(nil)
		if request.MoreItems <= 0 {
			return ret, nil
		}
		for _, uid := range uids {
			cInfo := dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo)
			contentType := mdb.CONTENT_TYPE_REGISTRY.ByID[cInfo.TypeId].Name
			ret = append(ret, core.ContentItem{UID: uid, Date: cInfo.Date, CreatedAt: cInfo.CreatedAt, ContentType: contentType, Suggester: strings.Join(suggesterNameParts, "")})
			if len(ret) >= request.MoreItems {
				break
			}
		}
		return ret, nil
	}
}

func init() {
	core.RegisterSuggester("DataContentUnitsSuggester", func(suggesterContext core.SuggesterContext) core.Suggester {
		return MakeDataContentUnitsSuggester(suggesterContext, []core.SuggesterFilter(nil), core.Last)
	})
	core.RegisterSuggester("DataCollectionsSuggester", func(suggesterContext core.SuggesterContext) core.Suggester {
		return MakeDataCollectionsSuggester(suggesterContext, []core.SuggesterFilter(nil), core.Last)
	})
}
