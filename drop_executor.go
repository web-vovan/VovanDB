package vovanDB

import (
	"fmt"
	"os"
)

func dropExecutor(s DropQuery) error {
	tableName := s.Table

	err := validateDropExecutor(s)

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

// Валидация
func validateDropExecutor(s DropQuery) error {
	// Существование файлов таблицы
	if !fileExists(getPathTableSchema(s.Table)) || !fileExists(getPathTableData(s.Table)) {
		return fmt.Errorf("невозможно удалить таблицу %s, она не существует", s.Table)
	}

	return nil
}
