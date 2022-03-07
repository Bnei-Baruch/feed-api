package learn

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tma15/gonline"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	elastic "gopkg.in/olivere/elastic.v6"

	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/databases/mdb"
	mdbModels "github.com/Bnei-Baruch/feed-api/databases/mdb/models"
	"github.com/Bnei-Baruch/feed-api/utils"
)

var (
	unitsCache      map[string]interface{}
	termvectorCache map[string]map[string]float64
	tagsCache       map[string][]string
	sourcesCache    map[string][]string
)

func InitFeatures() error {
	unitsCache = make(map[string]interface{})
	termvectorCache = make(map[string]map[string]float64)
	var err error
	tagsCache, err = loadTags()
	if err != nil {
		return err
	}
	sourcesCache, err = loadSources()
	log.Infof("tags: %d", len(tagsCache))
	log.Infof("sources: %d", len(sourcesCache))
	return err
}

type Classifier struct {
	Learner gonline.LearnerInterface
}

func UnitToFeatures(unit *mdbModels.ContentUnit) (map[string]float64, error) {
	// log.Infof("Unit: %+v", unit)
	cache, ok := unitsCache[unit.UID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unit %s not found in cache.", unit.UID))
	}
	cu, ok := cache.(*mdbModels.ContentUnit)
	if !ok {
		return nil, errors.New("Unexpected cache type, expected ContentUnit")
	}
	// log.Infof("Cache: %+v", cu)
	// log.Infof("Cache.R: %+v", cu.R)

	features := make(map[string]float64)

	features[fmt.Sprintf("uid:%s", cu.UID)] = 1.0
	features[fmt.Sprintf("content_unit:%s", cu.UID)] = 1.0
	features[fmt.Sprintf("content_type:%s", mdb.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name)] = 1.0

	for _, ccu := range cu.R.CollectionsContentUnits {
		features[fmt.Sprintf("collection:%s", ccu.R.Collection.UID)] = 1.0
		features[fmt.Sprintf("collection_content_type:%s", mdb.CONTENT_TYPE_REGISTRY.ByID[ccu.R.Collection.TypeID].Name)] = 1.0
	}

	for _, tagUid := range tagsCache[cu.UID] {
		features[fmt.Sprintf("tag:%s", tagUid)] = 1.0
	}

	for _, sourceUid := range sourcesCache[cu.UID] {
		features[fmt.Sprintf("source:%s", sourceUid)] = 1.0
	}

	if cu.Properties.Valid {
		var props map[string]interface{}
		if err := json.Unmarshal(cu.Properties.JSON, &props); err != nil {
			return nil, errors.New(fmt.Sprintf("Failed unmarshling props for %s", cu.UID))
		}
		if filmDate, ok := props["film_date"]; ok {
			dateStr := strings.Split(filmDate.(string), "T")[0] // remove the 'time' part
			if val, err := time.Parse("2006-01-02", dateStr); err != nil {
				return nil, errors.New(fmt.Sprintf("Failed parsing time %s for %s", dateStr, cu.UID))
			} else {
				features["effective_date_year"] = float64(val.Year())
				features["effective_date_month"] = float64(val.Month())
				features["effective_date_day"] = float64(val.Day())
				features["effective_date_weekday"] = float64(val.Weekday())
				features["effective_date_timestamp"] = float64(val.Unix())
			}
		}
		if originalLanguage, ok := props["original_language"]; ok {
			features[fmt.Sprintf("original_language:%s", originalLanguage.(string))] = 1.0
		}
		if duration, ok := props["duration"]; ok {
			features["duration"] = duration.(float64)
		}
	}

	termvector, ok := termvectorCache[unit.UID]
	if ok {
		for term, score := range termvector {
			features[term] = score
		}
	}

	parts := []string(nil)
	for k, v := range features {
		parts = append(parts, fmt.Sprintf("%s:%.1f", k, v))
	}
	sort.Strings(parts)
	log.Infof("example: %s", strings.Join(parts, ","))
	return features, nil
}

