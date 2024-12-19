package vovanDB

import (
	"fmt"
)

func insertExecutor(s InsertQuery) error {
	// tableName := s.Table

	err := validateInsertExecutor(s)

	if err != nil {
		return err
	}

	return nil
}

// Валидация
func validateInsertExecutor(s InsertQuery) error {
	tableName := s.Table

	// Существование таблицы
	if !fileExists(getPathTableSchema(tableName)) || !fileExists(getPathTableData(tableName)) {
		return fmt.Errorf("невозможно выполнить запрос, таблица %s не существует", tableName)
	}

	// Загружаем схему
	schema, err := getSchema(tableName)

	if err != nil {
		return err
	}

	// Сравниваем количество колонок
	if len(*schema.Columns) != len(s.Columns) {
		return fmt.Errorf("запрос невалиден, количество колонок в запросе и в таблице %s неравно", tableName)
	}

	// Сравниваем названия колонок
	for _, column := range s.Columns {
		if !schema.hasColumnInSchema(column) {
			return fmt.Errorf("запрос невалиден, в таблице %s нет колонки %s", tableName, column)
		}
	}

	// Сравниваем типы вставляемых значений
	for i, rowValues := range s.Values {
		// Проверяем количество
		if len(rowValues) != len(*schema.Columns) {
			return fmt.Errorf("запрос невалиден, в строке %d неверное количество элементов", i+1)
		}

		for j, value := range rowValues {
			if value.Type != (*schema.Columns)[j].Type {
				return fmt.Errorf("запрос невалиден, в строке %d неверный тип ", i+1)
			}
		}
	}

	return nil
}
