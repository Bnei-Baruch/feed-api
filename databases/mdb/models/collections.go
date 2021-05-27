// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Collection is an object representing the database table.
type Collection struct {
	ID         int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
	UID        string    `boil:"uid" json:"uid" toml:"uid" yaml:"uid"`
	TypeID     int64     `boil:"type_id" json:"type_id" toml:"type_id" yaml:"type_id"`
	CreatedAt  time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	Properties null.JSON `boil:"properties" json:"properties,omitempty" toml:"properties" yaml:"properties,omitempty"`
	Secure     int16     `boil:"secure" json:"secure" toml:"secure" yaml:"secure"`
	Published  bool      `boil:"published" json:"published" toml:"published" yaml:"published"`

	R *collectionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L collectionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CollectionColumns = struct {
	ID         string
	UID        string
	TypeID     string
	CreatedAt  string
	Properties string
	Secure     string
	Published  string
}{
	ID:         "id",
	UID:        "uid",
	TypeID:     "type_id",
	CreatedAt:  "created_at",
	Properties: "properties",
	Secure:     "secure",
	Published:  "published",
}

// Generated where

type whereHelpernull_JSON struct{ field string }

func (w whereHelpernull_JSON) EQ(x null.JSON) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_JSON) NEQ(x null.JSON) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_JSON) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_JSON) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_JSON) LT(x null.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_JSON) LTE(x null.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_JSON) GT(x null.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_JSON) GTE(x null.JSON) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelperint16 struct{ field string }

func (w whereHelperint16) EQ(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint16) NEQ(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint16) LT(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint16) LTE(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint16) GT(x int16) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint16) GTE(x int16) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint16) IN(slice []int16) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}

var CollectionWhere = struct {
	ID         whereHelperint64
	UID        whereHelperstring
	TypeID     whereHelperint64
	CreatedAt  whereHelpertime_Time
	Properties whereHelpernull_JSON
	Secure     whereHelperint16
	Published  whereHelperbool
}{
	ID:         whereHelperint64{field: "\"collections\".\"id\""},
	UID:        whereHelperstring{field: "\"collections\".\"uid\""},
	TypeID:     whereHelperint64{field: "\"collections\".\"type_id\""},
	CreatedAt:  whereHelpertime_Time{field: "\"collections\".\"created_at\""},
	Properties: whereHelpernull_JSON{field: "\"collections\".\"properties\""},
	Secure:     whereHelperint16{field: "\"collections\".\"secure\""},
	Published:  whereHelperbool{field: "\"collections\".\"published\""},
}

// CollectionRels is where relationship names are stored.
var CollectionRels = struct {
	Type                    string
	CollectionI18ns         string
	CollectionsContentUnits string
}{
	Type:                    "Type",
	CollectionI18ns:         "CollectionI18ns",
	CollectionsContentUnits: "CollectionsContentUnits",
}

// collectionR is where relationships are stored.
type collectionR struct {
	Type                    *ContentType
	CollectionI18ns         CollectionI18nSlice
	CollectionsContentUnits CollectionsContentUnitSlice
}

// NewStruct creates a new relationship struct
func (*collectionR) NewStruct() *collectionR {
	return &collectionR{}
}

// collectionL is where Load methods for each relationship are stored.
type collectionL struct{}

var (
	collectionAllColumns            = []string{"id", "uid", "type_id", "created_at", "properties", "secure", "published"}
	collectionColumnsWithoutDefault = []string{"uid", "type_id", "properties"}
	collectionColumnsWithDefault    = []string{"id", "created_at", "secure", "published"}
	collectionPrimaryKeyColumns     = []string{"id"}
)

type (
	// CollectionSlice is an alias for a slice of pointers to Collection.
	// This should generally be used opposed to []Collection.
	CollectionSlice []*Collection

	collectionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	collectionType                 = reflect.TypeOf(&Collection{})
	collectionMapping              = queries.MakeStructMapping(collectionType)
	collectionPrimaryKeyMapping, _ = queries.BindMapping(collectionType, collectionMapping, collectionPrimaryKeyColumns)
	collectionInsertCacheMut       sync.RWMutex
	collectionInsertCache          = make(map[string]insertCache)
	collectionUpdateCacheMut       sync.RWMutex
	collectionUpdateCache          = make(map[string]updateCache)
	collectionUpsertCacheMut       sync.RWMutex
	collectionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single collection record from the query.
func (q collectionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Collection, error) {
	o := &Collection{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for collections")
	}

	return o, nil
}

// All returns all Collection records from the query.
func (q collectionQuery) All(ctx context.Context, exec boil.ContextExecutor) (CollectionSlice, error) {
	var o []*Collection

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Collection slice")
	}

	return o, nil
}

