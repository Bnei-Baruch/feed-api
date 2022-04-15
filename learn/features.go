package learn

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
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

const TERMVECTOR_DAT = "/tmp/termvector.dat"

var (
	unitsCache      map[string]interface{}
	termvectorCache map[string]map[string]float64
	tagsCache       map[string][]string
	sourcesCache    map[string][]string
)

func InitFeatures() error {
	unitsCache = make(map[string]interface{})

	var err error
	termvectorCache, err = termVectorFromFile(TERMVECTOR_DAT)
	if err != nil {
		return err
	}

	tagsCache, err = loadTags()
	if err != nil {
		return err
	}
	log.Infof("tags: %d", len(tagsCache))

	sourcesCache, err = loadSources()
	if err != nil {
		return err
	}
	log.Infof("sources: %d", len(sourcesCache))

	return err
}

func termVectorFromFile(filename string) (map[string]map[string]float64, error) {
	cache := make(map[string]map[string]float64)

	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReaderSize(fp, 4096*64)
	lineNum := 0
	for {
		line, _, err := reader.ReadLine()
		lineNum += 1
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if strings.HasPrefix(string(line), "#") {
			continue /* ignore comment */
		}
		// log.Infof("Parsing: [%s]", line)
		uidTv := strings.SplitN(string(line), " ", 2)
		if len(uidTv) != 2 {
			log.Infof("line:%d has no termvector. This line is ignored.", lineNum)
			continue
		}
		uid := uidTv[0]
		// log.Infof("Uid: %s", uid)
		termsValues := strings.Split(uidTv[1], " ")
		terms := make(map[string]float64)
		for _, termValue := range termsValues {
			// log.Infof("Term value: %s", termValue)
			tv := strings.Split(termValue, ":")
			if len(tv) < 2 {
				return nil, errors.New(fmt.Sprintf("Expected term value %s to have 2 (or more) parts.", termValue))
			} else if len(tv) > 2 {
				tv[0] = strings.Join(tv[:len(tv)-2], ":")
				tv[1] = tv[len(tv)-1]
			}
			if valueFloat64, err := strconv.ParseFloat(tv[1], 64); err != nil {
				return nil, errors.New(fmt.Sprintf("Failed parsing float %s: %s", valueFloat64, err))
			} else {
				terms[tv[0]] = valueFloat64
			}
		}
		// log.Infof("%s %+v", uid, terms)
		cache[uid] = terms
	}
	return cache, nil
}

func termVectorToFile(filename string, cache map[string]map[string]float64) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	datawriter := bufio.NewWriter(file)

	keys := utils.StringKeys(cache)
	sort.Strings(keys)

	count := 0
	for _, key := range keys {
		features := cache[key]
		featuresKeys := utils.StringKeys(features)
		sort.Strings(featuresKeys)
		featuresParts := []string(nil)
		for _, featureKey := range featuresKeys {
			featuresParts = append(featuresParts, fmt.Sprintf("%s:%f", featureKey, features[featureKey]))
		}
		if len(featuresParts) > 0 {
			line := fmt.Sprintf("%s %s\n", key, strings.Join(featuresParts, " "))
			_, err = datawriter.WriteString(line)
			if err != nil {
				return err
			}
			count += 1
		}
	}
	if err := datawriter.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	log.Infof("Writing %d lines to termvector cache.", count)
	return nil
}

type Classifier struct {
	Learner gonline.LearnerInterface
}

func UnitToFeatures(prefix string, unit *mdbModels.ContentUnit) (map[string]float64, error) {
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

	features[fmt.Sprintf("%suid_%s", prefix, cu.UID)] = 1.0
	features[fmt.Sprintf("%scontent_unit_%s", prefix, cu.UID)] = 1.0
	features[fmt.Sprintf("%scontent_type_%s", prefix, mdb.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name)] = 1.0

	for _, ccu := range cu.R.CollectionsContentUnits {
		features[fmt.Sprintf("%scollection_%s", prefix, ccu.R.Collection.UID)] = 1.0
		features[fmt.Sprintf("%scollection_content_type_%s", prefix, mdb.CONTENT_TYPE_REGISTRY.ByID[ccu.R.Collection.TypeID].Name)] = 1.0
	}

	for _, tagUid := range tagsCache[cu.UID] {
		features[fmt.Sprintf("%stag_%s", prefix, tagUid)] = 1.0
	}

	for _, sourceUid := range sourcesCache[cu.UID] {
		features[fmt.Sprintf("%ssource_%s", prefix, sourceUid)] = 1.0
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
				features[fmt.Sprintf("%seffective_date_year", prefix)] = float64(val.Year())
				features[fmt.Sprintf("%seffective_date_month", prefix)] = float64(val.Month())
				features[fmt.Sprintf("%seffective_date_day", prefix)] = float64(val.Day())
				features[fmt.Sprintf("%seffective_date_weekday", prefix)] = float64(val.Weekday())
				features[fmt.Sprintf("%seffective_date_timestamp", prefix)] = float64(val.Unix())
			}
		}
		if originalLanguage, ok := props["original_language"]; ok {
			features[fmt.Sprintf("%soriginal_language_%s", prefix, originalLanguage.(string))] = 1.0
		}
		if duration, ok := props["duration"]; ok {
			features[fmt.Sprintf("%sduration", prefix)] = duration.(float64)
		}
	}

	termvector, ok := termvectorCache[unit.UID]
	if ok {
		for term, score := range termvector {
			features[fmt.Sprintf("%s%s", prefix, term)] = score
		}
	}

	parts := []string(nil)
	for k, v := range features {
		parts = append(parts, fmt.Sprintf("%s_%.1f", k, v))
	}
	sort.Strings(parts)
	// log.Infof("example: %s", strings.Join(parts, ","))
	return features, nil
}

