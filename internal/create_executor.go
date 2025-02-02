package internal

import (
	"fmt"
	"vovanDB/internal/helpers"
	"vovanDB/internal/schema"
)

func createExecutor(s CreateQuery) (string, error) {
	tableName := s.Table

	err := validateCreateQuery(s)

	if err != nil {
		return "", err
	}

	// Создаем файлы таблицы
	err = helpers.CreateTableFiles(tableName)

	if err != nil {
		return "", fmt.Errorf("не удалось создать файлы для таблицы: %s", tableName)
	}

	// Пишем мета-данные в файл схемы
	var columns []schema.ColumnSchema

	for _, column := range s.Columns {
		columns = append(columns, schema.ColumnSchema{
			Name:          column.Name,
			Type:          column.Type,
			AutoIncrement: column.AutoIncrement,
			NotNull:       column.NotNull,
		})
	}

	// Устанавливаем значение для колонок с auto_increment
	autoIncrementValues := make(map[string]int)

	for _, column := range columns {
		if column.AutoIncrement {
			autoIncrementValues[column.Name] = 0
		}
	}

	schema := schema.TableSchema{
		TableName:      tableName,
		Columns:        &columns,
		AutoIncrements: autoIncrementValues,
	}

	err = schema.WriteToFile()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("таблица %s успешно создана", tableName), nil
}
