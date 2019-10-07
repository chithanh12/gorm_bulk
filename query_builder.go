package gorm_bulk

import (
	"fmt"
	"strings"

	"bitbucket.org/hoiio/x-go/v3/xerrconv"
)

type queryBuilder struct{}

func QueryBuilder() *queryBuilder {
	return &queryBuilder{}
}

func (q *queryBuilder) BuildInsertQuery(tableName string, rows []interface{}) (string, []interface{}) {
	cols := Mapper().GetInsertColumns(rows[0])
	values := make([]string, 0, len(cols))
	for i := 0; i < len(cols); i++ {
		values = append(values, "?")
	}

	query := fmt.Sprintf("insert into `%v` (`%v`) values", tableName, strings.Join(cols, "`,`"))
	valuesHolder := fmt.Sprintf("(%v)", strings.Join(values, ",")) // (?,?,?)

	var params []interface{}
	rowValues := make([]string, 0)

	for _, row := range rows {
		rowMap, error := Mapper().GetValues(row)
		xerrconv.PanicIfError(xerrconv.DoPanic{Err: error})
		rowValues = append(rowValues, valuesHolder)

		for _, prop := range cols {
			params = append(params, rowMap[prop])
		}
	}

	query = fmt.Sprintf("%v %v", query, strings.Join(rowValues, ", "))
	return query, params
}

func (q *queryBuilder) BuildInsertOnDuplicateUpdate(tableName string, rows []interface{}) (string, []interface{}) {
	if len(rows) <= 0 {
		return "", nil
	}

	query, params := q.BuildInsertQuery(tableName, rows)
	updateCols := Mapper().GetUpdateColumns(rows[0])
	var update []string

	for _, col := range updateCols {
		update = append(update, fmt.Sprintf("`%v`=values(`%v`)", col, col))
	}

	onUpdateQuery := fmt.Sprintf("%v on duplicate key update %v", query, strings.Join(update, ","))
	return onUpdateQuery, params
}
