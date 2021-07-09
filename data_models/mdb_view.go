package data_models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/feed-api/databases/mdb/models"
	"github.com/Bnei-Baruch/feed-api/events"
	"github.com/Bnei-Baruch/feed-api/utils"
)

const ()

type MdbView struct {
	remote   *sql.DB
	local    *sql.DB
	name     string
	interval time.Duration
	tables   []TableInfo
}

func MakeMdbView(local *sql.DB, remote *sql.DB) *MdbView {
	tables := createTablesInfo()
	if !events.DebugMode {
		utils.Must(syncLocalMdb(tables, local, remote))
	}

	return &MdbView{
		remote,
		local,
		"MdbView",
		time.Duration(time.Minute),
		tables,
	}
}

func (m *MdbView) Name() string {
	return m.name
}

func (m *MdbView) Interval() time.Duration {
	return m.interval
}

type IdsTuple struct {
	Id1      int64  `boil:"id1"`
	Id2      int64  `boil:"id2"`
	Language string `boil:"language"`
}

type BoilRow interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
	Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error)
}

func ToFileSlice(s interface{}) []*models.File {
	slice := reflect.ValueOf(s).Elem()
	if slice.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Expected slice! got: %+v %+v", slice.Kind(), slice))
	}
	c := slice.Len()
	out := make([]*models.File, c)
	for i := 0; i < c; i++ {
		//log.Infof("%d: %+v", i, slice.Index(i))
		out[i] = slice.Index(i).Interface().(*models.File)
	}
	return out
}

func ToBoilRowSlice(s interface{}) []BoilRow {
	slice := reflect.ValueOf(s).Elem()
	if slice.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Expected slice! got: %+v %+v", slice.Kind(), slice))
	}
	c := slice.Len()
	out := make([]BoilRow, c)
	for i := 0; i < c; i++ {
		//log.Infof("%d: %+v", i, slice.Index(i))
		out[i] = slice.Index(i).Interface().(BoilRow)
	}
	return out
}

type Bindable interface {
	Bind(ctx context.Context, exec boil.Executor, obj interface{}) error
}

type FuncToInterface func() interface{}
type BoilQueryFunc func(mods ...qm.QueryMod) Bindable

type TableInfo struct {
	Name       string
	KeyColumns []string
	Units      FuncToInterface
	Query      BoilQueryFunc
}

func KeyColumns(info TableInfo) []string {
	columns := []string(nil)
	for i, column := range info.KeyColumns {
		if column == "language" {
			columns = append(columns, column)
		} else {
			columns = append(columns, fmt.Sprintf("%s as id%d", column, i+1))
		}
	}
	return columns
}

func fetchIds(info TableInfo, exec boil.ContextExecutor) ([]IdsTuple, error) {
	ids := []IdsTuple(nil)
	if err := info.Query(qm.Select(KeyColumns(info)...)).Bind(context.TODO(), exec, &ids); err != nil {
		log.Infof("Failed binding ids: %+v", info)
		return []IdsTuple(nil), err
	}
	return ids, nil
}

func whereIds(info TableInfo, ids []IdsTuple) string {
	values := []string(nil)
	for _, id := range ids {
		clauses := []string(nil)
		if id.Id1 != 0 {
			clauses = append(clauses, fmt.Sprintf("%s = %d", info.KeyColumns[0], id.Id1))
		}
		if id.Id2 != 0 {
			clauses = append(clauses, fmt.Sprintf("%s = %d", info.KeyColumns[1], id.Id2))
		}
		if id.Language != "" {
			clauses = append(clauses, fmt.Sprintf("%s = '%s'", info.KeyColumns[1], id.Language))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(clauses, " and ")))
	}
	return strings.Join(values, " or ")
}

func deleteRows(info TableInfo, ids []IdsTuple, exec boil.ContextExecutor) (int64, error) {
	if result, err := queries.Raw(fmt.Sprintf("delete from %s where %s", info.Name, whereIds(info, ids))).Exec(exec); err != nil {
		log.Infof("Error deleting from %s, ids: %s", info.Name, IdsToString(ids))
		return 0, err
	} else {
		if affected, err := result.RowsAffected(); err != nil {
			return 0, err
		} else {
			return affected, nil
		}
	}
}

type move struct {
	From int64
	To   int64
}

func appendF(f *models.File, filesMap map[int64]*models.File, sortedFiles *[]*models.File) {
	if f.ParentID.Valid {
		if _, found := filesMap[f.ParentID.Int64]; found {
			appendF(filesMap[f.ParentID.Int64], filesMap, sortedFiles)
		}
	}
	delete(filesMap, f.ID)
	*sortedFiles = append(*sortedFiles, f)
}

func fetchRows(info TableInfo, ids []IdsTuple, exec boil.ContextExecutor) ([]BoilRow, error) {
	slice := info.Units()
	err := info.Query(qm.Select("*"), qm.Where(whereIds(info, ids))).Bind(context.TODO(), exec, slice)
	if err != nil {
		// log.Infof("Failed binding rows: %+v Err: %+v", info, err)
		return nil, err
	}
	if info.Name != T_FILES {
		return ToBoilRowSlice(slice), nil
	} else {
		// Order rows by parent first, so insert will work properly.
		files := ToFileSlice(slice)
		filesMap := make(map[int64]*models.File)
		for _, f := range files {
			filesMap[f.ID] = f
		}
		sortedFiles := []*models.File(nil)
		for _, f := range filesMap {
			appendF(f, filesMap, &sortedFiles)
		}
		return ToBoilRowSlice(&sortedFiles), nil
	}
}

func CompareIdsTuple(a, b *IdsTuple) int {
	if a.Id1 == b.Id1 {
		if a.Language != "" && b.Language != "" {
			return strings.Compare(a.Language, b.Language)
		} else {
			if a.Id2 < b.Id2 {
				return -1
			} else if a.Id2 > b.Id2 {
				return 1
			}
			return 0
		}
	}
	if a.Id1 < b.Id1 {
		return -1
	}
	return 1
}

