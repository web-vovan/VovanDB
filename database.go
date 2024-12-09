package vovanDB

func Execute(sql string) error {
	// Лексический анализатор
	lexer := NewLexer(sql)

	// Токены лексического анализа
	tokens, err := lexer.Analyze()

	if err != nil {
		return err
	}

	// Парсер
	parser := NewParser(tokens)

	// Подготовленный запрос
	sqlQuery, err := parser.parse()

	if err != nil {
		return err
	}

	// Executor
	executor := NewExecutor(sqlQuery)

	// Выполняем запрос
	err = executor.executeQuery()

	if err != nil {
		return err
	}

	return nil
}
