package vovanDB

import (
	"fmt"
)

func createExecutor(s CreateQuery) error { 
    tableName := s.Table

	err := createValidator(s)

	if err != nil {
		return err
	}

	// Создаем файлы таблицы
	err = createTableFiles(tableName)

	if err != nil {
		return fmt.Errorf("не удалось создать файлы для таблицы: %s", tableName)
	}

	
	return nil
}

func createValidator(s CreateQuery) error {
	// Существование таблицы
	if fileExists(getPathTableSchema(s.Table)) || fileExists(getPathTableData(s.Table)) {
		return fmt.Errorf("невозможно создать таблицу %s, она уже существует", s.Table)
	}

	// Уникальность имен колонок
	nameColumns := make(map[string]bool)

	for _, columns := range s.Columns {
		if nameColumns[columns.Name] {
			return fmt.Errorf("дубль колонки %s", columns.Name)
		}

		nameColumns[columns.Name] = true
	}

	return nil
}
