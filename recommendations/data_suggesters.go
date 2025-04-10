package recommendations

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries"

	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/data_models"
	"github.com/Bnei-Baruch/feed-api/databases/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

const (
	CONTENT_UNITS_UIDS_KEY = "uids"
	COLLECTIONS_UIDS_KEY   = "collection-uids"
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
	if _, ok = s.suggesterContext.Cache[CONTENT_UNITS_UIDS_KEY]; !ok {
		prefiltered_start := time.Now()
		dm := s.suggesterContext.DataModels
		uids := utils.CopyStringMap(dm.ContentUnitsInfo.Prefiltered())
		log.Debugf("prefilter uids: %d", len(uids))
		utils.Profile("DataContentUnitsSuggester.More.Prefiltered", time.Now().Sub(prefiltered_start))

		language_start := time.Now()
		utils.IntersectMaps(uids, dm.LanguagesContentUnitsFilter.FilterValues(request.Options.Languages))
		utils.Profile("DataContentUnitsSuggester.More.Language", time.Now().Sub(language_start))

		// We don't want to filter only Rav units, as many are without Rav, such as FRIENDS_GATHERING and other
		// are lessons when friends teach.
		// persons_start := time.Now()
		// utils.IntersectMaps(uids, dm.PersonsContentUnitsFilter.FilterValues([]string{data_models.RAV_PERSON_UID}))
		// utils.Profile("DataContentUnitsSuggester.More.Persons", time.Now().Sub(persons_start))

		if request.Options.WithPosts {
			// Add blog posts.
			blog_start := time.Now()
			blogPosts := utils.CopyStringMap(dm.BlogPostsInfo.Prefiltered())
			utils.Profile("DataContentUnitsSuggester.More.BlogPrefiltered", time.Now().Sub(blog_start))

			// Filter blogs by languages.
			blog_language := time.Now()
			utils.FilterMap(blogPosts, func(blogPostId string) bool {
				blogId := dm.BlogPostsInfo.Data(blogPostId).(*data_models.BlogPostInfo).BlogId
				if blog, ok := mdb.BLOGS_REGISTRY.ByID[blogId]; !ok {
					return false
				} else {
					lang, ok := consts.BLOGS_LANG[blog.Name]
					return ok && utils.StringInSlice(lang, request.Options.Languages)
				}
			})
			utils.Profile("DataContentUnitsSuggester.More.BlogLanguage", time.Now().Sub(blog_language))

			blog_union := time.Now()
			utils.UnionMaps(uids, blogPosts)
			utils.Profile("DataContentUnitsSuggester.More.BlogUnion", time.Now().Sub(blog_union))
		}

		s.suggesterContext.Cache[CONTENT_UNITS_UIDS_KEY] = uids
	}
	return utils.CopyStringMap(s.suggesterContext.Cache[CONTENT_UNITS_UIDS_KEY].(map[string]bool))
}

func (s *DataContentUnitsSuggester) ContentUnitType(uid string) int64 {
	dm := s.suggesterContext.DataModels
	if data := dm.ContentUnitsInfo.Data(uid); data != nil {
		return data.(*data_models.ContentUnitInfo).TypeId
	}
	return mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_BLOG_POST].ID
}

func (s *DataContentUnitsSuggester) ContentUnitInfo(uid string) *data_models.ContentUnitInfo {
	dm := s.suggesterContext.DataModels
	if data := dm.ContentUnitsInfo.Data(uid); data != nil {
		return data.(*data_models.ContentUnitInfo)
	}
	bpi := dm.BlogPostsInfo.Data(uid).(*data_models.BlogPostInfo)
	return &data_models.ContentUnitInfo{
		mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_BLOG_POST].ID,
		data_models.BlogPostKey(bpi),
		bpi.PostedAt,
		bpi.PostedAt,
		true,
		false,
		0,
	}
}

