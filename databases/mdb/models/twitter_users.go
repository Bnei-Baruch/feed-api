// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// TwitterUser is an object representing the database table.
type TwitterUser struct {
	ID          int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	Username    string `boil:"username" json:"username" toml:"username" yaml:"username"`
	AccountID   string `boil:"account_id" json:"account_id" toml:"account_id" yaml:"account_id"`
	DisplayName string `boil:"display_name" json:"display_name" toml:"display_name" yaml:"display_name"`

	R *twitterUserR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L twitterUserL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var TwitterUserColumns = struct {
	ID          string
	Username    string
	AccountID   string
	DisplayName string
}{
	ID:          "id",
	Username:    "username",
	AccountID:   "account_id",
	DisplayName: "display_name",
}

// Generated where

var TwitterUserWhere = struct {
	ID          whereHelperint64
	Username    whereHelperstring
	AccountID   whereHelperstring
	DisplayName whereHelperstring
}{
	ID:          whereHelperint64{field: "\"twitter_users\".\"id\""},
	Username:    whereHelperstring{field: "\"twitter_users\".\"username\""},
	AccountID:   whereHelperstring{field: "\"twitter_users\".\"account_id\""},
	DisplayName: whereHelperstring{field: "\"twitter_users\".\"display_name\""},
}

// TwitterUserRels is where relationship names are stored.
var TwitterUserRels = struct {
	UserTwitterTweets string
}{
	UserTwitterTweets: "UserTwitterTweets",
}

// twitterUserR is where relationships are stored.
type twitterUserR struct {
	UserTwitterTweets TwitterTweetSlice
}

// NewStruct creates a new relationship struct
func (*twitterUserR) NewStruct() *twitterUserR {
	return &twitterUserR{}
}

// twitterUserL is where Load methods for each relationship are stored.
type twitterUserL struct{}

var (
	twitterUserAllColumns            = []string{"id", "username", "account_id", "display_name"}
	twitterUserColumnsWithoutDefault = []string{"username", "account_id", "display_name"}
	twitterUserColumnsWithDefault    = []string{"id"}
	twitterUserPrimaryKeyColumns     = []string{"id"}
)

type (
	// TwitterUserSlice is an alias for a slice of pointers to TwitterUser.
	// This should generally be used opposed to []TwitterUser.
	TwitterUserSlice []*TwitterUser

	twitterUserQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	twitterUserType                 = reflect.TypeOf(&TwitterUser{})
	twitterUserMapping              = queries.MakeStructMapping(twitterUserType)
	twitterUserPrimaryKeyMapping, _ = queries.BindMapping(twitterUserType, twitterUserMapping, twitterUserPrimaryKeyColumns)
	twitterUserInsertCacheMut       sync.RWMutex
	twitterUserInsertCache          = make(map[string]insertCache)
	twitterUserUpdateCacheMut       sync.RWMutex
	twitterUserUpdateCache          = make(map[string]updateCache)
	twitterUserUpsertCacheMut       sync.RWMutex
	twitterUserUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single twitterUser record from the query.
func (q twitterUserQuery) One(exec boil.Executor) (*TwitterUser, error) {
	o := &TwitterUser{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for twitter_users")
	}

	return o, nil
}

// All returns all TwitterUser records from the query.
func (q twitterUserQuery) All(exec boil.Executor) (TwitterUserSlice, error) {
	var o []*TwitterUser

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to TwitterUser slice")
	}

	return o, nil
}

// Count returns the count of all TwitterUser records in the query.
func (q twitterUserQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count twitter_users rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q twitterUserQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if twitter_users exists")
	}

	return count > 0, nil
}

// UserTwitterTweets retrieves all the twitter_tweet's TwitterTweets with an executor via user_id column.
func (o *TwitterUser) UserTwitterTweets(mods ...qm.QueryMod) twitterTweetQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"twitter_tweets\".\"user_id\"=?", o.ID),
	)

	query := TwitterTweets(queryMods...)
	queries.SetFrom(query.Query, "\"twitter_tweets\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"twitter_tweets\".*"})
	}

	return query
}

