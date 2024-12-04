package parser

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

func selectParser(p *Parser) (SQLQuery, error) {
	var columns []string
	var table string
	var conditions []Condition

	if p.isSelectQuery() {
		p.next()
	} else {
		return nil, fmt.Errorf("неверная структура запроса для select")
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
		return nil, fmt.Errorf("неверная структура запроса")
	}

	p.next()

	// Таблица
	if !p.isIdentifier() {
		return nil, fmt.Errorf("неверная структура запроса")
	}

	table = p.next().Value

	if p.isEnd() {
		return SelectQuery{
			Table:   table,
			Columns: columns,
		}, nil
	}

	if p.current().Value != "WHERE" {
		return nil, fmt.Errorf("неверная структура запроса")
	}

	p.next()

	// Условия
	for {
		if p.isEnd() {
			break
		}

		if p.isAndKeyword() {
			p.next()
			continue
		}

		condition, err := p.getCondition()

		if err != nil {
			return nil, err
		}

		conditions = append(conditions, condition)
	}

	if len(conditions) == 0 {
		return nil, fmt.Errorf("неверная структура запроса")
	}

	if p.isSemicolon() {
		p.next()
	}

	if !p.isEnd() {
		return nil, fmt.Errorf("неверная структура запроса")
	}

	return SelectQuery{
		Table:      table,
		Columns:    columns,
		Conditions: conditions,
	}, nil
}
