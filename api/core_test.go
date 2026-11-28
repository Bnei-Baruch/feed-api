package api

// Core use cases, called the way prod calls feed-api, on temporary databases built from the migrations and
// testdata/*_fixtures.sql:
//   - kmedia-mdb site: /recommend with the A/B variants current/new/random (src/sagas/recommended.js), /views, /watchingnow.
//   - kmedia mobile app: /feed with its namespaces and paging.
//   - archive-backend: /views.
// Only request-path models are built; ingestion (MDB sync, chronicles, SQL refresh) is not started.
//
// Run: FEED_TEST_PG_URL=postgres://user:pass@localhost/postgres?sslmode=disable go test ./api/ (skips when unset).

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/common/testdb"
	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/data_models"
	"github.com/Bnei-Baruch/feed-api/databases/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

var (
	testRouter *gin.Engine
	testMDB    *sql.DB
	testDM     *data_models.DataModels
)

func TestMain(m *testing.M) {
	pgUrl := os.Getenv("FEED_TEST_PG_URL")
	if pgUrl == "" {
		fmt.Println("FEED_TEST_PG_URL not set, skipping.")
		os.Exit(0)
	}
	os.Exit(run(m, pgUrl))
}

// Drops the temporary databases even when setup or a test panics.
func run(m *testing.M, pgUrl string) int {
	mdbTest, err := testdb.Make(pgUrl, "../databases/mdb/migrations", "testdata/mdb_fixtures.sql")
	utils.Must(err)
	defer mdbTest.Destroy()
	modelsTest, err := testdb.Make(pgUrl, "../databases/data_models/migrations", "testdata/data_models_fixtures.sql")
	utils.Must(err)
	defer modelsTest.Destroy()
	testMDB = mdbTest.DB
	utils.Must(mdb.InitTypeRegistries(testMDB))
	boil.SetDB(testMDB)
	modelsDb := common.MakeConnection(func() (*sql.DB, map[string]string, error) {
		return modelsTest.DB, map[string]string{}, nil
	}, func(db *sql.DB) error { return nil })

	testDM = makeRequestDataModels(testMDB, modelsDb)
	refreshAll()

	gin.SetMode(gin.TestMode)
	testRouter = gin.New()
	testRouter.Use(
		utils.DataStoresMiddleware(testMDB, testMDB, nil, modelsDb, testDM),
		utils.ErrorHandlingMiddleware(),
		utils.RecoveryMiddleware())
	SetupRoutes(testRouter)
	return m.Run()
}

// Same request-path models as data_models.MakeDataModels.
func makeRequestDataModels(db *sql.DB, modelsDb *common.Connection) *data_models.DataModels {
	ct := mdb.CONTENT_TYPE_REGISTRY.ByName
	return &data_models.DataModels{
		LanguagesContentUnitsFilter:   data_models.MakeMDBFilterModel(db, "LanguagesContentUnitsFilter", 0, fmt.Sprintf(data_models.LANGUAGES_CONTENT_UNITS_SQL, ct[consts.CT_ARTICLE].ID), [][]string{{"en"}, {"he"}, {"ru"}}),
		TagsContentUnitsFilter:        data_models.MakeMDBFilterModel(db, "TagsContentUnitsFilter", 0, data_models.TAGS_CONTENT_UNITS_SQL, nil),
		SourcesContentUnitsFilter:     data_models.MakeMDBFilterModel(db, "SourcesContentUnitsFilter", 0, data_models.SOURCES_CONTENT_UNITS_SQL, nil),
		PersonsContentUnitsFilter:     data_models.MakeMDBFilterModel(db, "PersonsContentUnitsFilter", 0, data_models.PERSONS_CONTENT_UNITS_SQL, [][]string{{data_models.RAV_PERSON_UID}}),
		CollectionsContentUnitsFilter: data_models.MakeMDBFilterModel(db, "CollectionsContentUnitsFilter", 0, data_models.COLLECTIONS_CONTENT_UNITS_SQL, nil),
		ContentUnitsCollectionsFilter: data_models.MakeMDBFilterModel(db, "ContentUnitsCollectionsFilter", 0, data_models.CONTENT_UNITS_COLLECTIONS_SQL, nil),
		ContentUnitsInfo:              data_models.MakeMDBDataModel(db, "ContentUnitsInfo", 0, fmt.Sprintf(data_models.CONTENT_UNITS_INFO_SQL, ct[consts.CT_LESSON_PART].ID), data_models.ScanContentUnitInfo, data_models.ContentUnitsPrefilter),
		CollectionsInfo:               data_models.MakeMDBDataModel(db, "CollectionsInfo", 0, data_models.COLLECTIONS_INFO_SQL, data_models.ScanCollectionInfo, data_models.CollectionsPrefilter),
		BlogPostsInfo:                 data_models.MakeMDBDataModel(db, "BlogPostsInfo", 0, data_models.BLOG_POST_INFO_SQL, data_models.ScanBlogPostInfo, data_models.BlogPostsPrefilter),
		SqlDataModel:                  data_models.MakeSqlDataModel(modelsDb),
	}
}