func IdsTupleSliceLess(s []IdsTuple) func(i, j int) bool {
	return func(i, j int) bool {
		return CompareIdsTuple(&s[i], &s[j]) == -1
	}
}

// Calculates difference a - b. Slices assumed to be sorted lower to higher.
func diffIdsTupleSlices(a []IdsTuple, b []IdsTuple) []IdsTuple {
	index := 0
	diff := []IdsTuple(nil)
	for _, n := range a {
		for index < len(b) && CompareIdsTuple(&n, &b[index]) == 1 {
			index++
		}
		if index >= len(b) || CompareIdsTuple(&n, &b[index]) == -1 {
			diff = append(diff, n)
		}
	}
	return diff
}

const (
	CHUNK_SIZE = 5000
)

// Syncs one table from remote to local.
// localIds - to be removed.
// remoteIds - to be added after local removed.
// localIds and removeIds should overlap most of the times.
func DeleteFromTable(info TableInfo, ids []IdsTuple, local *sql.DB) (int64, error) {
	start := time.Now()
	defer func() {
		utils.Profile("DeleteFromTable", time.Now().Sub(start))
	}()
	if affected, err := deleteRows(info, ids, local); err != nil {
		return 0, err
	} else {
		return affected, nil
	}
}

func InsertToTable(info TableInfo, ids []IdsTuple, local *sql.DB, remote *sql.DB) (int64, error) {
	start := time.Now()
	defer func() {
		utils.Profile("InsertToTable", time.Now().Sub(start))
	}()
	inserted := int64(0)
	for i := 0; i < len(ids); i += CHUNK_SIZE {
		chunk := ids[i:utils.MinInt(len(ids), i+CHUNK_SIZE)]
		if rows, err := fetchRows(info, chunk, remote); err != nil {
			return 0, err
		} else {
			log.Infof("Fetched %d rows. Chunk %d of total rows %d.", len(rows), i, len(ids))
			chunkInserted := int64(0)
			for _, row := range rows {
				if err := row.Insert(context.TODO(), local, boil.Infer()); err != nil {
					log.Warnf("Unable to insert, skipping: %+v\nRow: %+v", err, row)
				} else {
					chunkInserted++
					// log.Infof("Inserted")
				}
			}
			inserted += chunkInserted
		}
	}
	return inserted, nil
}

func RelationQuery(table, field1, field2 string) BoilQueryFunc {
	return func(mods ...qm.QueryMod) Bindable {
		mods = append([]qm.QueryMod{qm.Select(fmt.Sprintf("%s as id1", field1), fmt.Sprintf("%s as id2", field2)), qm.From(table)}, mods...)
		return models.NewQuery(mods...)
	}
}

type ContentUnitTag struct {
	ContentUnitID int64 `boil:"content_unit_id"`
	TagID         int64 `boil:"tag_id"`
}

func (cut *ContentUnitTag) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	query := fmt.Sprintf("INSERT INTO \"content_units_tags\" (\"content_unit_id\", \"tag_id\") VALUES (%d, %d)", cut.ContentUnitID, cut.TagID)
	if _, err := exec.ExecContext(ctx, query); err != nil {
		return errors.Wrap(err, "models: unable to insert into content_units_tags")
	}
	return nil
}

func (cut *ContentUnitTag) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if result, err := exec.ExecContext(ctx, "DELETE FROM \"content_units_tags\" WHERE \"content_unit_id\"=$1 AND \"tag_id\"=$2", cut.ContentUnitID, cut.TagID); err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from content_units_tags")
	} else {
		if rowsAff, err := result.RowsAffected(); err != nil {
			return 0, errors.Wrap(err, "models: failed to get rows affected by delete for content_units_tags")
		} else {
			return rowsAff, nil
		}
	}
}

type ContentUnitSource struct {
	ContentUnitID int64 `boil:"content_unit_id"`
	SourceID      int64 `boil:"source_id"`
}

func (cut *ContentUnitSource) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	query := fmt.Sprintf("INSERT INTO \"content_units_sources\" (\"content_unit_id\", \"source_id\") VALUES (%d, %d)", cut.ContentUnitID, cut.SourceID)
	if _, err := exec.ExecContext(ctx, query); err != nil {
		return errors.Wrap(err, "models: unable to insert into content_units_sources")
	}
	return nil
}

func (cut *ContentUnitSource) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if result, err := exec.ExecContext(ctx, "DELETE FROM \"content_units_sources\" WHERE \"content_unit_id\"=$1 AND \"source_id\"=$2", cut.ContentUnitID, cut.SourceID); err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from content_units_sources")
	} else {
		if rowsAff, err := result.RowsAffected(); err != nil {
			return 0, errors.Wrap(err, "models: failed to get rows affected by delete for content_units_sources")
		} else {
			return rowsAff, nil
		}
	}
}

type ContentUnitPublisher struct {
	ContentUnitID int64 `boil:"content_unit_id"`
	PublisherID   int64 `boil:"publisher_id"`
}

func (cut *ContentUnitPublisher) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	query := fmt.Sprintf("INSERT INTO \"content_units_publishers\" (\"content_unit_id\", \"publisher_id\") VALUES (%d, %d)", cut.ContentUnitID, cut.PublisherID)
	if _, err := exec.ExecContext(ctx, query); err != nil {
		return errors.Wrap(err, "models: unable to insert into content_units_publishers")
	}
	return nil
}

