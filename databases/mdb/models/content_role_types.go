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

// ContentRoleType is an object representing the database table.
type ContentRoleType struct {
	ID          int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name        string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`

	R *contentRoleTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L contentRoleTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ContentRoleTypeColumns = struct {
	ID          string
	Name        string
	Description string
}{
	ID:          "id",
	Name:        "name",
	Description: "description",
}

// Generated where

var ContentRoleTypeWhere = struct {
	ID          whereHelperint64
	Name        whereHelperstring
	Description whereHelpernull_String
}{
	ID:          whereHelperint64{field: "\"content_role_types\".\"id\""},
	Name:        whereHelperstring{field: "\"content_role_types\".\"name\""},
	Description: whereHelpernull_String{field: "\"content_role_types\".\"description\""},
}

// ContentRoleTypeRels is where relationship names are stored.
var ContentRoleTypeRels = struct {
	RoleContentUnitsPersons string
}{
	RoleContentUnitsPersons: "RoleContentUnitsPersons",
}

// contentRoleTypeR is where relationships are stored.
type contentRoleTypeR struct {
	RoleContentUnitsPersons ContentUnitsPersonSlice
}

// NewStruct creates a new relationship struct
func (*contentRoleTypeR) NewStruct() *contentRoleTypeR {
	return &contentRoleTypeR{}
}

// contentRoleTypeL is where Load methods for each relationship are stored.
type contentRoleTypeL struct{}

var (
	contentRoleTypeAllColumns            = []string{"id", "name", "description"}
	contentRoleTypeColumnsWithoutDefault = []string{"name", "description"}
	contentRoleTypeColumnsWithDefault    = []string{"id"}
	contentRoleTypePrimaryKeyColumns     = []string{"id"}
)

type (
	// ContentRoleTypeSlice is an alias for a slice of pointers to ContentRoleType.
	// This should generally be used opposed to []ContentRoleType.
	ContentRoleTypeSlice []*ContentRoleType

	contentRoleTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	contentRoleTypeType                 = reflect.TypeOf(&ContentRoleType{})
	contentRoleTypeMapping              = queries.MakeStructMapping(contentRoleTypeType)
	contentRoleTypePrimaryKeyMapping, _ = queries.BindMapping(contentRoleTypeType, contentRoleTypeMapping, contentRoleTypePrimaryKeyColumns)
	contentRoleTypeInsertCacheMut       sync.RWMutex
	contentRoleTypeInsertCache          = make(map[string]insertCache)
	contentRoleTypeUpdateCacheMut       sync.RWMutex
	contentRoleTypeUpdateCache          = make(map[string]updateCache)
	contentRoleTypeUpsertCacheMut       sync.RWMutex
	contentRoleTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single contentRoleType record from the query.
func (q contentRoleTypeQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ContentRoleType, error) {
	o := &ContentRoleType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for content_role_types")
	}

	return o, nil
}

// All returns all ContentRoleType records from the query.
func (q contentRoleTypeQuery) All(ctx context.Context, exec boil.ContextExecutor) (ContentRoleTypeSlice, error) {
	var o []*ContentRoleType

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ContentRoleType slice")
	}

	return o, nil
}

// Count returns the count of all ContentRoleType records in the query.
func (q contentRoleTypeQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count content_role_types rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q contentRoleTypeQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if content_role_types exists")
	}

	return count > 0, nil
}

// RoleContentUnitsPersons retrieves all the content_units_person's ContentUnitsPersons with an executor via role_id column.
func (o *ContentRoleType) RoleContentUnitsPersons(mods ...qm.QueryMod) contentUnitsPersonQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"content_units_persons\".\"role_id\"=?", o.ID),
	)

	query := ContentUnitsPersons(queryMods...)
	queries.SetFrom(query.Query, "\"content_units_persons\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"content_units_persons\".*"})
	}

	return query
}

// LoadRoleContentUnitsPersons allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (contentRoleTypeL) LoadRoleContentUnitsPersons(ctx context.Context, e boil.ContextExecutor, singular bool, maybeContentRoleType interface{}, mods queries.Applicator) error {
	var slice []*ContentRoleType
	var object *ContentRoleType

	if singular {
		object = maybeContentRoleType.(*ContentRoleType)
	} else {
		slice = *maybeContentRoleType.(*[]*ContentRoleType)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &contentRoleTypeR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &contentRoleTypeR{}
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

	query := NewQuery(qm.From(`content_units_persons`), qm.WhereIn(`content_units_persons.role_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load content_units_persons")
	}

	var resultSlice []*ContentUnitsPerson
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice content_units_persons")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on content_units_persons")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for content_units_persons")
	}

	if singular {
		object.R.RoleContentUnitsPersons = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &contentUnitsPersonR{}
			}
			foreign.R.Role = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.RoleID {
				local.R.RoleContentUnitsPersons = append(local.R.RoleContentUnitsPersons, foreign)
				if foreign.R == nil {
					foreign.R = &contentUnitsPersonR{}
				}
				foreign.R.Role = local
				break
			}
		}
	}

	return nil
}