func refreshAll() {
	for _, m := range []interface{ Refresh() error }{
		testDM.LanguagesContentUnitsFilter, testDM.TagsContentUnitsFilter, testDM.SourcesContentUnitsFilter,
		testDM.PersonsContentUnitsFilter, testDM.CollectionsContentUnitsFilter, testDM.ContentUnitsCollectionsFilter,
		testDM.ContentUnitsInfo, testDM.CollectionsInfo, testDM.BlogPostsInfo,
	} {
		utils.Must(m.Refresh())
	}
}

func post(t *testing.T, path string, body interface{}, resp interface{}) int {
	t.Helper()
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json") // as the site, app and archive-backend send it
	testRouter.ServeHTTP(w, req)
	if os.Getenv("FEED_TEST_DEBUG") != "" {
		t.Logf("DEBUG %s %s -> %d %s", path, b, w.Code, w.Body.String())
	}
	if w.Code == http.StatusOK && resp != nil {
		if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
			t.Fatalf("%s: bad json: %v: %s", path, err, w.Body.String())
		}
	}
	return w.Code
}

// -- Cases from the data -- //

type unitCase struct {
	Name        string
	Uid         string
	ContentType string
	Tags        []string
	Sources     []string
	Collections []string
	// Children of each source's parent (the site's "source collection").
	SourceCollections [][]string
}

func queryStrings(t *testing.T, q string, args ...interface{}) []string {
	t.Helper()
	ret := []string(nil)
	rows, err := testMDB.Query(q, args...)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			t.Fatal(err)
		}
		ret = append(ret, s)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return ret
}

func loadUnit(t *testing.T, name, uid string) unitCase {
	t.Helper()
	c := unitCase{Name: name, Uid: uid}
	var typeId int64
	if err := testMDB.QueryRow(`select type_id from content_units where uid = $1`, uid).Scan(&typeId); err != nil {
		t.Fatalf("%s %s: %v", name, uid, err)
	}
	c.ContentType = mdb.CONTENT_TYPE_REGISTRY.ByID[typeId].Name
	c.Tags = queryStrings(t, `select t.uid from tags t join content_units_tags x on x.tag_id = t.id join content_units cu on cu.id = x.content_unit_id where cu.uid = $1`, uid)
	c.Sources = queryStrings(t, `select s.uid from sources s join content_units_sources x on x.source_id = s.id join content_units cu on cu.id = x.content_unit_id where cu.uid = $1`, uid)
	c.Collections = queryStrings(t, `select c.uid from collections c join collections_content_units x on x.collection_id = c.id join content_units cu on cu.id = x.content_unit_id where cu.uid = $1 and c.published and c.secure = 0`, uid)
	for _, s := range c.Sources {
		c.SourceCollections = append(c.SourceCollections, queryStrings(t, `select s.uid from sources s where s.parent_id = (select parent_id from sources where uid = $1)`, s))
	}
	return c
}

