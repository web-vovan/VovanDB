package internal

import (
	"fmt"
)

type InsertValue struct {
	Type int
	Value string
}

type InsertQuery struct {
	Table   string
	Columns []string
	Values [][]InsertValue
}

func (q InsertQuery) QueryType() string {
	return "INSERT"
}

func insertParser(p *Parser) (SQLQuery, error) {
	var table string
	var values [][]InsertValue
	var columns []string

	if p.isInsertQuery() {
		p.next()
	} else {
		return nil, fmt.Errorf("неверная структура запроса для insert1")
	}

	if !p.isAndKeyword() && p.current().Value != "INTO" {
		return nil, fmt.Errorf("неверная структура запроса для insert2")
	}

	p.next()

	if !p.isIdentifier() {
		return nil, fmt.Errorf("неверная структура запроса для insert3")
	}

	table = p.next().Value

	if !p.isOpenParen() {
		return nil, fmt.Errorf("неверная структура запроса для insert4")
	}

	p.next()

	hasCloseParen := false

	// Собираем колонки
	for {
		if p.isComma() {
			p.next()
			continue
		}

		if p.isCloseParen() {
			hasCloseParen = true
			p.next()
			break
		}

		if !p.isIdentifier() {
			return nil, fmt.Errorf("неверная структура запроса для insert5")
		}

		columns = append(columns, p.next().Value)
	}

	if !hasCloseParen {
		return nil, fmt.Errorf("неверная структура запроса для insert6")
	}

	if !p.isAndKeyword() && p.current().Value != "VALUES" {
		return nil, fmt.Errorf("неверная структура запроса для insert7")
	}

	p.next()

	for {
		if p.isEnd() || p.isSemicolon(){
			break
		}

		if p.isComma() {
			p.next()
		}

		rowValues, err := p.getInsertRow()

		if err != nil {
			return nil, err
		}

		values = append(values, rowValues)
	}

	if p.isSemicolon() {
		p.next()
	}

	if !p.isEnd() {
		return nil, fmt.Errorf("неверная структура запроса insert8")
	}

	return InsertQuery{
		Table: table,
		Columns: columns,
		Values: values,
	}, nil
}
