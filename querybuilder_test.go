package esq

import (
	"encoding/json"
	"testing"

	"github.com/onsi/gomega"
)

func Test_search_byToken(t *testing.T) {
	gomega.RegisterTestingT(t)
	r := NewRequest(
		NewQuery().
			Must(Term("form1234", "form_id.keyword"))).
		Sort(ByField("form1234").Desc().Nested("answers").
			Filter(NewQuery().Must(
				Term("stuff", "nested_id.keyword")))).SetFrom(0).SetPageSize(25)

	q, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte(`
		{
			"query": {
				"bool": {
					"must": [
						{
							"term": {
								"form_id.keyword": "form1234"
							}
						}
					]
				}
			},
			"size": 25,
			"timeout": "10s"
		}
		`)

	gomega.Expect(q).Should(gomega.MatchJSON(expected))

}

// func Test_search_byLandingDateRange(t *testing.T) {
// 	gomega.RegisterTestingT(t)

// 	since, _ := time.Parse(time.RFC3339, "2017-10-18T09:48:12Z")
// 	until, _ := time.Parse(time.RFC3339, "2017-10-18T10:48:12Z")

// 	s := newSearch("form1234").byLandingDateRange(since, until)
// 	q, err := json.Marshal(s)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := []byte(`
// 		{
// 			"query": {
// 				"bool": {
// 					"must": [
// 						{
// 							"term": {
// 								"form_id.keyword": "form1234"
// 							}
// 						},
// 						{
// 							"range": {
// 								"landed_at": {
// 									"gte": "2017-10-18T09:48:12Z",
// 									"lte": "2017-10-18T10:48:12Z"
// 								}
// 							}
// 						}
// 					]
// 				}
// 			},
// 			"size": 25,
// 			"timeout": "10s"
// 		}
// 		`)

// 	gomega.Expect(q).Should(gomega.MatchJSON(expected))

// }

// func Test_search_bySubmitDateRange(t *testing.T) {
// 	gomega.RegisterTestingT(t)

// 	since, _ := time.Parse(time.RFC3339, "2017-10-18T09:48:12Z")
// 	until, _ := time.Parse(time.RFC3339, "2017-10-18T10:48:12Z")

// 	s := newSearch("form1234").bySubmitDateRange(since, until)
// 	q, err := json.Marshal(s)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := []byte(`
// 			{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							},
// 							{
// 								"range": {
// 									"submitted_at": {
// 										"gte": "2017-10-18T09:48:12Z",
// 										"lte": "2017-10-18T10:48:12Z"
// 									}
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"size": 25,
// 				"timeout": "10s"
// 			}
// 			`)

// 	gomega.Expect(q).Should(gomega.MatchJSON(expected))
// }

// func Test_search_inAnswers(t *testing.T) {
// 	gomega.RegisterTestingT(t)
// 	date, _ := time.Parse(time.RFC3339, "2017-10-18T09:48:12Z")

// 	tests := []struct {
// 		name string
// 		s    *search
// 		want string
// 	}{
// 		{
// 			"search for text", // also searches within multiple choices
// 			newSearch("form1234").inAnswers("hello world"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						],
// 						"should": [
// 							{
// 								"nested": {
// 									"path": "answers",
// 									"query": {
// 										"multi_match": {
// 											"fields": [
// 												"answers.text",
// 												"answers.url",
// 												"answers.file_url",
// 												"answers.email",
// 												"answers.choices.label"
// 											],
// 											"query": "hello world"
// 										}
// 									}
// 								}
// 							},
// 							{
// 								"match": {
// 									"_all": "hello world"
// 								}
// 							}
// 						],
// 						"minimum_should_match": 1
// 					}
// 				},
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		}, {
// 			"search for number",
// 			newSearch("form1234").inAnswers(1234),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						],
// 						"should": [
// 							{
// 								"nested": {
// 									"path": "answers",
// 									"query": {
// 										"multi_match": {
// 											"fields": [
// 												"answers.text",
// 												"answers.url",
// 												"answers.file_url",
// 												"answers.email",
// 												"answers.choices.label",
// 												"answers.number"
// 											],
// 											"query": 1234
// 										}
// 									}
// 								}
// 							},
// 							{
// 								"match": {
// 									"_all": 1234
// 								}
// 							}
// 						],
// 						"minimum_should_match": 1
// 					}
// 				},
// 				"size": 25,
// 				"timeout": "10s"
// 			}`},
// 		{
// 			"search for boolean",
// 			newSearch("form1234").inAnswers(true),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						],
// 						"should": [
// 							{
// 								"nested": {
// 									"path": "answers",
// 									"query": {
// 										"multi_match": {
// 											"fields": [
// 												"answers.text",
// 												"answers.url",
// 												"answers.file_url",
// 												"answers.email",
// 												"answers.choices.label",
// 												"answers.boolean"
// 											],
// 											"query": true
// 										}
// 									}
// 								}
// 							},
// 							{
// 								"match": {
// 									"_all": true
// 								}
// 							}
// 						],
// 						"minimum_should_match": 1
// 					}
// 				},
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		}, {
// 			"search for date",
// 			newSearch("form1234").inAnswers(date),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						],
// 						"should": [
// 							{
// 								"nested": {
// 									"path": "answers",
// 									"query": {
// 										"multi_match": {
// 											"fields": [
// 												"answers.date"
// 											],
// 											"query": "2017-10-18T09:48:12Z"
// 										}
// 									}
// 								}
// 							},
// 							{
// 								"match": {
// 									"_all": "2017-10-18T09:48:12Z"
// 								}
// 							}
// 						],
// 						"minimum_should_match": 1
// 					}
// 				},
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 	}

