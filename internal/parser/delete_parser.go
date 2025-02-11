package parser

import (
	"fmt"
	"vovanDB/internal/condition"
)

type DeleteQuery struct {
	Table      string
	Conditions []condition.Condition
}

func (q DeleteQuery) QueryType() string {
	return "DELETE"
}

func deleteParser(p *Parser) (SQLQuery, error) {
	var table string
	var conditions []condition.Condition

	if p.isDeleteQuery() {
		p.next()
	} else {
		return nil, fmt.Errorf("парсер поддерживает только delete запрос")
	}

	if p.current().Value != "FROM" {
		return nil, fmt.Errorf("неверная структура delete запроса, ожидается - FROM, указано - %s", p.current().Value)
	}

	p.next()

	// Таблица
	if !p.isIdentifier() {
		return nil, fmt.Errorf("неверная структура delete запроса, ожидается идентификатор")
	}

	table = p.next().Value

	if p.isSemicolon() {
		p.next()

		if p.isEnd() {
			return DeleteQuery{
				Table:  table,
			}, nil
		} else {
			return nil, fmt.Errorf("неверная структура delete запроса, указаны данные после ;")
		}
	}

	if p.isEnd() {
		return DeleteQuery{
			Table:  table,
		}, nil
	}

	if p.current().Value != "WHERE" {
		return nil, fmt.Errorf("неверная структура delete запроса, ожидается - WHERE, указано - %s", p.current().Value)
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
		return nil, fmt.Errorf("неверная структура delete запроса, ожидается конец запроса")
	}

	return DeleteQuery{
		Table:      table,
		Conditions: conditions,
	}, nil
}
