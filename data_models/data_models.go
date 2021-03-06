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
			cu.type_id = %d and (cu.properties->>'part')::int = 0 and (cu.properties->>'duration')::int < 1200,
			(cu.properties->>'duration')::int
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

	CONTENT_UNITS_WATCH_DURATION_SQL = `
		select 
			a.data->>'unit_uid' as event_unit_uid,
			case when
				cast(b.data->>'current_time' as float) > cast(a.data->>'current_time' as float) then
				cast(b.data->>'current_time' as float) - cast(a.data->>'current_time' as float)
			else 0 end as current_duration_seconds
		from entries a 
		left join entries b on (a.client_event_id=b.client_flow_id and b.client_event_type='player-stop')
		where 
			cast(a.created_at as date) >= '2021-04-01'
			and a.client_event_type='player-play' 
		order by a.created_at;`
)

type ContentUnitInfo struct {
	TypeId             int64
	Uid                string
	Date               time.Time
	CreatedAt          time.Time
	SecureAndPublished bool
	IsLessonPrep       bool
	Duration           time.Duration
}

func ScanContentUnitInfo(rows *sql.Rows, datas map[string]interface{}) error {
	cu := ContentUnitInfo{}
	var isLessonPrep null.Bool
	var duration null.Int64
	if err := rows.Scan(&cu.TypeId, &cu.Uid, &cu.Date, &cu.CreatedAt, &cu.SecureAndPublished, &isLessonPrep, &duration); err != nil {
		return err
	} else {
		cu.IsLessonPrep = isLessonPrep.Valid && isLessonPrep.Bool
		cu.Duration = time.Duration(duration.Int64) * time.Second
		datas[cu.Uid] = &cu
		return nil
	}
}

type CollectionInfo struct {
	TypeId    int64
	Uid       string
	Date      time.Time
	CreatedAt time.Time
	SourceUid string
}

func ScanCollectionInfo(rows *sql.Rows, datas map[string]interface{}) error {
	c := CollectionInfo{}
	var sourceUid null.String
	if err := rows.Scan(&c.TypeId, &c.Uid, &c.Date, &c.CreatedAt, &sourceUid); err != nil {
		return err
	} else {
		c.SourceUid = sourceUid.String
		datas[c.Uid] = &c
		return nil
	}
}

type RefreshModel interface {
	Name() string
	Refresh() error
	Interval() time.Duration
}

type DataModels struct {
	ticker      *time.Ticker
	nextRefresh map[string]time.Time

	LanguagesContentUnitsFilter   *MDBFilterModel
	TagsContentUnitsFilter        *MDBFilterModel
	SourcesContentUnitsFilter     *MDBFilterModel
	PersonsContentUnitsFilter     *MDBFilterModel
	CollectionsContentUnitsFilter *MDBFilterModel
	ContentUnitsCollectionsFilter *MDBFilterModel

	ContentUnitsInfo *MDBDataModel
	//ContentUnitsPopularity *MDBDataModel

	CollectionsInfo *MDBDataModel

	models []RefreshModel
}

func MakeDataModels(db *sql.DB) *DataModels {
	lcuf := MakeMDBFilterModel(db, "LanguagesContentUnitsFilter", time.Duration(time.Minute*10), LANGUAGES_CONTENT_UNITS_SQL)
	tcuf := MakeMDBFilterModel(db, "TagsContentUnitsFilter", time.Duration(time.Minute*10), TAGS_CONTENT_UNITS_SQL)
	scuf := MakeMDBFilterModel(db, "SourcesContentUnitsFilter", time.Duration(time.Minute*10), SOURCES_CONTENT_UNITS_SQL)
	pcuf := MakeMDBFilterModel(db, "PersonsContentUnitsFilter", time.Duration(time.Minute*10), PERSONS_CONTENT_UNITS_SQL)
	ccuf := MakeMDBFilterModel(db, "CollectionsContentUnitsFilter", time.Duration(time.Minute*10), COLLECTIONS_CONTENT_UNITS_SQL)
	cucf := MakeMDBFilterModel(db, "ContentUnitsCollectionsFilter", time.Duration(time.Minute*10), CONTENT_UNITS_COLLECTIONS_SQL)
	cui := MakeMDBDataModel(db, "ContentUnitsInfo", time.Duration(time.Minute*10), fmt.Sprintf(CONTENT_UNITS_INFO_SQL, mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_LESSON_PART].ID), ScanContentUnitInfo)
	ci := MakeMDBDataModel(db, "CollectionsInfo", time.Duration(time.Minute*10), COLLECTIONS_INFO_SQL, ScanCollectionInfo)
	models := []RefreshModel{lcuf, tcuf, scuf, pcuf, ccuf, cucf, cui, ci}

	dataModels := &DataModels{
		ticker:      time.NewTicker(time.Second),
		nextRefresh: make(map[string]time.Time),

		LanguagesContentUnitsFilter:   lcuf,
		TagsContentUnitsFilter:        tcuf,
		SourcesContentUnitsFilter:     scuf,
		PersonsContentUnitsFilter:     pcuf,
		CollectionsContentUnitsFilter: ccuf,
		ContentUnitsCollectionsFilter: cucf,
		ContentUnitsInfo:              cui,
		CollectionsInfo:               ci,

		models: models,
	}

	go func() {
		refresh := func() {
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
		log.Infof("Refreshed %s in %s", model.Name(), end.Sub(start))
		if err != nil {
			return errors.Wrap(err, model.Name())
		}
		return err
	}

	for _, dataModel := range dataModels.models {
		if when, ok := dataModels.nextRefresh[dataModel.Name()]; !ok || when.Before(time.Now()) {
			if err := refreshModel(dataModel); err != nil {
				return err
			}
			dataModels.nextRefresh[dataModel.Name()] = time.Now().Add(dataModel.Interval())
		}
	}
	return nil
}

// --- Data Models --- //

// --- MDB Filter Model --- //
type MDBFilterModel struct {
	db          *sql.DB
	name        string
	interval    time.Duration
	sql         string
	Values      map[string][]string
	valuesMutex sync.RWMutex
}

func MakeMDBFilterModel(db *sql.DB, name string, interval time.Duration, sql string) *MDBFilterModel {
	return &MDBFilterModel{db, name, interval, sql, make(map[string][]string), sync.RWMutex{}}
}

func (model *MDBFilterModel) Name() string {
	return model.name
}

func (model *MDBFilterModel) Interval() time.Duration {
	return model.interval
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
type ScanRows func(rows *sql.Rows, datas map[string]interface{}) error

type MDBDataModel struct {
	db         *sql.DB
	name       string
	interval   time.Duration
	sql        string
	scanRows   ScanRows
	Datas      map[string]interface{}
	datasMutex sync.RWMutex
}

func MakeMDBDataModel(db *sql.DB, name string, interval time.Duration, sql string, scanRows ScanRows) *MDBDataModel {
	return &MDBDataModel{db, name, interval, sql, scanRows, make(map[string]interface{}), sync.RWMutex{}}
}

func (model *MDBDataModel) Name() string {
	return model.name
}

func (model *MDBDataModel) Interval() time.Duration {
	return model.interval
}

func (model *MDBDataModel) Refresh() error {
	rows, err := queries.Raw(model.db, model.sql).Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	tmp := make(map[string]interface{})
	for rows.Next() {
		if err := model.scanRows(rows, tmp); err != nil {
			return err
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
	return model.Datas[key]
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
