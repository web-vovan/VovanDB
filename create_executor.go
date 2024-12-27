package vovanDB

import (
	"encoding/json"
	"fmt"
	"os"
)

func createExecutor(s CreateQuery) error {
	tableName := s.Table

	err := validateCreateQuery(s)

	if err != nil {
		return err
	}

	// Создаем файлы таблицы
	err = createTableFiles(tableName)

	if err != nil {
		return fmt.Errorf("не удалось создать файлы для таблицы: %s", tableName)
	}

	// Пишем мета-данные в файл схемы
	var columns []ColumnSchema

	for _, column := range s.Columns {
		columns = append(columns, ColumnSchema{
			Name: column.Name,
			Type: column.Type,
		})
	}

	schema := TableSchema{
		TableName: tableName,
		Columns:   &columns,
	}

	schemaData, err := json.MarshalIndent(schema, "", "  ")

	if err != nil {
		return fmt.Errorf("ошибка при сериализации файла схемы в таблице %s: %w", tableName, err)
	}

	err = os.WriteFile(getPathTableSchema(tableName), schemaData, 0644)

	if err != nil {
		return fmt.Errorf("не удалось записать данные в файл схемы для таблицы: %s", tableName)
	}

	return nil
}
