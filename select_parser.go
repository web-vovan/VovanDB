package vovanDB

import (
	"fmt"
)

type SelectQuery struct {
	Table      string
	Columns    []string
	Conditions []Condition
}

func (q SelectQuery) String() string {
	return fmt.Sprintf("Table: %s\nColumns: %s\nConditions: \n%s", q.Table, q.Columns, q.Conditions)
}

func (q SelectQuery) QueryType() string {
	return "SELECT"
}

func (q SelectQuery) showAllColumns() bool {
	return q.Columns[0] == "*"
}

func selectParser(p *Parser) (SQLQuery, error) {
	var columns []string
	var table string
	var conditions []Condition

	if p.isSelectQuery() {
		p.next()
	} else {
		return nil, fmt.Errorf("неверная структура запроса1")
	}

	// Колонки
	for p.isIdentifier() || p.isComma() || p.isAsterisk() {
		if p.isComma() {
			p.next()
			continue
		}

		columns = append(columns, p.next().Value)
	}

	if p.isEnd() || p.current().Value != "FROM" {
		return nil, fmt.Errorf("неверная структура запроса2")
	}

	p.next()

	// Таблица
	if !p.isIdentifier() {
		return nil, fmt.Errorf("неверная структура запроса3")
	}

	table = p.next().Value

	if p.isSemicolon() {
		p.next()

		if p.isEnd() {
			return SelectQuery{
				Table:   table,
				Columns: columns,
			}, nil
		} else {
			return nil, fmt.Errorf("неверная структура запроса4")
		}
	}

	if p.isEnd() {
		return SelectQuery{
			Table:   table,
			Columns: columns,
		}, nil
	}

	if p.current().Value != "WHERE" {
		return nil, fmt.Errorf("неверная структура запроса5")
	}

	p.next()

	// Условия
	conditions, err := p.getArrayConditions()

	if err != nil {
		return nil, err
	}

	if p.isSemicolon() {
		p.next()
	}

	if !p.isEnd() {
		return nil, fmt.Errorf("неверная структура запроса6")
	}

	return SelectQuery{
		Table:      table,
		Columns:    columns,
		Conditions: conditions,
	}, nil
}