func PreloadUnits(units []interface{}) error {
	uids := []string(nil)
	for i := range units {
		cu, ok := units[i].(*mdbModels.ContentUnit)
		if !ok {
			continue
		}
		if _, exist := unitsCache[cu.UID]; exist {
			continue
		}
		uids = append(uids, cu.UID)
	}
	if len(uids) == 0 {
		return nil
	}
	var contentUnits []*mdbModels.ContentUnit
	if err := mdbModels.NewQuery(
		qm.From("content_units as cu"),
		qm.Load("ContentUnitI18ns"),
		qm.Load("CollectionsContentUnits"),
		qm.Load("CollectionsContentUnits.Collection"),
		qm.WhereIn("uid in ?", utils.ToInterfaceSlice(uids)...),
	).Bind(nil, common.LocalMdb, &contentUnits); err != nil {
		return err
	}
	log.Infof("Loaded %d units", len(units))
	for i := range contentUnits {
		unitsCache[contentUnits[i].UID] = contentUnits[i]
	}
	return nil
}

func ExampleFeatures(example interface{}) (map[string]float64, error) {
	cu, ok := example.(*mdbModels.ContentUnit)
	if !ok {
		return nil, errors.New("Implemented for content unit only.")
	}
	return UnitToFeatures(cu)
}

func EvalLearner(epoch int, learner *gonline.LearnerInterface, x_test *[]map[string]float64, y_test *[]string) {
	numCorr := 0
	numTotal := 0
	cls := gonline.Classifier{}
	w := (*learner).GetParam()
	ft, label := (*learner).GetDics()
	cls.Weight = *w
	cls.FtDict = *ft
	cls.LabelDict = *label
	for i, x_i := range *x_test {
		j := cls.Predict(&x_i)
		if cls.LabelDict.Id2elem[j] == (*y_test)[i] {
			numCorr += 1
		}
		numTotal += 1
	}
	acc := float64(numCorr) / float64(numTotal)
	fmt.Printf("epoch:%d test accuracy: %f (%d/%d)\n", epoch, acc, numCorr, numTotal)
}

func Examples(examples []interface{}, label string) (*[]map[string]float64, *[]string, error) {
	x := []map[string]float64(nil)
	y := []string(nil)
	for i := range examples {
		features, err := ExampleFeatures(examples[i])
		if err != nil {
			// return nil, nil, err
			log.Warnf("Skipping %+v, not supported yet: %+v.", examples[i], err)
			continue
		}
		x = append(x, features)
		y = append(y, label)
	}
	return &x, &y, nil
}

func PrepareExamples(positives []interface{}, negatives []interface{}) (*[]map[string]float64, *[]string, *[]map[string]float64, *[]string, error) {
	units := append(positives, negatives...)
	if err := PreloadUnits(units); err != nil {
		return nil, nil, nil, nil, err
	}
	if err := CacheTermVectors(units); err != nil {
		return nil, nil, nil, nil, err
	}
	x_pos, y_pos, err := Examples(positives, "p")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	x_neg, y_neg, err := Examples(negatives, "n")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return x_pos, y_pos, x_neg, y_neg, nil
}

