package data_models

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/databases/data_models/models"
	"github.com/Bnei-Baruch/feed-api/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type SqlRefreshModel struct {
	modelsDb *common.Connection

	sqlFiles []string
	sqls     []string
}

func LoadSqls(path string, files []string) ([]string, error) {
	sqls := []string(nil)
	for _, file := range files {
		if b, err := ioutil.ReadFile(filepath.Join(path, file)); err != nil {
			return nil, err
		} else {
			sqls = append(sqls, string(b))
		}
	}
	return sqls, nil
}

func MakeSqlRefreshModel(sqlFiles []string, modelsDb *common.Connection) *SqlRefreshModel {
	sqlsPath := viper.GetString("data_models.sqls_path")
	sqls, err := LoadSqls(sqlsPath, sqlFiles)
	utils.Must(err)

	return &SqlRefreshModel{modelsDb, sqlFiles, sqls}
}

func (cm *SqlRefreshModel) Name() string {
	return fmt.Sprintf("SqlRefreshModel-%s", strings.Join(cm.sqlFiles, "-"))
}

func (cm *SqlRefreshModel) Refresh() error {
	log.Debug("Update sql models.")
	defer log.Debugf("Update sql models done.")
	params := make(map[string]string)
	cm.modelsDb.FillParams(params)

	// Play units by minutes.
	minutesPrevEndReadId := []struct {
		IdMax null.String `boil:"id_max"`
	}(nil)
	if err := cm.modelsDb.With(models.NewQuery(qm.Select("max(event_end_id_max) as id_max"), qm.From("dwh_fact_play_units_by_minutes"))).Bind(context.TODO(), &minutesPrevEndReadId); err != nil {
		return err
	}
	if len(minutesPrevEndReadId) == 1 && minutesPrevEndReadId[0].IdMax.Valid {
		params["$minutes-prev-read-id"] = minutesPrevEndReadId[0].IdMax.String
	} else {
		params["$minutes-prev-read-id"] = ""
	}

	// Downloads by minutes.
	downloadsMinutesPrevReadId := []struct {
		IdMax null.String `boil:"id_max"`
	}(nil)
	if err := cm.modelsDb.With(models.NewQuery(qm.Select("max(event_id_max) as id_max"), qm.From("dwh_fact_download_units_by_minutes"))).Bind(context.TODO(), &downloadsMinutesPrevReadId); err != nil {
		return err
	}
	if len(downloadsMinutesPrevReadId) == 1 && downloadsMinutesPrevReadId[0].IdMax.Valid {
		params["$download-minutes-prev-read-id"] = downloadsMinutesPrevReadId[0].IdMax.String
	} else {
		params["$download-minutes-prev-read-id"] = ""
	}

	// Page enter by minutes.
	pageEnterMinutesPrevReadId := []struct {
		IdMax null.String `boil:"id_max"`
	}(nil)
	if err := cm.modelsDb.With(models.NewQuery(qm.Select("max(event_id_max) as id_max"), qm.From("dwh_fact_page_enter_units_by_minutes"))).Bind(context.TODO(), &pageEnterMinutesPrevReadId); err != nil {
		return err
	}
	if len(pageEnterMinutesPrevReadId) == 1 && pageEnterMinutesPrevReadId[0].IdMax.Valid {
		params["$page-enter-minutes-prev-read-id"] = pageEnterMinutesPrevReadId[0].IdMax.String
	} else {
		params["$page-enter-minutes-prev-read-id"] = ""
	}

	// 2018 - 2021 odl content user measures.
	contentUnitMeasuresOldCount := []struct {
		Count null.Int64 `boil:"count"`
	}(nil)
	if err := cm.modelsDb.With(models.NewQuery(qm.Select("count(*) as count"), qm.From("dwh_content_units_measures_2018_2021"))).Bind(context.TODO(), &contentUnitMeasuresOldCount); err != nil {
		return err
	}

	for i, sql := range cm.sqls {
		// Skip adding old measures more than once.
		if cm.sqlFiles[i] == INSERT_CONTENT_UNITS_2018_2021_MEASURES &&
			len(contentUnitMeasuresOldCount) == 1 &&
			contentUnitMeasuresOldCount[0].Count.Valid &&
			contentUnitMeasuresOldCount[0].Count.Int64 > 0 {
			log.Debugf("Skipping %s", cm.sqlFiles[i])
			continue
		}
		for param, value := range params {
			sql = strings.ReplaceAll(sql, param, value)
		}
		start := time.Now()
		log.Debugf("Running %s with %+v", cm.sqlFiles[i], params)
		if result, err := cm.modelsDb.With(queries.Raw(sql)).Exec(); err != nil {
			log.Warnf("Error running sql %s: %+v", cm.sqlFiles[i], err)
		} else {
			log.Debugf("Updated sql %s, result: %+v", cm.sqlFiles[i], result)
		}
		utils.Profile(fmt.Sprintf("SqlDataModel: %s", cm.sqlFiles[i]), time.Now().Sub(start))
	}
	return nil
}

func (cm *SqlRefreshModel) Interval() time.Duration {
	return time.Duration(time.Minute)
}