func (cut *ContentUnitPublisher) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if result, err := exec.ExecContext(ctx, "DELETE FROM \"content_units_publishers\" WHERE \"content_unit_id\"=$1 AND \"publisher_id\"=$2", cut.ContentUnitID, cut.PublisherID); err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from content_units_publishers")
	} else {
		if rowsAff, err := result.RowsAffected(); err != nil {
			return 0, errors.Wrap(err, "models: failed to get rows affected by delete for content_units_publishers")
		} else {
			return rowsAff, nil
		}
	}
}

type FileOperation struct {
	FileID      int64 `boil:"file_id"`
	OperationID int64 `boil:"operation_id"`
}

func (cut *FileOperation) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	query := fmt.Sprintf("INSERT INTO \"files_operations\" (\"file_id\", \"operation_id\") VALUES (%d, %d)", cut.FileID, cut.OperationID)
	if _, err := exec.ExecContext(ctx, query); err != nil {
		return errors.Wrap(err, "models: unable to insert into files_operations")
	}
	return nil
}

func (cut *FileOperation) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if result, err := exec.ExecContext(ctx, "DELETE FROM \"files_operations\" WHERE \"file_id\"=$1 AND \"operation_id\"=$2", cut.FileID, cut.OperationID); err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from files_operations")
	} else {
		if rowsAff, err := result.RowsAffected(); err != nil {
			return 0, errors.Wrap(err, "models: failed to get rows affected by delete for files_operations")
		} else {
			return rowsAff, nil
		}
	}
}

type FileStorage struct {
	FileID    int64 `boil:"file_id"`
	StorageID int64 `boil:"storage_id"`
}

func (cut *FileStorage) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	query := fmt.Sprintf("INSERT INTO \"files_storages\" (\"file_id\", \"storage_id\") VALUES (%d, %d)", cut.FileID, cut.StorageID)
	if _, err := exec.ExecContext(ctx, query); err != nil {
		return errors.Wrap(err, "models: unable to insert into files_storages")
	}
	return nil
}

func (cut *FileStorage) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if result, err := exec.ExecContext(ctx, "DELETE FROM \"files_storages\" WHERE \"file_id\"=$1 AND \"storage_id\"=$2", cut.FileID, cut.StorageID); err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from files_storages")
	} else {
		if rowsAff, err := result.RowsAffected(); err != nil {
			return 0, errors.Wrap(err, "models: failed to get rows affected by delete for files_storages")
		} else {
			return rowsAff, nil
		}
	}
}

type AuthorSource struct {
	AuthorID int64 `boil:"author_id"`
	SourceID int64 `boil:"source_id"`
}

func (cut *AuthorSource) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	query := fmt.Sprintf("INSERT INTO \"authors_sources\" (\"author_id\", \"source_id\") VALUES (%d, %d)", cut.AuthorID, cut.SourceID)
	if _, err := exec.ExecContext(ctx, query); err != nil {
		return errors.Wrap(err, "models: unable to insert into authors_sources")
	}
	return nil
}

func (cut *AuthorSource) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if result, err := exec.ExecContext(ctx, "DELETE FROM \"authors_sources\" WHERE \"author_id\"=$1 AND \"source_id\"=$2", cut.AuthorID, cut.SourceID); err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from authors_sources")
	} else {
		if rowsAff, err := result.RowsAffected(); err != nil {
			return 0, errors.Wrap(err, "models: failed to get rows affected by delete for authors_sources")
		} else {
			return rowsAff, nil
		}
	}
}

const (
	T_AUTHOR_I18N               = "author_i18n"
	T_AUTHORS                   = "authors"
	T_AUTHORS_SOURCES           = "authors_sources"
	T_BLOG_POSTS                = "blog_posts"
	T_BLOGS                     = "blogs"
	T_COLLECTION_I18N           = "collection_i18n"
	T_COLLECTIONS               = "collections"
	T_COLLECTIONS_CONTENT_UNITS = "collections_content_units"
	T_CONTENT_ROLE_TYPES        = "content_role_types"
	T_CONTENT_TYPES             = "content_types"
	T_CONTENT_UNIT_DERIVATIONS  = "content_unit_derivations"
	T_CONTENT_UNIT_I18N         = "content_unit_i18n"
	T_CONTENT_UNITS             = "content_units"
	T_CONTENT_UNITS_PERSONS     = "content_units_persons"
	T_CONTENT_UNITS_PUBLISHERS  = "content_units_publishers"
	T_CONTENT_UNITS_SOURCES     = "content_units_sources"
	T_CONTENT_UNITS_TAGS        = "content_units_tags"
	T_FILES                     = "files"
	T_FILES_OPERATIONS          = "files_operations"
	T_FILES_STORAGES            = "files_storages"
	T_MIGRATIONS                = "migrations"
	T_OPERATION_TYPES           = "operation_types"
	T_OPERATIONS                = "operations"
	T_PERSON_I18N               = "person_i18n"
	T_PERSONS                   = "persons"
	T_PUBLISHER_I18N            = "publisher_i18n"
	T_PUBLISHERS                = "publishers"
	T_SOURCE_I18N               = "source_i18n"
	T_SOURCE_TYPES              = "source_types"
	T_SOURCES                   = "sources"
	T_STORAGES                  = "storages"
	T_TAG_I18N                  = "tag_i18n"
	T_TAGS                      = "tags"
	T_TWITTER_TWEETS            = "twitter_tweets"
	T_TWITTER_USERS             = "twitter_users"
	T_USERS                     = "users"

	E_COLLECTION_CREATE               = "COLLECTION_CREATE"
	E_COLLECTION_UPDATE               = "COLLECTION_UPDATE"
	E_COLLECTION_DELETE               = "COLLECTION_DELETE"
	E_COLLECTION_PUBLISHED_CHANGE     = "COLLECTION_PUBLISHED_CHANGE"
	E_COLLECTION_CONTENT_UNITS_CHANGE = "COLLECTION_CONTENT_UNITS_CHANGE"

	E_CONTENT_UNIT_CREATE             = "CONTENT_UNIT_CREATE"
	E_CONTENT_UNIT_UPDATE             = "CONTENT_UNIT_UPDATE"
	E_CONTENT_UNIT_DELETE             = "CONTENT_UNIT_DELETE"
	E_CONTENT_UNIT_PUBLISHED_CHANGE   = "CONTENT_UNIT_PUBLISHED_CHANGE"
	E_CONTENT_UNIT_DERIVATIVES_CHANGE = "CONTENT_UNIT_DERIVATIVES_CHANGE"
	E_CONTENT_UNIT_SOURCES_CHANGE     = "CONTENT_UNIT_SOURCES_CHANGE"
	E_CONTENT_UNIT_TAGS_CHANGE        = "CONTENT_UNIT_TAGS_CHANGE"
	E_CONTENT_UNIT_PERSONS_CHANGE     = "CONTENT_UNIT_PERSONS_CHANGE"
	E_CONTENT_UNIT_PUBLISHERS_CHANGE  = "CONTENT_UNIT_PUBLISHERS_CHANGE"

	E_FILE_UPDATE    = "FILE_UPDATE"
	E_FILE_PUBLISHED = "FILE_PUBLISHED"
	E_FILE_INSERT    = "FILE_INSERT"
	E_FILE_REPLACE   = "FILE_REPLACE"
	E_FILE_REMOVE    = "FILE_REMOVE"

	E_SOURCE_CREATE = "SOURCE_CREATE"
	E_SOURCE_UPDATE = "SOURCE_UPDATE"

	E_TAG_CREATE = "TAG_CREATE"
	E_TAG_UPDATE = "TAG_UPDATE"

	E_PERSON_CREATE = "PERSON_CREATE"
	E_PERSON_UPDATE = "PERSON_UPDATE"
	E_PERSON_DELETE = "PERSON_DELETE"

	E_PUBLISHER_CREATE = "PUBLISHER_CREATE"
	E_PUBLISHER_UPDATE = "PUBLISHER_UPDATE"

	E_BLOG_POST_CREATE = "BLOG_POST_CREATE"
	E_BLOG_POST_UPDATE = "BLOG_POST_UPDATE"
	E_BLOG_POST_DELETE = "BLOG_POST_DELETE"

	E_TWEET_CREATE = "TWEET_CREATE"
	E_TWEET_UPDATE = "TWEET_UPDATE"
	E_TWEET_DELETE = "TWEET_DELETE"
)