func cases(t *testing.T) []unitCase {
	return []unitCase{
		loadUnit(t, "lesson-latest", "lesson11"),
		loadUnit(t, "lesson-oldest-preparation", "lesson40"),
		loadUnit(t, "program-chapter", "chaptr21"),
		loadUnit(t, "clip", "clipunt1"),
	}
}

// -- Site requests (kmedia-mdb src/sagas/recommended.js) -- //

func spec(filters []map[string]interface{}, order int) map[string]interface{} {
	return map[string]interface{}{"name": "DataContentUnitsSuggester", "filters": filters, "order_selector": order}
}

func roundRobin(selector int, args []string) map[string]interface{} {
	return map[string]interface{}{"name": "RoundRobinSuggester", "specs": []interface{}{
		spec([]map[string]interface{}{{"filter_selector": selector, "args": args}, {"filter_selector": 8}}, 5), // WatchingNow
		spec([]map[string]interface{}{{"filter_selector": selector, "args": args}}, 3),                         // Random
		spec([]map[string]interface{}{{"filter_selector": selector, "args": args}}, 0),                         // Last
	}}
}

func siteSpecs(c unitCase, variant string) []interface{} {
	specs := []interface{}(nil)
	switch variant {
	case "random":
		return []interface{}{map[string]interface{}{"name": "DataContentUnitsSuggester", "order_selector": 3}}
	case "new":
		specs = append(specs, spec([]map[string]interface{}{{"filter_selector": 1, "args": []string{"VIDEO_PROGRAM"}}}, 3))
		if utils.StringInSlice(c.ContentType, []string{consts.CT_LESSON_PART, consts.CT_FULL_LESSON, consts.CT_WOMEN_LESSON, consts.CT_VIRTUAL_LESSON, consts.CT_LECTURE}) {
			for _, s := range c.Sources {
				specs = append(specs, roundRobin(3, []string{s}))
			}
			for _, sc := range c.SourceCollections {
				if len(sc) > 0 { // A source without a parent has no source collection.
					specs = append(specs, roundRobin(3, sc))
				}
			}
		} else {
			for _, tag := range c.Tags {
				specs = append(specs, roundRobin(2, []string{tag}))
			}
			for _, col := range c.Collections {
				specs = append(specs, roundRobin(4, []string{col}))
			}
		}
	}
	return append(specs, map[string]interface{}{"name": "Default"})
}

func recommendRequest(c unitCase, specs []interface{}, languages, skip []string) map[string]interface{} {
	return map[string]interface{}{
		"more_items":   12,
		"current_feed": []interface{}{},
		"options": map[string]interface{}{
			"recommend":        map[string]interface{}{"uid": c.Uid},
			"languages":        languages,
			"skip_uids":        skip,
			"specs":            specs,
			"watching_now_min": 50,
			"popular_min":      100,
		},
	}
}

var siteLanguages = [][]string{{"he"}, {"ru"}, {"en"}, {"he", "en"}, {"ru", "en"}, {"es"}}

// -- Checks -- //

type itemInfo struct {
	exists, visible bool
	languages       map[string]bool
}

// Each item must be a published, non-secure unit or collection; for units, the languages they have files in.
func lookupItem(t *testing.T, ci core.ContentItem) itemInfo {
	t.Helper()
	info := itemInfo{languages: map[string]bool{}}
	var published bool
	var secure int
	var id int64
	if err := testMDB.QueryRow(`select id, published, secure from content_units where uid = $1`, ci.UID).Scan(&id, &published, &secure); err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	} else if err == nil {
		info.exists, info.visible = true, published && secure == 0
		for _, l := range queryStrings(t, `select distinct language from files where content_unit_id = $1 and language is not null`, id) {
			info.languages[l] = true
		}
		return info
	}
	if err := testMDB.QueryRow(`select published, secure from collections where uid = $1`, ci.UID).Scan(&published, &secure); err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	} else if err == nil {
		info.exists, info.visible, info.languages = true, published && secure == 0, nil
	}
	return info
}

