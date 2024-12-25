package vovanDB

import (
	"fmt"
	"strings"
)

func selectExecutor(s SelectQuery) error {
	tableName := s.Table

	err := validateSelectExecutor(s)

	if err != nil {
		return err
	}

	// Загружаем схему
	tableSchema, err := getSchema(tableName)

	if err != nil {
		return err
	}

	// Загружаем данные таблицы
	tableData, err := getTableData(tableName)

	if err != nil {
		return err
	}

	var builder strings.Builder

	builder.WriteString("[")

	for i, line := range tableData {
		builder.WriteString("{")

		for j, data := range line {
			builder.WriteString("\"" + (*tableSchema.Columns)[j].Name + "\"" + ":\"" + data + "\"")

			if j < len(line)-1 {
				builder.WriteString(",")
			}
		}

		builder.WriteString("}")

		if i < len(tableData)-1 {
			builder.WriteString(",")
		}
	}

	fmt.Println(builder.String()) 

	return nil
}

// Валидация
func validateSelectExecutor(s SelectQuery) error {
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

	// Сравниваем названия колонок в select
	if !(len(s.Columns) == 1 && s.Columns[0] == "*") {
		for _, column := range s.Columns {
			if !schema.hasColumnInSchema(column) {
				return fmt.Errorf("запрос невалиден, в таблице %s нет колонки %s", tableName, column)
			}
		}
	}

	// Сравниваем названия колонок в where
	for _, c := range s.Conditions {
		if !schema.hasColumnInSchema(c.Column) {
			return fmt.Errorf("запрос невалиден, в таблице %s нет колонки %s", tableName, c.Column)
		}
	}

	// Сравниваем типы в where
	for _, c := range s.Conditions {
		columnType, err := schema.getColumnType(c.Column)

		if err != nil {
			return err
		}

		if columnType != c.ValueType {
			return fmt.Errorf("неверный тип колонки %s в условии where", c.Column)
		}
	}

	return nil
}
