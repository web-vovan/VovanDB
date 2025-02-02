package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"vovanDB/internal/schema"
	"vovanDB/internal/helpers"
)

// Получение схемы таблицы
func GetSchema(tableName string) (schema.TableSchema, error) {
	var schema schema.TableSchema
	schemaData, err := os.ReadFile(helpers.GetPathTableSchema(tableName))

	if err != nil {
		return schema, fmt.Errorf("не удалось загрузить файл схемы для таблицы %s", tableName)
	}

	err = json.Unmarshal(schemaData, &schema)

	if err != nil {
		return schema, fmt.Errorf("ошибка при декодировании файла схемы для таблицы %s", tableName)
	}

	return schema, nil
}
