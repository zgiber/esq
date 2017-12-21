package esq

import (
	"encoding/json"
)

// Request can be executed with a search DSL, which includes the Query DSL, within its body.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-request-body.html
type Request struct {
	query    *Query
	sortings []*Sorting
	from     int
	size     int
	timeout  string
}

// Query element within the search request body allows to define a query.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-request-query.html
type Query struct {
	leafQuery     *leafQuery
	nestedQuery   *nestedQuery
	compoundQuery *compoundQuery
}

type compoundQuery struct {
	Bool struct {
		Must               []*Query    `json:"must,omitempty"`
		MustNot            []*Query    `json:"must_not,omitempty"`
		Should             []*Query    `json:"should,omitempty"`
		Filter             []*Query    `json:"filter,omitempty"`
		MinimumShouldMatch interface{} `json:"minimum_should_match,omitempty"`
	} `json:"bool,omitempty"`
}

type nestedQuery struct {
	Nested struct {
		Path  string `json:"path,omitempty"`
		Query *Query `json:"query,omitempty"`
	} `json:"nested,omitempty"`
}

type leafQuery struct {
	queryType  string
	parameters map[string]interface{}
}

// Sorting has parameters to define which order
// a query should retrieve the results.
type Sorting struct {
	byField      string
	order        string
	mode         string
	nestedPath   string
	nestedFilter *Query
}

// NewQuery returns an initialized *Query
// which does not hold any value. It is supposed
// to be used as a 'frame' for sub-queries.
func NewQuery() *Query {
	return &Query{}
}

// NewRequest constructs a valid request from a Query.
func NewRequest(q *Query) *Request {
	return &Request{
		query: q,
	}
}

// SetFrom allows to retrieve hits from a certain offset. Defaults to 0.
func (r *Request) SetFrom(value int) *Request {
	r.from = value
	return r
}

// SetPageSize sets number of hits to return.
// Defaults to 10. If you do not care about getting some hits back
// but only about the number of matches and/or aggregations,
// setting the value to 0 will help performance.
func (r *Request) SetPageSize(value int) *Request {
	r.size = value
	return r
}

// Sort applies the given sortings to the Request.
func (r *Request) Sort(sortings ...*Sorting) *Request {

	for _, srt := range sortings {
		r.sortings = append(r.sortings, srt)
	}
	return r
}

// MarshalJSON serializes the Request in the expected format.
func (r *Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Query    *Query     `json:"query,omitempty"`
		Sortings []*Sorting `json:"sort,omitempty"`
		From     int        `json:"from,omitempty"`
		Size     int        `json:"size,omitempty"`
		Timeout  string     `json:"timeout,omitempty"`
	}{
		Query:    r.query,
		Sortings: r.sortings,
		From:     r.from,
		Size:     r.size,
		Timeout:  r.timeout,
	})
}

func (q *Query) isLeaf() bool {
	return q.leafQuery != nil
}

func (q *Query) isNested() bool {
	return q.nestedQuery != nil
}

// MarshalJSON serializes the Query in the expected format.
func (q *Query) MarshalJSON() ([]byte, error) {
	if q.isLeaf() {
		return json.Marshal(q.leafQuery)
	}

	if q.isNested() {
		return json.Marshal(q.nestedQuery)
	}

	return json.Marshal(q.compoundQuery)
}

func (lq *leafQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		lq.queryType: lq.parameters,
	})
}

// Must query must appear in matching documents and will contribute to the score
// The returned query is always a compound query. If the original query is not a
// compound query, it will become the first subquery in the returned query.
func (q *Query) Must(subQuery ...*Query) *Query {
	if q.compoundQuery == nil {
		q.compoundQuery = &compoundQuery{}
	}

	if q.isLeaf() {
		q.compoundQuery.Bool.Must = append(q.compoundQuery.Bool.Must, &Query{leafQuery: q.leafQuery})
		q.leafQuery = nil
	}

	if q.isNested() {
		q.compoundQuery.Bool.Must = append(q.compoundQuery.Bool.Must, &Query{nestedQuery: q.nestedQuery})
		q.nestedQuery = nil
	}

	q.compoundQuery.Bool.Must = append(q.compoundQuery.Bool.Must, subQuery...)
	return q
}

