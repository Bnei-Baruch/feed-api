// Code generated by SQLBoiler (https://github.com/Bnei-Baruch/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package mdbmdbmodels

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/boil"
	"github.com/Bnei-Baruch/sqlboiler/queries"
	"github.com/Bnei-Baruch/sqlboiler/queries/qm"
	"github.com/Bnei-Baruch/sqlboiler/strmangle"
	"github.com/pkg/errors"
	"gopkg.in/volatiletech/null.v6"
)

// ContentType is an object representing the database table.
type ContentType struct {
	ID          int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name        string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`

	R *contentTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L contentTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ContentTypeColumns = struct {
	ID          string
	Name        string
	Description string
}{
	ID:          "id",
	Name:        "name",
	Description: "description",
}

// contentTypeR is where relationships are stored.
type contentTypeR struct {
	TypeCollections  CollectionSlice
	TypeContentUnits ContentUnitSlice
}

// contentTypeL is where Load methods for each relationship are stored.
type contentTypeL struct{}

var (
	contentTypeColumns               = []string{"id", "name", "description"}
	contentTypeColumnsWithoutDefault = []string{"name", "description"}
	contentTypeColumnsWithDefault    = []string{"id"}
	contentTypePrimaryKeyColumns     = []string{"id"}
)

type (
	// ContentTypeSlice is an alias for a slice of pointers to ContentType.
	// This should generally be used opposed to []ContentType.
	ContentTypeSlice []*ContentType

	contentTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	contentTypeType                 = reflect.TypeOf(&ContentType{})
	contentTypeMapping              = queries.MakeStructMapping(contentTypeType)
	contentTypePrimaryKeyMapping, _ = queries.BindMapping(contentTypeType, contentTypeMapping, contentTypePrimaryKeyColumns)
	contentTypeInsertCacheMut       sync.RWMutex
	contentTypeInsertCache          = make(map[string]insertCache)
	contentTypeUpdateCacheMut       sync.RWMutex
	contentTypeUpdateCache          = make(map[string]updateCache)
	contentTypeUpsertCacheMut       sync.RWMutex
	contentTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single contentType record from the query, and panics on error.
func (q contentTypeQuery) OneP() *ContentType {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single contentType record from the query.
func (q contentTypeQuery) One() (*ContentType, error) {
	o := &ContentType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mdbmdbmodels: failed to execute a one query for content_types")
	}

	return o, nil
}

// AllP returns all ContentType records from the query, and panics on error.
func (q contentTypeQuery) AllP() ContentTypeSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all ContentType records from the query.
func (q contentTypeQuery) All() (ContentTypeSlice, error) {
	var o []*ContentType

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "mdbmdbmodels: failed to assign all query results to ContentType slice")
	}

	return o, nil
}

// CountP returns the count of all ContentType records in the query, and panics on error.
func (q contentTypeQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all ContentType records in the query.
func (q contentTypeQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mdbmdbmodels: failed to count content_types rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q contentTypeQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q contentTypeQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mdbmdbmodels: failed to check if content_types exists")
	}

	return count > 0, nil
}

// TypeCollectionsG retrieves all the collection's collections via type_id column.
func (o *ContentType) TypeCollectionsG(mods ...qm.QueryMod) collectionQuery {
	return o.TypeCollections(boil.GetDB(), mods...)
}

// TypeCollections retrieves all the collection's collections with an executor via type_id column.
func (o *ContentType) TypeCollections(exec boil.Executor, mods ...qm.QueryMod) collectionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"collections\".\"type_id\"=?", o.ID),
	)

	query := Collections(exec, queryMods...)
	queries.SetFrom(query.Query, "\"collections\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"collections\".*"})
	}

	return query
}

// TypeContentUnitsG retrieves all the content_unit's content units via type_id column.
func (o *ContentType) TypeContentUnitsG(mods ...qm.QueryMod) contentUnitQuery {
	return o.TypeContentUnits(boil.GetDB(), mods...)
}

// TypeContentUnits retrieves all the content_unit's content units with an executor via type_id column.
func (o *ContentType) TypeContentUnits(exec boil.Executor, mods ...qm.QueryMod) contentUnitQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"content_units\".\"type_id\"=?", o.ID),
	)

	query := ContentUnits(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_units\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"content_units\".*"})
	}

	return query
}

// LoadTypeCollections allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (contentTypeL) LoadTypeCollections(e boil.Executor, singular bool, maybeContentType interface{}) error {
	var slice []*ContentType
	var object *ContentType

	count := 1
	if singular {
		object = maybeContentType.(*ContentType)
	} else {
		slice = *maybeContentType.(*[]*ContentType)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentTypeR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentTypeR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"collections\" where \"type_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load collections")
	}
	defer results.Close()

	var resultSlice []*Collection
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice collections")
	}

	if singular {
		object.R.TypeCollections = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.TypeID {
				local.R.TypeCollections = append(local.R.TypeCollections, foreign)
				break
			}
		}
	}

	return nil
}

// LoadTypeContentUnits allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (contentTypeL) LoadTypeContentUnits(e boil.Executor, singular bool, maybeContentType interface{}) error {
	var slice []*ContentType
	var object *ContentType

	count := 1
	if singular {
		object = maybeContentType.(*ContentType)
	} else {
		slice = *maybeContentType.(*[]*ContentType)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentTypeR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentTypeR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_units\" where \"type_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load content_units")
	}
	defer results.Close()

	var resultSlice []*ContentUnit
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice content_units")
	}

	if singular {
		object.R.TypeContentUnits = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.TypeID {
				local.R.TypeContentUnits = append(local.R.TypeContentUnits, foreign)
				break
			}
		}
	}

	return nil
}

// AddTypeCollectionsG adds the given related objects to the existing relationships
// of the content_type, optionally inserting them as new records.
// Appends related to o.R.TypeCollections.
// Sets related.R.Type appropriately.
// Uses the global database handle.
func (o *ContentType) AddTypeCollectionsG(insert bool, related ...*Collection) error {
	return o.AddTypeCollections(boil.GetDB(), insert, related...)
}

// AddTypeCollectionsP adds the given related objects to the existing relationships
// of the content_type, optionally inserting them as new records.
// Appends related to o.R.TypeCollections.
// Sets related.R.Type appropriately.
// Panics on error.
func (o *ContentType) AddTypeCollectionsP(exec boil.Executor, insert bool, related ...*Collection) {
	if err := o.AddTypeCollections(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddTypeCollectionsGP adds the given related objects to the existing relationships
// of the content_type, optionally inserting them as new records.
// Appends related to o.R.TypeCollections.
// Sets related.R.Type appropriately.
// Uses the global database handle and panics on error.
func (o *ContentType) AddTypeCollectionsGP(insert bool, related ...*Collection) {
	if err := o.AddTypeCollections(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddTypeCollections adds the given related objects to the existing relationships
// of the content_type, optionally inserting them as new records.
// Appends related to o.R.TypeCollections.
// Sets related.R.Type appropriately.
func (o *ContentType) AddTypeCollections(exec boil.Executor, insert bool, related ...*Collection) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.TypeID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"collections\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"type_id"}),
				strmangle.WhereClause("\"", "\"", 2, collectionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.TypeID = o.ID
		}
	}

	if o.R == nil {
		o.R = &contentTypeR{
			TypeCollections: related,
		}
	} else {
		o.R.TypeCollections = append(o.R.TypeCollections, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &collectionR{
				Type: o,
			}
		} else {
			rel.R.Type = o
		}
	}
	return nil
}

// AddTypeContentUnitsG adds the given related objects to the existing relationships
// of the content_type, optionally inserting them as new records.
// Appends related to o.R.TypeContentUnits.
// Sets related.R.Type appropriately.
// Uses the global database handle.
func (o *ContentType) AddTypeContentUnitsG(insert bool, related ...*ContentUnit) error {
	return o.AddTypeContentUnits(boil.GetDB(), insert, related...)
}

// AddTypeContentUnitsP adds the given related objects to the existing relationships
// of the content_type, optionally inserting them as new records.
// Appends related to o.R.TypeContentUnits.
// Sets related.R.Type appropriately.
// Panics on error.
func (o *ContentType) AddTypeContentUnitsP(exec boil.Executor, insert bool, related ...*ContentUnit) {
	if err := o.AddTypeContentUnits(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddTypeContentUnitsGP adds the given related objects to the existing relationships
// of the content_type, optionally inserting them as new records.
// Appends related to o.R.TypeContentUnits.
// Sets related.R.Type appropriately.
// Uses the global database handle and panics on error.
func (o *ContentType) AddTypeContentUnitsGP(insert bool, related ...*ContentUnit) {
	if err := o.AddTypeContentUnits(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddTypeContentUnits adds the given related objects to the existing relationships
// of the content_type, optionally inserting them as new records.
// Appends related to o.R.TypeContentUnits.
// Sets related.R.Type appropriately.
func (o *ContentType) AddTypeContentUnits(exec boil.Executor, insert bool, related ...*ContentUnit) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.TypeID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"content_units\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"type_id"}),
				strmangle.WhereClause("\"", "\"", 2, contentUnitPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.TypeID = o.ID
		}
	}

	if o.R == nil {
		o.R = &contentTypeR{
			TypeContentUnits: related,
		}
	} else {
		o.R.TypeContentUnits = append(o.R.TypeContentUnits, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &contentUnitR{
				Type: o,
			}
		} else {
			rel.R.Type = o
		}
	}
	return nil
}

// ContentTypesG retrieves all records.
func ContentTypesG(mods ...qm.QueryMod) contentTypeQuery {
	return ContentTypes(boil.GetDB(), mods...)
}

// ContentTypes retrieves all the records using an executor.
func ContentTypes(exec boil.Executor, mods ...qm.QueryMod) contentTypeQuery {
	mods = append(mods, qm.From("\"content_types\""))
	return contentTypeQuery{NewQuery(exec, mods...)}
}

// FindContentTypeG retrieves a single record by ID.
func FindContentTypeG(id int64, selectCols ...string) (*ContentType, error) {
	return FindContentType(boil.GetDB(), id, selectCols...)
}

// FindContentTypeGP retrieves a single record by ID, and panics on error.
func FindContentTypeGP(id int64, selectCols ...string) *ContentType {
	retobj, err := FindContentType(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindContentType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindContentType(exec boil.Executor, id int64, selectCols ...string) (*ContentType, error) {
	contentTypeObj := &ContentType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"content_types\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(contentTypeObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mdbmdbmodels: unable to select from content_types")
	}

	return contentTypeObj, nil
}

// FindContentTypeP retrieves a single record by ID with an executor, and panics on error.
func FindContentTypeP(exec boil.Executor, id int64, selectCols ...string) *ContentType {
	retobj, err := FindContentType(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *ContentType) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *ContentType) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *ContentType) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *ContentType) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("mdbmdbmodels: no content_types provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(contentTypeColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	contentTypeInsertCacheMut.RLock()
	cache, cached := contentTypeInsertCache[key]
	contentTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			contentTypeColumns,
			contentTypeColumnsWithDefault,
			contentTypeColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(contentTypeType, contentTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(contentTypeType, contentTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"content_types\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"content_types\" DEFAULT VALUES"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to insert into content_types")
	}

	if !cached {
		contentTypeInsertCacheMut.Lock()
		contentTypeInsertCache[key] = cache
		contentTypeInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single ContentType record. See Update for
// whitelist behavior description.
func (o *ContentType) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single ContentType record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *ContentType) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the ContentType, and panics on error.
// See Update for whitelist behavior description.
func (o *ContentType) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the ContentType.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *ContentType) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	contentTypeUpdateCacheMut.RLock()
	cache, cached := contentTypeUpdateCache[key]
	contentTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			contentTypeColumns,
			contentTypePrimaryKeyColumns,
			whitelist,
		)

		if len(wl) == 0 {
			return errors.New("mdbmdbmodels: unable to update content_types, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"content_types\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, contentTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(contentTypeType, contentTypeMapping, append(wl, contentTypePrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to update content_types row")
	}

	if !cached {
		contentTypeUpdateCacheMut.Lock()
		contentTypeUpdateCache[key] = cache
		contentTypeUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q contentTypeQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q contentTypeQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to update all for content_types")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ContentTypeSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o ContentTypeSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o ContentTypeSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ContentTypeSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("mdbmdbmodels: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"content_types\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, contentTypePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to update all in contentType slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *ContentType) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *ContentType) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *ContentType) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *ContentType) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("mdbmdbmodels: no content_types provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(contentTypeColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
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
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	contentTypeUpsertCacheMut.RLock()
	cache, cached := contentTypeUpsertCache[key]
	contentTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertColumnSet(
			contentTypeColumns,
			contentTypeColumnsWithDefault,
			contentTypeColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		update := strmangle.UpdateColumnSet(
			contentTypeColumns,
			contentTypePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("mdbmdbmodels: unable to upsert content_types, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(contentTypePrimaryKeyColumns))
			copy(conflict, contentTypePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"content_types\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(contentTypeType, contentTypeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(contentTypeType, contentTypeMapping, ret)
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

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to upsert content_types")
	}

	if !cached {
		contentTypeUpsertCacheMut.Lock()
		contentTypeUpsertCache[key] = cache
		contentTypeUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single ContentType record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContentType) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single ContentType record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *ContentType) DeleteG() error {
	if o == nil {
		return errors.New("mdbmdbmodels: no ContentType provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single ContentType record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContentType) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single ContentType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ContentType) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("mdbmdbmodels: no ContentType provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), contentTypePrimaryKeyMapping)
	sql := "DELETE FROM \"content_types\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to delete from content_types")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q contentTypeQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q contentTypeQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("mdbmdbmodels: no contentTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to delete all from content_types")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o ContentTypeSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o ContentTypeSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("mdbmdbmodels: no ContentType slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o ContentTypeSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ContentTypeSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("mdbmdbmodels: no ContentType slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"content_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, contentTypePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to delete all from contentType slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *ContentType) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *ContentType) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *ContentType) ReloadG() error {
	if o == nil {
		return errors.New("mdbmdbmodels: no ContentType provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ContentType) Reload(exec boil.Executor) error {
	ret, err := FindContentType(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContentTypeSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContentTypeSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContentTypeSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("mdbmdbmodels: empty ContentTypeSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContentTypeSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	contentTypes := ContentTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"content_types\".* FROM \"content_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, contentTypePrimaryKeyColumns, len(*o))

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&contentTypes)
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to reload all in ContentTypeSlice")
	}

	*o = contentTypes

	return nil
}

// ContentTypeExists checks if the ContentType row exists.
func ContentTypeExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"content_types\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mdbmdbmodels: unable to check if content_types exists")
	}

	return exists, nil
}

// ContentTypeExistsG checks if the ContentType row exists.
func ContentTypeExistsG(id int64) (bool, error) {
	return ContentTypeExists(boil.GetDB(), id)
}

// ContentTypeExistsGP checks if the ContentType row exists. Panics on error.
func ContentTypeExistsGP(id int64) bool {
	e, err := ContentTypeExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ContentTypeExistsP checks if the ContentType row exists. Panics on error.
func ContentTypeExistsP(exec boil.Executor, id int64) bool {
	e, err := ContentTypeExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
