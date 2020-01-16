package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Bnei-Baruch/sqlboiler/queries/qm"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/mdb"
	mdbmodels "github.com/Bnei-Baruch/feed-api/mdb/models"
	"github.com/Bnei-Baruch/feed-api/utils"
)

var SECURE_PUBLISHED_MOD = qm.Where(fmt.Sprintf("secure=%d AND published IS TRUE", consts.SEC_PUBLIC))

// Responds with JSON of given response or aborts the request with the given error.
func concludeRequest(c *gin.Context, resp interface{}, err *HttpError) {
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

func ItemsHandler(c *gin.Context) {
	r := ItemsRequest{}
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleItems(c.MustGet("MDB_DB").(*sql.DB), r)
	concludeRequest(c, resp, err)
}

func handleItems(db *sql.DB, r ItemsRequest) (*ItemsResponse, *HttpError) {
	log.Infof("r: %+v", r)
	mods := []qm.QueryMod{SECURE_PUBLISHED_MOD}

	var total int64
	// Count total query.
	countMods := append([]qm.QueryMod{qm.Select("count(DISTINCT id)")}, mods...)
	err := mdbmodels.ContentUnits(db, countMods...).QueryRow().Scan(&total)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if total == 0 {
		return &ItemsResponse{Items: make([]interface{}, 0)}, nil
	}

	mods = append(mods, qm.GroupBy("id"))
	mods = append(mods, qm.OrderBy("created_at desc"))
	mods = append(mods, qm.Limit(10))
	mods = append(mods, qm.Offset(r.Offset))

	// Eager loading.
	loadTables := []string{
		"CollectionsContentUnits",
		"CollectionsContentUnits.Collection",
	}
	mods = append(mods, qm.Load(loadTables...))

	// Data query.
	units, err := mdbmodels.ContentUnits(db, mods...).All()
	if err != nil {
		return nil, NewInternalError(err)
	}

	// response
	cus, ex := prepareCUs(db, units, "he")
	if ex != nil {
		return nil, ex
	}

	resp := &ItemsResponse{
		Total: total,
		Items: utils.ToInterfaceSlice(cus),
	}

	return resp, nil
}

// units must be loaded with their CCUs loaded with their collections
func prepareCUs(db *sql.DB, units []*mdbmodels.ContentUnit, language string) ([]*ContentUnit, *HttpError) {
	// Filter secure published collections
	// Load i18n for all content units and all collections - total 2 DB round trips
	cuids := make([]int64, len(units))
	cids := make([]int64, 0)
	for i, x := range units {
		cuids[i] = x.ID
		b := x.R.CollectionsContentUnits[:0]
		for _, y := range x.R.CollectionsContentUnits {
			if consts.SEC_PUBLIC == y.R.Collection.Secure && y.R.Collection.Published {
				b = append(b, y)
				cids = append(cids, y.CollectionID)
			}
			x.R.CollectionsContentUnits = b
		}
	}

	cui18nsMap, err := loadCUI18ns(db, language, cuids)
	if err != nil {
		return nil, NewInternalError(err)
	}
	ci18nsMap, err := loadCI18ns(db, language, cids)
	if err != nil {
		return nil, NewInternalError(err)
	}

	cus := make([]*ContentUnit, len(units))
	for i, x := range units {
		cu, err := mdbToCU(x)
		if err != nil {
			return nil, NewInternalError(err)
		}
		if i18ns, ok := cui18nsMap[x.ID]; ok {
			setCUI18n(cu, language, i18ns)
		}

		// collections
		cu.Collections = make(map[string]*Collection, 0)
		for _, ccu := range x.R.CollectionsContentUnits {
			cl := ccu.R.Collection

			cc, err := mdbToC(cl)
			if err != nil {
				return nil, NewInternalError(err)
			}
			if i18ns, ok := ci18nsMap[cl.ID]; ok {
				setCI18n(cc, language, i18ns)
			}

			// Dirty hack for unique mapping - needs to parse in client...
			key := fmt.Sprintf("%s____%s", cl.UID, ccu.Name)
			cu.Collections[key] = cc
		}

		cus[i] = cu
	}

	return cus, nil
}

func loadCUI18ns(db *sql.DB, language string, ids []int64) (map[int64]map[string]*mdbmodels.ContentUnitI18n, error) {
	i18nsMap := make(map[int64]map[string]*mdbmodels.ContentUnitI18n, len(ids))
	if len(ids) == 0 {
		return i18nsMap, nil
	}

	// Load from DB
	i18ns, err := mdbmodels.ContentUnitI18ns(db,
		qm.WhereIn("content_unit_id in ?", utils.ConvertArgsInt64(ids)...),
		qm.AndIn("language in ?", utils.ConvertArgsString(consts.I18N_LANG_ORDER[language])...)).
		All()
	if err != nil {
		return nil, errors.Wrap(err, "Load content units i18ns from DB")
	}

	// Group by content unit and language
	for _, x := range i18ns {
		v, ok := i18nsMap[x.ContentUnitID]
		if !ok {
			v = make(map[string]*mdbmodels.ContentUnitI18n, 1)
			i18nsMap[x.ContentUnitID] = v
		}
		v[x.Language] = x
	}

	return i18nsMap, nil
}

func loadCI18ns(db *sql.DB, language string, ids []int64) (map[int64]map[string]*mdbmodels.CollectionI18n, error) {
	i18nsMap := make(map[int64]map[string]*mdbmodels.CollectionI18n, len(ids))
	if len(ids) == 0 {
		return i18nsMap, nil
	}

	// Load from DB
	i18ns, err := mdbmodels.CollectionI18ns(db,
		qm.WhereIn("collection_id in ?", utils.ConvertArgsInt64(ids)...),
		qm.AndIn("language in ?", utils.ConvertArgsString(consts.I18N_LANG_ORDER[language])...)).
		All()
	if err != nil {
		return nil, errors.Wrap(err, "Load collections i18ns from DB")
	}

	// Group by collection and language

	for _, x := range i18ns {
		v, ok := i18nsMap[x.CollectionID]
		if !ok {
			v = make(map[string]*mdbmodels.CollectionI18n, 1)
			i18nsMap[x.CollectionID] = v
		}
		v[x.Language] = x
	}

	return i18nsMap, nil
}

func mdbToCU(cu *mdbmodels.ContentUnit) (*ContentUnit, error) {
	var props mdb.ContentUnitProperties
	if err := cu.Properties.Unmarshal(&props); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal properties")
	}

	u := &ContentUnit{
		Item: Item{
			ID:          cu.UID,
			ContentType: mdb.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name,
			CreatedAt:   cu.CreatedAt,
		},
		mdbID:            cu.ID,
		Duration:         props.Duration,
		OriginalLanguage: props.OriginalLanguage,
	}

	if !props.FilmDate.IsZero() {
		u.FilmDate = &utils.Date{Time: props.FilmDate.Time}
	}

	return u, nil
}

