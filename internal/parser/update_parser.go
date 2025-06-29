package parser

import (
	"fmt"
	"vovanDB/internal/condition"
)

type UpdateValue struct {
	ColumnName string
	Value      string
	Type       string
}

type UpdateQuery struct {
	Table      string
	Values     []UpdateValue
	Conditions []condition.Condition
}

func (q UpdateQuery) QueryType() string {
	return "UPDATE"
}

func updateParser(p *Parser) (SQLQuery, error) {
	var table string
	var values []UpdateValue
	var conditions []condition.Condition

	if p.isUpdateQuery() {
		p.next()
	} else {
		return nil, fmt.Errorf("парсер поддерживает только update запрос")
	}

	// Таблица
	if !p.isIdentifier() {
		return nil, fmt.Errorf("неверная структура запроса2")
	}

	table = p.next().Value

	if !p.isOperator() && p.current().Value != "SET" {
		return nil, fmt.Errorf("неверная структура запроса3")
	}

	p.next()

	// Значения для обновления
	values, err := p.getUpdateValues()

	if err != nil {
		return nil, err
	}

	if p.isSemicolon() {
		p.next()

		if p.isEnd() {
			return UpdateQuery{
				Table:  table,
				Values: values,
			}, nil
		} else {
			return nil, fmt.Errorf("неверная структура запроса7")
		}
	}

	if p.isEnd() {
		return UpdateQuery{
			Table:  table,
			Values: values,
		}, nil
	}

	if p.current().Value != "WHERE" {
		return nil, fmt.Errorf("неверная структура запроса5")
	}

	p.next()

	// Условия
	conditions, err = p.getArrayConditions()

	if err != nil {
		return nil, err
	}

	if p.isSemicolon() {
		p.next()
	}

	if !p.isEnd() {
		return nil, fmt.Errorf("неверная структура запроса6")
	}

	return UpdateQuery{
		Table:      table,
		Values:     values,
		Conditions: conditions,
	}, nil
}