func checkItems(t *testing.T, where string, items []core.ContentItem, languages []string, notUids map[string]bool) {
	t.Helper()
	for _, ci := range items {
		if notUids[ci.UID] {
			t.Errorf("%s: %s repeated, current or skipped", where, ci.UID)
		}
		notUids[ci.UID] = true
		info := lookupItem(t, ci)
		if !info.exists || !info.visible {
			t.Errorf("%s: %s (%s) exists=%v published&non-secure=%v", where, ci.UID, ci.ContentType, info.exists, info.visible)
		}
		if info.languages != nil && len(languages) > 0 {
			ok := false
			for _, l := range languages {
				ok = ok || info.languages[l]
			}
			if !ok {
				t.Errorf("%s: unit %s (%s) has no file in %v", where, ci.UID, ci.ContentType, languages)
			}
		}
		if ci.ContentType == "" || ci.Suggester == "" {
			t.Errorf("%s: %s missing content_type or suggester: %+v", where, ci.UID, ci)
		}
	}
}

func uids(items []core.ContentItem) []string {
	ret := []string(nil)
	for _, ci := range items {
		ret = append(ret, ci.UID)
	}
	return ret
}

// -- Tests -- //

func TestRecommendSiteVariants(t *testing.T) {
	for _, c := range cases(t) {
		for _, variant := range []string{"current", "new", "random"} {
			for _, languages := range siteLanguages {
				name := fmt.Sprintf("%s/%s/%s", c.Name, variant, strings.Join(languages, "+"))
				t.Run(name, func(t *testing.T) {
					specs := siteSpecs(c, variant)
					resp := MoreResponse{}
					if code := post(t, "/recommend", recommendRequest(c, specs, languages, nil), &resp); code != http.StatusOK {
						t.Fatalf("status %d", code)
					}
					if len(resp.Feeds) != len(specs) {
						t.Fatalf("feeds %d, want one per spec %d", len(resp.Feeds), len(specs))
					}
					total := 0
					seen := map[string]bool{c.Uid: true}
					for i, feed := range resp.Feeds {
						if len(feed) > 12 {
							t.Errorf("feed %d has %d items > more_items 12", i, len(feed))
						}
						total += len(feed)
						checkItems(t, fmt.Sprintf("feed %d", i), feed, languages, seen)
					}
					// The Default spec always has candidates for a common page in a site language.
					if last := resp.Feeds[len(resp.Feeds)-1]; variant != "random" && languages[0] != "es" && len(last) == 0 {
						t.Errorf("Default spec returned nothing")
					}
					t.Logf("%d feeds, %d items", len(resp.Feeds), total)
				})
			}
		}
	}
}

// The site passes skip uids (units already shown on the page); they never come back.
func TestRecommendSkipUids(t *testing.T) {
	c := cases(t)[0]
	first := MoreResponse{}
	if code := post(t, "/recommend", recommendRequest(c, siteSpecs(c, "current"), []string{"he"}, nil), &first); code != http.StatusOK {
		t.Fatalf("status %d", code)
	}
	skip := uids(first.Feeds[0])
	if len(skip) == 0 {
		t.Fatal("nothing recommended")
	}
	second := MoreResponse{}
	if code := post(t, "/recommend", recommendRequest(c, siteSpecs(c, "current"), []string{"he"}, skip), &second); code != http.StatusOK {
		t.Fatalf("status %d", code)
	}
	notUids := map[string]bool{c.Uid: true}
	for _, uid := range skip {
		notUids[uid] = true
	}
	if len(second.Feeds) == 0 || len(second.Feeds[0]) == 0 {
		t.Fatal("nothing recommended after skipping the first page")
	}
	for i, feed := range second.Feeds {
		checkItems(t, fmt.Sprintf("feed %d", i), feed, []string{"he"}, notUids)
	}
}

