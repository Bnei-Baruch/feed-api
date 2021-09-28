package data_models

import (
	"context"
	"sort"
	"time"

	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/databases/data_models/models"
	"github.com/Bnei-Baruch/feed-api/utils"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type SqlDataModel struct {
	modelsDb *common.Connection
}

func MakeSqlDataModel(modelsDb *common.Connection) *SqlDataModel {
	return &SqlDataModel{modelsDb}
}

type Count struct {
	Uid   string `boil:"uid"`
	Count int64  `boil:"count"`
}

func (dm *SqlDataModel) WatchingNow(uids []string) ([]int64, error) {
	if count, err := dm.AllWatchingNow(); err != nil {
		return nil, err
	} else {
		ret := []int64(nil)
		for _, uid := range uids {
			ret = append(ret, count[uid])
		}
		return ret, nil
	}
}

func (dm *SqlDataModel) AllWatchingNow() (map[string]int64, error) {
	start := time.Now()
	defer func() {
		utils.Profile("AllWatchingNow", time.Now().Sub(start))
	}()
	count := []Count(nil)
	if err := dm.modelsDb.With(models.NewQuery(
		qm.Select("event_unit_uid as uid, unique_users_watching_now_count as count"),
		qm.From("dwh_content_units_measures"),
		qm.Where("unique_users_watching_now_count > 0"),
	)).Bind(context.TODO(), &count); err != nil {
		return nil, err
	}
	utils.Profile("AllWatchingNow.Sql", time.Now().Sub(start))
	cMap := make(map[string]int64, len(count))
	for _, c := range count {
		cMap[c.Uid] = c.Count
	}
	return cMap, nil
}

func (dm *SqlDataModel) SortWatchingNow(uids []string) error {
	if count, err := dm.AllWatchingNow(); err != nil {
		return err
	} else {
		sort.SliceStable(uids, func(i, j int) bool {
			return count[uids[i]] > count[uids[j]]
		})
	}
	return nil
}

func (dm *SqlDataModel) Views(uids []string) ([]int64, error) {
	if count, err := dm.AllViews(); err != nil {
		return nil, err
	} else {
		ret := []int64(nil)
		for _, uid := range uids {
			ret = append(ret, count[uid])
		}
		return ret, nil
	}
}

func (dm *SqlDataModel) AllViews() (map[string]int64, error) {
	count := []Count(nil)
	if err := dm.modelsDb.With(models.NewQuery(
		qm.Select("event_unit_uid as uid, unique_users_count as count"),
		qm.From("dwh_content_units_measures"),
		qm.Where("unique_users_count > 0"),
	)).Bind(context.TODO(), &count); err != nil {
		return nil, err
	}
	cMap := make(map[string]int64, len(count))
	for _, c := range count {
		cMap[c.Uid] = c.Count
	}
	return cMap, nil
}

func (dm *SqlDataModel) SortPopular(uids []string) error {
	if count, err := dm.AllViews(); err != nil {
		return err
	} else {
		/*for _, uid := range uids {
			if c, ok := count[uid]; ok {
				log.Infof("%s: %d", uid, c)
			}
		}*/
		sort.SliceStable(uids, func(i, j int) bool {
			return count[uids[i]] > count[uids[j]]
		})
		/*if len(uids) > 30 {
			log.Infof("After Uids: %+v...", uids[0:30])
		}*/
	}
	return nil
}