func createTablesInfo() []TableInfo {
	return []TableInfo{
		// Types
		TableInfo{T_CONTENT_TYPES, []string{"id"}, func() interface{} { return &[]*models.ContentType{} }, func(mods ...qm.QueryMod) Bindable { return models.ContentTypes(mods...) }},
		TableInfo{T_CONTENT_ROLE_TYPES, []string{"id"}, func() interface{} { return &[]*models.ContentRoleType{} }, func(mods ...qm.QueryMod) Bindable { return models.ContentRoleTypes(mods...) }},
		TableInfo{T_OPERATION_TYPES, []string{"id"}, func() interface{} { return &[]*models.OperationType{} }, func(mods ...qm.QueryMod) Bindable { return models.OperationTypes(mods...) }},

		// Authors
		TableInfo{T_AUTHORS, []string{"id"}, func() interface{} { return &[]*models.Author{} }, func(mods ...qm.QueryMod) Bindable { return models.Authors(mods...) }},
		TableInfo{T_AUTHOR_I18N, []string{"author_id", "language"}, func() interface{} { return &[]*models.AuthorI18n{} }, func(mods ...qm.QueryMod) Bindable { return models.AuthorI18ns(mods...) }},

		// Blogs
		TableInfo{T_BLOGS, []string{"id"}, func() interface{} { return &[]*models.Blog{} }, func(mods ...qm.QueryMod) Bindable { return models.Blogs(mods...) }},
		TableInfo{T_BLOG_POSTS, []string{"id"}, func() interface{} { return &[]*models.BlogPost{} }, func(mods ...qm.QueryMod) Bindable { return models.BlogPosts(mods...) }},

		// Collections
		TableInfo{T_COLLECTIONS, []string{"id"}, func() interface{} { return &[]*models.Collection{} }, func(mods ...qm.QueryMod) Bindable { return models.Collections(mods...) }},
		TableInfo{T_COLLECTION_I18N, []string{"collection_id", "language"}, func() interface{} { return &[]*models.CollectionI18n{} }, func(mods ...qm.QueryMod) Bindable { return models.CollectionI18ns(mods...) }},

		// ContentUnits
		TableInfo{T_CONTENT_UNITS, []string{"id"}, func() interface{} { return &[]*models.ContentUnit{} }, func(mods ...qm.QueryMod) Bindable { return models.ContentUnits(mods...) }},
		TableInfo{T_CONTENT_UNIT_I18N, []string{"content_unit_id", "language"}, func() interface{} { return &[]*models.ContentUnitI18n{} }, func(mods ...qm.QueryMod) Bindable { return models.ContentUnitI18ns(mods...) }},

		// Users
		TableInfo{T_USERS, []string{"id"}, func() interface{} { return &[]*models.User{} }, func(mods ...qm.QueryMod) Bindable { return models.Users(mods...) }},

		// Files
		TableInfo{T_FILES, []string{"id"}, func() interface{} { return &[]*models.File{} }, func(mods ...qm.QueryMod) Bindable { return models.Files(mods...) }},

		// Operations
		TableInfo{T_OPERATIONS, []string{"id"}, func() interface{} { return &[]*models.Operation{} }, func(mods ...qm.QueryMod) Bindable { return models.Operations(mods...) }},

		// Storages
		TableInfo{T_STORAGES, []string{"id"}, func() interface{} { return &[]*models.Storage{} }, func(mods ...qm.QueryMod) Bindable { return models.Storages(mods...) }},

		// Publishers
		TableInfo{T_PUBLISHERS, []string{"id"}, func() interface{} { return &[]*models.Publisher{} }, func(mods ...qm.QueryMod) Bindable { return models.Publishers(mods...) }},
		TableInfo{T_PUBLISHER_I18N, []string{"publisher_id", "language"}, func() interface{} { return &[]*models.PublisherI18n{} }, func(mods ...qm.QueryMod) Bindable { return models.PublisherI18ns(mods...) }},

		// Persons
		TableInfo{T_PERSONS, []string{"id"}, func() interface{} { return &[]*models.Person{} }, func(mods ...qm.QueryMod) Bindable { return models.Persons(mods...) }},
		TableInfo{T_PERSON_I18N, []string{"person_id", "language"}, func() interface{} { return &[]*models.PersonI18n{} }, func(mods ...qm.QueryMod) Bindable { return models.PersonI18ns(mods...) }},

		// Sources
		TableInfo{T_SOURCES, []string{"id"}, func() interface{} { return &[]*models.Source{} }, func(mods ...qm.QueryMod) Bindable { return models.Sources(mods...) }},
		TableInfo{T_SOURCE_I18N, []string{"source_id", "language"}, func() interface{} { return &[]*models.SourceI18n{} }, func(mods ...qm.QueryMod) Bindable { return models.SourceI18ns(mods...) }},
		TableInfo{T_SOURCE_TYPES, []string{"id"}, func() interface{} { return &[]*models.SourceType{} }, func(mods ...qm.QueryMod) Bindable { return models.SourceTypes(mods...) }},

		// Tags
		TableInfo{T_TAGS, []string{"id"}, func() interface{} { return &[]*models.Tag{} }, func(mods ...qm.QueryMod) Bindable { return models.Tags(mods...) }},
		TableInfo{T_TAG_I18N, []string{"tag_id", "language"}, func() interface{} { return &[]*models.TagI18n{} }, func(mods ...qm.QueryMod) Bindable { return models.TagI18ns(mods...) }},

		// Twitter
		TableInfo{T_TWITTER_USERS, []string{"id"}, func() interface{} { return &[]*models.TwitterUser{} }, func(mods ...qm.QueryMod) Bindable { return models.TwitterUsers(mods...) }},
		TableInfo{T_TWITTER_TWEETS, []string{"id"}, func() interface{} { return &[]*models.TwitterTweet{} }, func(mods ...qm.QueryMod) Bindable { return models.TwitterTweets(mods...) }},

		// Content Unit Relations
		TableInfo{T_COLLECTIONS_CONTENT_UNITS, []string{"collection_id", "content_unit_id"}, func() interface{} { return &[]*models.CollectionsContentUnit{} }, func(mods ...qm.QueryMod) Bindable { return models.CollectionsContentUnits(mods...) }},
		TableInfo{T_CONTENT_UNIT_DERIVATIONS, []string{"source_id", "derived_id"}, func() interface{} { return &[]*models.ContentUnitDerivation{} }, func(mods ...qm.QueryMod) Bindable { return models.ContentUnitDerivations(mods...) }},
		TableInfo{T_CONTENT_UNITS_TAGS, []string{"content_unit_id", "tag_id"}, func() interface{} { return &[]*ContentUnitTag{} }, RelationQuery("content_units_tags", "content_unit_id", "tag_id")},
		TableInfo{T_CONTENT_UNITS_PERSONS, []string{"content_unit_id", "person_id"}, func() interface{} { return &[]*models.ContentUnitsPerson{} }, func(mods ...qm.QueryMod) Bindable { return models.ContentUnitsPersons(mods...) }},
		TableInfo{T_CONTENT_UNITS_SOURCES, []string{"content_unit_id", "source_id"}, func() interface{} { return &[]*ContentUnitSource{} }, RelationQuery("content_units_sources", "content_unit_id", "source_id")},
		TableInfo{T_CONTENT_UNITS_PUBLISHERS, []string{"content_unit_id", "publisher_id"}, func() interface{} { return &[]*ContentUnitPublisher{} }, RelationQuery("content_units_publishers", "content_unit_id", "publisher_id")},

		// File Relations
		TableInfo{T_FILES_OPERATIONS, []string{"file_id", "operation_id"}, func() interface{} { return &[]*FileOperation{} }, RelationQuery("files_operations", "file_id", "operation_id")},
		TableInfo{T_FILES_STORAGES, []string{"file_id", "storage_id"}, func() interface{} { return &[]*FileStorage{} }, RelationQuery("files_storages", "file_id", "storage_id")},

		// Author Sources
		TableInfo{T_AUTHORS_SOURCES, []string{"author_id", "source_id"}, func() interface{} { return &[]*AuthorSource{} }, RelationQuery("authors_sources", "author_id", "source_id")},
	}
}