// Old clients send one spec (or none) and read "feed".
func TestRecommendSingleSpecAndDefault(t *testing.T) {
	c := cases(t)[0]
	for name, options := range map[string]map[string]interface{}{
		"no spec":      {"recommend": map[string]interface{}{"uid": c.Uid}, "languages": []string{"he"}},
		"default spec": {"recommend": map[string]interface{}{"uid": c.Uid}, "languages": []string{"he"}, "spec": map[string]interface{}{"name": "Default"}},
	} {
		t.Run(name, func(t *testing.T) {
			resp := MoreResponse{}
			if code := post(t, "/recommend", map[string]interface{}{"more_items": 12, "current_feed": []interface{}{}, "options": options}, &resp); code != http.StatusOK {
				t.Fatalf("status %d", code)
			}
			if len(resp.Feed) == 0 || len(resp.Feed) > 12 {
				t.Errorf("feed has %d items", len(resp.Feed))
			}
			checkItems(t, "feed", resp.Feed, []string{"he"}, map[string]bool{c.Uid: true})
		})
	}
}

func TestRecommendBadRequests(t *testing.T) {
	// An unknown suggester is an error from the spec, an unknown uid only has nothing to relate to.
	for _, tc := range []struct {
		name, uid, suggester string
		want                 int
	}{{"unknown suggester", "lesson11", "NoSuchSuggester", http.StatusInternalServerError}, {"unknown uid", "zzzzzzzz", "Default", http.StatusOK}} {
		t.Run(tc.name, func(t *testing.T) {
			body := map[string]interface{}{"more_items": 12, "options": map[string]interface{}{"recommend": map[string]interface{}{"uid": tc.uid}, "languages": []string{"he"}, "specs": []interface{}{map[string]interface{}{"name": tc.suggester}}}}
			if code := post(t, "/recommend", body, nil); code != tc.want {
				t.Errorf("status %d, want %d", code, tc.want)
			}
		})
	}
}

// Mobile app: /feed with its namespaces, paging with current_feed.
func TestFeedMobileNamespacesAndPaging(t *testing.T) {
	for _, ns := range []string{"", "kmedia-app-1", "kmedia-app-1.31", "kmedia-app-1.4", "kmedia-app-1.84"} {
		for _, languages := range [][]string{{"he"}, {"ru"}, {"en"}} {
			name := fmt.Sprintf("ns=%s/%s", ns, languages[0])
			t.Run(name, func(t *testing.T) {
				feed := []core.ContentItem(nil)
				seen := map[string]bool{}
				for page := 0; page < 2; page++ {
					resp := MoreResponse{}
					req := map[string]interface{}{"namespace": ns, "more_items": 10, "current_feed": feed, "options": map[string]interface{}{"languages": languages}}
					if code := post(t, "/feed", req, &resp); code != http.StatusOK {
						t.Fatalf("page %d: status %d", page, code)
					}
					if len(resp.Feed) == 0 {
						t.Fatalf("page %d: empty", page)
					}
					for i, ci := range resp.Feed {
						if ci.FeedOrder != int64(len(feed)+i) {
							t.Errorf("page %d: item %d feed_order %d", page, i, ci.FeedOrder)
						}
					}
					checkItems(t, fmt.Sprintf("page %d", page), resp.Feed, languages, seen)
					feed = append(feed, resp.Feed...)
				}
			})
		}
	}
}

// /views (site, archive-backend), /uniqueviews, /watchingnow (site): one count per uid, in order, 0 if unknown.
func TestViewsAndWatchingNow(t *testing.T) {
	req := map[string]interface{}{"uids": []string{"lesson11", "lesson12", "chaptr11", "clipunt1", "zzzzzzzz"}}
	for path, want := range map[string]map[string][]int64{
		"/views":       {"views": {500, 400, 300, 150, 0}},
		"/uniqueviews": {"views": {300, 250, 200, 100, 0}},
		"/watchingnow": {"watching_now": {80, 60, 0, 55, 0}},
	} {
		t.Run(path, func(t *testing.T) {
			resp := map[string][]int64{}
			if code := post(t, path, req, &resp); code != http.StatusOK {
				t.Fatalf("status %d", code)
			}
			if !reflect.DeepEqual(resp, want) {
				t.Errorf("got %v, want %v", resp, want)
			}
		})
	}
}
