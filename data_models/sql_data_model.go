package data_models

import (
	"context"
	"sort"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/databases/data_models/models"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type SqlDataModel struct {
	modelsDb *common.Connection
}

func MakeSqlDataModel(modelsDb *common.Connection) *SqlDataModel {
	return &SqlDataModel{modelsDb}
}

type Count struct {
	Uid   null.String `boil:"uid"`
	Count null.Int64  `boil:"count"`
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
		if c.Uid.Valid && c.Count.Valid {
			cMap[c.Uid.String] = c.Count.Int64
		}
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
		qm.Select("event_unit_uid as uid, total_page_enter_count as count"),
		qm.From("dwh_content_units_measures"),
		qm.Where("total_page_enter_count > 0"),
	)).Bind(context.TODO(), &count); err != nil {
		return nil, err
	}
	cMap := make(map[string]int64, len(count))
	for _, c := range count {
		if c.Uid.Valid && c.Count.Valid {
			cMap[c.Uid.String] = c.Count.Int64
		}
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

func (dm *SqlDataModel) UniqueViews(uids []string) ([]int64, error) {
	if count, err := dm.AllUniqueViews(); err != nil {
		return nil, err
	} else {
		ret := []int64(nil)
		for _, uid := range uids {
			ret = append(ret, count[uid])
		}
		return ret, nil
	}
}

func (dm *SqlDataModel) AllUniqueViews() (map[string]int64, error) {
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
		if c.Uid.Valid && c.Count.Valid {
			cMap[c.Uid.String] = c.Count.Int64
		}
	}
	return cMap, nil
}

func (dm *SqlDataModel) SortUniquePopular(uids []string) error {
	if count, err := dm.AllUniqueViews(); err != nil {
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
