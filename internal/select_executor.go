package internal

import (
	"strings"
	"vovanDB/internal/constants"
)

func selectExecutor(s SelectQuery) (string, error) {
	tableName := s.Table

	err := validateSelectQuery(s)

	if err != nil {
		return "", err
	}

	// Загружаем схему
	tableSchema, _ := getSchema(tableName)

	// Загружаем данные таблицы
	tableData, err := getTableData(tableName)

	if err != nil {
		return "", err
	}

	// Индексы неподходящих строк
	notMatchingRowIndices, err := getNotMatchingRowIndices(&tableData, &tableSchema, &s.Conditions)

	if err != nil {
		return "", err
	}

	// Индексы подходящих колонок
	matchingColumnIndices, err := getMatchingColumnIndices(&tableSchema, s.Columns)

	if err != nil {
		return "", err
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

			schemaColumn := (*tableSchema.Columns)[j]

			if schemaColumn.Type == constants.TYPE_DIGIT ||
				schemaColumn.Type == constants.TYPE_BOOL ||
				(!schemaColumn.NotNull && data == "NULL") {
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

	return builder.String(), nil
}