// 	for _, test := range tests {
// 		q, err := json.Marshal(test.s)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		gomega.Expect(q).Should(gomega.MatchJSON(test.want))
// 	}
// }

// func Test_search_inAnswerField(t *testing.T) {
// 	// Not implemented yet
// 	// The use case would be when search is restricted to a given answer id
// 	// e.g. search "awesome" inside the answers for "Why do you like our product?" question
// 	// but omit submissions where "awesome" is in an answer for a different question.
// }

// func Test_setCompleted(t *testing.T) {
// 	gomega.RegisterTestingT(t)
// 	s := newSearch("form1234").setCompleted(true)
// 	q, err := json.Marshal(s)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := []byte(`{
// 		"query": {
// 			"bool": {
// 				"must": [
// 					{
// 						"term": {
// 							"form_id.keyword": "form1234"
// 						}
// 					},
// 					{
// 						"exists": {
// 							"field": "submitted_at"
// 						}
// 					}
// 				]
// 			}
// 		},
// 		"size": 25,
// 		"timeout": "10s"
// 	}`)

// 	gomega.Expect(q).Should(gomega.MatchJSON(expected))
// }

// func Test_setFrom(t *testing.T) {
// 	gomega.RegisterTestingT(t)
// 	s := newSearch("form1234").setCompleted(true)
// 	q, err := json.Marshal(s)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := []byte(`{
// 		"query": {
// 			"bool": {
// 				"must": [
// 					{
// 						"term": {
// 							"form_id.keyword": "form1234"
// 						}
// 					},
// 					{
// 						"exists": {
// 							"field": "submitted_at"
// 						}
// 					}
// 				]
// 			}
// 		},
// 		"size": 25,
// 		"timeout": "10s"
// 	}`)

// 	gomega.Expect(q).Should(gomega.MatchJSON(expected))
// }

// func Test_setPageSize(t *testing.T) {
// 	gomega.RegisterTestingT(t)
// 	s := newSearch("form1234").setPageSize(100)
// 	q, err := json.Marshal(s)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := []byte(`{
// 		"query": {
// 		  "bool": {
// 			"must": [
// 			  {
// 				"term": {
// 				  "form_id.keyword": "form1234"
// 				}
// 			  }
// 			]
// 		  }
// 		},
// 		"size": 100,
// 		"timeout": "10s"
// 	  }`)

// 	gomega.Expect(q).Should(gomega.MatchJSON(expected))
// }

// func Test_byToken(t *testing.T) {
// 	gomega.RegisterTestingT(t)
// 	s := newSearch("form1234").byToken("token6789")
// 	q, err := json.Marshal(s)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expected := []byte(`{
// 		"query": {
// 		  "bool": {
// 			"must": [
// 			  {
// 				"term": {
// 				  "form_id.keyword": "form1234"
// 				}
// 			  },
// 			  {
// 				"term": {
// 				  "token.keyword": "token6789"
// 				}
// 			  }
// 			]
// 		  }
// 		},
// 		"size": 25,
// 		"timeout": "10s"
// 	  }`)

// 	gomega.Expect(q).Should(gomega.MatchJSON(expected))
// }
// func Test_query_mustNot(t *testing.T) {
// 	gomega.RegisterTestingT(t)

// 	tests := []struct {
// 		name string
// 		q    *query
// 		want string
// 	}{
// 		{
// 			"from new query",
// 			newQuery().mustNot(newTermQuery("form1234", "form_id.keyword")),
// 			`{
// 				"bool": {
// 					"must_not": [
// 						{
// 							"term": {
// 								"form_id.keyword": "form1234"
// 							}
// 						}
// 					]
// 				}
// 			}`,
// 		}, {
// 			"from existing leaf query",
// 			newTermQuery("form1234", "form_id.keyword").mustNot(),
// 			`{
// 				"bool": {
// 					"must_not": [
// 						{
// 							"term": {
// 								"form_id.keyword": "form1234"
// 							}
// 						}
// 					]
// 				}
// 			}`,
// 		},
// 		{
// 			"from existing compound query",
// 			newTermQuery("form1234", "form_id.keyword").mustNot(
// 				newMatchQuery(true, "some.field"),
// 				newMatchQuery("blue", "some_other.field"),
// 			),
// 			`{
// 				"bool": {
// 					"must_not": [
// 						{
// 							"term": {
// 								"form_id.keyword": "form1234"
// 							}
// 						},
// 						{
// 							"match": {
// 								"some.field": true
// 							}
// 						},
// 						{
// 							"match": {
// 								"some_other.field": "blue"
// 							}
// 						}
// 					]
// 				}
// 			}`,
// 		},
// 		{
// 			"from existing nested query",
// 			newNestedQuery("answers", newMatchQuery(true, "answers.boolean")).mustNot(
// 				newMatchQuery("blue", "some_other.field"),
// 			),
// 			`{
// 				"bool": {
// 					"must_not": [
// 						{
// 							"nested": {
// 								"path": "answers",
// 								"query": {
// 									"match": {
// 										"answers.boolean": true
// 									}
// 								}
// 							}
// 						},
// 						{
// 							"match": {
// 								"some_other.field": "blue"
// 							}
// 						}
// 					]
// 				}
// 			}`,
// 		},
// 	}

