// Code generated by SQLBoiler 4.8.6 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// SourceI18n is an object representing the database table.
type SourceI18n struct {
	SourceID    int64       `boil:"source_id" json:"source_id" toml:"source_id" yaml:"source_id"`
	Language    string      `boil:"language" json:"language" toml:"language" yaml:"language"`
	Name        null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	CreatedAt   time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *sourceI18nR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L sourceI18nL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SourceI18nColumns = struct {
	SourceID    string
	Language    string
	Name        string
	Description string
	CreatedAt   string
}{
	SourceID:    "source_id",
	Language:    "language",
	Name:        "name",
	Description: "description",
	CreatedAt:   "created_at",
}

var SourceI18nTableColumns = struct {
	SourceID    string
	Language    string
	Name        string
	Description string
	CreatedAt   string
}{
	SourceID:    "source_i18n.source_id",
	Language:    "source_i18n.language",
	Name:        "source_i18n.name",
	Description: "source_i18n.description",
	CreatedAt:   "source_i18n.created_at",
}

// Generated where

var SourceI18nWhere = struct {
	SourceID    whereHelperint64
	Language    whereHelperstring
	Name        whereHelpernull_String
	Description whereHelpernull_String
	CreatedAt   whereHelpertime_Time
}{
	SourceID:    whereHelperint64{field: "\"source_i18n\".\"source_id\""},
	Language:    whereHelperstring{field: "\"source_i18n\".\"language\""},
	Name:        whereHelpernull_String{field: "\"source_i18n\".\"name\""},
	Description: whereHelpernull_String{field: "\"source_i18n\".\"description\""},
	CreatedAt:   whereHelpertime_Time{field: "\"source_i18n\".\"created_at\""},
}

// SourceI18nRels is where relationship names are stored.
var SourceI18nRels = struct {
	Source string
}{
	Source: "Source",
}

// sourceI18nR is where relationships are stored.
type sourceI18nR struct {
	Source *Source `boil:"Source" json:"Source" toml:"Source" yaml:"Source"`
}

// NewStruct creates a new relationship struct
func (*sourceI18nR) NewStruct() *sourceI18nR {
	return &sourceI18nR{}
}

// sourceI18nL is where Load methods for each relationship are stored.
type sourceI18nL struct{}

var (
	sourceI18nAllColumns            = []string{"source_id", "language", "name", "description", "created_at"}
	sourceI18nColumnsWithoutDefault = []string{"source_id", "language"}
	sourceI18nColumnsWithDefault    = []string{"name", "description", "created_at"}
	sourceI18nPrimaryKeyColumns     = []string{"source_id", "language"}
	sourceI18nGeneratedColumns      = []string{}
)

type (
	// SourceI18nSlice is an alias for a slice of pointers to SourceI18n.
	// This should almost always be used instead of []SourceI18n.
	SourceI18nSlice []*SourceI18n

	sourceI18nQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	sourceI18nType                 = reflect.TypeOf(&SourceI18n{})
	sourceI18nMapping              = queries.MakeStructMapping(sourceI18nType)
	sourceI18nPrimaryKeyMapping, _ = queries.BindMapping(sourceI18nType, sourceI18nMapping, sourceI18nPrimaryKeyColumns)
	sourceI18nInsertCacheMut       sync.RWMutex
	sourceI18nInsertCache          = make(map[string]insertCache)
	sourceI18nUpdateCacheMut       sync.RWMutex
	sourceI18nUpdateCache          = make(map[string]updateCache)
	sourceI18nUpsertCacheMut       sync.RWMutex
	sourceI18nUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single sourceI18n record from the query.
func (q sourceI18nQuery) One(exec boil.Executor) (*SourceI18n, error) {
	o := &SourceI18n{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for source_i18n")
	}

	return o, nil
}

// All returns all SourceI18n records from the query.
func (q sourceI18nQuery) All(exec boil.Executor) (SourceI18nSlice, error) {
	var o []*SourceI18n

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to SourceI18n slice")
	}

	return o, nil
}

// Count returns the count of all SourceI18n records in the query.
func (q sourceI18nQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count source_i18n rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q sourceI18nQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if source_i18n exists")
	}

	return count > 0, nil
}

// Source pointed to by the foreign key.
func (o *SourceI18n) Source(mods ...qm.QueryMod) sourceQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SourceID),
	}

	queryMods = append(queryMods, mods...)

	query := Sources(queryMods...)
	queries.SetFrom(query.Query, "\"sources\"")

	return query
}