// Count returns the count of all Collection records in the query.
func (q collectionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count collections rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q collectionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if collections exists")
	}

	return count > 0, nil
}

// Type pointed to by the foreign key.
func (o *Collection) Type(mods ...qm.QueryMod) contentTypeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.TypeID),
	}

	queryMods = append(queryMods, mods...)

	query := ContentTypes(queryMods...)
	queries.SetFrom(query.Query, "\"content_types\"")

	return query
}

// CollectionI18ns retrieves all the collection_i18n's CollectionI18ns with an executor.
func (o *Collection) CollectionI18ns(mods ...qm.QueryMod) collectionI18nQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"collection_i18n\".\"collection_id\"=?", o.ID),
	)

	query := CollectionI18ns(queryMods...)
	queries.SetFrom(query.Query, "\"collection_i18n\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"collection_i18n\".*"})
	}

	return query
}

// CollectionsContentUnits retrieves all the collections_content_unit's CollectionsContentUnits with an executor.
func (o *Collection) CollectionsContentUnits(mods ...qm.QueryMod) collectionsContentUnitQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"collections_content_units\".\"collection_id\"=?", o.ID),
	)

	query := CollectionsContentUnits(queryMods...)
	queries.SetFrom(query.Query, "\"collections_content_units\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"collections_content_units\".*"})
	}

	return query
}

// LoadType allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (collectionL) LoadType(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCollection interface{}, mods queries.Applicator) error {
	var slice []*Collection
	var object *Collection

	if singular {
		object = maybeCollection.(*Collection)
	} else {
		slice = *maybeCollection.(*[]*Collection)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &collectionR{}
		}
		args = append(args, object.TypeID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &collectionR{}
			}

			for _, a := range args {
				if a == obj.TypeID {
					continue Outer
				}
			}

			args = append(args, obj.TypeID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`content_types`), qm.WhereIn(`content_types.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ContentType")
	}

	var resultSlice []*ContentType
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ContentType")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for content_types")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for content_types")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Type = foreign
		if foreign.R == nil {
			foreign.R = &contentTypeR{}
		}
		foreign.R.TypeCollections = append(foreign.R.TypeCollections, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.TypeID == foreign.ID {
				local.R.Type = foreign
				if foreign.R == nil {
					foreign.R = &contentTypeR{}
				}
				foreign.R.TypeCollections = append(foreign.R.TypeCollections, local)
				break
			}
		}
	}

	return nil
}

// LoadCollectionI18ns allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (collectionL) LoadCollectionI18ns(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCollection interface{}, mods queries.Applicator) error {
	var slice []*Collection
	var object *Collection

	if singular {
		object = maybeCollection.(*Collection)
	} else {
		slice = *maybeCollection.(*[]*Collection)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &collectionR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &collectionR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`collection_i18n`), qm.WhereIn(`collection_i18n.collection_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load collection_i18n")
	}

	var resultSlice []*CollectionI18n
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice collection_i18n")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on collection_i18n")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for collection_i18n")
	}

	if singular {
		object.R.CollectionI18ns = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &collectionI18nR{}
			}
			foreign.R.Collection = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.CollectionID {
				local.R.CollectionI18ns = append(local.R.CollectionI18ns, foreign)
				if foreign.R == nil {
					foreign.R = &collectionI18nR{}
				}
				foreign.R.Collection = local
				break
			}
		}
	}

	return nil
}

