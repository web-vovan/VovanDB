package parser

import (
	"fmt"
	"vovanDB/internal/condition"
)

type Sorting struct {
	Field     string
	Direction string
}

type SelectQuery struct {
	Table      string
	Columns    []string
	Conditions []condition.Condition
	Sorting    Sorting
}

func (q SelectQuery) String() string {
	return fmt.Sprintf("Table: %s\nColumns: %s\nConditions: \n%s\nSorting: %s\n", q.Table, q.Columns, q.Conditions, q.Sorting)
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
	var orderField string
	var orderDirection = "ASC"
	conditions  := []condition.Condition{}

	if p.isSelectQuery() {
		p.next()
	} else {
		return nil, fmt.Errorf("парсер поддерживает только select запрос")
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
		return nil, fmt.Errorf("неверная структура select запроса, ожидается - FROM, указано - %s", p.current().Value)
	}

	p.next()

	// Таблица
	if !p.isIdentifier() {
		return nil, fmt.Errorf("неверная структура select запроса, ожидается идентификатор")
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
			return nil, fmt.Errorf("неверная структура select запроса, указаны данные после ;")
		}
	}

	if p.isEnd() {
		return SelectQuery{
			Table:   table,
			Columns: columns,
		}, nil
	}

	if p.current().Value == "WHERE" {

		p.next()

		// Условия
		conditionsList, err := p.getArrayConditions()

		if err != nil {
			return nil, err
		}

		conditions = conditionsList
	}

	if p.current().Value == "ORDER" {
		p.next()

		if p.current().Value != "BY" {
			return nil, fmt.Errorf("неверная структура select запроса, ожидается - BY, указано - %s", p.current().Value)
		}

		p.next()

		if !p.isIdentifier() {
			return nil, fmt.Errorf("неверная структура select запроса, в ORDER BY ожидается идентификатор")
		}

		orderField = p.next().Value

		if p.current().Value == "ASC" || p.current().Value == "DESC" {
			orderDirection = p.next().Value
		}
	}

	if p.isSemicolon() {
		p.next()
	}

	if !p.isEnd() {
		return nil, fmt.Errorf("неверная структура select запроса, ожидается - конец запроса, указано - %s", p.current().Value)
	}

	return SelectQuery{
		Table:      table,
		Columns:    columns,
		Conditions: conditions,
		Sorting: Sorting{
			Field: orderField,
			Direction: orderDirection,
		},
	}, nil
}
