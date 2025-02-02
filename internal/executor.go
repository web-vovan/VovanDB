package internal

import "fmt"

type Executor struct {
	sqlQuery SQLQuery
}

func NewExecutor(s SQLQuery) *Executor {
	return &Executor{
		sqlQuery: s,
	}
}

func (e *Executor) ExecuteQuery() (string, error) {
	switch e.sqlQuery.(type) {
	case CreateQuery:
		createQuery, ok := e.sqlQuery.(CreateQuery)

		if !ok {
			return "", fmt.Errorf("ошибка при преобразовании типа CreateQuery")
		}

		data, err := createExecutor(createQuery)

		if err != nil {
			return "", err
		}

		return data, nil
	case DropQuery:
		dropQuery, ok := e.sqlQuery.(DropQuery)

		if !ok {
			return "", fmt.Errorf("ошибка при преобразовании типа DropQuery")
		}

		data, err := dropExecutor(dropQuery)

		if err != nil {
			return "", err
		}

		return data, nil
	case InsertQuery:
		insertQuery, ok := e.sqlQuery.(InsertQuery)

		if !ok {
			return "", fmt.Errorf("ошибка при преобразовании типа InsertQuery")
		}

		data, err := insertExecutor(insertQuery)

		if err != nil {
			return "", err
		}

		return data, nil
	case SelectQuery:
		selectQuery, ok := e.sqlQuery.(SelectQuery)

		if !ok {
			return "", fmt.Errorf("ошибка при преобразовании типа InsertQuery")
		}

		data, err := selectExecutor(selectQuery)

		if err != nil {
			return "", err
		}

		return data, err
	case UpdateQuery:
		updateQuery, ok := e.sqlQuery.(UpdateQuery)

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