func syncLocalMdb(tables []TableInfo, local *sql.DB, remote *sql.DB) error {
	scope := make(map[string]*ScopeIds)
	for _, info := range tables {
		var remoteIds, localIds []IdsTuple
		var err error
		if remoteIds, err = fetchIds(info, remote); err != nil {
			return err
		}
		if localIds, err = fetchIds(info, local); err != nil {
			return err
		}
		sort.Slice(localIds, IdsTupleSliceLess(localIds))
		sort.Slice(remoteIds, IdsTupleSliceLess(remoteIds))
		ids := append(diffIdsTupleSlices(localIds, remoteIds), diffIdsTupleSlices(remoteIds, localIds)...)
		log.Infof("%-30s local ids: %-7d remote ids: %-7d diff: %d", fmt.Sprintf("[%s]", info.Name), len(localIds), len(remoteIds), len(ids))
		for _, id := range ids {
			addScopeId(info.Name, id, scope, SCOPE_LOCAL)
			addScopeId(info.Name, id, scope, SCOPE_REMOTE)
		}
	}

	if err := applyScope(scope, tables, local, remote); err != nil {
		return err
	}
	log.Infof("Finished syncing all tables.")

	return nil
}

func InterfaceToInt64(n interface{}) (int64, error) {
	if f, ok := n.(float64); ok {
		return int64(f), nil
	}
	if i, ok := n.(int64); ok {
		return i, nil
	}
	return 0, errors.New(fmt.Sprintf("Expected int64 or float64, got %+v", reflect.TypeOf(n)))
}

