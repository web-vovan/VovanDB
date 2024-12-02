package parser

import (
	"fmt"
)

type SelectQuery struct {
	Table     string
	Columns   []string
	Conditions []Condition
}

func (q SelectQuery) String() string {
	return fmt.Sprintf("Table: %s, Columns: %s, Condition: (%s)", q.Table, q.Columns, q.Conditions)
}

func (q SelectQuery) QueryType() string {
	return "SELECT"
}

func selectParser(p *Parser) (SQLQuery, error) {
	var columns []string
	var table string
	var conditions []Condition

	return SelectQuery{
		Table:     table,
		Columns:   columns,
		Conditions: conditions,
	}, nil
}
