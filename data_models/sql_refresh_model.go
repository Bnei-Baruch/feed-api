package data_models

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/Bnei-Baruch/feed-api/databases/data_models/models"
	"github.com/Bnei-Baruch/feed-api/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type SqlRefreshModel struct {
	modelsDb *sql.DB

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

func MakeSqlRefreshModel(sqlFiles []string, modelsDb *sql.DB) *SqlRefreshModel {
	sqlsPath := viper.GetString("data_models.sqls_path")
	sqls, err := LoadSqls(sqlsPath, sqlFiles)
	utils.Must(err)

	return &SqlRefreshModel{modelsDb, sqlFiles, sqls}
}

func (cm *SqlRefreshModel) Name() string {
	return "SqlRefreshModel"
}

func (cm *SqlRefreshModel) Refresh() error {
	log.Info("Update sql models.")
	params := make(map[string]string)

	minutesPrevEndReadId := []struct {
		IdMax null.String `boil:"id_max"`
	}(nil)
	if err := models.NewQuery(qm.Select("max(event_end_id_max) as id_max"), qm.From("dwh_fact_play_units_by_minutes")).Bind(context.TODO(), cm.modelsDb, &minutesPrevEndReadId); err != nil {
		return err
	}
	if len(minutesPrevEndReadId) == 1 && minutesPrevEndReadId[0].IdMax.Valid {
		params["$minutes-prev-read-id"] = minutesPrevEndReadId[0].IdMax.String
	} else {
		params["$minutes-read-id"] = ""
	}

	for i, sql := range cm.sqls {
		//log.Infof("Before %s", sql)
		for param, value := range params {
			sql = strings.ReplaceAll(sql, param, value)
		}
		//log.Infof("After %s", sql)
		start := time.Now()
		log.Infof("Running %s", cm.sqlFiles[i])
		if result, err := queries.Raw(sql).Exec(cm.modelsDb); err != nil {
			log.Warnf("Error running sqls: %+v", err)
			// return err
		} else {
			log.Infof("Updated sql %s, result: %+v", cm.sqlFiles[i], result)
		}
		utils.Profile(fmt.Sprintf("SqlDataModel: %s", cm.sqlFiles[i]), time.Now().Sub(start))
	}
	return nil
}

func (cm *SqlRefreshModel) Interval() time.Duration {
	// This should not be used - run automatically after chronicles window.
	return time.Second
}