const (
	SCOPE_LOCAL  = true
	SCOPE_REMOTE = false
)

func addScopeId(table string, id IdsTuple, scope map[string]*ScopeIds, isLocal bool) {
	if _, ok := scope[table]; !ok {
		scope[table] = &ScopeIds{nil, nil}
	}
	scopeIds := scope[table]
	if isLocal {
		for _, localId := range scopeIds.local {
			if CompareIdsTuple(&localId, &id) == 0 {
				return
			}
		}
		scopeIds.local = append(scopeIds.local, id)
		sort.Slice(scopeIds.local, IdsTupleSliceLess(scopeIds.local))
	} else {
		for _, remoteId := range scopeIds.remote {
			if CompareIdsTuple(&remoteId, &id) == 0 {
				return
			}
		}
		scopeIds.remote = append(scopeIds.remote, id)
		sort.Slice(scopeIds.remote, IdsTupleSliceLess(scopeIds.remote))
	}
}

func readIdFromEvent(key string, payload map[string]interface{}) (int64, error) {
	if id, ok := payload[key]; !ok {
		return 0, errors.New(fmt.Sprintf("Failed extracting %s from payload: %+v. Skipping", key, payload))
	} else {
		if int64Id, err := InterfaceToInt64(id); err != nil {
			return 0, errors.New(fmt.Sprintf("Failed converting %s for: %+v. Skipping", key, payload))
		} else {
			return int64Id, nil
		}
	}
}

func addCollectionScope(collectionId int64, scope map[string]*ScopeIds, local, remote *sql.DB) error {
	addScopeId(T_COLLECTIONS, IdsTuple{collectionId, 0, ""}, scope, SCOPE_LOCAL)
	addScopeId(T_COLLECTION_I18N, IdsTuple{collectionId, 0, ""}, scope, SCOPE_LOCAL)
	addScopeId(T_COLLECTIONS_CONTENT_UNITS, IdsTuple{collectionId, 0, ""}, scope, SCOPE_LOCAL)
	addScopeId(T_COLLECTIONS, IdsTuple{collectionId, 0, ""}, scope, SCOPE_REMOTE)
	addScopeId(T_COLLECTION_I18N, IdsTuple{collectionId, 0, ""}, scope, SCOPE_REMOTE)
	addScopeId(T_COLLECTIONS_CONTENT_UNITS, IdsTuple{collectionId, 0, ""}, scope, SCOPE_REMOTE)
	return nil
}

func addContentUnitScope(contentUnitIds []int64, scope map[string]*ScopeIds, local, remote *sql.DB) error {
	if err := addContentUnitScopeSingleSide(contentUnitIds, scope, local, SCOPE_LOCAL, true /*withFiles*/); err != nil {
		return err
	}
	if err := addContentUnitScopeSingleSide(contentUnitIds, scope, remote, SCOPE_REMOTE, true /*withFiles*/); err != nil {
		return err
	}
	return nil
}

func addContentUnitScopeSingleSide(contentUnitIds []int64, scope map[string]*ScopeIds, exec *sql.DB, isLocal bool, withFiles bool) error {
	for _, contentUnitId := range contentUnitIds {
		addScopeId(T_CONTENT_UNITS, IdsTuple{contentUnitId, 0, ""}, scope, isLocal)
		addScopeId(T_CONTENT_UNIT_I18N, IdsTuple{contentUnitId, 0, ""}, scope, isLocal)
		addScopeId(T_COLLECTIONS_CONTENT_UNITS, IdsTuple{0, contentUnitId, ""}, scope, isLocal)
		addScopeId(T_CONTENT_UNIT_DERIVATIONS, IdsTuple{contentUnitId, 0, ""}, scope, isLocal)
		addScopeId(T_CONTENT_UNIT_DERIVATIONS, IdsTuple{0, contentUnitId, ""}, scope, isLocal)
		addScopeId(T_CONTENT_UNITS_PERSONS, IdsTuple{contentUnitId, 0, ""}, scope, isLocal)
		addScopeId(T_CONTENT_UNITS_TAGS, IdsTuple{contentUnitId, 0, ""}, scope, isLocal)
		addScopeId(T_CONTENT_UNITS_SOURCES, IdsTuple{contentUnitId, 0, ""}, scope, isLocal)
		addScopeId(T_CONTENT_UNITS_PUBLISHERS, IdsTuple{contentUnitId, 0, ""}, scope, isLocal)
	}

	if withFiles {
		if files, err := models.Files(qm.Select("id"), qm.WhereIn("content_unit_id in ?", utils.ToInterfaceSlice(contentUnitIds)...)).All(context.TODO(), exec); err != nil {
			return err
		} else {
			fileIds := []int64(nil)
			for _, f := range files {
				fileIds = append(fileIds, f.ID)
			}
			if _, _, err := addFileScopeSingleSide(fileIds, scope, exec, isLocal, false /*withContentUnits*/); err != nil {
				return err
			}
		}
	}

	return nil
}

