package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Bnei-Baruch/feed-api/consts"
)

type Feed struct {
	Suggester Suggester
}

func MakeFeedFromSuggester(suggester Suggester, suggesterContext SuggesterContext) *Feed {
	return &Feed{Suggester: suggester}
}

func MakeFeed(suggesterContext SuggesterContext) *Feed {
	return MakeFeedFromSuggester(MakeSortSuggester([]Suggester{MakeRoundRobinSuggester([]Suggester{
		// 1. Morning lesson.
		MakeCollectionSuggester(suggesterContext.DB, consts.CT_DAILY_LESSON),
		// 2. Additional lessons.
		MakeRoundRobinSuggester([]Suggester{
			MakeContentUnitsSuggester(suggesterContext.DB, []string{consts.CT_LECTURE}),
			MakeContentUnitsSuggester(suggesterContext.DB, []string{consts.CT_VIRTUAL_LESSON}),
			MakeContentUnitsSuggester(suggesterContext.DB, []string{consts.CT_WOMEN_LESSON}),
			MakeContentUnitsSuggester(suggesterContext.DB, []string{consts.CT_EVENT_PART}),
		}),
		// 3. TODO: Twitter.
		// 4. Programs.
		MakeContentUnitsSuggester(suggesterContext.DB, []string{
			consts.CT_VIDEO_PROGRAM_CHAPTER,
		}),
		// 5. Blog (Article?, Publication?).
		MakeRoundRobinSuggester([]Suggester{
			MakeContentUnitsSuggester(suggesterContext.DB, []string{consts.CT_BLOG_POST}),
			MakeContentUnitsSuggester(suggesterContext.DB, []string{consts.CT_ARTICLE}),
			MakeContentUnitsSuggester(suggesterContext.DB, []string{consts.CT_PUBLICATION}),
		}),
		// 6. Yeshivat + Mean.
		MakeRoundRobinSuggester([]Suggester{
			MakeContentUnitsSuggester(suggesterContext.DB, []string{consts.CT_FRIENDS_GATHERING}),
			MakeContentUnitsSuggester(suggesterContext.DB, []string{consts.CT_MEAL}),
		}),
		// 7. Clip.
		MakeContentUnitsSuggester(suggesterContext.DB, []string{
			consts.CT_CLIP,
		}),
	})}), suggesterContext)
}

func (f *Feed) More(r MoreRequest) ([]ContentItem, error) {
	return f.Suggester.More(r)
}

/*
func Merge(r MoreRequest, suggestions [][]ContentItem) ([]ContentItem, error) {
	mergedFeed := append([]ContentItem(nil), r.CurrentFeed...)
	for _, s := range suggestions {
		mergedFeed = append(mergedFeed, s...)
	}
	sort.SliceStable(mergedFeed, func(i, j int) bool {
		return mergedFeed[i].CreatedAt.After(mergedFeed[j].CreatedAt)
	})
	return mergedFeed[0:utils.MinInt(len(r.CurrentFeed)+r.MoreItems, len(mergedFeed))], nil
}
*/

func MakeDefaultSuggester(suggesterContext SuggesterContext) (Suggester, error) {
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
								},
								{
									"filter_selector": FPopularFilter
								}
							],
							"order_selector": OPopular
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
						},
								{
									"filter_selector": FPopularFilter
								}
					],
					"order_selector": OPopular
				},
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
									"filter_selector": FPopularFilter
								}
					],
					"order_selector": OPopular
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
						},
								{
									"filter_selector": FPopularFilter
								}
					],
					"order_selector": OPopular
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

	for filterName, filterValue := range FILTER_STRING_TO_VALUE {
		rootJSON = strings.ReplaceAll(rootJSON, filterName, fmt.Sprintf("%d", filterValue))
	}
	for orderName, orderValue := range ORDER_STRING_TO_VALUE {
		rootJSON = strings.ReplaceAll(rootJSON, orderName, fmt.Sprintf("%d", orderValue))
	}

	return MakeSuggesterFromJson(suggesterContext, rootJSON)

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

func MakeSuggesterFromJson(suggesterContext SuggesterContext, jsonStr string) (Suggester, error) {
	var spec SuggesterSpec
	if err := json.Unmarshal([]byte(jsonStr), &spec); err != nil {
		return nil, err
	}

	if s, err := MakeSuggesterFromName(suggesterContext, spec.Name); err != nil {
		return nil, err
	} else {
		if err := s.UnmarshalSpec(suggesterContext, spec); err != nil {
			return nil, err
		} else {
			return s, nil
		}
	}
}
