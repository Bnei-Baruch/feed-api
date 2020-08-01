// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
)

// Blog is an object representing the database table.
type Blog struct {
	ID   int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`
	URL  string `boil:"url" json:"url" toml:"url" yaml:"url"`

	R *blogR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L blogL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var BlogColumns = struct {
	ID   string
	Name string
	URL  string
}{
	ID:   "id",
	Name: "name",
	URL:  "url",
}

// blogR is where relationships are stored.
type blogR struct {
	BlogPosts BlogPostSlice
}

// blogL is where Load methods for each relationship are stored.
type blogL struct{}

var (
	blogColumns               = []string{"id", "name", "url"}
	blogColumnsWithoutDefault = []string{"name", "url"}
	blogColumnsWithDefault    = []string{"id"}
	blogPrimaryKeyColumns     = []string{"id"}
)

type (
	// BlogSlice is an alias for a slice of pointers to Blog.
	// This should generally be used opposed to []Blog.
	BlogSlice []*Blog

	blogQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	blogType                 = reflect.TypeOf(&Blog{})
	blogMapping              = queries.MakeStructMapping(blogType)
	blogPrimaryKeyMapping, _ = queries.BindMapping(blogType, blogMapping, blogPrimaryKeyColumns)
	blogInsertCacheMut       sync.RWMutex
	blogInsertCache          = make(map[string]insertCache)
	blogUpdateCacheMut       sync.RWMutex
	blogUpdateCache          = make(map[string]updateCache)
	blogUpsertCacheMut       sync.RWMutex
	blogUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single blog record from the query, and panics on error.
func (q blogQuery) OneP() *Blog {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single blog record from the query.
func (q blogQuery) One() (*Blog, error) {
	o := &Blog{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mdbmdbmodels: failed to execute a one query for blogs")
	}

	return o, nil
}

// AllP returns all Blog records from the query, and panics on error.
func (q blogQuery) AllP() BlogSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Blog records from the query.
func (q blogQuery) All() (BlogSlice, error) {
	var o []*Blog

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "mdbmdbmodels: failed to assign all query results to Blog slice")
	}

	return o, nil
}

// CountP returns the count of all Blog records in the query, and panics on error.
func (q blogQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Blog records in the query.
func (q blogQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "mdbmdbmodels: failed to count blogs rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q blogQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q blogQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "mdbmdbmodels: failed to check if blogs exists")
	}

	return count > 0, nil
}

// BlogPostsG retrieves all the blog_post's blog posts.
func (o *Blog) BlogPostsG(mods ...qm.QueryMod) blogPostQuery {
	return o.BlogPosts(boil.GetDB(), mods...)
}

// BlogPosts retrieves all the blog_post's blog posts with an executor.
func (o *Blog) BlogPosts(exec boil.Executor, mods ...qm.QueryMod) blogPostQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"blog_posts\".\"blog_id\"=?", o.ID),
	)

	query := BlogPosts(exec, queryMods...)
	queries.SetFrom(query.Query, "\"blog_posts\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"blog_posts\".*"})
	}

	return query
}

// LoadBlogPosts allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (blogL) LoadBlogPosts(e boil.Executor, singular bool, maybeBlog interface{}) error {
	var slice []*Blog
	var object *Blog

	count := 1
	if singular {
		object = maybeBlog.(*Blog)
	} else {
		slice = *maybeBlog.(*[]*Blog)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &blogR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &blogR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"blog_posts\" where \"blog_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load blog_posts")
	}
	defer results.Close()

	var resultSlice []*BlogPost
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice blog_posts")
	}

	if singular {
		object.R.BlogPosts = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.BlogID {
				local.R.BlogPosts = append(local.R.BlogPosts, foreign)
				break
			}
		}
	}

	return nil
}

// AddBlogPostsG adds the given related objects to the existing relationships
// of the blog, optionally inserting them as new records.
// Appends related to o.R.BlogPosts.
// Sets related.R.Blog appropriately.
// Uses the global database handle.
func (o *Blog) AddBlogPostsG(insert bool, related ...*BlogPost) error {
	return o.AddBlogPosts(boil.GetDB(), insert, related...)
}

// AddBlogPostsP adds the given related objects to the existing relationships
// of the blog, optionally inserting them as new records.
// Appends related to o.R.BlogPosts.
// Sets related.R.Blog appropriately.
// Panics on error.
func (o *Blog) AddBlogPostsP(exec boil.Executor, insert bool, related ...*BlogPost) {
	if err := o.AddBlogPosts(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddBlogPostsGP adds the given related objects to the existing relationships
// of the blog, optionally inserting them as new records.
// Appends related to o.R.BlogPosts.
// Sets related.R.Blog appropriately.
// Uses the global database handle and panics on error.
func (o *Blog) AddBlogPostsGP(insert bool, related ...*BlogPost) {
	if err := o.AddBlogPosts(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddBlogPosts adds the given related objects to the existing relationships
// of the blog, optionally inserting them as new records.
// Appends related to o.R.BlogPosts.
// Sets related.R.Blog appropriately.
func (o *Blog) AddBlogPosts(exec boil.Executor, insert bool, related ...*BlogPost) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.BlogID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"blog_posts\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"blog_id"}),
				strmangle.WhereClause("\"", "\"", 2, blogPostPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.BlogID = o.ID
		}
	}

	if o.R == nil {
		o.R = &blogR{
			BlogPosts: related,
		}
	} else {
		o.R.BlogPosts = append(o.R.BlogPosts, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &blogPostR{
				Blog: o,
			}
		} else {
			rel.R.Blog = o
		}
	}
	return nil
}

// BlogsG retrieves all records.
func BlogsG(mods ...qm.QueryMod) blogQuery {
	return Blogs(boil.GetDB(), mods...)
}

// Blogs retrieves all the records using an executor.
func Blogs(exec boil.Executor, mods ...qm.QueryMod) blogQuery {
	mods = append(mods, qm.From("\"blogs\""))
	return blogQuery{NewQuery(exec, mods...)}
}

// FindBlogG retrieves a single record by ID.
func FindBlogG(id int64, selectCols ...string) (*Blog, error) {
	return FindBlog(boil.GetDB(), id, selectCols...)
}

// FindBlogGP retrieves a single record by ID, and panics on error.
func FindBlogGP(id int64, selectCols ...string) *Blog {
	retobj, err := FindBlog(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindBlog retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindBlog(exec boil.Executor, id int64, selectCols ...string) (*Blog, error) {
	blogObj := &Blog{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"blogs\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(blogObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "mdbmdbmodels: unable to select from blogs")
	}

	return blogObj, nil
}

// FindBlogP retrieves a single record by ID with an executor, and panics on error.
func FindBlogP(exec boil.Executor, id int64, selectCols ...string) *Blog {
	retobj, err := FindBlog(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Blog) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Blog) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Blog) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Blog) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("mdbmdbmodels: no blogs provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(blogColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	blogInsertCacheMut.RLock()
	cache, cached := blogInsertCache[key]
	blogInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			blogColumns,
			blogColumnsWithDefault,
			blogColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(blogType, blogMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(blogType, blogMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"blogs\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"blogs\" DEFAULT VALUES"
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
		return errors.Wrap(err, "mdbmdbmodels: unable to insert into blogs")
	}

	if !cached {
		blogInsertCacheMut.Lock()
		blogInsertCache[key] = cache
		blogInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single Blog record. See Update for
// whitelist behavior description.
func (o *Blog) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Blog record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Blog) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Blog, and panics on error.
// See Update for whitelist behavior description.
func (o *Blog) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Blog.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Blog) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	blogUpdateCacheMut.RLock()
	cache, cached := blogUpdateCache[key]
	blogUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			blogColumns,
			blogPrimaryKeyColumns,
			whitelist,
		)

		if len(wl) == 0 {
			return errors.New("mdbmdbmodels: unable to update blogs, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"blogs\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, blogPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(blogType, blogMapping, append(wl, blogPrimaryKeyColumns...))
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
		return errors.Wrap(err, "mdbmdbmodels: unable to update blogs row")
	}

	if !cached {
		blogUpdateCacheMut.Lock()
		blogUpdateCache[key] = cache
		blogUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q blogQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q blogQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to update all for blogs")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o BlogSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o BlogSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o BlogSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o BlogSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), blogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"blogs\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, blogPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to update all in blog slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Blog) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Blog) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Blog) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Blog) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("mdbmdbmodels: no blogs provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(blogColumnsWithDefault, o)

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

	blogUpsertCacheMut.RLock()
	cache, cached := blogUpsertCache[key]
	blogUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertColumnSet(
			blogColumns,
			blogColumnsWithDefault,
			blogColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		update := strmangle.UpdateColumnSet(
			blogColumns,
			blogPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("mdbmdbmodels: unable to upsert blogs, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(blogPrimaryKeyColumns))
			copy(conflict, blogPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"blogs\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(blogType, blogMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(blogType, blogMapping, ret)
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
		return errors.Wrap(err, "mdbmdbmodels: unable to upsert blogs")
	}

	if !cached {
		blogUpsertCacheMut.Lock()
		blogUpsertCache[key] = cache
		blogUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single Blog record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Blog) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Blog record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Blog) DeleteG() error {
	if o == nil {
		return errors.New("mdbmdbmodels: no Blog provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Blog record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Blog) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Blog record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Blog) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("mdbmdbmodels: no Blog provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), blogPrimaryKeyMapping)
	sql := "DELETE FROM \"blogs\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to delete from blogs")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q blogQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q blogQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("mdbmdbmodels: no blogQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to delete all from blogs")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o BlogSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o BlogSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("mdbmdbmodels: no Blog slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o BlogSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o BlogSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("mdbmdbmodels: no Blog slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), blogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"blogs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, blogPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to delete all from blog slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Blog) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Blog) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Blog) ReloadG() error {
	if o == nil {
		return errors.New("mdbmdbmodels: no Blog provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Blog) Reload(exec boil.Executor) error {
	ret, err := FindBlog(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *BlogSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *BlogSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BlogSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("mdbmdbmodels: empty BlogSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BlogSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	blogs := BlogSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), blogPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"blogs\".* FROM \"blogs\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, blogPrimaryKeyColumns, len(*o))

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&blogs)
	if err != nil {
		return errors.Wrap(err, "mdbmdbmodels: unable to reload all in BlogSlice")
	}

	*o = blogs

	return nil
}

// BlogExists checks if the Blog row exists.
func BlogExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"blogs\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "mdbmdbmodels: unable to check if blogs exists")
	}

	return exists, nil
}

// BlogExistsG checks if the Blog row exists.
func BlogExistsG(id int64) (bool, error) {
	return BlogExists(boil.GetDB(), id)
}

// BlogExistsGP checks if the Blog row exists. Panics on error.
func BlogExistsGP(id int64) bool {
	e, err := BlogExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// BlogExistsP checks if the Blog row exists. Panics on error.
func BlogExistsP(exec boil.Executor, id int64) bool {
	e, err := BlogExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
