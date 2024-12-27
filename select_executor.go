package vovanDB

import (
	"strings"
)

func selectExecutor(s SelectQuery) error {
	tableName := s.Table

	err := validateSelectQuery(s)

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

	// Индексы строк для фильтрации
	var filterRowIndex = make(map[int]bool)

	// Фильтруем данные
	if len(s.Conditions) > 0 {
		// Условия для фильтрации с индексами
		var indexCondition = make(map[int]Condition)

		for _, condition := range s.Conditions {
			index, err := tableSchema.getColumnIndex(condition.Column)

			if err != nil {
				return err
			}

			indexCondition[index] = condition
		}

		for i, line := range tableData {
			hasFiltered := false
			for j, condition := range indexCondition {
				if condition.Value != line[j] {
					hasFiltered = true
					break
				}
			}

			if hasFiltered {
				filterRowIndex[i] = true
			}
		}
	}

	// Индексы колонок для фильтрации
	var filterColumnIndex = make(map[int]bool)

	if s.showAllColumns() {
		for i := range s.Columns {
			filterColumnIndex[i] = true
		}
	} else {
		for _, column := range s.Columns {
			index, err := tableSchema.getColumnIndex(column)

			if err != nil {
				return err
			}

			filterColumnIndex[index] = true
		}
	}

	var builder strings.Builder

	countRows := len(tableData) - len(filterRowIndex)
	countColumns := len(filterColumnIndex)

	builder.WriteString("[")

	for i, line := range tableData {
		// Пропускаем отфильтрованную строку
		if _, ok := filterRowIndex[i]; ok {
			continue
		}

		builder.WriteString("{")

		for j, data := range line {
			// Пропускаем колонки
			if _, ok := filterColumnIndex[j]; !ok {
				continue
			}

			if (*tableSchema.Columns)[j].Type == TYPE_DIGIT || (*tableSchema.Columns)[j].Type == TYPE_BOOL {
				builder.WriteString("\"" + (*tableSchema.Columns)[j].Name + "\"" + ":" + data)
			} else {
				builder.WriteString("\"" + (*tableSchema.Columns)[j].Name + "\"" + ":\"" + data + "\"")
			}

			if j < countColumns-1 {
				builder.WriteString(",")
			}
		}

		builder.WriteString("}")

		if i < countRows-1 {
			builder.WriteString(",")
		}
	}

	builder.WriteString("]")

	return nil
}