func (s *DataContentUnitsSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	log.Debugf("DataContentUnitsSuggester MoreRequest:%d  %+v  %+v", request.MoreItems, core.CurrentFeedsToUidsString(request.CurrentFeed), request.Options.SkipUids)
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
		for i := range request.CurrentFeed {
			delete(uids, request.CurrentFeed[i].UID)
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
					return utils.Int64InSlice(s.ContentUnitType(uid), int64Args)
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
					watchingNowMin := int64(0)
					if request.Options.WatchingNowMin != nil {
						watchingNowMin = *request.Options.WatchingNowMin
					}
					watchingNowUids := make(map[string]bool, len(watchingNow))
					for uid, count := range watchingNow {
						if count > watchingNowMin {
							watchingNowUids[uid] = true
						}
					}
					utils.IntersectMaps(uids, watchingNowUids)
				}
			case core.PopularFilter:
				suggesterNameParts = append(suggesterNameParts, ";PopularFilter:[", strings.Join(filter.Args, ","), "]")
				if views, err := dm.SqlDataModel.AllUniqueViews(); err != nil {
					return nil, err
				} else {
					popularMin := int64(0)
					if request.Options.PopularMin != nil {
						popularMin = *request.Options.PopularMin
					}
					popularUids := make(map[string]bool, len(views))
					for uid, count := range views {
						if count > popularMin {
							popularUids[uid] = true
						}
					}
					utils.IntersectMaps(uids, popularUids)
				}
			case core.AgeFilter:
				suggesterNameParts = append(suggesterNameParts, ";Age:[", strings.Join(filter.Args, ","), "]")
				if len(filter.Args) != 1 {
					return nil, errors.New(fmt.Sprintf("Expected age args be of length 1, got: '%d'.", len(filter.Args)))
				} else if ageSeconds, err := strconv.Atoi(filter.Args[0]); err != nil {
					return nil, errors.New(fmt.Sprintf("Expected age arg to be number as string, got: '%s'.", filter.Args[0]))
				} else {
					age := time.Now().Add(time.Duration(-ageSeconds) * time.Second)
					utils.FilterMap(uids, func(uid string) bool {
						return age.Before(s.ContentUnitInfo(uid).Date)
					})
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
				return from.After(s.ContentUnitInfo(uid).Date)
			})
			log.Debugf("Debug timestamp from %d => %d", l, len(uids))
		}
		utils.Profile("DataContentUnitsSuggester.More.Filters", time.Now().Sub(filters_start))
		sort_start := time.Now()
		uidsSlice := utils.SliceFromMap(uids)
		log.Debugf("uidsSlice: %d", len(uidsSlice))
		sortLastPrev := func(i, j int) bool {
			icu := s.ContentUnitInfo(uidsSlice[i])
			jcu := s.ContentUnitInfo(uidsSlice[j])
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
				info := s.ContentUnitInfo(uid)
				return recommendInfo.Date.Before(info.Date) || (recommendInfo.Date.Equal(info.Date) && recommendInfo.CreatedAt.Before(info.CreatedAt))
			})
			sort.SliceStable(uidsSlice, func(i, j int) bool {
				icu := s.ContentUnitInfo(uidsSlice[i])
				jcu := s.ContentUnitInfo(uidsSlice[j])
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
				info := s.ContentUnitInfo(uid)
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
			cuInfo := s.ContentUnitInfo(uid)
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

func (s *DataCollectionsSuggester) initialUids(request core.MoreRequest) map[string]bool {
	init_start := time.Now()
	defer func() {
		utils.Profile("DataCollectionsSuggester.More.Init", time.Now().Sub(init_start))
	}()
	var ok bool
	if _, ok = s.suggesterContext.Cache[COLLECTIONS_UIDS_KEY]; !ok {
		log.Debugf("Initialize collections cache")
		prefiltered_start := time.Now()
		dm := s.suggesterContext.DataModels
		uids := utils.CopyStringMap(dm.CollectionsInfo.Prefiltered())
		log.Debugf("prefilter uids: %d", len(uids))
		utils.Profile("DataCollectionsSuggester.More.Prefiltered", time.Now().Sub(prefiltered_start))
		s.suggesterContext.Cache[COLLECTIONS_UIDS_KEY] = uids
	}
	return utils.CopyStringMap(s.suggesterContext.Cache[COLLECTIONS_UIDS_KEY].(map[string]bool))
}

func (s *DataCollectionsSuggester) More(request core.MoreRequest) ([]core.ContentItem, error) {
	log.Debugf("DataCollectionsSuggester MoreRequest:%d  %+v %+v", request.MoreItems, core.CurrentFeedsToUidsString(request.CurrentFeed), request.Options.SkipUids)
	if recommendInfo, err := LoadContentUnitRecommendInfo(request.Options.Recommend.Uid, s.suggesterContext); err != nil {
		return nil, err
	} else {
		dm := s.suggesterContext.DataModels
		uids := s.initialUids(request)
		for _, uid := range request.Options.SkipUids {
			delete(uids, uid)
		}
		for i := range request.CurrentFeed {
			delete(uids, request.CurrentFeed[i].UID)
		}
		delete(uids, recommendInfo.Uid)
		suggesterNameParts := []string{"DataCollectionsSuggester"}
		for _, filter := range s.filters {
			switch filter.FilterSelector {
			case core.UnitContentTypes:
				log.Errorf("Did not expect UnitContentType filter for DataCollectionsSuggester")
			case core.CollectionContentTypes:
				suggesterNameParts = append(suggesterNameParts, ";CollectionContentTypes:[", strings.Join(filter.Args, ","), "]")
				int64Args := core.ContentTypesToInt64Ids(filter.Args)
				utils.FilterMap(uids, func(uid string) bool {
					return utils.Int64InSlice(dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo).TypeId, int64Args)
				})
			case core.Tags:
				log.Errorf("Did not expect Tags filter for DataCollectionsSuggester")
			case core.Sources:
				log.Errorf("Did not expect Sources filter for DataCollectionsSuggester")
			case core.Collections:
				suggesterNameParts = append(suggesterNameParts, ";Collections:[", strings.Join(filter.Args, ","), "]")
				utils.FilterMap(uids, func(uid string) bool {
					return utils.StringInSlice(uid, filter.Args)
				})
			case core.SameTag:
				log.Errorf("Did not expect SameTag filter for DataCollectionsSuggester")
			case core.SameCollection:
				suggesterNameParts = append(suggesterNameParts, ";SameCollection:[", strings.Join(filter.Args, ","), "]")
				utils.IntersectMaps(uids, dm.ContentUnitsCollectionsFilter.FilterValues([]string{request.Options.Recommend.Uid}))
			case core.SameSource:
				suggesterNameParts = append(suggesterNameParts, ";SameSource:[", strings.Join(filter.Args, ","), "]")
				utils.FilterMap(uids, func(uid string) bool {
					return utils.StringInSlice(dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo).SourceUid, recommendInfo.Sources)
				})
			case core.AgeFilter:
				suggesterNameParts = append(suggesterNameParts, ";Age:[", strings.Join(filter.Args, ","), "]")
				if len(filter.Args) != 1 {
					return nil, errors.New(fmt.Sprintf("Expected age args be of length 1, got: '%d'.", len(filter.Args)))
				} else if ageSeconds, err := strconv.Atoi(filter.Args[0]); err != nil {
					return nil, errors.New(fmt.Sprintf("Expected age arg to be number as string, got: '%s'.", filter.Args[0]))
				} else {
					age := time.Now().Add(time.Duration(-ageSeconds) * time.Second)
					utils.FilterMap(uids, func(uid string) bool {
						return age.Before(dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo).Date)
					})
				}
			default:
				log.Errorf("Did not expect filter selector enum %d", filter.FilterSelector)
			}
		}
		log.Debugf("Debug timestamp: %+v", request.Options.DebugTimestamp)
		if request.Options.DebugTimestamp != nil {
			from := time.Unix(*request.Options.DebugTimestamp, 0)
			log.Debugf("From: %+v", from)
			l := len(uids)
			utils.FilterMap(uids, func(uid string) bool {
				return from.After(dm.CollectionsInfo.Data(uid).(*data_models.CollectionInfo).Date)
			})
			log.Debugf("Debug timestamp from %d => %d", l, len(uids))
		}
		uidsSlice := utils.SliceFromMap(uids)
		sortLastPrev := func(i, j int) bool {
			icu := dm.CollectionsInfo.Data(uidsSlice[i]).(*data_models.CollectionInfo)
			jcu := dm.CollectionsInfo.Data(uidsSlice[j]).(*data_models.CollectionInfo)
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
			log.Errorf("Did not expect Next order for DataCollectionsSuggester")
		case core.Prev:
			log.Errorf("Did not expect Prev order for DataCollectionsSuggester")
		case core.Rand:
			suggesterNameParts = append(suggesterNameParts, ";Rand")
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(uidsSlice), func(i, j int) { uidsSlice[i], uidsSlice[j] = uidsSlice[j], uidsSlice[i] })
		}
		ret := []core.ContentItem(nil)
		if request.MoreItems <= 0 {
			return ret, nil
		}
		for _, uid := range uidsSlice {
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
