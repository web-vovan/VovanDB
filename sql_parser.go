package parser

import (
	"fmt"
)

type SQLQuery interface {
	QueryType() string
}

// Условие
type Condition struct {
	Column   string
	Operator string
	Value    string
}

func (c Condition) String() string {
	return fmt.Sprintf("Column: %s, Operator: %s, Value: %s", c.Column, c.Operator, c.Value)
}

func Parse(sql string) (SQLQuery, error) {
	// Лексический анализатор
	lexer := NewLexer(sql)

	// Токены лексического анализа
	tokens, err := lexer.Analyze()

	if err != nil {
		return nil, err
	}
	
	printTokens(tokens)

	return nil, nil
	// if tokens[0].Value == "SELECT" {
	// 	result, err := selectParser(tokens)

	// 	fmt.Println(result)
	// 	fmt.Println(err)

	// 	return result, err
	// }

	// return nil, fmt.Errorf("данный тип запроса пока не поддерживается %s", tokens[0].Value)
}
