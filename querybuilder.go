package esq

import (
	"encoding/json"
)

type Request struct {
	query    *Query
	sortings []*Sorting
	from     int
	size     int
	timeout  string
}

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

type Sorting struct {
	byField      string
	order        string
	mode         string
	nestedPath   string
	nestedFilter *Query
}

func NewQuery() *Query {
	return &Query{}
}

func NewRequest(q *Query) *Request {
	return &Request{
		query:   q,
		from:    0,
		size:    25,
		timeout: "10s",
	}
}

func (r *Request) SetFrom(value int) *Request {
	r.from = value
	return r
}

func (r *Request) SetPageSize(value int) *Request {
	r.size = value
	return r
}

func (r *Request) Sort(sortings ...*Sorting) *Request {

	for _, srt := range sortings {
		r.sortings = append(r.sortings, srt)
	}
	return r
}

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

// If q is empty, a new compound 'bool' query is returned.
// If q is not empty and not a compound query, then q becomes the first entry
// in the returned compound query's 'must' entries, followed by subQuery items.
// If q is already a compound query, then subQuery items are appended to the 'must' clause of it.
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

func NestedQuery(path string, q *Query) *Query {
	nq := &nestedQuery{}
	nq.Nested.Path = path
	nq.Nested.Query = q
	return &Query{nestedQuery: nq}
}

// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-minimum-should-match.html
// it controls how many should clauses must be met by ES
// for a result to be displayed. E.g. the value is 1 then
// a document which does not match any 'should' clauses is not
// returned by the search.
func (q *Query) minimumShouldMatch(value interface{}) *Query {
	if q.isLeaf() {
		q.leafQuery.parameters["minimum_should_match"] = value
		return q
	}

	q.compoundQuery.Bool.MinimumShouldMatch = value
	return q
}

// Match queries accept text/numerics/dates
func Exists(field string) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "exists"
	q.leafQuery.parameters = map[string]interface{}{
		"field": field,
	}
	return q
}

// Match queries accept text/numerics/date/etc values
func Match(value interface{}, field string) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "match"
	q.leafQuery.parameters = map[string]interface{}{
		field: value,
	}
	return q
}

// The multi_match query builds on the match query to allow multi-field queries
func MultiMatch(value interface{}, fields ...string) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "multi_match"
	q.leafQuery.parameters = map[string]interface{}{
		"query":  value,
		"fields": fields,
	}
	return q
}

func Term(value interface{}, field string) *Query {
	q := &Query{leafQuery: &leafQuery{}}
	q.leafQuery.queryType = "term"
	q.leafQuery.parameters = map[string]interface{}{
		field: value,
	}
	return q
}

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

func ByField(fieldname string) *Sorting {
	if len(fieldname) == 0 {
		return nil
	}

	return &Sorting{
		byField: fieldname,
	}
}

func (s *Sorting) Desc() *Sorting {
	s.order = "desc"
	return s
}

func (s *Sorting) Asc() *Sorting {
	s.order = "desc"
	return s
}

func (s *Sorting) Min() *Sorting {
	s.mode = "min"
	return s
}

func (s *Sorting) Max() *Sorting {
	s.mode = "max"
	return s
}

func (s *Sorting) Sum() *Sorting {
	s.mode = "sum"
	return s
}

func (s *Sorting) Avg() *Sorting {
	s.mode = "avg"
	return s
}

func (s *Sorting) Median() *Sorting {
	s.mode = "median"
	return s
}

func (s *Sorting) Nested(path string) *Sorting {
	s.nestedPath = path
	return s
}

func (s *Sorting) Filter(q *Query) *Sorting {
	s.nestedFilter = q
	return s
}
