package sql

import (
	"fmt"
	"strings"
)

// A PostgreSQL query
type Query struct {
	baseQuery string
	filters   []QueryFilter
	limit     int
	offset    int
}

// A filter (where clause) for a PostgreSQL query
type QueryFilter = string

type FilterOperator = string

// Operators for where clauses
const (
	OpEq FilterOperator = "="
	OpLt FilterOperator = "<"
	OpGt FilterOperator = ">"
)

// Create a new Query
func NewQuery(baseQuery string) *Query {
	return &Query{baseQuery: baseQuery, filters: []QueryFilter{}}
}

// Create a new QueryFilter
// varchar values need to be enclosed in "'" manually
// Example:
// for varchar/string value: NewQueryFilter("table.field", "'myVarchar'", OpEq)
// for integer/numeric value: NewQueryFilter("table.field", "4.6", OpEq)
func NewQueryFilter(key string, value string, operator FilterOperator) QueryFilter {
	return fmt.Sprintf("%s %s %s", key, operator, value)
}

// Add a filter with operator "=" to the query
func (q *Query) AddEqFilter(key string, value string) {
	q.filters = append(q.filters, NewQueryFilter(key, value, OpEq))
}

// Add a filter with operator "<" to the query
func (q *Query) AddLtFilter(key string, value string) {
	q.filters = append(q.filters, NewQueryFilter(key, value, OpLt))
}

// Add a filter with operator ">" to the query
func (q *Query) AddGtFilter(key string, value string) {
	q.filters = append(q.filters, NewQueryFilter(key, value, OpGt))
}

// Add a limit to the query
func (q *Query) AddLimit(limit int) {
	q.limit = limit
}

// Add a offset to the query
func (q *Query) AddOffset(offset int) {
	q.offset = offset
}

// Render the complete query with all appended clauses
func (q *Query) String() string {
	fullQuery := q.baseQuery
	// cut group by clauses to add where clauses first
	groupByIndex := strings.LastIndex(fullQuery, "group by")
	groupClause := ""
	if groupByIndex >= 0 {
		groupClause = fullQuery[groupByIndex:]
		fullQuery = fullQuery[:groupByIndex]
	}

	// add where clauses and limit/offset
	for i, filter := range q.filters {
		if i == 0 {
			// first filter is appended with "WHERE"
			fullQuery = fmt.Sprintf("%s WHERE %s", fullQuery, filter)
			continue
		}
		// subsequent filters are appended with "AND"
		fullQuery = fmt.Sprintf("%s AND %s", fullQuery, filter)
	}
	if q.limit > 0 {
		fullQuery = fmt.Sprintf("%s LIMIT %d", fullQuery, q.limit)
	}
	if q.offset > 0 {
		fullQuery = fmt.Sprintf("%s OFFSET %d", fullQuery, q.offset)
	}

	if groupClause != "" {
		// re-add group by clause
		fullQuery = fmt.Sprintf("%s %s", fullQuery, groupClause)
	}

	return fullQuery
}