// MustNot query must not appear in the matching documents.
// Clauses are executed in filter context meaning that scoring is ignored
// and clauses are considered for caching. Because scoring is ignored,
// a score of 0 for all documents is returned.
//
// The returned query is always a compound query. If the original query is not a
// compound query, it will become the first subquery in the returned query.
func (q *Query) MustNot(subQuery ...*Query) *Query {
	if q.compoundQuery == nil {
		q.compoundQuery = &compoundQuery{}
	}

	if q.isLeaf() {
		q.compoundQuery.Bool.MustNot = append(q.compoundQuery.Bool.MustNot, &Query{leafQuery: q.leafQuery})
		q.leafQuery = nil
	}

	if q.isNested() {
		q.compoundQuery.Bool.MustNot = append(q.compoundQuery.Bool.MustNot, &Query{nestedQuery: q.nestedQuery})
		q.nestedQuery = nil
	}

	q.compoundQuery.Bool.MustNot = append(q.compoundQuery.Bool.MustNot, subQuery...)
	return q
}

// Should query should appear in the matching document.
// If the bool query is in a query context and has a must or filter clause
// then a document will match the bool query even if none of the should queries match.
// In this case these clauses are only used to influence the score.
// If the bool query is a filter context or has neither must or filter
// then at least one of the should queries must match a document for it to match the bool query.
// This behavior may be explicitly controlled by settings the minimum_should_match parameter.
func (q *Query) Should(subQuery ...*Query) *Query {
	if q.compoundQuery == nil {
		q.compoundQuery = &compoundQuery{}
	}

	if q.isLeaf() {
		q.compoundQuery.Bool.Should = append(q.compoundQuery.Bool.Should, &Query{leafQuery: q.leafQuery})
		q.leafQuery = nil
	}

	if q.isNested() {
		q.compoundQuery.Bool.Should = append(q.compoundQuery.Bool.Should, &Query{nestedQuery: q.nestedQuery})
		q.nestedQuery = nil
	}

	q.compoundQuery.Bool.Should = append(q.compoundQuery.Bool.Should, subQuery...)
	return q
}

// Filter clause (query) must appear in matching documents.
// However unlike must the score of the query will be ignored.
// Filter clauses are executed in filter context, meaning that
// scoring is ignored and clauses are considered for caching.
func (q *Query) Filter(subQuery ...*Query) *Query {
	if q.compoundQuery == nil {
		q.compoundQuery = &compoundQuery{}
	}

	if q.isLeaf() {
		q.compoundQuery.Bool.Filter = append(q.compoundQuery.Bool.Filter, &Query{leafQuery: q.leafQuery})
		q.leafQuery = nil
	}

	if q.isNested() {
		q.compoundQuery.Bool.Filter = append(q.compoundQuery.Bool.Filter, &Query{nestedQuery: q.nestedQuery})
		q.nestedQuery = nil
	}

	q.compoundQuery.Bool.Filter = append(q.compoundQuery.Bool.Filter, subQuery...)
	return q
}

// NestedQuery allows to query nested objects / docs (see nested mapping).
// The query is executed against the nested objects / docs
// as if they were indexed as separate docs (they are, internally) and resulting
// in the root parent doc (or parent nested mapping).
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-nested-query.html
func NestedQuery(path string, q *Query) *Query {
	nq := &nestedQuery{}
	nq.Nested.Path = path
	nq.Nested.Query = q
	return &Query{nestedQuery: nq}
}

// minimumShouldMatch controls how many should clauses must be met by ES
// for a result to be displayed. E.g. the value is 1 then
// a document which does not match any 'should' clauses is not
// returned by the search.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-minimum-should-match.html
func (q *Query) minimumShouldMatch(value interface{}) *Query {
	if q.isLeaf() {
		q.leafQuery.parameters["minimum_should_match"] = value
		return q
	}

	q.compoundQuery.Bool.MinimumShouldMatch = value
	return q
}

// Exists returns documents that have at least one non-null value in the original field.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-exists-query.html
func Exists(field string) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "exists"
	q.leafQuery.parameters = map[string]interface{}{
		"field": field,
	}
	return q
}

