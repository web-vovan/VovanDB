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
	tableSchema, _ := getSchema(tableName)

	// Загружаем данные таблицы
	tableData, err := getTableData(tableName)

	if err != nil {
		return err
	}

	// Индексы неподходящих строк
	notMatchingRowIndices, err := getNotMatchingRowIndices(&tableData, &tableSchema, &s.Conditions)

	if err != nil {
		return err
	}

	// Индексы подходящих колонок
	matchingColumnIndices, err := getMatchingColumnIndices(&tableSchema, s.Columns)

	if err != nil {
		return err
	}

	var builder strings.Builder

	countRows := len(tableData) - len(notMatchingRowIndices)
	countColumns := len(matchingColumnIndices)

	builder.WriteString("[")

	for i, line := range tableData {
		// Пропускаем строку, не подходящую под фильтр
		if _, ok := notMatchingRowIndices[i]; ok {
			continue
		}

		builder.WriteString("{")

		for j, data := range line {
			// Пропускаем колонки
			if _, ok := matchingColumnIndices[j]; !ok {
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