// LoadSource allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (sourceI18nL) LoadSource(e boil.Executor, singular bool, maybeSourceI18n interface{}, mods queries.Applicator) error {
	var slice []*SourceI18n
	var object *SourceI18n

	if singular {
		object = maybeSourceI18n.(*SourceI18n)
	} else {
		slice = *maybeSourceI18n.(*[]*SourceI18n)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &sourceI18nR{}
		}
		args = append(args, object.SourceID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &sourceI18nR{}
			}

			for _, a := range args {
				if a == obj.SourceID {
					continue Outer
				}
			}

			args = append(args, obj.SourceID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`sources`),
		qm.WhereIn(`sources.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Source")
	}

	var resultSlice []*Source
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Source")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for sources")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for sources")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Source = foreign
		if foreign.R == nil {
			foreign.R = &sourceR{}
		}
		foreign.R.SourceI18ns = append(foreign.R.SourceI18ns, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SourceID == foreign.ID {
				local.R.Source = foreign
				if foreign.R == nil {
					foreign.R = &sourceR{}
				}
				foreign.R.SourceI18ns = append(foreign.R.SourceI18ns, local)
				break
			}
		}
	}

	return nil
}

// SetSource of the sourceI18n to the related item.
// Sets o.R.Source to related.
// Adds o to related.R.SourceI18ns.
func (o *SourceI18n) SetSource(exec boil.Executor, insert bool, related *Source) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"source_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"source_id"}),
		strmangle.WhereClause("\"", "\"", 2, sourceI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.SourceID, o.Language}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.SourceID = related.ID
	if o.R == nil {
		o.R = &sourceI18nR{
			Source: related,
		}
	} else {
		o.R.Source = related
	}

	if related.R == nil {
		related.R = &sourceR{
			SourceI18ns: SourceI18nSlice{o},
		}
	} else {
		related.R.SourceI18ns = append(related.R.SourceI18ns, o)
	}

	return nil
}

// SourceI18ns retrieves all the records using an executor.
func SourceI18ns(mods ...qm.QueryMod) sourceI18nQuery {
	mods = append(mods, qm.From("\"source_i18n\""))
	return sourceI18nQuery{NewQuery(mods...)}
}

// FindSourceI18n retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSourceI18n(exec boil.Executor, sourceID int64, language string, selectCols ...string) (*SourceI18n, error) {
	sourceI18nObj := &SourceI18n{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"source_i18n\" where \"source_id\"=$1 AND \"language\"=$2", sel,
	)

	q := queries.Raw(query, sourceID, language)

	err := q.Bind(nil, exec, sourceI18nObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from source_i18n")
	}

	return sourceI18nObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SourceI18n) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no source_i18n provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(sourceI18nColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	sourceI18nInsertCacheMut.RLock()
	cache, cached := sourceI18nInsertCache[key]
	sourceI18nInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			sourceI18nAllColumns,
			sourceI18nColumnsWithDefault,
			sourceI18nColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"source_i18n\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"source_i18n\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
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
		return errors.Wrap(err, "models: unable to insert into source_i18n")
	}

	if !cached {
		sourceI18nInsertCacheMut.Lock()
		sourceI18nInsertCache[key] = cache
		sourceI18nInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the SourceI18n.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SourceI18n) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	sourceI18nUpdateCacheMut.RLock()
	cache, cached := sourceI18nUpdateCache[key]
	sourceI18nUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			sourceI18nAllColumns,
			sourceI18nPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update source_i18n, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"source_i18n\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sourceI18nPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, append(wl, sourceI18nPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	var result sql.Result
	result, err = exec.Exec(cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update source_i18n row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for source_i18n")
	}

	if !cached {
		sourceI18nUpdateCacheMut.Lock()
		sourceI18nUpdateCache[key] = cache
		sourceI18nUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q sourceI18nQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for source_i18n")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for source_i18n")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SourceI18nSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourceI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"source_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, sourceI18nPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in sourceI18n slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all sourceI18n")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SourceI18n) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no source_i18n provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(sourceI18nColumnsWithDefault, o)

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

	sourceI18nUpsertCacheMut.RLock()
	cache, cached := sourceI18nUpsertCache[key]
	sourceI18nUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			sourceI18nAllColumns,
			sourceI18nColumnsWithDefault,
			sourceI18nColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			sourceI18nAllColumns,
			sourceI18nPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert source_i18n, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(sourceI18nPrimaryKeyColumns))
			copy(conflict, sourceI18nPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"source_i18n\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(sourceI18nType, sourceI18nMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert source_i18n")
	}

	if !cached {
		sourceI18nUpsertCacheMut.Lock()
		sourceI18nUpsertCache[key] = cache
		sourceI18nUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single SourceI18n record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SourceI18n) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no SourceI18n provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sourceI18nPrimaryKeyMapping)
	sql := "DELETE FROM \"source_i18n\" WHERE \"source_id\"=$1 AND \"language\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from source_i18n")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for source_i18n")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q sourceI18nQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no sourceI18nQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from source_i18n")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for source_i18n")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SourceI18nSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourceI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"source_i18n\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sourceI18nPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from sourceI18n slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for source_i18n")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SourceI18n) Reload(exec boil.Executor) error {
	ret, err := FindSourceI18n(exec, o.SourceID, o.Language)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SourceI18nSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SourceI18nSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sourceI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"source_i18n\".* FROM \"source_i18n\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sourceI18nPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SourceI18nSlice")
	}

	*o = slice

	return nil
}

// SourceI18nExists checks if the SourceI18n row exists.
func SourceI18nExists(exec boil.Executor, sourceID int64, language string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"source_i18n\" where \"source_id\"=$1 AND \"language\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, sourceID, language)
	}
	row := exec.QueryRow(sql, sourceID, language)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if source_i18n exists")
	}

	return exists, nil
}
