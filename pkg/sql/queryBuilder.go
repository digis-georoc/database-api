package sql

import (
	"fmt"
	"strings"
)

// A PostgreSQL query
// param: baseQuery - the sql query text to be populated with parameters and filter options
// param: filterValues - list of filter values to be safely applied as parameters later
// param: limit - limit param for basic pagination
// param: offset - offset param for basic pagination
type Query struct {
	baseQuery    string
	filterValues []interface{}
	limit        int
	offset       int
}

// A filter (e.g. where clause) for a PostgreSQL query
type QueryFilter = string

type FilterOperator = string

type FilterJunctor = string

// Operators for sql filters
const (
	OpEq        FilterOperator = "="
	OpLt        FilterOperator = "<"
	OpGt        FilterOperator = ">"
	OpIn        FilterOperator = "IN"
	OpLike      FilterOperator = "LIKE"
	OpBetween   FilterOperator = "BETWEEN"
	OpInPolygon FilterOperator = "INPOLYGON"

	OpAnd   FilterJunctor = "AND"
	OpOr    FilterJunctor = "OR"
	OpWhere FilterJunctor = "WHERE"

	SEPARATOR string = ","

	MIN_LOWER_BOUND = "0"
	MAX_UPPER_BOUND = "9999999"
)

var OperatorMap map[string]FilterOperator = map[string]FilterOperator{
	"eq":  OpEq,
	"in":  OpIn,
	"gt":  OpGt,
	"lt":  OpLt,
	"lk":  OpLike,
	"btw": OpBetween,
}

// Create a new Query
func NewQuery(baseQuery string) *Query {
	return &Query{baseQuery: baseQuery, filterValues: []interface{}{}}
}

// Create a new QueryFilter
func NewQueryFilter(key string, value string, operator FilterOperator, junctor FilterJunctor) QueryFilter {
	return fmt.Sprintf("%s %s %s %s", junctor, key, operator, value)
}

// Add a filter depending on the operator
func (q *Query) AddFilter(key string, value string, operator FilterOperator, junctor FilterJunctor) {
	switch operator {
	case OpEq:
		q.AddEqFilter(key, value, junctor)
	case OpGt:
		q.AddGtFilter(key, value, junctor)
	case OpLt:
		q.AddLtFilter(key, value, junctor)
	case OpIn:
		q.AddInFilter(key, value, junctor)
	case OpLike:
		q.AddLikeFilter(key, value, junctor)
	case OpBetween:
		q.AddBetweenFilter(key, value, junctor)
	case OpInPolygon:
		q.AddInPolygonFilter(key, value, junctor)
	}
}

// Add a filter with operator "=" to the query
func (q *Query) AddEqFilter(key string, value string, junctor FilterJunctor) {
	q.filterValues = append(q.filterValues, value)
	placeholder := fmt.Sprintf("$%d", len(q.filterValues))
	filterString := NewQueryFilter(key, placeholder, OpEq, junctor)
	q.baseQuery = fmt.Sprintf("%s %s", q.baseQuery, filterString)
}

// Add a filter with operator "<" to the query
func (q *Query) AddLtFilter(key string, value string, junctor FilterJunctor) {
	q.filterValues = append(q.filterValues, value)
	placeholder := fmt.Sprintf("$%d", len(q.filterValues))
	filterString := NewQueryFilter(key, placeholder, OpLt, junctor)
	q.baseQuery = fmt.Sprintf("%s %s", q.baseQuery, filterString)
}

// Add a filter with operator ">" to the query
func (q *Query) AddGtFilter(key string, value string, junctor FilterJunctor) {
	q.filterValues = append(q.filterValues, value)
	placeholder := fmt.Sprintf("$%d", len(q.filterValues))
	filterString := NewQueryFilter(key, placeholder, OpGt, junctor)
	q.baseQuery = fmt.Sprintf("%s %s", q.baseQuery, filterString)
}

