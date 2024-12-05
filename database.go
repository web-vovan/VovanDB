package vovanDB

import (
	"fmt"
)

func Execute(sql string) error {
	// Лексический анализатор
	lexer := NewLexer(sql)

	// Токены лексического анализа
	tokens, err := lexer.Analyze()

	if err != nil {
		return err
	}

	// printTokens(tokens)

	// Парсер
	parser := NewParser(tokens)

	var sqlQuery SQLQuery

	if parser.isCreateQuery() {
		sqlQuery, err = createParser(parser)		
	} else if parser.isSelectQuery() {
		sqlQuery, err = selectParser(parser)
	} else {
		return fmt.Errorf("данный тип запроса пока не поддерживается %s", tokens[0].Value)
	}

	if err != nil {
		return err
	}

	// fmt.Println(sqlQuery)

	// Executor
	executor := NewExecutor(sqlQuery)

	// Выполняем запрос
	err = executor.executeQuery()

	if err != nil {
		return err
	}

	return nil
}