func addFileScopeParentsOnly(fileIds []int64, scope map[string]*ScopeIds, exec *sql.DB) ([]int64, error) {
	children := make([]int64, len(fileIds))
	copy(children, fileIds)

	fileIdsSlice := []int64(nil)
	fileIdsMap := make(map[int64]bool)
	for len(children) > 0 {
		log.Infof("Children: %+v", children)
		if files, err := models.Files(qm.Select("id, parent_id"), qm.WhereIn("id in ?", utils.ToInterfaceSlice(children)...)).All(context.TODO(), exec); err != nil {
			return nil, err
		} else {
			for _, file := range files {
				if !fileIdsMap[file.ID] {
					fileIdsSlice = append(fileIdsSlice, file.ID)
				}
				fileIdsMap[file.ID] = true
			}
			children = []int64(nil)
			for _, file := range files {
				if file.ParentID.Valid && !fileIdsMap[file.ParentID.Int64] {
					fileIdsSlice = append(fileIdsSlice, file.ParentID.Int64)
					children = append(children, file.ParentID.Int64)
					fileIdsMap[file.ParentID.Int64] = true
				}
			}
		}
	}
	log.Infof("Closure: %+v", fileIdsSlice)
	return fileIdsSlice, nil
}

func addFileScope(fileIds []int64, scope map[string]*ScopeIds, local, remote *sql.DB) error {
	var localUnits, remoteUnits map[int64]bool
	var withParentsRemote []int64
	var err error
	if withParentsRemote, remoteUnits, err = addFileScopeSingleSide(fileIds, scope, remote, SCOPE_REMOTE, true /*withContentUnits*/); err != nil {
		return err
	}
	if _, localUnits, err = addFileScopeSingleSide(withParentsRemote, scope, local, SCOPE_LOCAL, true /*withContentUnits*/); err != nil {
		return err
	}

	contentUnitIds := []int64(nil)
	for id, _ := range remoteUnits {
		if !localUnits[id] {
			contentUnitIds = append(contentUnitIds, id)
		}
	}

	if len(contentUnitIds) > 0 {
		return addContentUnitScopeSingleSide(contentUnitIds, scope, remote, SCOPE_REMOTE, false /*withFiles*/)
	}

	return nil
}

func addFileScopeSingleSide(fileIds []int64, scope map[string]*ScopeIds, exec *sql.DB, isLocal bool, withContentUnits bool) ([]int64, map[int64]bool, error) {
	var withParents []int64
	var err error
	if withParents, err = addFileScopeParentsOnly(fileIds, scope, exec); err != nil {
		return nil, nil, err
	}
	for _, fileId := range withParents {
		addScopeId(T_FILES, IdsTuple{fileId, 0, ""}, scope, isLocal)
		addScopeId(T_FILES_OPERATIONS, IdsTuple{fileId, 0, ""}, scope, isLocal)
		addScopeId(T_FILES_STORAGES, IdsTuple{fileId, 0, ""}, scope, isLocal)
	}

	if !withContentUnits {
		return withParents, nil, nil
	}

	units := make(map[int64]bool)
	if files, err := models.Files(qm.Select("content_unit_id"), qm.WhereIn("id in ?", utils.ToInterfaceSlice(withParents)...)).All(context.TODO(), exec); err != nil {
		return withParents, nil, errors.Wrap(err, "Faile fetching local files")
	} else {
		for _, file := range files {
			if file.ContentUnitID.Valid {
				units[file.ContentUnitID.Int64] = true
			} else {
				log.Warnf("Non valid content unit id in file id: %d", file.ID)
			}
		}
	}
	return withParents, units, nil
}

type ScopeIds struct {
	local  []IdsTuple
	remote []IdsTuple
}

