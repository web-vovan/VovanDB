package vovanDB

import (
	"fmt"
	"os"
)

func dropExecutor(s DropQuery) error {
	tableName := s.Table

	err := validateDropQuery(s)

	if err != nil {
		return err
	}

	// Удаляем файл схемы
	err = os.Remove(getPathTableSchema(tableName))

	if err != nil {
		return fmt.Errorf("не удалось удалить файл схемы для таблицы: %s", tableName)
	}

	err = os.Remove(getPathTableData(tableName))

	if err != nil {
		return fmt.Errorf("не удалось удалить файл с данными для таблицы: %s", tableName)
	}

	return nil
}
