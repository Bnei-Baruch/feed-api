package data_models

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/Bnei-Baruch/feed-api/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/queries"
)

type ReadIdFunc func() string

type SqlModels struct {
	localChroniclesDb *sql.DB
	modelsDb          *sql.DB
	prevReadIdFunc    ReadIdFunc

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

func MakeSqlModels(sqlFiles []string, cDb *sql.DB, modelsDb *sql.DB, prevReadIdFunc ReadIdFunc) *SqlModels {
	sqlsPath := viper.GetString("data_models.sqls_path")
	sqls, err := LoadSqls(sqlsPath, sqlFiles)
	utils.Must(err)

	return &SqlModels{cDb, modelsDb, prevReadIdFunc, sqlFiles, sqls}
}

func (cm *SqlModels) Name() string {
	return "SqlModels"
}

func (cm *SqlModels) Refresh() error {
	log.Info("Update sql models.")
	params := make(map[string]string)
	if cm.prevReadIdFunc != nil {
		params["$prev-read-id"] = cm.prevReadIdFunc()
	}
	for i, sql := range cm.sqls {
		//log.Infof("Before %s", sql)
		for param, value := range params {
			sql = strings.ReplaceAll(sql, param, value)
		}
		//log.Infof("After %s", sql)
		log.Infof("Running %s", cm.sqlFiles[i])
		if result, err := queries.Raw(sql).Exec(cm.modelsDb); err != nil {
			log.Warnf("Error running sqls: %+v", err)
			// return err
		} else {
			log.Infof("Updated sql %s, result: %+v", cm.sqlFiles[i], result)
		}
	}
	return nil
}

func (cm *SqlModels) Interval() time.Duration {
	// This should not be used - run automatically after chronicles window.
	return time.Second
}