func CrossValidate(positives []interface{}, negatives []interface{}) error {
	x_pos, y_pos, x_neg, y_neg, err := PrepareExamples(positives, negatives)
	if err != nil {
		return err
	}
	x_examples := append(*x_pos, *x_neg...)
	y_examples := append(*y_pos, *y_neg...)
	gonline.ShuffleData(&x_examples, &y_examples)

	folds := 5
	size := len(x_examples) / folds
	log.Infof("Cross validating %d examples. Fold size: %d.", len(x_examples), size)
	for i := 0; i < folds; i++ {
		x_test := x_examples[size*i : utils.MinInt(len(x_examples), size*(1+i))]
		y_test := y_examples[size*i : utils.MinInt(len(x_examples), size*(1+i))]

		x_train := []map[string]float64(nil)
		y_train := []string(nil)
		if i > 0 {
			x_train = append(x_train, x_examples[0:size*i]...)
			y_train = append(y_train, y_examples[0:size*i]...)
		}
		if i < folds-1 {
			x_train = append(x_train, x_examples[size*(1+i):len(x_examples)]...)
			y_train = append(y_train, y_examples[size*(1+i):len(y_examples)]...)
		}

		log.Infof("Train size: %d, Test size: %d.", len(x_train), len(x_test))

		var learner gonline.LearnerInterface
		learner = gonline.NewPA("II", 1)
		times := 5
		for j := 0; j < times; j++ {
			gonline.ShuffleData(&x_train, &y_train)
			learner.Fit(&x_train, &y_train)
		}
		EvalLearner(i, &learner, &x_test, &y_test)
	}
	return nil
}

func loadSources() (map[string][]string, error) {
	rows, err := queries.Raw(fmt.Sprintf(`
	WITH RECURSIVE rec_sources AS (
		SELECT
			s.id,
			s.uid,
			s.position,
			ARRAY [a.code, s.uid] "path"
		FROM sources s INNER JOIN authors_sources aas ON s.id = aas.source_id
			INNER JOIN authors a ON a.id = aas.author_id
		UNION
		SELECT
			s.id,
			s.uid,
			s.position,
			rs.path || s.uid
		FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
	)
	SELECT
		cu.uid,
		array_agg(DISTINCT item) FILTER (WHERE item IS NOT NULL AND item <> '')
	FROM content_units_sources cus
			INNER JOIN rec_sources AS rs ON cus.source_id = rs.id
			INNER JOIN content_units AS cu ON cus.content_unit_id = cu.id
			, unnest(rs.path) item
	GROUP BY cu.uid;`)).Query(common.LocalMdb)

	if err != nil {
		return nil, errors.Wrap(err, "Load sources")
	}
	defer rows.Close()

	return rowsToUIDToValues(rows)
}

func loadTags() (map[string][]string, error) {
	rows, err := queries.Raw(fmt.Sprintf(`
	WITH RECURSIVE rec_tags AS (
		SELECT
			t.id,
			t.uid,
			ARRAY [t.uid] :: CHAR(8) [] "path"
		FROM tags t
		WHERE parent_id IS NULL
		UNION
		SELECT
			t.id,
			t.uid,
			(rt.path || t.uid) :: CHAR(8) []
		FROM tags t INNER JOIN rec_tags rt ON t.parent_id = rt.id
	)
	SELECT
		cu.uid,
		array_agg(DISTINCT item)
	FROM content_units_tags cut
			INNER JOIN rec_tags AS rt ON cut.tag_id = rt.id
			INNER JOIN content_units AS cu ON cut.content_unit_id = cu.id
			, unnest(rt.path) item
	GROUP BY cu.uid;`)).Query(common.LocalMdb)

	if err != nil {
		return nil, errors.Wrap(err, "Load tags")
	}
	defer rows.Close()

	return rowsToUIDToValues(rows)
}

func rowsToUIDToValues(rows *sql.Rows) (map[string][]string, error) {
	m := make(map[string][]string)

	for rows.Next() {
		var cuUID string
		var values pq.StringArray
		err := rows.Scan(&cuUID, &values)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		m[cuUID] = values
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err()")
	}

	return m, nil
}

type MdbUid struct {
	MdbUid string `json:"mdb_uid"`
}

