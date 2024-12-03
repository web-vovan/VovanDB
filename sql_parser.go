package parser

import (
	"fmt"
)

type SQLQuery interface {
	QueryType() string
}


func Parse(sql string) (SQLQuery, error) {
	// Лексический анализатор
	lexer := NewLexer(sql)

	// Токены лексического анализа
	tokens, err := lexer.Analyze()

	if err != nil {
		return nil, err
	}

	// printTokens(tokens)

	parser := NewParser(tokens)

	if parser.isSelectQuery() {
		result, err := selectParser(parser)

		fmt.Println(result)
		fmt.Println(err)

		return result, err
	}

	return nil, fmt.Errorf("данный тип запроса пока не поддерживается %s", tokens[0].Value)
}