// Match queries accept text/numerics/date/etc values.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-match-query.html
func Match(value interface{}, field string) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "match"
	q.leafQuery.parameters = map[string]interface{}{
		field: value,
	}
	return q
}

// MatchPhrase query analyzes the text and creates a phrase query out of the analyzed text.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-match-query-phrase.html
func MatchPhrase(value string, field string) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "match_phrase"
	q.leafQuery.parameters = map[string]interface{}{
		field: value,
	}
	return q
}

// MultiMatch query builds on the match query to allow multi-field queries.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-multi-match-query.html
func MultiMatch(value interface{}, fields ...string) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "multi_match"
	q.leafQuery.parameters = map[string]interface{}{
		"query":  value,
		"fields": fields,
	}
	return q
}

// Term query finds documents that contain the exact term specified in the inverted index.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-term-query.html
func Term(value interface{}, field string, boost float64) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "term"
	q.leafQuery.parameters = map[string]interface{}{
		field: value,
	}

	if boost > 0.0 {
		q.leafQuery.parameters["boost"] = boost
	}

	return q
}

// Terms query Filters documents that have fields that match any of the provided terms (not analyzed).
// Value is expected to be a slice of values.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-terms-query.html
func Terms(value interface{}, field string, boost float64) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "terms"
	q.leafQuery.parameters = map[string]interface{}{
		field: value,
	}

	if boost > 0.0 {
		q.leafQuery.parameters["boost"] = boost
	}

	return q
}

// Range matches documents with fields that have terms within a certain range.
// The type of the Lucene query depends on the field type,
// for string fields, the TermRangeQuery,
// while for number/date fields, the query is a NumericRangeQuery.
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-range-query.html
func Range(gte, lte interface{}, field string) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "range"
	q.leafQuery.parameters = map[string]interface{}{
		field: struct {
			Gte interface{} `json:"gte,omitempty"`
			Lte interface{} `json:"lte,omitempty"`
		}{
			gte,
			lte,
		},
	}
	return q
}

// MarshalJSON serializes the Sorting in the expected format.
func (s *Sorting) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("{}"), nil
	}

	return json.Marshal(map[string]interface{}{
		s.byField: struct {
			Order        string `json:"order,omitempty"`
			NestedPath   string `json:"nested_path,omitempty"`
			NestedFilter *Query `json:"nested_filter,omitempty"`
		}{s.order, s.nestedPath, s.nestedFilter},
	})
}

// ByField returns sorting by the given field's value.
func ByField(fieldname string) *Sorting {
	if len(fieldname) == 0 {
		return nil
	}

	return &Sorting{
		byField: fieldname,
	}
}

// Desc returns sorting with descending order.
func (s *Sorting) Desc() *Sorting {
	s.order = "desc"
	return s
}

// Asc returns sorting with ascending order.
func (s *Sorting) Asc() *Sorting {
	s.order = "desc"
	return s
}

// Min returns sorting with ascending order
// when sorting by array or multi-valued fields.
func (s *Sorting) Min() *Sorting {
	s.mode = "min"
	return s
}

// Max returns a sorting which picks the highest value
// when sorting by array or multi-valued fields.
func (s *Sorting) Max() *Sorting {
	s.mode = "max"
	return s
}

// Sum returns a sorting which uses the sum of all values
// when sorting by array or multi-valued fields.
func (s *Sorting) Sum() *Sorting {
	s.mode = "sum"
	return s
}

// Avg returns a sorting which uses the average of all values
// when sorting by array or multi-valued fields.
func (s *Sorting) Avg() *Sorting {
	s.mode = "avg"
	return s
}

// Median returns a sorting which uses the median of all values
// when sorting by array or multi-valued fields.
func (s *Sorting) Median() *Sorting {
	s.mode = "median"
	return s
}

// Nested returns a sorting which uses a nested field.
// The path parameter is the path for the nested field.
// If filter is not nil, only matching values are considered
// for the sorting.
func (s *Sorting) Nested(path string, filter *Query) *Sorting {
	s.nestedPath = path
	if filter != nil {
		s.nestedFilter = filter
	}
	return s
}
