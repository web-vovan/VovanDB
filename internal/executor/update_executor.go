package executor

import (
	"bytes"
	"fmt"
	"os"
	"vovanDB/internal/helpers"
	"vovanDB/internal/parser"
	"vovanDB/internal/validator"
	schemaHelpers "vovanDB/internal/schema/helpers"
	executorHelpers "vovanDB/internal/executor/helpers"
)

func updateExecutor(s parser.UpdateQuery) (string, error) {
	tableName := s.Table

	err := validator.ValidateUpdateQuery(s)

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

	// Нет строк для обновления
	if len(matchingRowIndices) == 0 {
		return "успешно обновлено 0 строк",nil
	}

	// Значения для колонок
	var columnValues = make(map[int]parser.UpdateValue)

	for _, value := range s.Values {
		i, err := tableSchema.GetColumnIndex(value.ColumnName)

		if err != nil {
			return "", err
		}

		columnValues[i] = value
	}

	for i := range matchingRowIndices {
		// Меняем значения в строках
		for j := range tableData[i] {
			// Пропускаем колонки
			if _, ok := columnValues[j]; !ok {
				continue
			}

			// Меняем значение колонки в строке
			tableData[i][j] = columnValues[j].Value
		}
	}

	// Новые данные для вставки
	var updateData bytes.Buffer

	for _, row := range tableData {
		updateData.Write(executorHelpers.TransformArrStringToRowData(row))
	}

	file, err := os.Create(helpers.GetPathTableData(tableName))

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = file.Write(updateData.Bytes())

	if err != nil {
		return "", fmt.Errorf("не удалось записать данные в файл: %w", err)
	}

	return fmt.Sprintf("успешно обновлено %d строк", len(matchingRowIndices)), nil
}