// 	for _, test := range tests {
// 		q, err := json.Marshal(test.q)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		gomega.Expect(q).Should(gomega.MatchJSON(test.want))
// 	}
// }

// func Test_newSort(t *testing.T) {
// 	gomega.RegisterTestingT(t)

// 	tests := []struct {
// 		name string
// 		s    *search
// 		want string
// 	}{
// 		{
// 			"sort by text field",
// 			newSearch("form1234").sortBy("textfield_123", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"answers.text.keyword": {
// 							"order": "desc",
// 							"nested_path": "answers",
// 							"nested_filter": {
// 								"term": {
// 									"answers.field.id.keyword": "123"
// 								}
// 							}
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by number field",
// 			newSearch("form1234").sortBy("number_123", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"answers.number": {
// 							"order": "desc",
// 							"nested_path": "answers",
// 							"nested_filter": {
// 								"term": {
// 									"answers.field.id.keyword": "123"
// 								}
// 							}
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by date field",
// 			newSearch("form1234").sortBy("date_123", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"answers.date": {
// 							"order": "desc",
// 							"nested_path": "answers",
// 							"nested_filter": {
// 								"term": {
// 									"answers.field.id.keyword": "123"
// 								}
// 							}
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by boolean field",
// 			newSearch("form1234").sortBy("yesno_123", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"answers.boolean": {
// 							"order": "desc",
// 							"nested_path": "answers",
// 							"nested_filter": {
// 								"term": {
// 									"answers.field.id.keyword": "123"
// 								}
// 							}
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by landing date",
// 			newSearch("form1234").sortBy("landed_at", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"landed_at": {
// 							"order": "desc"
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by submit date",
// 			newSearch("form1234").sortBy("submitted_at", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"submitted_at": {
// 							"order": "desc"
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by platform",
// 			newSearch("form1234").sortBy("platform", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"metadata.platform.keyword": {
// 							"order": "desc"
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by referer",
// 			newSearch("form1234").sortBy("referer", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"metadata.referer.keyword": {
// 							"order": "desc"
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by ip",
// 			newSearch("form1234").sortBy("ip", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"metadata.network_id.keyword": {
// 							"order": "desc"
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by ip",
// 			newSearch("form1234").sortBy("ip", "desc"),
// 			`{
// 				"query": {
// 					"bool": {
// 						"must": [
// 							{
// 								"term": {
// 									"form_id.keyword": "form1234"
// 								}
// 							}
// 						]
// 					}
// 				},
// 				"sort": [
// 					{
// 						"metadata.network_id.keyword": {
// 							"order": "desc"
// 						}
// 					}
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			}`,
// 		},
// 		{
// 			"sort by score",
// 			newSearch("form1234").sortBy("score", "desc"),
// 			`{
// 				"query": {
// 				  "bool": {
// 					"must": [
// 					  {
// 						"term": {
// 						  "form_id.keyword": "form1234"
// 						}
// 					  }
// 					]
// 				  }
// 				},
// 				"sort": [
// 				  {
// 					"calculated.score": {
// 					  "order": "desc"
// 					}
// 				  }
// 				],
// 				"size": 25,
// 				"timeout": "10s"
// 			  }`,
// 		},
// 		{
// 			"sort by text field",
// 			newSearch("form1234").sortBy("payment_123_price", "desc"),
// 			`{
// 					"query": {
// 						"bool": {
// 							"must": [
// 								{
// 									"term": {
// 										"form_id.keyword": "form1234"
// 									}
// 								}
// 							]
// 						}
// 					},
// 					"sort": [
// 						{
// 							"answers.payment.amount.keyword": {
// 								"order": "desc",
// 								"nested_path": "answers",
// 								"nested_filter": {
// 									"term": {
// 										"answers.field.id.keyword": "123"
// 									}
// 								}
// 							}
// 						}
// 					],
// 					"size": 25,
// 					"timeout": "10s"
// 				}`,
// 		},
// 	}

// 	for _, test := range tests {
// 		s, err := json.Marshal(test.s)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		gomega.Expect(s).Should(gomega.MatchJSON(test.want))
// 	}

// }
