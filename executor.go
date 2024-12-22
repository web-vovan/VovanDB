package vovanDB

import "fmt"

type Executor struct {
	sqlQuery SQLQuery
}

func NewExecutor(s SQLQuery) *Executor {
	return &Executor{
		sqlQuery: s,
	}
}

func (e *Executor) executeQuery() error {
	var err error

	switch e.sqlQuery.(type) {
	case CreateQuery:
		createQuery, ok := e.sqlQuery.(CreateQuery)

		if !ok {
			return fmt.Errorf("ошибка при преобразовании типа CreateQuery")
		}

		err = createExecutor(createQuery)

		if err != nil {
			return err
		}
	case DropQuery:
		dropQuery, ok := e.sqlQuery.(DropQuery)

		if !ok {
			return fmt.Errorf("ошибка при преобразовании типа DropQuery")
		}

		err = dropExecutor(dropQuery)

		if err != nil {
			return err
		}
	case InsertQuery:
		insertQuery, ok := e.sqlQuery.(InsertQuery)

		if !ok {
			return fmt.Errorf("ошибка при преобразовании типа InsertQuery")
		}

		err = insertExecutor(insertQuery)

		if err != nil {
			return err
		}
	case SelectQuery:
		selectQuery, ok := e.sqlQuery.(SelectQuery)

		if !ok {
			return fmt.Errorf("ошибка при преобразовании типа InsertQuery")
		}

		err = selectExecutor(selectQuery)

		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("данный тип запросов не поддерживается")
	}

	return nil
}