// LoadCollectionsContentUnits allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (collectionL) LoadCollectionsContentUnits(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCollection interface{}, mods queries.Applicator) error {
	var slice []*Collection
	var object *Collection

	if singular {
		object = maybeCollection.(*Collection)
	} else {
		slice = *maybeCollection.(*[]*Collection)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &collectionR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &collectionR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`collections_content_units`), qm.WhereIn(`collections_content_units.collection_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load collections_content_units")
	}

	var resultSlice []*CollectionsContentUnit
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice collections_content_units")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on collections_content_units")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for collections_content_units")
	}

	if singular {
		object.R.CollectionsContentUnits = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &collectionsContentUnitR{}
			}
			foreign.R.Collection = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.CollectionID {
				local.R.CollectionsContentUnits = append(local.R.CollectionsContentUnits, foreign)
				if foreign.R == nil {
					foreign.R = &collectionsContentUnitR{}
				}
				foreign.R.Collection = local
				break
			}
		}
	}

	return nil
}

// SetType of the collection to the related item.
// Sets o.R.Type to related.
// Adds o to related.R.TypeCollections.
func (o *Collection) SetType(ctx context.Context, exec boil.ContextExecutor, insert bool, related *ContentType) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"collections\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"type_id"}),
		strmangle.WhereClause("\"", "\"", 2, collectionPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.TypeID = related.ID
	if o.R == nil {
		o.R = &collectionR{
			Type: related,
		}
	} else {
		o.R.Type = related
	}

	if related.R == nil {
		related.R = &contentTypeR{
			TypeCollections: CollectionSlice{o},
		}
	} else {
		related.R.TypeCollections = append(related.R.TypeCollections, o)
	}

	return nil
}

// AddCollectionI18ns adds the given related objects to the existing relationships
// of the collection, optionally inserting them as new records.
// Appends related to o.R.CollectionI18ns.
// Sets related.R.Collection appropriately.
func (o *Collection) AddCollectionI18ns(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*CollectionI18n) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.CollectionID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"collection_i18n\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"collection_id"}),
				strmangle.WhereClause("\"", "\"", 2, collectionI18nPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.CollectionID, rel.Language}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.CollectionID = o.ID
		}
	}

	if o.R == nil {
		o.R = &collectionR{
			CollectionI18ns: related,
		}
	} else {
		o.R.CollectionI18ns = append(o.R.CollectionI18ns, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &collectionI18nR{
				Collection: o,
			}
		} else {
			rel.R.Collection = o
		}
	}
	return nil
}

// AddCollectionsContentUnits adds the given related objects to the existing relationships
// of the collection, optionally inserting them as new records.
// Appends related to o.R.CollectionsContentUnits.
// Sets related.R.Collection appropriately.
func (o *Collection) AddCollectionsContentUnits(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*CollectionsContentUnit) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.CollectionID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"collections_content_units\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"collection_id"}),
				strmangle.WhereClause("\"", "\"", 2, collectionsContentUnitPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.CollectionID, rel.ContentUnitID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.CollectionID = o.ID
		}
	}

	if o.R == nil {
		o.R = &collectionR{
			CollectionsContentUnits: related,
		}
	} else {
		o.R.CollectionsContentUnits = append(o.R.CollectionsContentUnits, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &collectionsContentUnitR{
				Collection: o,
			}
		} else {
			rel.R.Collection = o
		}
	}
	return nil
}

// Collections retrieves all the records using an executor.
func Collections(mods ...qm.QueryMod) collectionQuery {
	mods = append(mods, qm.From("\"collections\""))
	return collectionQuery{NewQuery(mods...)}
}

// FindCollection retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCollection(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Collection, error) {
	collectionObj := &Collection{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"collections\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, collectionObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from collections")
	}

	return collectionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Collection) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no collections provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(collectionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	collectionInsertCacheMut.RLock()
	cache, cached := collectionInsertCache[key]
	collectionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			collectionAllColumns,
			collectionColumnsWithDefault,
			collectionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(collectionType, collectionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(collectionType, collectionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"collections\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"collections\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into collections")
	}

	if !cached {
		collectionInsertCacheMut.Lock()
		collectionInsertCache[key] = cache
		collectionInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Collection.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Collection) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	collectionUpdateCacheMut.RLock()
	cache, cached := collectionUpdateCache[key]
	collectionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			collectionAllColumns,
			collectionPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update collections, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"collections\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, collectionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(collectionType, collectionMapping, append(wl, collectionPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update collections row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for collections")
	}

	if !cached {
		collectionUpdateCacheMut.Lock()
		collectionUpdateCache[key] = cache
		collectionUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q collectionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for collections")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for collections")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CollectionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), collectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"collections\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, collectionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in collection slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all collection")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Collection) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no collections provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(collectionColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	collectionUpsertCacheMut.RLock()
	cache, cached := collectionUpsertCache[key]
	collectionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			collectionAllColumns,
			collectionColumnsWithDefault,
			collectionColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			collectionAllColumns,
			collectionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert collections, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(collectionPrimaryKeyColumns))
			copy(conflict, collectionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"collections\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(collectionType, collectionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(collectionType, collectionMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert collections")
	}

	if !cached {
		collectionUpsertCacheMut.Lock()
		collectionUpsertCache[key] = cache
		collectionUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Collection record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Collection) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Collection provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), collectionPrimaryKeyMapping)
	sql := "DELETE FROM \"collections\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from collections")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for collections")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q collectionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no collectionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from collections")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for collections")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CollectionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), collectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"collections\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, collectionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from collection slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for collections")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Collection) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindCollection(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CollectionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CollectionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), collectionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"collections\".* FROM \"collections\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, collectionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in CollectionSlice")
	}

	*o = slice

	return nil
}

// CollectionExists checks if the Collection row exists.
func CollectionExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"collections\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if collections exists")
	}

	return exists, nil
}
