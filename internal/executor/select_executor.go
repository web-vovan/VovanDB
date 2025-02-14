package executor

import (
	"strings"
	"vovanDB/internal/constants"
	"vovanDB/internal/helpers"
	"vovanDB/internal/parser"
	"vovanDB/internal/validator"
	executorHelpers "vovanDB/internal/executor/helpers"
	schemaHelpers "vovanDB/internal/schema/helpers"
)

func selectExecutor(s parser.SelectQuery) (string, error) {
	tableName := s.Table

	err := validator.ValidateSelectQuery(s)

	if err != nil {
		return "", err
	}

	// Загружаем схему
	tableSchema, _ := schemaHelpers.GetSchema(tableName)

	// Загружаем данные таблицы
	tableData, err := helpers.GetTableData(tableName)

	if err != nil {
		return "", err
	}

	// Индексы подходящих строк
	matchingRowIndices, err := executorHelpers.GetMatchingRowIndices(&tableData, &tableSchema, &s.Conditions)

	if err != nil {
		return "", err
	}

	// Индексы подходящих колонок
	matchingColumnIndices, err := executorHelpers.GetMatchingColumnIndices(&tableSchema, s.Columns)

	if err != nil {
		return "", err
	}

	var builder strings.Builder

	countRows := len(tableData) - len(matchingRowIndices)
	countColumns := len(matchingColumnIndices)

	builder.WriteString("[")

	for i, line := range tableData {
		// Пропускаем строку, не подходящую под фильтр
		if _, ok := matchingRowIndices[i]; !ok {
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
