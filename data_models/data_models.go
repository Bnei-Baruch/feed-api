package data_models

import (
	"database/sql"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"

	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
	"github.com/Bnei-Baruch/sqlboiler/queries"
)

const (
	DATA_MODELS_REFRESH_SECONDS = 600

	LANGUAGES_CONTENT_UNITS_SQL = `
		select
			f.language,
			array_agg(distinct cu.uid)
		from
			content_units as cu,
			files as f
		where
			f.content_unit_id = cu.id and
			f.mime_type in ('video/mp4', 'audio/mpeg')
		group by
			f.language;`

	TAGS_CONTENT_UNITS_SQL = `
		select
			t.uid,
			array_agg(distinct cu.uid)
		from
			content_units as cu,
			tags as t,
			content_units_tags as cut
		where
			cut.content_unit_id = cu.id and
			t.id = cut.tag_id
		group by
			t.uid;`

	SOURCES_CONTENT_UNITS_SQL = `
		select
			s.uid,
			array_agg(distinct cu.uid)
		from
			content_units as cu,
			sources as s,
			content_units_sources as cus
		where
			cus.content_unit_id = cu.id and
			s.id = cus.source_id
		group by
			s.uid;`

	PERSONS_CONTENT_UNITS_SQL = `
		select
			p.uid,
			array_agg(distinct cu.uid)
		from
			content_units as cu,
			content_units_persons as cup,
			persons as p
		where
			p.id = cup.person_id and
			cup.content_unit_id = cu.id
		group by
			p.uid;`

	RAV_PERSON_UID = "abcdefgh"

	COLLECTIONS_CONTENT_UNITS_SQL = `
		select
			c.uid,
			array_agg(distinct cu.uid)
		from
			content_units as cu,
			collections_content_units as ccu,
			collections as c
		where
			c.id = ccu.collection_id and
			ccu.content_unit_id = cu.id
		group by
			c.uid;`

	CONTENT_UNITS_COLLECTIONS_SQL = `
		select
			cu.uid,
			array_agg(distinct c.uid)
		from
			content_units as cu,
			collections_content_units as ccu,
			collections as c
		where
			c.id = ccu.collection_id and
			ccu.content_unit_id = cu.id
		group by
			cu.uid;`

	CONTENT_UNITS_INFO_SQL = `
		select
			cu.type_id,
			cu.uid as uid,
			coalesce(cu.properties->>'film_date', cu.properties->>'start_date', cu.created_at::text)::date as date,
			cu.created_at as created_at,
			cu.secure = 0 AND cu.published IS TRUE as secure_and_published,
			cu.type_id = %d and (cu.properties->>'part')::int = 0 and (cu.properties->>'duration')::int < 1200
		from
			content_units as cu;`

	COLLECTIONS_INFO_SQL = `
		select
			c.type_id,
			c.uid as uid,
			coalesce(c.properties->>'start_date', c.created_at::text)::date as date,
			c.created_at as created_at,
			c.properties->>'source'
		from
			collections as c;`
)

type ContentUnitInfo struct {
	TypeId             int64
	Uid                string
	Date               time.Time
	CreatedAt          time.Time
	SecureAndPublished bool
	IsLessonPrep       bool
}

func ScanContentUnitInfo(rows *sql.Rows) (string, interface{}, error) {
	cu := ContentUnitInfo{}
	var isLessonPrep null.Bool
	if err := rows.Scan(&cu.TypeId, &cu.Uid, &cu.Date, &cu.CreatedAt, &cu.SecureAndPublished, &isLessonPrep); err != nil {
		return "", nil, err
	} else {
		cu.IsLessonPrep = isLessonPrep.Valid && isLessonPrep.Bool
		return cu.Uid, &cu, err
	}
}

type CollectionInfo struct {
	TypeId    int64
	Uid       string
	Date      time.Time
	CreatedAt time.Time
	SourceUid string
}

func ScanCollectionInfo(rows *sql.Rows) (string, interface{}, error) {
	c := CollectionInfo{}
	var sourceUid null.String
	if err := rows.Scan(&c.TypeId, &c.Uid, &c.Date, &c.CreatedAt, &sourceUid); err != nil {
		return "", nil, err
	} else {
		c.SourceUid = sourceUid.String
		return c.Uid, &c, err
	}
}

type RefreshModel interface {
	Name() string
	Refresh() error
}

type DataModels struct {
	ticker                        *time.Ticker
	LanguagesContentUnitsFilter   *MDBFilterModel
	TagsContentUnitsFilter        *MDBFilterModel
	SourcesContentUnitsFilter     *MDBFilterModel
	PersonsContentUnitsFilter     *MDBFilterModel
	CollectionsContentUnitsFilter *MDBFilterModel
	ContentUnitsCollectionsFilter *MDBFilterModel

	ContentUnitsInfo *MDBDataModel
	//ContentUnitsPopularity *MDBDataModel

	CollectionsInfo *MDBDataModel
}

