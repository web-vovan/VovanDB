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

	switch t := e.sqlQuery.(type) {
	case CreateQuery:
        createQuery, ok := e.sqlQuery.(CreateQuery);

        if !ok {
            return fmt.Errorf("ошибка при преобразовании типа CreateQuery")
        }

		err = createExecutor(createQuery)

        if err != nil {
            return err
        }
    default:
        return fmt.Errorf("не поддерживается тип %s", t)
	}

	return nil
}
