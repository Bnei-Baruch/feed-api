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

// PersonI18n is an object representing the database table.
type PersonI18n struct {
	PersonID         int64       `boil:"person_id" json:"person_id" toml:"person_id" yaml:"person_id"`
	Language         string      `boil:"language" json:"language" toml:"language" yaml:"language"`
	OriginalLanguage null.String `boil:"original_language" json:"original_language,omitempty" toml:"original_language" yaml:"original_language,omitempty"`
	Name             null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	Description      null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	UserID           null.Int64  `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	CreatedAt        time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *personI18nR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L personI18nL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PersonI18nColumns = struct {
	PersonID         string
	Language         string
	OriginalLanguage string
	Name             string
	Description      string
	UserID           string
	CreatedAt        string
}{
	PersonID:         "person_id",
	Language:         "language",
	OriginalLanguage: "original_language",
	Name:             "name",
	Description:      "description",
	UserID:           "user_id",
	CreatedAt:        "created_at",
}

// Generated where

var PersonI18nWhere = struct {
	PersonID         whereHelperint64
	Language         whereHelperstring
	OriginalLanguage whereHelpernull_String
	Name             whereHelpernull_String
	Description      whereHelpernull_String
	UserID           whereHelpernull_Int64
	CreatedAt        whereHelpertime_Time
}{
	PersonID:         whereHelperint64{field: "\"person_i18n\".\"person_id\""},
	Language:         whereHelperstring{field: "\"person_i18n\".\"language\""},
	OriginalLanguage: whereHelpernull_String{field: "\"person_i18n\".\"original_language\""},
	Name:             whereHelpernull_String{field: "\"person_i18n\".\"name\""},
	Description:      whereHelpernull_String{field: "\"person_i18n\".\"description\""},
	UserID:           whereHelpernull_Int64{field: "\"person_i18n\".\"user_id\""},
	CreatedAt:        whereHelpertime_Time{field: "\"person_i18n\".\"created_at\""},
}

// PersonI18nRels is where relationship names are stored.
var PersonI18nRels = struct {
	Person string
	User   string
}{
	Person: "Person",
	User:   "User",
}

// personI18nR is where relationships are stored.
type personI18nR struct {
	Person *Person
	User   *User
}

// NewStruct creates a new relationship struct
func (*personI18nR) NewStruct() *personI18nR {
	return &personI18nR{}
}

// personI18nL is where Load methods for each relationship are stored.
type personI18nL struct{}

var (
	personI18nAllColumns            = []string{"person_id", "language", "original_language", "name", "description", "user_id", "created_at"}
	personI18nColumnsWithoutDefault = []string{"person_id", "language", "original_language", "name", "description", "user_id"}
	personI18nColumnsWithDefault    = []string{"created_at"}
	personI18nPrimaryKeyColumns     = []string{"person_id", "language"}
)

type (
	// PersonI18nSlice is an alias for a slice of pointers to PersonI18n.
	// This should generally be used opposed to []PersonI18n.
	PersonI18nSlice []*PersonI18n

	personI18nQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	personI18nType                 = reflect.TypeOf(&PersonI18n{})
	personI18nMapping              = queries.MakeStructMapping(personI18nType)
	personI18nPrimaryKeyMapping, _ = queries.BindMapping(personI18nType, personI18nMapping, personI18nPrimaryKeyColumns)
	personI18nInsertCacheMut       sync.RWMutex
	personI18nInsertCache          = make(map[string]insertCache)
	personI18nUpdateCacheMut       sync.RWMutex
	personI18nUpdateCache          = make(map[string]updateCache)
	personI18nUpsertCacheMut       sync.RWMutex
	personI18nUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single personI18n record from the query.
func (q personI18nQuery) One(ctx context.Context, exec boil.ContextExecutor) (*PersonI18n, error) {
	o := &PersonI18n{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for person_i18n")
	}

	return o, nil
}

// All returns all PersonI18n records from the query.
func (q personI18nQuery) All(ctx context.Context, exec boil.ContextExecutor) (PersonI18nSlice, error) {
	var o []*PersonI18n

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to PersonI18n slice")
	}

	return o, nil
}

// Count returns the count of all PersonI18n records in the query.
func (q personI18nQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count person_i18n rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q personI18nQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if person_i18n exists")
	}

	return count > 0, nil
}

// Person pointed to by the foreign key.
func (o *PersonI18n) Person(mods ...qm.QueryMod) personQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.PersonID),
	}

	queryMods = append(queryMods, mods...)

	query := Persons(queryMods...)
	queries.SetFrom(query.Query, "\"persons\"")

	return query
}

// User pointed to by the foreign key.
func (o *PersonI18n) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadPerson allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (personI18nL) LoadPerson(ctx context.Context, e boil.ContextExecutor, singular bool, maybePersonI18n interface{}, mods queries.Applicator) error {
	var slice []*PersonI18n
	var object *PersonI18n

	if singular {
		object = maybePersonI18n.(*PersonI18n)
	} else {
		slice = *maybePersonI18n.(*[]*PersonI18n)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &personI18nR{}
		}
		args = append(args, object.PersonID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &personI18nR{}
			}

			for _, a := range args {
				if a == obj.PersonID {
					continue Outer
				}
			}

			args = append(args, obj.PersonID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`persons`), qm.WhereIn(`persons.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Person")
	}

	var resultSlice []*Person
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Person")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for persons")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for persons")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Person = foreign
		if foreign.R == nil {
			foreign.R = &personR{}
		}
		foreign.R.PersonI18ns = append(foreign.R.PersonI18ns, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.PersonID == foreign.ID {
				local.R.Person = foreign
				if foreign.R == nil {
					foreign.R = &personR{}
				}
				foreign.R.PersonI18ns = append(foreign.R.PersonI18ns, local)
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (personI18nL) LoadUser(ctx context.Context, e boil.ContextExecutor, singular bool, maybePersonI18n interface{}, mods queries.Applicator) error {
	var slice []*PersonI18n
	var object *PersonI18n

	if singular {
		object = maybePersonI18n.(*PersonI18n)
	} else {
		slice = *maybePersonI18n.(*[]*PersonI18n)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &personI18nR{}
		}
		if !queries.IsNil(object.UserID) {
			args = append(args, object.UserID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &personI18nR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.UserID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.UserID) {
				args = append(args, obj.UserID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`users`), qm.WhereIn(`users.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for users")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for users")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.PersonI18ns = append(foreign.R.PersonI18ns, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.UserID, foreign.ID) {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.PersonI18ns = append(foreign.R.PersonI18ns, local)
				break
			}
		}
	}

	return nil
}

// SetPerson of the personI18n to the related item.
// Sets o.R.Person to related.
// Adds o to related.R.PersonI18ns.
func (o *PersonI18n) SetPerson(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Person) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"person_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"person_id"}),
		strmangle.WhereClause("\"", "\"", 2, personI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.PersonID, o.Language}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.PersonID = related.ID
	if o.R == nil {
		o.R = &personI18nR{
			Person: related,
		}
	} else {
		o.R.Person = related
	}

	if related.R == nil {
		related.R = &personR{
			PersonI18ns: PersonI18nSlice{o},
		}
	} else {
		related.R.PersonI18ns = append(related.R.PersonI18ns, o)
	}

	return nil
}

// SetUser of the personI18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.PersonI18ns.
func (o *PersonI18n) SetUser(ctx context.Context, exec boil.ContextExecutor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"person_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, personI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.PersonID, o.Language}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.UserID, related.ID)
	if o.R == nil {
		o.R = &personI18nR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			PersonI18ns: PersonI18nSlice{o},
		}
	} else {
		related.R.PersonI18ns = append(related.R.PersonI18ns, o)
	}

	return nil
}

// RemoveUser relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *PersonI18n) RemoveUser(ctx context.Context, exec boil.ContextExecutor, related *User) error {
	var err error

	queries.SetScanner(&o.UserID, nil)
	if _, err = o.Update(ctx, exec, boil.Whitelist("user_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.User = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.PersonI18ns {
		if queries.Equal(o.UserID, ri.UserID) {
			continue
		}

		ln := len(related.R.PersonI18ns)
		if ln > 1 && i < ln-1 {
			related.R.PersonI18ns[i] = related.R.PersonI18ns[ln-1]
		}
		related.R.PersonI18ns = related.R.PersonI18ns[:ln-1]
		break
	}
	return nil
}

// PersonI18ns retrieves all the records using an executor.
func PersonI18ns(mods ...qm.QueryMod) personI18nQuery {
	mods = append(mods, qm.From("\"person_i18n\""))
	return personI18nQuery{NewQuery(mods...)}
}

// FindPersonI18n retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPersonI18n(ctx context.Context, exec boil.ContextExecutor, personID int64, language string, selectCols ...string) (*PersonI18n, error) {
	personI18nObj := &PersonI18n{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"person_i18n\" where \"person_id\"=$1 AND \"language\"=$2", sel,
	)

	q := queries.Raw(query, personID, language)

	err := q.Bind(ctx, exec, personI18nObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from person_i18n")
	}

	return personI18nObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PersonI18n) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no person_i18n provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(personI18nColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	personI18nInsertCacheMut.RLock()
	cache, cached := personI18nInsertCache[key]
	personI18nInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			personI18nAllColumns,
			personI18nColumnsWithDefault,
			personI18nColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(personI18nType, personI18nMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(personI18nType, personI18nMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"person_i18n\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"person_i18n\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into person_i18n")
	}

	if !cached {
		personI18nInsertCacheMut.Lock()
		personI18nInsertCache[key] = cache
		personI18nInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the PersonI18n.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PersonI18n) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	personI18nUpdateCacheMut.RLock()
	cache, cached := personI18nUpdateCache[key]
	personI18nUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			personI18nAllColumns,
			personI18nPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update person_i18n, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"person_i18n\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, personI18nPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(personI18nType, personI18nMapping, append(wl, personI18nPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update person_i18n row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for person_i18n")
	}

	if !cached {
		personI18nUpdateCacheMut.Lock()
		personI18nUpdateCache[key] = cache
		personI18nUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q personI18nQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for person_i18n")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for person_i18n")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PersonI18nSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), personI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"person_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, personI18nPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in personI18n slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all personI18n")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PersonI18n) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no person_i18n provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(personI18nColumnsWithDefault, o)

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

	personI18nUpsertCacheMut.RLock()
	cache, cached := personI18nUpsertCache[key]
	personI18nUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			personI18nAllColumns,
			personI18nColumnsWithDefault,
			personI18nColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			personI18nAllColumns,
			personI18nPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert person_i18n, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(personI18nPrimaryKeyColumns))
			copy(conflict, personI18nPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"person_i18n\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(personI18nType, personI18nMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(personI18nType, personI18nMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert person_i18n")
	}

	if !cached {
		personI18nUpsertCacheMut.Lock()
		personI18nUpsertCache[key] = cache
		personI18nUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single PersonI18n record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PersonI18n) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no PersonI18n provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), personI18nPrimaryKeyMapping)
	sql := "DELETE FROM \"person_i18n\" WHERE \"person_id\"=$1 AND \"language\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from person_i18n")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for person_i18n")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q personI18nQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no personI18nQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from person_i18n")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for person_i18n")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PersonI18nSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), personI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"person_i18n\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, personI18nPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from personI18n slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for person_i18n")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PersonI18n) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPersonI18n(ctx, exec, o.PersonID, o.Language)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PersonI18nSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PersonI18nSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), personI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"person_i18n\".* FROM \"person_i18n\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, personI18nPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PersonI18nSlice")
	}

	*o = slice

	return nil
}

// PersonI18nExists checks if the PersonI18n row exists.
func PersonI18nExists(ctx context.Context, exec boil.ContextExecutor, personID int64, language string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"person_i18n\" where \"person_id\"=$1 AND \"language\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, personID, language)
	}
	row := exec.QueryRowContext(ctx, sql, personID, language)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if person_i18n exists")
	}

	return exists, nil
}