func eventsScope(datas []events.Data, local, remote *sql.DB) (map[string]*ScopeIds, error) {
	start := time.Now()
	defer func() {
		utils.Profile("eventsScope", time.Now().Sub(start))
	}()
	scope := make(map[string]*ScopeIds)
	contentUnitIds := []int64(nil)
	fileIds := []int64(nil)
	for _, data := range datas {
		switch data.Type {
		case E_BLOG_POST_CREATE, E_BLOG_POST_UPDATE, E_BLOG_POST_DELETE:
			if id, err := readIdFromEvent("wpId", data.Payload); err != nil {
				log.Warnf("%+v", err)
			} else {
				addScopeId(T_BLOG_POSTS, IdsTuple{id, 0, ""}, scope, SCOPE_LOCAL)
				addScopeId(T_BLOG_POSTS, IdsTuple{id, 0, ""}, scope, SCOPE_REMOTE)
			}

		case E_COLLECTION_CREATE, E_COLLECTION_UPDATE, E_COLLECTION_DELETE, E_COLLECTION_PUBLISHED_CHANGE, E_COLLECTION_CONTENT_UNITS_CHANGE:
			if id, err := readIdFromEvent("id", data.Payload); err != nil {
				log.Warnf("Failed fetching id: %+v", err)
			} else {
				if err := addCollectionScope(id, scope, local, remote); err != nil {
					return nil, err
				}
			}

		case E_CONTENT_UNIT_CREATE, E_CONTENT_UNIT_UPDATE, E_CONTENT_UNIT_DELETE, E_CONTENT_UNIT_PUBLISHED_CHANGE, E_CONTENT_UNIT_DERIVATIVES_CHANGE, E_CONTENT_UNIT_SOURCES_CHANGE, E_CONTENT_UNIT_TAGS_CHANGE, E_CONTENT_UNIT_PERSONS_CHANGE, E_CONTENT_UNIT_PUBLISHERS_CHANGE:
			if id, err := readIdFromEvent("id", data.Payload); err != nil {
				log.Warnf("%+v", err)
			} else {
				contentUnitIds = append(contentUnitIds, id)
			}

		case E_FILE_REPLACE:
			for _, key := range []string{"old", "new"} {
				if p, ok := data.Payload[key]; !ok {
					log.Warnf("Failed extracting %s from %s: %+v. Skipping.", key, E_FILE_REPLACE, data)
				} else {
					if mapP, ok := p.(map[string]interface{}); !ok {
						log.Warnf("Failed converting payload to map: %+v. Skipping.", p)
					} else {
						if id, err := readIdFromEvent("id", mapP); err != nil {
							log.Warnf("%+v", err)
						} else {
							fileIds = append(fileIds, id)
						}
					}
				}
			}

		case E_FILE_UPDATE, E_FILE_PUBLISHED, E_FILE_INSERT, E_FILE_REMOVE:
			if id, err := readIdFromEvent("id", data.Payload); err != nil {
				log.Warnf("%+v", err)
			} else {
				fileIds = append(fileIds, id)
			}

		case E_SOURCE_CREATE, E_SOURCE_UPDATE:
			if id, err := readIdFromEvent("id", data.Payload); err != nil {
				log.Warnf("%+v", err)
			} else {
				addScopeId(T_SOURCES, IdsTuple{id, 0, ""}, scope, SCOPE_LOCAL)
				addScopeId(T_SOURCES, IdsTuple{id, 0, ""}, scope, SCOPE_REMOTE)
			}

		case E_TAG_CREATE, E_TAG_UPDATE:
			if id, err := readIdFromEvent("id", data.Payload); err != nil {
				log.Warnf("%+v", err)
			} else {
				addScopeId(T_TAGS, IdsTuple{id, 0, ""}, scope, SCOPE_LOCAL)
				addScopeId(T_TAGS, IdsTuple{id, 0, ""}, scope, SCOPE_REMOTE)
			}

		case E_PERSON_CREATE, E_PERSON_UPDATE, E_PERSON_DELETE:
			if id, err := readIdFromEvent("id", data.Payload); err != nil {
				log.Warnf("%+v", err)
			} else {
				addScopeId(T_PERSONS, IdsTuple{id, 0, ""}, scope, SCOPE_LOCAL)
				addScopeId(T_PERSONS, IdsTuple{id, 0, ""}, scope, SCOPE_REMOTE)
			}

		case E_PUBLISHER_CREATE, E_PUBLISHER_UPDATE:
			if id, err := readIdFromEvent("id", data.Payload); err != nil {
				log.Warnf("%+v", err)
			} else {
				addScopeId(T_PUBLISHERS, IdsTuple{id, 0, ""}, scope, SCOPE_LOCAL)
				addScopeId(T_PUBLISHERS, IdsTuple{id, 0, ""}, scope, SCOPE_REMOTE)
			}

		case E_TWEET_CREATE, E_TWEET_UPDATE, E_TWEET_DELETE:
			if id, err := readIdFromEvent("tid", data.Payload); err != nil {
				log.Warnf("%+v", err)
			} else {
				addScopeId(T_TWITTER_TWEETS, IdsTuple{id, 0, ""}, scope, SCOPE_LOCAL)
				addScopeId(T_TWITTER_TWEETS, IdsTuple{id, 0, ""}, scope, SCOPE_REMOTE)
			}

		default:
			log.Warnf("Did not expect event of type: %+v. Ignoring.", data)
		}
	}

	if err := addContentUnitScope(contentUnitIds, scope, local, remote); err != nil {
		return nil, err
	}
	if err := addFileScope(fileIds, scope, local, remote); err != nil {
		return nil, err
	}
	return scope, nil
}

func (m *MdbView) Refresh() error {
	start := time.Now()
	defer func() {
		utils.Profile("Refresh", time.Now().Sub(start))
	}()
	datas := events.ReadAndClearEvents()
	log.Infof("New events to handle: %+v", datas)
	if scope, err := eventsScope(datas, m.local, m.remote); err != nil {
		return errors.Wrap(err, "Error generating scope from events")
	} else {
		if err := applyScope(scope, m.tables, m.local, m.remote); err != nil {
			return errors.Wrap(err, "Error applying scope")
		}
	}
	return nil
}

func IdsToString(ids []IdsTuple) string {
	parts := []string(nil)
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("(%d,%d,%s)", id.Id1, id.Id2, id.Language))
	}
	return strings.Join(parts, ",")
}

func PrintScope(scope map[string]*ScopeIds) {
	log.Info("Scope:")
	for table, info := range scope {
		log.Infof("[%s]:", table)
		log.Infof("  Local: %s", IdsToString(info.local))
		log.Infof("  Remote: %s", IdsToString(info.remote))
	}
}

func applyScope(scope map[string]*ScopeIds, tables []TableInfo, local *sql.DB, remote *sql.DB) error {
	start := time.Now()
	defer func() {
		utils.Profile("applyScope", time.Now().Sub(start))
	}()
	defer utils.PrintProfile(true)
	PrintScope(scope)
	// First delete in reverse order.
	log.Infof("Deleting...")
	for i := len(tables) - 1; i >= 0; i-- {
		info := tables[i]
		if scopeIds, ok := scope[info.Name]; ok && len(scopeIds.local) > 0 {
			// log.Infof("Deleting from %s with ids %+v", info.Name, scopeIds)
			if deleted, err := DeleteFromTable(info, scopeIds.local, local); err != nil {
				log.Warnf("Error deleting from table %s: %+v", info.Name, err)
			} else {
				log.Infof("Deleted %d from %s", deleted, info.Name)
			}
		} else {
			log.Infof("No ids for %s", info.Name)
		}
	}
	log.Infof("Inserting...")
	for _, info := range tables {
		if scopeIds, ok := scope[info.Name]; ok {
			// log.Infof("Inserting to %s with ids %+v", info.Name, scopeIds)
			if inserted, err := InsertToTable(info, scopeIds.remote, local, remote); err != nil {
				return err
			} else {
				log.Infof("Inserted %d to %s", inserted, info.Name)
			}
		} else {
			log.Infof("No ids for %s", info.Name)
		}
	}
	log.Infof("Finished syncing all tables.")
	return nil
}
