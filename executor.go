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
		err = createExecutor(e)
    default:
        return fmt.Errorf("не поддерживается тип %s", t)
	}

	if err != nil {
		return err
	}

	return nil
}