// AddRoleContentUnitsPersons adds the given related objects to the existing relationships
// of the content_role_type, optionally inserting them as new records.
// Appends related to o.R.RoleContentUnitsPersons.
// Sets related.R.Role appropriately.
func (o *ContentRoleType) AddRoleContentUnitsPersons(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*ContentUnitsPerson) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.RoleID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"content_units_persons\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"role_id"}),
				strmangle.WhereClause("\"", "\"", 2, contentUnitsPersonPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ContentUnitID, rel.PersonID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.RoleID = o.ID
		}
	}

	if o.R == nil {
		o.R = &contentRoleTypeR{
			RoleContentUnitsPersons: related,
		}
	} else {
		o.R.RoleContentUnitsPersons = append(o.R.RoleContentUnitsPersons, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &contentUnitsPersonR{
				Role: o,
			}
		} else {
			rel.R.Role = o
		}
	}
	return nil
}

// ContentRoleTypes retrieves all the records using an executor.
func ContentRoleTypes(mods ...qm.QueryMod) contentRoleTypeQuery {
	mods = append(mods, qm.From("\"content_role_types\""))
	return contentRoleTypeQuery{NewQuery(mods...)}
}

// FindContentRoleType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindContentRoleType(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*ContentRoleType, error) {
	contentRoleTypeObj := &ContentRoleType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"content_role_types\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, contentRoleTypeObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from content_role_types")
	}

	return contentRoleTypeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ContentRoleType) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no content_role_types provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(contentRoleTypeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	contentRoleTypeInsertCacheMut.RLock()
	cache, cached := contentRoleTypeInsertCache[key]
	contentRoleTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			contentRoleTypeAllColumns,
			contentRoleTypeColumnsWithDefault,
			contentRoleTypeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(contentRoleTypeType, contentRoleTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(contentRoleTypeType, contentRoleTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"content_role_types\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"content_role_types\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into content_role_types")
	}

	if !cached {
		contentRoleTypeInsertCacheMut.Lock()
		contentRoleTypeInsertCache[key] = cache
		contentRoleTypeInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the ContentRoleType.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ContentRoleType) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	contentRoleTypeUpdateCacheMut.RLock()
	cache, cached := contentRoleTypeUpdateCache[key]
	contentRoleTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			contentRoleTypeAllColumns,
			contentRoleTypePrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update content_role_types, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"content_role_types\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, contentRoleTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(contentRoleTypeType, contentRoleTypeMapping, append(wl, contentRoleTypePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update content_role_types row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for content_role_types")
	}

	if !cached {
		contentRoleTypeUpdateCacheMut.Lock()
		contentRoleTypeUpdateCache[key] = cache
		contentRoleTypeUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q contentRoleTypeQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for content_role_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for content_role_types")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ContentRoleTypeSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentRoleTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"content_role_types\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, contentRoleTypePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in contentRoleType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all contentRoleType")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ContentRoleType) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no content_role_types provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(contentRoleTypeColumnsWithDefault, o)

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

	contentRoleTypeUpsertCacheMut.RLock()
	cache, cached := contentRoleTypeUpsertCache[key]
	contentRoleTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			contentRoleTypeAllColumns,
			contentRoleTypeColumnsWithDefault,
			contentRoleTypeColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			contentRoleTypeAllColumns,
			contentRoleTypePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert content_role_types, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(contentRoleTypePrimaryKeyColumns))
			copy(conflict, contentRoleTypePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"content_role_types\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(contentRoleTypeType, contentRoleTypeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(contentRoleTypeType, contentRoleTypeMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert content_role_types")
	}

	if !cached {
		contentRoleTypeUpsertCacheMut.Lock()
		contentRoleTypeUpsertCache[key] = cache
		contentRoleTypeUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single ContentRoleType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ContentRoleType) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ContentRoleType provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), contentRoleTypePrimaryKeyMapping)
	sql := "DELETE FROM \"content_role_types\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from content_role_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for content_role_types")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q contentRoleTypeQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no contentRoleTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from content_role_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for content_role_types")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ContentRoleTypeSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentRoleTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"content_role_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, contentRoleTypePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from contentRoleType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for content_role_types")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ContentRoleType) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindContentRoleType(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContentRoleTypeSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ContentRoleTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentRoleTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"content_role_types\".* FROM \"content_role_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, contentRoleTypePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ContentRoleTypeSlice")
	}

	*o = slice

	return nil
}

// ContentRoleTypeExists checks if the ContentRoleType row exists.
func ContentRoleTypeExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"content_role_types\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if content_role_types exists")
	}

	return exists, nil
}
