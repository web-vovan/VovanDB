package main

import (
	"fmt"
)

type DropQuery struct {
	Table string
}

func (q DropQuery) String() string {
	return fmt.Sprintf("Table: %s\n", q.Table)
}

func (q DropQuery) QueryType() string {
	return "DROP"
}

func dropParser(p *Parser) (SQLQuery, error) {
	var table string

	if p.isDropQuery() {
		p.next()
	} else {
		return nil, fmt.Errorf("неверная структура запроса для drop")
	}

	if !p.isAndKeyword() && p.current().Value != "TABLE" {
		return nil, fmt.Errorf("неверная структура запроса для drop")
	}

	p.next()

	if !p.isIdentifier() {
		return nil, fmt.Errorf("неверная структура запроса для drop")
	}

	table = p.next().Value

	if p.isSemicolon() {
		p.next()
	}

	if !p.isEnd() {
		return nil, fmt.Errorf("неверная структура запроса для drop4")
	}

	return DropQuery{
		Table: table,
	}, nil
}