func CacheTermVectors(units []interface{}) error {
	log.Infof("Caching term vectors for %d units.", len(units))
	uidsByLanguage := make(map[string][]string)
	for i := range units {
		cu, ok := units[i].(*mdbModels.ContentUnit)
		if !ok {
			continue
		}
		var exist bool
		if cu, exist = unitsCache[cu.UID].(*mdbModels.ContentUnit); !exist {
			continue
		}
		if _, ok := termvectorCache[cu.UID]; ok {
			continue
		}
		if cu.R != nil {
			for _, i18n := range cu.R.ContentUnitI18ns {
				if i18n.Language != "" {
					uidsByLanguage[i18n.Language] = append(uidsByLanguage[i18n.Language], cu.UID)
				}
			}
		}
	}
	if len(uidsByLanguage) == 0 {
		return nil
	}
	log.Infof("Loaded %d languages.", len(uidsByLanguage))
	for lang, uids := range uidsByLanguage {
		log.Infof("\t%s: %d", lang, len(uids))
	}
	indexTemplate := viper.GetString("elasticsearch.index_template")
	for lang, uids := range uidsByLanguage {
		langIndex := strings.ReplaceAll(indexTemplate, "<LANG>", lang)

		var searchResult *elastic.SearchResult

		log.Infof("Fetching document ids.")
		ids := make(map[string]string)
		total := -1
		for true {
			if searchResult != nil && searchResult.Hits != nil {
				if total == -1 {
					log.Infof("lang: %s, total: %d", lang, searchResult.Hits.TotalHits)
					total = int(searchResult.Hits.TotalHits)
				}
				for _, h := range searchResult.Hits.Hits {
					// log.Infof("Hits: %+v", h.Source)
					var mdbUid MdbUid
					if err := json.Unmarshal(*h.Source, &mdbUid); err != nil {
						return err
					}
					// log.Infof("%s: %s", h.Id, mdbUid.MdbUid)
					ids[h.Id] = mdbUid.MdbUid
				}
			}
			query := elastic.NewBoolQuery().Filter(elastic.NewTermsQuery("mdb_uid", utils.ToInterfaceSlice(uids)...))
			scrollClient := common.ESC.Scroll().
				Index(langIndex).
				Sort("_doc", true).
				KeepAlive("5m").
				Size(1000).
				FetchSourceContext(elastic.NewFetchSourceContext(true).Include("mdb_uid")).
				Query(query)
			if searchResult != nil {
				scrollClient = scrollClient.ScrollId(searchResult.ScrollId)
			}
			var err error
			searchResult, err = scrollClient.Do(context.TODO())
			if err != nil {
				if err == io.EOF {
					break
				}
				if elastic.IsNotFound(err) {
					log.Infof("Failed finding index for %s, skipping.", lang)
					break
				}
				return err
			}
		}
		log.Infof("%s: %d", lang, len(ids))

		if len(ids) == 0 {
			continue
		}

		docs := []interface{}(nil)
		for id, _ := range ids {
			doc := make(map[string]interface{})
			doc["_index"] = langIndex
			doc["_id"] = id
			doc["_type"] = "_doc"
			doc["fields"] = []string{"title", "description", "content"}
			doc["offsets"] = false
			doc["positions"] = false
			doc["payloads"] = false
			doc["term_statistics"] = false
			filter := make(map[string]interface{})
			filter["max_doc_freq"] = 2000
			doc["filter"] = filter

			docs = append(docs, doc)
		}
		body := make(map[string]interface{})
		body["docs"] = docs

		log.Infof("Fetching term vectors...")
		results, err := common.ESC.MultiTermVectors().BodyJson(body).Do(context.TODO())
		if err != nil {
			return err
		}
		for _, doc := range results.Docs {
			uid := ids[doc.Id]
			if _, ok := termvectorCache[uid]; !ok {
				termvectorCache[uid] = make(map[string]float64)
			}
			tv := termvectorCache[uid]

			for field, terms := range doc.TermVectors {
				for term, score := range terms.Terms {
					tv[fmt.Sprintf("%s.%s.%s", lang, field, term)] = score.Score
				}
			}
			// log.Infof("doc: %+v", doc)
		}
		log.Infof("Done fetching terms for lang: %s", lang)
	}
	return nil
}