func MakeDataModels(db *sql.DB) *DataModels {
	dataModels := &DataModels{
		time.NewTicker(DATA_MODELS_REFRESH_SECONDS * time.Second),
		MakeMDBFilterModel(db, "LanguagesContentUnitsFilter", LANGUAGES_CONTENT_UNITS_SQL),
		MakeMDBFilterModel(db, "TagsContentUnitsFilter", TAGS_CONTENT_UNITS_SQL),
		MakeMDBFilterModel(db, "SourcesContentUnitsFilter", SOURCES_CONTENT_UNITS_SQL),
		MakeMDBFilterModel(db, "PersonsContentUnitsFilter", PERSONS_CONTENT_UNITS_SQL),
		MakeMDBFilterModel(db, "CollectionsContentUnitsFilter", COLLECTIONS_CONTENT_UNITS_SQL),
		MakeMDBFilterModel(db, "ContentUnitsCollectionsFilter", CONTENT_UNITS_COLLECTIONS_SQL),
		MakeMDBDataModel(db, "ContentUnitsInfo", fmt.Sprintf(CONTENT_UNITS_INFO_SQL, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID), ScanContentUnitInfo),
		//MakeMDBDataModel(db, "ContentUnitsPopularity", CONTENT_UNITS_POPULARITY_SQL, ScanContentUnitPopularity),
		MakeMDBDataModel(db, "CollectionsInfo", COLLECTIONS_INFO_SQL, ScanCollectionInfo),
	}

	go func() {
		refresh := func() {
			log.Info("Refreshing data models.")
			if err := dataModels.Refresh(); err != nil {
				log.Errorf("Error Refresh: %s", err.Error())
			}
		}
		refresh()
		for range dataModels.ticker.C {
			refresh()
		}
	}()

	return dataModels
}

func (dataModels *DataModels) Refresh() error {
	refreshModel := func(model RefreshModel) error {
		start := time.Now()
		err := model.Refresh()
		end := time.Now()
		log.Infof("Refreshing %s in %s", model.Name(), end.Sub(start))
		if err != nil {
			return errors.Wrap(err, model.Name())
		}
		return err
	}
	if err := refreshModel(dataModels.LanguagesContentUnitsFilter); err != nil {
		return err
	}
	if err := refreshModel(dataModels.TagsContentUnitsFilter); err != nil {
		return err
	}
	if err := refreshModel(dataModels.SourcesContentUnitsFilter); err != nil {
		return err
	}
	if err := refreshModel(dataModels.PersonsContentUnitsFilter); err != nil {
		return err
	}
	if err := refreshModel(dataModels.CollectionsContentUnitsFilter); err != nil {
		return err
	}
	if err := refreshModel(dataModels.ContentUnitsCollectionsFilter); err != nil {
		return err
	}
	if err := refreshModel(dataModels.ContentUnitsInfo); err != nil {
		return err
	}
	if err := refreshModel(dataModels.CollectionsInfo); err != nil {
		return err
	}
	return nil
}

// --- Data Models --- //

// --- MDB Filter Model --- //
type MDBFilterModel struct {
	db          *sql.DB
	name        string
	sql         string
	Values      map[string][]string
	valuesMutex sync.RWMutex
}

func MakeMDBFilterModel(db *sql.DB, name string, sql string) *MDBFilterModel {
	return &MDBFilterModel{db, name, sql, make(map[string][]string), sync.RWMutex{}}
}

func (model *MDBFilterModel) Name() string {
	return model.name
}

func (model *MDBFilterModel) Refresh() error {
	rows, err := queries.Raw(model.db, model.sql).Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	tmp := make(map[string][]string)
	for rows.Next() {
		var key null.String
		uids := []string(nil)
		if err := rows.Scan(&key, pq.Array(&uids)); err != nil {
			return err
		}
		sort.Strings(uids)
		// We ignore key.Valid as we want to use in non valid state the default empty string.
		tmp[key.String] = uids
	}

	model.valuesMutex.Lock()
	defer model.valuesMutex.Unlock()
	for k, v := range tmp {
		model.Values[k] = v
	}
	return nil
}

func (model *MDBFilterModel) FilterValues(keys []string) []string {
	model.valuesMutex.RLock()
	defer model.valuesMutex.RUnlock()
	ret := []string(nil)
	for _, key := range keys {
		if ret == nil {
			ret = model.Values[key]
		} else {
			ret = utils.UnionSorted(ret, model.Values[key])
		}
	}
	return ret
}

// -- MDB Data Model -- //
type ScanRows func(rows *sql.Rows) (string, interface{}, error)

type MDBDataModel struct {
	db         *sql.DB
	name       string
	sql        string
	scanRows   ScanRows
	Datas      map[string]interface{}
	datasMutex sync.RWMutex
}

func MakeMDBDataModel(db *sql.DB, name string, sql string, scanRows ScanRows) *MDBDataModel {
	return &MDBDataModel{db, name, sql, scanRows, make(map[string]interface{}), sync.RWMutex{}}
}

func (model *MDBDataModel) Name() string {
	return model.name
}

func (model *MDBDataModel) Refresh() error {
	rows, err := queries.Raw(model.db, model.sql).Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	tmp := make(map[string]interface{})
	for rows.Next() {
		if key, data, err := model.scanRows(rows); err != nil {
			return err
		} else {
			tmp[key] = data
		}
	}

	model.datasMutex.Lock()
	defer model.datasMutex.Unlock()
	for k, d := range tmp {
		model.Datas[k] = d
	}
	return nil
}

func (model *MDBDataModel) Data(key string) interface{} {
	model.datasMutex.RLock()
	defer model.datasMutex.RUnlock()
	if data, ok := model.Datas[key]; ok {
		return data
	}
	return nil
}

func (model *MDBDataModel) Keys() []string {
	model.datasMutex.RLock()
	defer model.datasMutex.RUnlock()
	keys := []string(nil)
	for k := range model.Datas {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
