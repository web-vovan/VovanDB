package schema

import (
	"encoding/json"
	"fmt"
	"os"

	"vovanDB/internal/helpers"
)

type Table struct {
	Schema *TableSchema
}

type TableSchema struct {
	TableName      string          `json:"tableName"`
	Columns        *[]ColumnSchema `json:"columns"`
	AutoIncrements map[string]int  `json:"autoIncrementValues"`
}

type ColumnSchema struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	AutoIncrement bool   `json:"autoIncrement"`
	NotNull       bool   `json:"notNull"`
}

// Проверка наличия колонки в схеме
func (s *TableSchema) HasColumnInSchema(columnName string) bool {
	for _, c := range *s.Columns {
		if c.Name == columnName {
			return true
		}
	}

	return false
}

// Колонка
func (s *TableSchema) GetColumn(columnName string) (ColumnSchema, error) {
	for _, c := range *s.Columns {
		if c.Name == columnName {
			return c, nil
		}
	}

	return ColumnSchema{}, fmt.Errorf("колонки %s нет в схеме таблицы %s", columnName, s.TableName)
}

// Тип колонки
func (s *TableSchema) GetColumnType(columnName string) (string, error) {
	for _, c := range *s.Columns {
		if c.Name == columnName {
			return c.Type, nil
		}
	}

	return "", fmt.Errorf("колонки %s нет в схеме таблицы %s", columnName, s.TableName)
}

// Индекс колонки
func (s *TableSchema) GetColumnIndex(columnName string) (int, error) {
	for i, c := range *s.Columns {
		if c.Name == columnName {
			return i, nil
		}
	}

	return -1, fmt.Errorf("колонки %s нет в схеме таблицы %s", columnName, s.TableName)
}

func (s *TableSchema) WriteToFile() error {
	schemaData, err := json.MarshalIndent(s, "", "  ")

	if err != nil {
		return fmt.Errorf("ошибка при сериализации файла схемы в таблице %s: %w", s.TableName, err)
	}

	err = os.WriteFile(helpers.GetPathTableSchema(s.TableName), schemaData, 0644)

	if err != nil {
		return fmt.Errorf("не удалось записать данные в файл схемы для таблицы: %s", s.TableName)
	}

	return nil
}

func (s *TableSchema) HasAutoIncrementColumn() bool {
	return len(s.AutoIncrements) > 0
}

func (s *TableSchema) GetAutoIncrementColumnName() string {
	for _, c := range *s.Columns {
		if c.AutoIncrement {
			return c.Name
		}
	}

	return ""
}

func (s *TableSchema) GetAutoIncrementColumnIndex() int {
	for i, c := range *s.Columns {
		if c.AutoIncrement {
			return i
		}
	}

	return -1
}

func (s *TableSchema) GetAutoIncrementColumnValue() int {
	column := s.GetAutoIncrementColumnName()

	if column == "" {
		return -1
	}

	r, ok := s.AutoIncrements[column]

	if !ok {
		return -1
	}

	return r
}

func (s *TableSchema) IncrementColumn(column string) {
	s.AutoIncrements[column]++
}
