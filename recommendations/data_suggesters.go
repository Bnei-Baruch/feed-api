package recommendations

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/queries"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/data_models"
	"github.com/Bnei-Baruch/feed-api/mdb"
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
		rows, err := queries.Raw(suggesterContext.DB, sqlStr).Query()
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

func (s *DataContentUnitsSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	if recommendInfo, err := LoadContentUnitRecommendInfo(request.Options.Recommend.Uid, s.suggesterContext); err != nil {
		return nil, err
	} else {
		dm := s.suggesterContext.DataModels
		uids := dm.ContentUnitsInfo.Keys()
		// Filter published and secure.
		uids = utils.Filter(uids, func(uid string) bool {
			return dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo).SecureAndPublished &&
				!dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo).IsLessonPrep &&
				uid != recommendInfo.Uid &&
				!utils.StringInSlice(uid, request.Options.SkipUids)
		})
		uids = utils.IntersectSorted(uids, dm.PersonsContentUnitsFilter.FilterValues([]string{"abcdefgh"}))
		uids = utils.IntersectSorted(uids, dm.LanguagesContentUnitsFilter.FilterValues(request.Options.Languages))
		suggesterNameParts := []string{"DataContentUnitsSuggester"}
		for _, filter := range s.filters {
			switch filter.FilterSelector {
			case core.UnitContentTypes:
				suggesterNameParts = append(suggesterNameParts, ";UnitContentTypes:[", strings.Join(filter.Args, ","), "]")
				int64Args := core.ContentTypesToInt64Ids(filter.Args)
				uids = utils.Filter(uids, func(uid string) bool {
					return utils.Int64InSlice(dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo).TypeId, int64Args)
				})
			case core.CollectionContentTypes:
				suggesterNameParts = append(suggesterNameParts, ";CollectionContentTypes:[", strings.Join(filter.Args, ","), "]")
				int64Args := core.ContentTypesToInt64Ids(filter.Args)
				uids = utils.Filter(uids, func(uid string) bool {
					collectionsUids := dm.ContentUnitsCollectionsFilter.FilterValues([]string{uid})
					for _, collectionUid := range collectionsUids {
						if utils.Int64InSlice(dm.CollectionsInfo.Data(collectionUid).(*data_models.CollectionInfo).TypeId, int64Args) {
							return true
						}
					}
					return false
				})
			case core.Tags:
				suggesterNameParts = append(suggesterNameParts, ";Tags:[", strings.Join(filter.Args, ","), "]")
				uids = utils.IntersectSorted(uids, dm.TagsContentUnitsFilter.FilterValues(filter.Args))
			case core.Sources:
				suggesterNameParts = append(suggesterNameParts, ";Sources:[", strings.Join(filter.Args, ","), "]")
				uids = utils.IntersectSorted(uids, dm.SourcesContentUnitsFilter.FilterValues(filter.Args))
			case core.Collections:
				suggesterNameParts = append(suggesterNameParts, ";Collections:[", strings.Join(filter.Args, ","), "]")
				uids = utils.IntersectSorted(uids, dm.CollectionsContentUnitsFilter.FilterValues(filter.Args))
			case core.SameTag:
				suggesterNameParts = append(suggesterNameParts, ";SameTag:[", strings.Join(filter.Args, ","), "]")
				uids = utils.IntersectSorted(uids, dm.TagsContentUnitsFilter.FilterValues(recommendInfo.Tags))
			case core.SameCollection:
				suggesterNameParts = append(suggesterNameParts, ";SameCollection:[", strings.Join(filter.Args, ","), "]")
				uids = utils.IntersectSorted(uids, dm.CollectionsContentUnitsFilter.FilterValues(dm.ContentUnitsCollectionsFilter.FilterValues([]string{request.Options.Recommend.Uid})))
			case core.SameSource:
				suggesterNameParts = append(suggesterNameParts, ";SameSource:[", strings.Join(filter.Args, ","), "]")
				uids = utils.IntersectSorted(uids, dm.SourcesContentUnitsFilter.FilterValues(recommendInfo.Sources))
			default:
				log.Errorf("Did not expect filter selector enum %d", filter.FilterSelector)
			}
		}
		sortLastPrev := func(i, j int) bool {
			icu := dm.ContentUnitsInfo.Data(uids[i]).(*data_models.ContentUnitInfo)
			jcu := dm.ContentUnitsInfo.Data(uids[j]).(*data_models.ContentUnitInfo)
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
			suggesterNameParts = append(suggesterNameParts, ";Next")
			uids = utils.Filter(uids, func(uid string) bool {
				info := dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo)
				return recommendInfo.Date.Before(info.Date) || (recommendInfo.Date.Equal(info.Date) && recommendInfo.CreatedAt.Before(info.CreatedAt))
			})
			sort.SliceStable(uids, func(i, j int) bool {
				icu := dm.ContentUnitsInfo.Data(uids[i]).(*data_models.ContentUnitInfo)
				jcu := dm.ContentUnitsInfo.Data(uids[j]).(*data_models.ContentUnitInfo)
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
			uids = utils.Filter(uids, func(uid string) bool {
				info := dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo)
				return recommendInfo.Date.After(info.Date) || (recommendInfo.Date.Equal(info.Date) && recommendInfo.CreatedAt.After(info.CreatedAt))
			})
			sort.SliceStable(uids, sortLastPrev)
		case core.Rand:
			suggesterNameParts = append(suggesterNameParts, ";Rand")
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(uids), func(i, j int) { uids[i], uids[j] = uids[j], uids[i] })
		case core.Popular:
			suggesterNameParts = append(suggesterNameParts, ";Popular")
			sort.SliceStable(uids, func(i, j int) bool {
				var iWatch *data_models.ContentUnitsWatchDuration
				if dm.ContentUnitsWatchDuration.Data(uids[i]) == nil {
					iWatch = &data_models.ContentUnitsWatchDuration{uids[i], 0, 0}
				} else {
					iWatch = dm.ContentUnitsWatchDuration.Data(uids[i]).(*data_models.ContentUnitsWatchDuration)
				}
				var jWatch *data_models.ContentUnitsWatchDuration
				if dm.ContentUnitsWatchDuration.Data(uids[j]) == nil {
					jWatch = &data_models.ContentUnitsWatchDuration{uids[j], 0, 0}
				} else {
					jWatch = dm.ContentUnitsWatchDuration.Data(uids[j]).(*data_models.ContentUnitsWatchDuration)
				}
				return iWatch.Count > jWatch.Count
			})
		}
		ret := []core.ContentItem(nil)
		if request.MoreItems <= 0 {
			return ret, nil
		}
		for _, uid := range uids {
			cuInfo := dm.ContentUnitsInfo.Data(uid).(*data_models.ContentUnitInfo)
			contentType := mdb.CONTENT_TYPE_REGISTRY.ByID[cuInfo.TypeId].Name
			ret = append(ret, core.ContentItem{UID: uid, Date: cuInfo.Date, CreatedAt: cuInfo.CreatedAt, ContentType: contentType, Suggester: strings.Join(suggesterNameParts, "")})
			if len(ret) >= request.MoreItems {
				break
			}
		}
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
		// Filter published and secure, remove skipped and current uids.
		uids = utils.Filter(uids, func(uid string) bool {
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
				uids = utils.Filter(uids, func(uid string) bool {
					return utils.Int64InSlice(dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo).TypeId, int64Args)
				})
			case core.Tags:
				log.Errorf("Did not expect Tags filter for DataCollectionsSuggester")
			case core.Sources:
				log.Errorf("Did not expect Sources filter for DataCollectionsSuggester")
			case core.Collections:
				suggesterNameParts = append(suggesterNameParts, ";Collections:[", strings.Join(filter.Args, ","), "]")
				uids = utils.Filter(uids, func(uid string) bool {
					return utils.StringInSlice(uid, filter.Args)
				})
			case core.SameTag:
				log.Errorf("Did not expect SameTag filter for DataCollectionsSuggester")
			case core.SameCollection:
				suggesterNameParts = append(suggesterNameParts, ";SameCollection:[", strings.Join(filter.Args, ","), "]")
				uids = utils.IntersectSorted(uids, dm.ContentUnitsCollectionsFilter.FilterValues([]string{request.Options.Recommend.Uid}))
			case core.SameSource:
				suggesterNameParts = append(suggesterNameParts, ";SameSource:[", strings.Join(filter.Args, ","), "]")
				uids = utils.Filter(uids, func(uid string) bool {
					return utils.StringInSlice(dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo).SourceUid, recommendInfo.Sources)
				})
			default:
				log.Errorf("Did not expect filter selector enum %d", filter.FilterSelector)
			}
		}
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