func PreloadUnits(units []interface{}) error {
	addedUids := make(map[string]bool)
	uids := []string(nil)
	for i := range units {
		cu, ok := units[i].(*mdbModels.ContentUnit)
		if !ok {
			continue
		}
		if _, exist := unitsCache[cu.UID]; exist {
			continue
		}
		if _, added := addedUids[cu.UID]; added {
			continue
		}
		uids = append(uids, cu.UID)
		addedUids[cu.UID] = true
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

func ExampleFeatures(example RecommendItemsPair) (map[string]float64, error) {
	rCu, rOk := example.Recommended.(*mdbModels.ContentUnit)
	if !rOk {
		return nil, errors.New("Implemented for content unit only.")
	}
	sCu, sOk := example.Selected.(*mdbModels.ContentUnit)
	if !sOk {
		return nil, errors.New("Implemented for content unit only.")
	}
	rFeatures, err := UnitToFeatures("r_", rCu)
	if err != nil {
		return nil, err
	}
	sFeatures, err := UnitToFeatures("s_", sCu)
	if err != nil {
		return nil, err
	}
	for k, v := range sFeatures {
		rFeatures[k] = v
	}
	return rFeatures, nil
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

func ExamplesFromPairs(examples []RecommendItemsPair, label string) (*[]map[string]float64, *[]string, error) {
	x := []map[string]float64(nil)
	y := []string(nil)
	skipped := 0
	for i := range examples {
		features, err := ExampleFeatures(examples[i])
		if err != nil {
			// return nil, nil, err
			if skipped < 3 {
				log.Warnf("Skipping %s example %+v, not supported yet: %+v.", label, examples[i], err)
			}
			skipped += 1
			continue
		}
		x = append(x, features)
		y = append(y, label)
	}
	if skipped >= 3 {
		log.Info("...")
		log.Info("...")
	}
	log.Infof("Skipped total of %d %s examples.", skipped, label)
	return &x, &y, nil
}

func PrepareExamples(positives []RecommendItemsPair, negatives []RecommendItemsPair) (*[]map[string]float64, *[]string, *[]map[string]float64, *[]string, error) {
	units := []interface{}(nil)
	for _, example := range append(positives, negatives...) {
		units = append(units, example.Recommended, example.Selected)
	}
	if err := PreloadUnits(units); err != nil {
		return nil, nil, nil, nil, err
	}
	if err := CacheTermVectors(units); err != nil {
		return nil, nil, nil, nil, err
	}
	x_pos, y_pos, err := ExamplesFromPairs(positives, "p")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	x_neg, y_neg, err := ExamplesFromPairs(negatives, "n")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return x_pos, y_pos, x_neg, y_neg, nil
}

func CrossValidate(positives []RecommendItemsPair, negatives []RecommendItemsPair) error {
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
		y_test := y_examples[size*i : utils.MinInt(len(y_examples), size*(1+i))]

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

		learners := make(map[string]gonline.LearnerInterface)
		learners["p"] = gonline.NewPerceptron()
		learners["pa"] = gonline.NewPA("", 0.01)
		learners["pa1"] = gonline.NewPA("I", 0.01)
		learners["pa2"] = gonline.NewPA("II", 0.01)
		learners["cw"] = gonline.NewCW(0.8)
		learners["arow"] = gonline.NewArow(10.)
		// learners["adam"] = gonline.NewAdam()

		for algorithm, learner := range learners {
			log.Infof("Algorithm: %s Fold: %d/%d", algorithm, i+1, folds)
			times := 5
			for j := 0; j < times; j++ {
				gonline.ShuffleData(&x_train, &y_train)
				learner.Fit(&x_train, &y_train)
			}
			EvalLearner(i, &learner, &x_test, &y_test)
		}
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
	log.Infof("Caching term vectors for %d units. Cache size: %d", len(units), len(termvectorCache))
	uidsByLanguage := make(map[string][]string)
	termvectorExist := 0
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
			termvectorExist += 1
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
	log.Infof("Loading %d languages (%d units were found).", len(uidsByLanguage), termvectorExist)
	for lang, uids := range uidsByLanguage {
		log.Infof("    %s: %d", lang, len(uids))
	}
	indexTemplate := viper.GetString("elasticsearch.index_template")
	for lang, uids := range uidsByLanguage {
		langIndex := strings.ReplaceAll(indexTemplate, "<LANG>", lang)

		var searchResult *elastic.SearchResult

		log.Infof("Fetching document ids for %s.", lang)
		ids := make(map[string]string)
		total := -1
		for true {
			if searchResult != nil && searchResult.Hits != nil {
				if total == -1 {
					log.Infof("    lang: %s, total: %d", lang, searchResult.Hits.TotalHits)
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
		log.Infof("Unique ids %s: %d", lang, len(ids))

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

		log.Infof("Fetching %d term vectors...", len(docs))
		results, err := common.ESC.MultiTermVectors().BodyJson(body).Do(context.TODO())
		if err != nil {
			return err
		}
		log.Infof("Fetched %d docs for %s", len(results.Docs), lang)
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
	return termVectorToFile(TERMVECTOR_DAT, termvectorCache)
}
