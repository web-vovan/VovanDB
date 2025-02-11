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

func deleteExecutor(s parser.DeleteQuery) (string, error) {
	tableName := s.Table

	err := validator.ValidateDeleteQuery(s)

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

	fmt.Println(matchingRowIndices)

	if err != nil {
		return "", err
	}

	// Нет строк для удаления
	if len(matchingRowIndices) == 0 {
		return "успешно удалено 0 строк",nil
	}

	// Новые данные для вставки
	var newData bytes.Buffer

	for i, row := range tableData {
		// Пропускаем строки для удаления
		if _, ok := matchingRowIndices[i]; ok {
			continue
		}

		newData.Write(executorHelpers.TransformArrStringToRowData(row))
	}

	file, err := os.Create(helpers.GetPathTableData(tableName))

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = file.Write(newData.Bytes())

	if err != nil {
		return "", fmt.Errorf("не удалось записать данные в файл: %w", err)
	}

	return fmt.Sprintf("успешно удалено %d строк", len(matchingRowIndices)), nil
}