// LoadUserTwitterTweets allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (twitterUserL) LoadUserTwitterTweets(e boil.Executor, singular bool, maybeTwitterUser interface{}, mods queries.Applicator) error {
	var slice []*TwitterUser
	var object *TwitterUser

	if singular {
		object = maybeTwitterUser.(*TwitterUser)
	} else {
		slice = *maybeTwitterUser.(*[]*TwitterUser)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &twitterUserR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &twitterUserR{}
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

	query := NewQuery(qm.From(`twitter_tweets`), qm.WhereIn(`twitter_tweets.user_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load twitter_tweets")
	}

	var resultSlice []*TwitterTweet
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice twitter_tweets")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on twitter_tweets")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for twitter_tweets")
	}

	if singular {
		object.R.UserTwitterTweets = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &twitterTweetR{}
			}
			foreign.R.User = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.UserID {
				local.R.UserTwitterTweets = append(local.R.UserTwitterTweets, foreign)
				if foreign.R == nil {
					foreign.R = &twitterTweetR{}
				}
				foreign.R.User = local
				break
			}
		}
	}

	return nil
}

// AddUserTwitterTweets adds the given related objects to the existing relationships
// of the twitter_user, optionally inserting them as new records.
// Appends related to o.R.UserTwitterTweets.
// Sets related.R.User appropriately.
func (o *TwitterUser) AddUserTwitterTweets(exec boil.Executor, insert bool, related ...*TwitterTweet) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.UserID = o.ID
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"twitter_tweets\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
				strmangle.WhereClause("\"", "\"", 2, twitterTweetPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}
			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.UserID = o.ID
		}
	}

	if o.R == nil {
		o.R = &twitterUserR{
			UserTwitterTweets: related,
		}
	} else {
		o.R.UserTwitterTweets = append(o.R.UserTwitterTweets, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &twitterTweetR{
				User: o,
			}
		} else {
			rel.R.User = o
		}
	}
	return nil
}

// TwitterUsers retrieves all the records using an executor.
func TwitterUsers(mods ...qm.QueryMod) twitterUserQuery {
	mods = append(mods, qm.From("\"twitter_users\""))
	return twitterUserQuery{NewQuery(mods...)}
}

// FindTwitterUser retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindTwitterUser(exec boil.Executor, iD int64, selectCols ...string) (*TwitterUser, error) {
	twitterUserObj := &TwitterUser{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"twitter_users\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, twitterUserObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from twitter_users")
	}

	return twitterUserObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *TwitterUser) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no twitter_users provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(twitterUserColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	twitterUserInsertCacheMut.RLock()
	cache, cached := twitterUserInsertCache[key]
	twitterUserInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			twitterUserAllColumns,
			twitterUserColumnsWithDefault,
			twitterUserColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(twitterUserType, twitterUserMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(twitterUserType, twitterUserMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"twitter_users\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"twitter_users\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into twitter_users")
	}

	if !cached {
		twitterUserInsertCacheMut.Lock()
		twitterUserInsertCache[key] = cache
		twitterUserInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the TwitterUser.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *TwitterUser) Update(exec boil.Executor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	twitterUserUpdateCacheMut.RLock()
	cache, cached := twitterUserUpdateCache[key]
	twitterUserUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			twitterUserAllColumns,
			twitterUserPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update twitter_users, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"twitter_users\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, twitterUserPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(twitterUserType, twitterUserMapping, append(wl, twitterUserPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update twitter_users row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for twitter_users")
	}

	if !cached {
		twitterUserUpdateCacheMut.Lock()
		twitterUserUpdateCache[key] = cache
		twitterUserUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q twitterUserQuery) UpdateAll(exec boil.Executor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for twitter_users")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for twitter_users")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o TwitterUserSlice) UpdateAll(exec boil.Executor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), twitterUserPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"twitter_users\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, twitterUserPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in twitterUser slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all twitterUser")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *TwitterUser) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no twitter_users provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(twitterUserColumnsWithDefault, o)

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

	twitterUserUpsertCacheMut.RLock()
	cache, cached := twitterUserUpsertCache[key]
	twitterUserUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			twitterUserAllColumns,
			twitterUserColumnsWithDefault,
			twitterUserColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			twitterUserAllColumns,
			twitterUserPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert twitter_users, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(twitterUserPrimaryKeyColumns))
			copy(conflict, twitterUserPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"twitter_users\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(twitterUserType, twitterUserMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(twitterUserType, twitterUserMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert twitter_users")
	}

	if !cached {
		twitterUserUpsertCacheMut.Lock()
		twitterUserUpsertCache[key] = cache
		twitterUserUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single TwitterUser record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *TwitterUser) Delete(exec boil.Executor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no TwitterUser provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), twitterUserPrimaryKeyMapping)
	sql := "DELETE FROM \"twitter_users\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from twitter_users")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for twitter_users")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q twitterUserQuery) DeleteAll(exec boil.Executor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no twitterUserQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.Exec(exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from twitter_users")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for twitter_users")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TwitterUserSlice) DeleteAll(exec boil.Executor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), twitterUserPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"twitter_users\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, twitterUserPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	result, err := exec.Exec(sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from twitterUser slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for twitter_users")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *TwitterUser) Reload(exec boil.Executor) error {
	ret, err := FindTwitterUser(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TwitterUserSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := TwitterUserSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), twitterUserPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"twitter_users\".* FROM \"twitter_users\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, twitterUserPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in TwitterUserSlice")
	}

	*o = slice

	return nil
}

// TwitterUserExists checks if the TwitterUser row exists.
func TwitterUserExists(exec boil.Executor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"twitter_users\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if twitter_users exists")
	}

	return exists, nil
}
