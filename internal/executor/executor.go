package executor

import (
	"fmt"
	"vovanDB/internal/parser"
)

type Executor struct {
	sqlQuery parser.SQLQuery
}

func NewExecutor(s parser.SQLQuery) *Executor {
	return &Executor{
		sqlQuery: s,
	}
}

func (e *Executor) ExecuteQuery() (string, error) {
	switch e.sqlQuery.(type) {
	case parser.CreateQuery:
		createQuery, ok := e.sqlQuery.(parser.CreateQuery)

		if !ok {
			return "", fmt.Errorf("ошибка при преобразовании типа CreateQuery")
		}

		data, err := createExecutor(createQuery)

		if err != nil {
			return "", err
		}

		return data, nil
	case parser.DropQuery:
		dropQuery, ok := e.sqlQuery.(parser.DropQuery)

		if !ok {
			return "", fmt.Errorf("ошибка при преобразовании типа DropQuery")
		}

		data, err := dropExecutor(dropQuery)

		if err != nil {
			return "", err
		}

		return data, nil
	case parser.InsertQuery:
		insertQuery, ok := e.sqlQuery.(parser.InsertQuery)

		if !ok {
			return "", fmt.Errorf("ошибка при преобразовании типа InsertQuery")
		}

		data, err := insertExecutor(insertQuery)

		if err != nil {
			return "", err
		}

		return data, nil
	case parser.SelectQuery:
		selectQuery, ok := e.sqlQuery.(parser.SelectQuery)

		if !ok {
			return "", fmt.Errorf("ошибка при преобразовании типа InsertQuery")
		}

		data, err := selectExecutor(selectQuery)

		if err != nil {
			return "", err
		}

		return data, err
	case parser.UpdateQuery:
		updateQuery, ok := e.sqlQuery.(parser.UpdateQuery)

		if !ok {
			return "", fmt.Errorf("ошибка при преобразовании типа InsertQuery")
		}

		data, err := updateExecutor(updateQuery)

		if err != nil {
			return "", err
		}

		return data, err
	}

	return "", fmt.Errorf("данный тип запросов не поддерживается")
}
