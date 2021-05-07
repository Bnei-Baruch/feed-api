package recommendations

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Bnei-Baruch/feed-api/core"
)

type Recommender struct {
	Suggester core.Suggester
}

func MakeRecommender(suggesterContext core.SuggesterContext) (*Recommender, error) {
	lessonsJSON := `
		{
			"name": "RoundRobinSuggester",
			"specs": [
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameSource
								}
							],
							"order_selector": OPrev
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": OPrev
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameSource
								}
							],
							"order_selector": ORand
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": ORand
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameSource
								}
							],
							"order_selector": ONext
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": ONext
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameSource
								}
							],
							"order_selector": ORand
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": ORand
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FCollectionContentTypes,
									"args": [
										"LESSONS_SERIES"
									]
								}
							],
							"order_selector": OLast
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": OLast
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameSource
								}
							],
							"order_selector": ORand
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LECTURE"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": OLast
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LECTURE"
									]
								},
								{
									"filter_selector": FSameSource
								}
							],
							"order_selector": ORand
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LECTURE"
									]
								}
							],
							"order_selector": ORand
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"VIDEO_PROGRAM_CHAPTER"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": OLast
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"VIDEO_PROGRAM_CHAPTER"
									]
								}
							],
							"order_selector": ORand
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"VIDEO_PROGRAM_CHAPTER"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": OLast
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"VIDEO_PROGRAM_CHAPTER"
									]
								}
							],
							"order_selector": ORand
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"CLIP"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": OLast
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"CLIP"
									]
								}
							],
							"order_selector": ORand
						}
					]
				},
				{
					"name": "CompletionSuggester",
					"specs": [
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"CLIP"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": OLast
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"CLIP"
									]
								}
							],
							"order_selector": ORand
						}
					]
				},
				{
					"name": "DataContentUnitsSuggester",
					"filters": [
						{
							"filter_selector": FUnitContentTypes,
							"args": [
								"CLIP"
							]
						}
					],
					"order_selector": ORand
				},
				{
					"name": "DataContentUnitsSuggester",
					"filters": [
						{
							"filter_selector": FUnitContentTypes,
							"args": [
								"CLIP"
							]
						}
					],
					"order_selector": ORand
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
								{
									"name": "DataContentUnitsSuggester",
									"filters": [
										{
											"filter_selector": FUnitContentTypes,
											"args": [
												"CLIP"
											]
										},
										{
											"filter_selector": FSameTag
										}
									],
									"order_selector": OLast
								},
								{
									"name": "DataContentUnitsSuggester",
									"filters": [
										{
											"filter_selector": FUnitContentTypes,
											"args": [
												"CLIP"
											]
										}
									],
									"order_selector": OLast
								}
							]
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FSameCollection
								}
							],
							"order_selector": OLast
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FSameCollection
								}
							],
							"order_selector": OPrev
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameSource
								}
							],
							"order_selector": ONext
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FUnitContentTypes,
									"args": [
										"LESSON_PART"
									]
								},
								{
									"filter_selector": FSameSource
								}
							],
							"order_selector": OPrev
						},
						{
							"name": "DataCollectionsSuggester",
							"filters": [
								{
									"filter_selector": FCollectionContentTypes,
									"args": [
										"LESSONS_SERIES"
									]
								},
								{
									"filter_selector": FSameSource
								}
							],
							"order_selector": OLast
						},
						{
							"name":"CompletionSuggester",
							"specs":[
								{
									"name": "RoundRobinSuggester",
									"specs": [
										{
											"name": "DataContentUnitsSuggester",
											"filters": [
												{
													"filter_selector": FUnitContentTypes,
													"args": [
														"LESSON_PART"
													]
												},
												{
													"filter_selector": FSameTag
												}
											],
											"order_selector": OLast
										},
										{
											"name": "DataContentUnitsSuggester",
											"filters": [
												{
													"filter_selector": FUnitContentTypes,
													"args": [
														"VIRTUAL_LESSON"
													]
												},
												{
													"filter_selector": FSameTag
												}
											],
											"order_selector": OLast
										},
										{
											"name": "DataContentUnitsSuggester",
											"filters": [
												{
													"filter_selector": FUnitContentTypes,
													"args": [
														"WOMEN_LESSON"
													]
												},
												{
													"filter_selector": FSameTag
												}
											],
											"order_selector": OLast
										},
										{
											"name": "DataContentUnitsSuggester",
											"filters": [
												{
													"filter_selector": FUnitContentTypes,
													"args": [
														"LECTURE"
													]
												},
												{
													"filter_selector": FSameTag
												}
											],
											"order_selector": OLast
										}
									]
								},
								{
									"name": "DataContentUnitsSuggester",
									"filters": [
										{
											"filter_selector": FUnitContentTypes,
											"args": [
												"LESSON_PART",
												"VIRTUAL_LESSON",
												"WOMEN_LESSON"
											]
										}
									],
									"order_selector": OLast
								}
							]
						},
						{
							"name":"CompletionSuggester",
							"specs":[
								{
									"name": "DataContentUnitsSuggester",
									"filters": [
										{
											"filter_selector": FUnitContentTypes,
											"args": [
												"VIDEO_PROGRAM_CHAPTER"
											]
										},
										{
											"filter_selector": FSameTag
										}
									],
									"order_selector": OLast
								},
								{
									"name": "DataContentUnitsSuggester",
									"filters": [
										{
											"filter_selector": FUnitContentTypes,
											"args": [
												"VIDEO_PROGRAM_CHAPTER"
											]
										}
									],
									"order_selector": OLast
								}
							]
						},
						{
							"name": "DataContentUnitsSuggester",
							"filters": [
								{
									"filter_selector": FCollectionContentTypes,
									"args": [
										"CONGRESS"
									]
								},
								{
									"filter_selector": FSameTag
								}
							],
							"order_selector": OLast
						}
					]
				},
				{
					"name": "DataContentUnitsSuggester",
					"filters": [
						{
							"filter_selector": FUnitContentTypes,
							"args": [
								"CLIP",
								"LESSON_PART",
								"VIDEO_PROGRAM_CHAPTER"
							]
						}
					],
					"order_selector": ORand
				}
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

	for filterName, filterValue := range core.FILTER_STRING_TO_VALUE {
		rootJSON = strings.ReplaceAll(rootJSON, filterName, fmt.Sprintf("%d", filterValue))
	}
	for orderName, orderValue := range core.ORDER_STRING_TO_VALUE {
		rootJSON = strings.ReplaceAll(rootJSON, orderName, fmt.Sprintf("%d", orderValue))
	}

	var spec core.SuggesterSpec
	if err := json.Unmarshal([]byte(rootJSON), &spec); err != nil {
		return nil, err
	}

	if s, err := core.MakeSuggesterFromName(suggesterContext, spec.Name); err != nil {
		return nil, err
	} else {
		if err := s.UnmarshalSpec(suggesterContext, spec); err != nil {
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
