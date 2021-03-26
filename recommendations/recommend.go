package recommendations

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Bnei-Baruch/feed-api/core"
)

type Recommender struct {
	Suggester core.Suggester
}

func MakeRecommender(db *sql.DB) (*Recommender, error) {
	lessonsJSON := `
		{
			"name": "RoundRobinSuggester",
			"specs": [
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "PrevContentUnitsSameSourceSuggester",
							"args": [
								"LESSON_PART"
							]
						},
						{
							"name": "ContentTypesSameTagSuggester",
							"args": [
								"LESSON_PART"
							],
							"time_selector": 2
						},
						{
							"name": "RandomContentUnitsSameSourceSuggester",
							"args": [
								"LESSON_PART"
							]
						},
						{
							"name": "ContentTypesSameTagSuggester",
							"args": [
								"LESSON_PART"
							],
							"time_selector": 3
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "NextContentUnitsSameSourceSuggester",
							"args": [
								"LESSON_PART"
							]
						},
						{
							"name": "ContentTypesSameTagSuggester",
							"time_selector": 1,
							"args": [
								"LESSON_PART"
							]
						},
						{
							"name": "RandomContentUnitsSameSourceSuggester",
							"args": [
								"LESSON_PART"
							]
						},
						{
							"name": "ContentTypesSameTagSuggester",
							"args": [
								"LESSON_PART"
							],
							"time_selector": 3
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "ContentUnitCollectionSuggester",
							"args": [
								"LESSONS_SERIES"
							]
						},
						{
							"name": "ContentTypesSameTagSuggester",
							"args": [
								"LESSON_PART"
							]
						},
						{
							"name": "RandomContentUnitsSameSourceSuggester",
							"args": [
								"LESSON_PART"
							]
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "ContentTypesSameTagSuggester",
							"args": [
								"LECTURE"
							]
						},
						{
							"name": "RandomContentUnitsSameSourceSuggester",
							"args": [
								"LECTURE"
							]
						},
						{
							"name": "RandomContentTypesSuggester",
							"args": [
								"LECTURE"
							]
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "ContentTypesSameTagSuggester",
							"args": [
								"VIDEO_PROGRAM_CHAPTER"
							]
						},
						{
							"name": "RandomContentTypesSuggester",
							"args": [
								"VIDEO_PROGRAM_CHAPTER"
							]
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "ContentTypesSameTagSuggester",
							"args": [
								"VIDEO_PROGRAM_CHAPTER"
							]
						},
						{
							"name": "RandomContentTypesSuggester",
							"args": [
								"VIDEO_PROGRAM_CHAPTER"
							]
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "ContentTypesSameTagSuggester",
							"args": [
								"CLIP"
							]
						},
						{
							"name": "RandomContentTypesSuggester",
							"args": [
								"CLIP"
							]
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "ContentTypesSameTagSuggester",
							"args": [
								"CLIP"
							]
						},
						{
							"name": "RandomContentTypesSuggester",
							"args": [
								"CLIP"
							]
						}
					]
				},
				{
					"name": "RandomContentTypesSuggester",
					"args": [
						"CLIP"
					]
				},
				{
					"name": "RandomContentTypesSuggester",
					"args": [
						"CLIP"
					]
				}
			]
		}
	`
	defaultJSON := `
		{
			"name":"CompletionSuggester",
			"specs":[
				{
					"name":"RoundRobinSuggester",
					"specs":[
						{
							"name":"CompletionSuggester",
							"specs":[
								{"name":"LastClipsSameTagSuggester"},
								{"name":"LastClipsSuggester"}
							]
						},
						{"name":"LastContentUnitsSameCollectionSuggester"},
						{"name":"PrevContentUnitsSameCollectionSuggester"},
						{"name":"NextContentUnitsSameSourceSuggester","args":["LESSON_PART"]},
						{"name":"PrevContentUnitsSameSourceSuggester","args":["LESSON_PART"]},
						{"name":"LastCollectionSameSourceSuggester","args":["LESSONS_SERIES"]},
						{
							"name":"CompletionSuggester",
							"specs":[
								{"name":"LastLessonsSameTagSuggester"},
								{"name":"LastLessonsSuggester"}
							]
						},
						{
							"name":"CompletionSuggester",
							"specs":[
								{"name":"LastProgramsSameTagSuggester"},
								{"name":"LastProgramsSuggester"}
							]
						},
						{"name":"LastCongressSameTagSuggester"}
					]
				},
				{"name":"RandomContentTypesSuggester","args":["CLIP","LESSON_PART","VIDEO_PROGRAM_CHAPTER"]}
			]
		}
	`

	rootJSON := fmt.Sprintf(`
		{
			"name":"ContentTypeSuggester",
			"args":["LESSON_PART", "*"],
			"specs":[%s,%s]
		}
	`, lessonsJSON, defaultJSON)

	var spec core.SuggesterSpec
	if err := json.Unmarshal([]byte(rootJSON), &spec); err != nil {
		return nil, err
	}

	fmt.Printf("?!?!?!??????????? %+v", spec)

	if s, err := core.MakeSuggesterFromName(db, spec.Name); err != nil {
		return nil, err
	} else {
		if err := s.UnmarshalSpec(db, spec); err != nil {
			return nil, err
		} else {
			return &Recommender{s}, nil
		}
	}

	/*return &Recommender{Suggester: core.MakeCompletionSuggester([]core.Suggester{
		core.MakeRoundRobinSuggester([]core.Suggester{
			core.MakeCompletionSuggester([]core.Suggester{MakeLastClipsSameTagSuggester(db), MakeLastClipsSuggester(db)}),
			MakeLastContentUnitsSameCollectionSuggester(db),
			MakePrevContentUnitsSameCollectionSuggester(db),
			MakeNextContentUnitsSameSourceSuggester(db, []string{consts.CT_LESSON_PART}),
			MakePrevContentUnitsSameSourceSuggester(db, []string{consts.CT_LESSON_PART}),
			MakeLastCollectionSameSourceSuggester(db, []string{consts.CT_LESSONS_SERIES}),
			core.MakeCompletionSuggester([]core.Suggester{MakeLastLessonsSameTagSuggester(db), MakeLastLessonsSuggester(db)}),
			core.MakeCompletionSuggester([]core.Suggester{MakeLastProgramsSameTagSuggester(db), MakeLastProgramsSuggester(db)}),
			MakeLastCongressSameTagSuggester(db),
		}),
		MakeRandomContentTypesSuggester(db, []string{consts.CT_CLIP, consts.CT_LESSON_PART, consts.CT_VIDEO_PROGRAM_CHAPTER}, []string(nil) /* tagUids * /),
	})}, nil*/

}

func (recommender *Recommender) Recommend(r core.MoreRequest) ([]core.ContentItem, error) {
	return recommender.Suggester.More(r)
}