func setCUI18n(cu *ContentUnit, language string, i18ns map[string]*mdbmodels.ContentUnitI18n) {
	for _, l := range consts.I18N_LANG_ORDER[language] {
		li18n, ok := i18ns[l]
		if ok {
			if cu.Name == "" && li18n.Name.Valid {
				cu.Name = li18n.Name.String
			}
			if cu.Description == "" && li18n.Description.Valid {
				cu.Description = li18n.Description.String
			}
		}
	}
}

func mdbToC(c *mdbmodels.Collection) (cl *Collection, err error) {
	var props mdb.CollectionProperties
	if err = c.Properties.Unmarshal(&props); err != nil {
		err = errors.Wrap(err, "json.Unmarshal properties")
		return
	}

	cl = &Collection{
		Item: Item{
			ID:          c.UID,
			ContentType: mdb.CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name,
		},
		Country:         props.Country,
		City:            props.City,
		FullAddress:     props.FullAddress,
		Genres:          props.Genres,
		DefaultLanguage: props.DefaultLanguage,
		HolidayID:       props.HolidayTag,
		SourceID:        props.Source,
		Number:          props.Number,
	}

	if !props.FilmDate.IsZero() {
		cl.FilmDate = &utils.Date{Time: props.FilmDate.Time}
	}
	if !props.StartDate.IsZero() {
		cl.StartDate = &utils.Date{Time: props.StartDate.Time}
	}
	if !props.EndDate.IsZero() {
		cl.EndDate = &utils.Date{Time: props.EndDate.Time}
	}

	return
}

func setCI18n(c *Collection, language string, i18ns map[string]*mdbmodels.CollectionI18n) {
	for _, l := range consts.I18N_LANG_ORDER[language] {
		li18n, ok := i18ns[l]
		if ok {
			if c.Name == "" && li18n.Name.Valid {
				c.Name = li18n.Name.String
			}
			if c.Description == "" && li18n.Description.Valid {
				c.Description = li18n.Description.String
			}
		}
	}
}
