package internal

import (
	"fmt"
	"os"
	"vovanDB/internal/helpers"
	"vovanDB/internal/parser"
)

func dropExecutor(s parser.DropQuery) (string, error) {
	tableName := s.Table

	err := validateDropQuery(s)

	if err != nil {
		return "", err
	}

	// Удаляем файл схемы
	err = os.Remove(helpers.GetPathTableSchema(tableName))

	if err != nil {
		return "", fmt.Errorf("не удалось удалить файл схемы для таблицы: %s", tableName)
	}

	err = os.Remove(helpers.GetPathTableData(tableName))

	if err != nil {
		return "", fmt.Errorf("не удалось удалить файл с данными для таблицы: %s", tableName)
	}

	return fmt.Sprintf("таблица %s успешно удалена", tableName), nil
}
