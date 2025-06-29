package executor

import (
	"sort"
	"strconv"
	"strings"
	"vovanDB/internal/constants"
	executorHelpers "vovanDB/internal/executor/helpers"
	"vovanDB/internal/helpers"
	"vovanDB/internal/parser"
	schemaHelpers "vovanDB/internal/schema/helpers"
	"vovanDB/internal/validator"
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

	// Фильтруем строки
	filterData := [][]string{}

	if len(matchingRowIndices) < len(tableData) {
		for i, line := range tableData {
			// Пропускаем строку, не подходящую под фильтр
			if _, ok := matchingRowIndices[i]; !ok {
				continue
			}

			filterData = append(filterData, line)
		}
	} else {
		filterData = tableData
	}

	// Сортируем строки
	if s.Sorting.Field != "" {
		orderField, _ := tableSchema.GetColumn(s.Sorting.Field)
		orderFieldIndex, _ := tableSchema.GetColumnIndex(s.Sorting.Field)

		sort.Slice(filterData, func(i, j int) bool {
			if orderField.Type == constants.COLUMN_INT {
				num1, _ := strconv.Atoi(filterData[i][orderFieldIndex])
				num2, _ := strconv.Atoi(filterData[j][orderFieldIndex])

				if s.Sorting.Direction == "ASC" {
					return num1 < num2
				}

				return num1 > num2
			}

			if s.Sorting.Direction == "ASC" {
				return filterData[i][orderFieldIndex] < filterData[j][orderFieldIndex]
			}

			return filterData[i][orderFieldIndex] > filterData[j][orderFieldIndex]
		})
	}

	// Готовим json c данными
	var builder strings.Builder

	countRows := len(matchingRowIndices)
	countColumns := len(matchingColumnIndices)

	builder.WriteString("[")
	addRows := 1

	for _, line := range filterData {
		builder.WriteString("{")

		for j, data := range line {
			// Пропускаем колонки
			if _, ok := matchingColumnIndices[j]; !ok {
				continue
			}

			schemaColumn := (*tableSchema.Columns)[j]

			if schemaColumn.Type == constants.COLUMN_INT ||
				schemaColumn.Type == constants.COLUMN_BOOL ||
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

		if addRows < countRows {
			builder.WriteString(",")
			addRows++
		}
	}

	builder.WriteString("]")

	return builder.String(), nil
}
