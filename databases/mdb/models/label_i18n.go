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

// LabelI18n is an object representing the database table.
type LabelI18n struct {
	LabelID   int64       `boil:"label_id" json:"label_id" toml:"label_id" yaml:"label_id"`
	Language  string      `boil:"language" json:"language" toml:"language" yaml:"language"`
	Name      null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	UserID    null.Int64  `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	CreatedAt time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *labelI18nR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L labelI18nL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var LabelI18nColumns = struct {
	LabelID   string
	Language  string
	Name      string
	UserID    string
	CreatedAt string
}{
	LabelID:   "label_id",
	Language:  "language",
	Name:      "name",
	UserID:    "user_id",
	CreatedAt: "created_at",
}

var LabelI18nTableColumns = struct {
	LabelID   string
	Language  string
	Name      string
	UserID    string
	CreatedAt string
}{
	LabelID:   "label_i18n.label_id",
	Language:  "label_i18n.language",
	Name:      "label_i18n.name",
	UserID:    "label_i18n.user_id",
	CreatedAt: "label_i18n.created_at",
}

// Generated where

var LabelI18nWhere = struct {
	LabelID   whereHelperint64
	Language  whereHelperstring
	Name      whereHelpernull_String
	UserID    whereHelpernull_Int64
	CreatedAt whereHelpertime_Time
}{
	LabelID:   whereHelperint64{field: "\"label_i18n\".\"label_id\""},
	Language:  whereHelperstring{field: "\"label_i18n\".\"language\""},
	Name:      whereHelpernull_String{field: "\"label_i18n\".\"name\""},
	UserID:    whereHelpernull_Int64{field: "\"label_i18n\".\"user_id\""},
	CreatedAt: whereHelpertime_Time{field: "\"label_i18n\".\"created_at\""},
}

// LabelI18nRels is where relationship names are stored.
var LabelI18nRels = struct {
	Label string
	User  string
}{
	Label: "Label",
	User:  "User",
}

// labelI18nR is where relationships are stored.
type labelI18nR struct {
	Label *Label `boil:"Label" json:"Label" toml:"Label" yaml:"Label"`
	User  *User  `boil:"User" json:"User" toml:"User" yaml:"User"`
}

// NewStruct creates a new relationship struct
func (*labelI18nR) NewStruct() *labelI18nR {
	return &labelI18nR{}
}

// labelI18nL is where Load methods for each relationship are stored.
type labelI18nL struct{}

var (
	labelI18nAllColumns            = []string{"label_id", "language", "name", "user_id", "created_at"}
	labelI18nColumnsWithoutDefault = []string{"label_id", "language"}
	labelI18nColumnsWithDefault    = []string{"name", "user_id", "created_at"}
	labelI18nPrimaryKeyColumns     = []string{"label_id", "language"}
	labelI18nGeneratedColumns      = []string{}
)

type (
	// LabelI18nSlice is an alias for a slice of pointers to LabelI18n.
	// This should almost always be used instead of []LabelI18n.
	LabelI18nSlice []*LabelI18n

	labelI18nQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	labelI18nType                 = reflect.TypeOf(&LabelI18n{})
	labelI18nMapping              = queries.MakeStructMapping(labelI18nType)
	labelI18nPrimaryKeyMapping, _ = queries.BindMapping(labelI18nType, labelI18nMapping, labelI18nPrimaryKeyColumns)
	labelI18nInsertCacheMut       sync.RWMutex
	labelI18nInsertCache          = make(map[string]insertCache)
	labelI18nUpdateCacheMut       sync.RWMutex
	labelI18nUpdateCache          = make(map[string]updateCache)
	labelI18nUpsertCacheMut       sync.RWMutex
	labelI18nUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single labelI18n record from the query.
func (q labelI18nQuery) One(exec boil.Executor) (*LabelI18n, error) {
	o := &LabelI18n{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for label_i18n")
	}

	return o, nil
}

// All returns all LabelI18n records from the query.
func (q labelI18nQuery) All(exec boil.Executor) (LabelI18nSlice, error) {
	var o []*LabelI18n

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to LabelI18n slice")
	}

	return o, nil
}

// Count returns the count of all LabelI18n records in the query.
func (q labelI18nQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count label_i18n rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q labelI18nQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if label_i18n exists")
	}

	return count > 0, nil
}

// Label pointed to by the foreign key.
func (o *LabelI18n) Label(mods ...qm.QueryMod) labelQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.LabelID),
	}

	queryMods = append(queryMods, mods...)

	query := Labels(queryMods...)
	queries.SetFrom(query.Query, "\"labels\"")

	return query
}

// User pointed to by the foreign key.
func (o *LabelI18n) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadLabel allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (labelI18nL) LoadLabel(e boil.Executor, singular bool, maybeLabelI18n interface{}, mods queries.Applicator) error {
	var slice []*LabelI18n
	var object *LabelI18n

	if singular {
		object = maybeLabelI18n.(*LabelI18n)
	} else {
		slice = *maybeLabelI18n.(*[]*LabelI18n)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &labelI18nR{}
		}
		args = append(args, object.LabelID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &labelI18nR{}
			}

			for _, a := range args {
				if a == obj.LabelID {
					continue Outer
				}
			}

			args = append(args, obj.LabelID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`labels`),
		qm.WhereIn(`labels.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Label")
	}

	var resultSlice []*Label
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Label")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for labels")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for labels")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Label = foreign
		if foreign.R == nil {
			foreign.R = &labelR{}
		}
		foreign.R.LabelI18ns = append(foreign.R.LabelI18ns, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.LabelID == foreign.ID {
				local.R.Label = foreign
				if foreign.R == nil {
					foreign.R = &labelR{}
				}
				foreign.R.LabelI18ns = append(foreign.R.LabelI18ns, local)
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (labelI18nL) LoadUser(e boil.Executor, singular bool, maybeLabelI18n interface{}, mods queries.Applicator) error {
	var slice []*LabelI18n
	var object *LabelI18n

	if singular {
		object = maybeLabelI18n.(*LabelI18n)
	} else {
		slice = *maybeLabelI18n.(*[]*LabelI18n)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &labelI18nR{}
		}
		if !queries.IsNil(object.UserID) {
			args = append(args, object.UserID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &labelI18nR{}
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

	query := NewQuery(
		qm.From(`users`),
		qm.WhereIn(`users.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
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
		foreign.R.LabelI18ns = append(foreign.R.LabelI18ns, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.UserID, foreign.ID) {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.LabelI18ns = append(foreign.R.LabelI18ns, local)
				break
			}
		}
	}

	return nil
}

// SetLabel of the labelI18n to the related item.
// Sets o.R.Label to related.
// Adds o to related.R.LabelI18ns.
func (o *LabelI18n) SetLabel(exec boil.Executor, insert bool, related *Label) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"label_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"label_id"}),
		strmangle.WhereClause("\"", "\"", 2, labelI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.LabelID, o.Language}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.LabelID = related.ID
	if o.R == nil {
		o.R = &labelI18nR{
			Label: related,
		}
	} else {
		o.R.Label = related
	}

	if related.R == nil {
		related.R = &labelR{
			LabelI18ns: LabelI18nSlice{o},
		}
	} else {
		related.R.LabelI18ns = append(related.R.LabelI18ns, o)
	}

	return nil
}

// SetUser of the labelI18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.LabelI18ns.
func (o *LabelI18n) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"label_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, labelI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.LabelID, o.Language}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.UserID, related.ID)
	if o.R == nil {
		o.R = &labelI18nR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			LabelI18ns: LabelI18nSlice{o},
		}
	} else {
		related.R.LabelI18ns = append(related.R.LabelI18ns, o)
	}

	return nil
}

// RemoveUser relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *LabelI18n) RemoveUser(exec boil.Executor, related *User) error {
	var err error

	queries.SetScanner(&o.UserID, nil)
	if _, err = o.Update(exec, boil.Whitelist("user_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.User = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.LabelI18ns {
		if queries.Equal(o.UserID, ri.UserID) {
			continue
		}

		ln := len(related.R.LabelI18ns)
		if ln > 1 && i < ln-1 {
			related.R.LabelI18ns[i] = related.R.LabelI18ns[ln-1]
		}
		related.R.LabelI18ns = related.R.LabelI18ns[:ln-1]
		break
	}
	return nil
}

// LabelI18ns retrieves all the records using an executor.
func LabelI18ns(mods ...qm.QueryMod) labelI18nQuery {
	mods = append(mods, qm.From("\"label_i18n\""))
	return labelI18nQuery{NewQuery(mods...)}
}

// FindLabelI18n retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindLabelI18n(exec boil.Executor, labelID int64, language string, selectCols ...string) (*LabelI18n, error) {
	labelI18nObj := &LabelI18n{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"label_i18n\" where \"label_id\"=$1 AND \"language\"=$2", sel,
	)

	q := queries.Raw(query, labelID, language)

	err := q.Bind(nil, exec, labelI18nObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from label_i18n")
	}

	return labelI18nObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *LabelI18n) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no label_i18n provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(labelI18nColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	labelI18nInsertCacheMut.RLock()
	cache, cached := labelI18nInsertCache[key]
	labelI18nInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			labelI18nAllColumns,
			labelI18nColumnsWithDefault,
			labelI18nColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(labelI18nType, labelI18nMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(labelI18nType, labelI18nMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"label_i18n\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"label_i18n\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into label_i18n")
	}

	if !cached {
		labelI18nInsertCacheMut.Lock()
		labelI18nInsertCache[key] = cache
		labelI18nInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the LabelI18n.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *LabelI18n) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	labelI18nUpdateCacheMut.RLock()
	cache, cached := labelI18nUpdateCache[key]
	labelI18nUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			labelI18nAllColumns,
			labelI18nPrimaryKeyColumns,
		)
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update label_i18n, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"label_i18n\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, labelI18nPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(labelI18nType, labelI18nMapping, append(wl, labelI18nPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update label_i18n row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for label_i18n")
	}

	if !cached {
		labelI18nUpdateCacheMut.Lock()
		labelI18nUpdateCache[key] = cache
		labelI18nUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q labelI18nQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for label_i18n")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for label_i18n")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o LabelI18nSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), labelI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"label_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, labelI18nPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in labelI18n slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all labelI18n")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *LabelI18n) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no label_i18n provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(labelI18nColumnsWithDefault, o)

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

	labelI18nUpsertCacheMut.RLock()
	cache, cached := labelI18nUpsertCache[key]
	labelI18nUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			labelI18nAllColumns,
			labelI18nColumnsWithDefault,
			labelI18nColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			labelI18nAllColumns,
			labelI18nPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert label_i18n, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(labelI18nPrimaryKeyColumns))
			copy(conflict, labelI18nPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"label_i18n\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(labelI18nType, labelI18nMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(labelI18nType, labelI18nMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert label_i18n")
	}

	if !cached {
		labelI18nUpsertCacheMut.Lock()
		labelI18nUpsertCache[key] = cache
		labelI18nUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single LabelI18n record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *LabelI18n) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no LabelI18n provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), labelI18nPrimaryKeyMapping)
	sql := "DELETE FROM \"label_i18n\" WHERE \"label_id\"=$1 AND \"language\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from label_i18n")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for label_i18n")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q labelI18nQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no labelI18nQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from label_i18n")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for label_i18n")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o LabelI18nSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), labelI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"label_i18n\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, labelI18nPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from labelI18n slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for label_i18n")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *LabelI18n) Reload(exec boil.Executor) error {
	ret, err := FindLabelI18n(exec, o.LabelID, o.Language)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *LabelI18nSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := LabelI18nSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), labelI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"label_i18n\".* FROM \"label_i18n\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, labelI18nPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in LabelI18nSlice")
	}

	*o = slice

	return nil
}

// LabelI18nExists checks if the LabelI18n row exists.
func LabelI18nExists(exec boil.Executor, labelID int64, language string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"label_i18n\" where \"label_id\"=$1 AND \"language\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, labelID, language)
	}
	row := exec.QueryRow(sql, labelID, language)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if label_i18n exists")
	}

	return exists, nil
}
