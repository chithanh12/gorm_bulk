package gorm_bulk

import (
	"errors"
	"fmt"
	"strings"
)

type queryBuilder struct{}
type QueryResult struct {
	Query      string
	Parameters []interface{}
}

func QueryBuilder() *queryBuilder {
	return &queryBuilder{}
}

func (q *queryBuilder) BuildInsertQuery(tableName string, rows []interface{}, ignoreCols ...string) (*QueryResult, error) {
	cols, err := Mapper().GetColumns(rows[0])
	if err != nil {
		return nil, err
	}

	if len(ignoreCols) > 0 {
		for _, c := range ignoreCols {
			cols = removeElement(cols, c)
		}
	}
	values := make([]string, 0, len(cols))
	for i := 0; i < len(cols); i++ {
		values = append(values, "?")
	}

	query := fmt.Sprintf("insert into `%v` (`%v`) values", tableName, strings.Join(cols, "`,`"))
	valuesHolder := fmt.Sprintf("(%v)", strings.Join(values, ",")) // (?,?,?)

	var params []interface{}
	rowValues := make([]string, 0)

	for _, row := range rows {
		rowMap, err := Mapper().GetValues(row)
		if err != nil {
			return nil, err
		}
		rowValues = append(rowValues, valuesHolder)

		for _, prop := range cols {
			params = append(params, rowMap[prop])
		}
	}

	query = fmt.Sprintf("%v %v", query, strings.Join(rowValues, ", "))

	return &QueryResult{
		Query:      query,
		Parameters: params,
	}, nil
}

func (q *queryBuilder) BuildInsertOnDuplicateUpdate(tableName string, rows []interface{}, ignoreColumns ...string) (*QueryResult, error) {
	if len(rows) <= 0 {
		return nil, errors.New("Empty rows parameters")
	}

	insertResult, err := q.BuildInsertQuery(tableName, rows)
	if err != nil {
		return nil, err
	}
	updateCols, err := Mapper().GetColumns(rows[0])
	if err != nil {
		return nil, err
	}

	var update []string

	for _, col := range updateCols {
		update = append(update, fmt.Sprintf("`%v`=values(`%v`)", col, col))
	}

	insertResult.Query = fmt.Sprintf("%v on duplicate key update %v", insertResult.Query, strings.Join(update, ","))
	return insertResult, nil
}

func removeElement(items []string, el string) []string {
	idx := -1
	for i, item := range items {
		if item == el {
			idx = i
			break
		}
	}

	if idx >= 0 {
		return append(items[:idx], items[idx+1:]...)
	}

	return items
}