// Add a filter to check if a field is in a set of  values
// param key: the table.field to check against
// param values: comma-separated string of values
func (q *Query) AddInFilter(key string, values string, junctor FilterJunctor) {
	valueSplit := strings.Split(values, SEPARATOR)
	placeholderString := ""
	for i, val := range valueSplit {
		q.filterValues = append(q.filterValues, val)
		if i == 0 {
			// add first value with bracket
			placeholderString = fmt.Sprintf("($%d", len(q.filterValues))
			continue
		}
		placeholderString = fmt.Sprintf("%s,$%d", placeholderString, len(q.filterValues))
	}
	// add closing bracket
	placeholderString = fmt.Sprintf("%s)", placeholderString)

	filterString := NewQueryFilter(key, placeholderString, OpIn, junctor)
	q.baseQuery = fmt.Sprintf("%s %s", q.baseQuery, filterString)
}

// Add a filter with operator "LIKE" to the query
func (q *Query) AddLikeFilter(key string, value string, junctor FilterJunctor) {
	q.filterValues = append(q.filterValues, value)
	placeholder := fmt.Sprintf("$%d", len(q.filterValues))
	filterString := NewQueryFilter(key, placeholder, OpLike, junctor)
	q.baseQuery = fmt.Sprintf("%s %s", q.baseQuery, filterString)
}

// Add a filter for value range to the query
func (q *Query) AddBetweenFilter(key string, value string, junctor FilterJunctor) {
	lower, upper, _ := strings.Cut(value, SEPARATOR)
	if upper == "" {
		// assume max value
		upper = MAX_UPPER_BOUND
	}
	if lower == "" {
		// assume min value
		lower = MIN_LOWER_BOUND
	}
	q.filterValues = append(q.filterValues, lower)
	lowerPlaceholder := fmt.Sprintf("$%d", len(q.filterValues))
	q.filterValues = append(q.filterValues, upper)
	upperPlaceholder := fmt.Sprintf("$%d", len(q.filterValues))
	operatorString := fmt.Sprintf("%s AND %s", lowerPlaceholder, upperPlaceholder)
	filterString := NewQueryFilter(key, operatorString, OpBetween, junctor)
	q.baseQuery = fmt.Sprintf("%s %s", q.baseQuery, filterString)
}

// Add a filter with to check for polygon to the query
func (q *Query) AddInPolygonFilter(key string, value string, junctor FilterJunctor) {
	q.filterValues = append(q.filterValues, fmt.Sprintf("POLYGON(%s)", value))
	placeholder := fmt.Sprintf("$%d", len(q.filterValues))
	filterString := fmt.Sprintf("%s ST_WITHIN(%s, ST_GEOMETRYFROMTEXT(%s))", junctor, key, placeholder)
	q.baseQuery = fmt.Sprintf("%s %s", q.baseQuery, filterString)
}

// Add a subquery / sql block to the query
// Do not enter user-provided values here as they are not sanitized.
// For user values, use filters
func (q *Query) AddSQLBlock(sql string) {
	q.baseQuery = fmt.Sprintf("%s %s", q.baseQuery, sql)
}

// Add a limit to the query
func (q *Query) AddLimit(limit int) {
	q.limit = limit
}

// Add a offset to the query
func (q *Query) AddOffset(offset int) {
	q.offset = offset
}

// Retrieve the list of filterValues
func (q *Query) GetFilterValues() []interface{} {
	return q.filterValues
}

// Retrieve the full query string, including the limit and offset if set
func (q *Query) GetQueryString() string {
	if q.limit > 0 {
		q.baseQuery = fmt.Sprintf("%s LIMIT %d", q.baseQuery, q.limit)
	}
	if q.offset > 0 {
		q.baseQuery = fmt.Sprintf("%s OFFSET %d", q.baseQuery, q.offset)
	}
	return q.baseQuery
}
