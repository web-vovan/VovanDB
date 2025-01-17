package vovanDB

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	Type          int    `json:"type"`
	AutoIncrement bool   `json:"autoIncrement"`
	NotNull       bool   `json:"notNull"`
}

func (s TableSchema) String() string {
	var result string

	result = "Table: " + s.TableName + "\n"
	result += "Columns: \n"

	for _, c := range *s.Columns {
		result += c.Name + " : " + typeNames[c.Type] + "\n"
	}

	result += "AutoIncrements: \n"

	for i, v := range s.AutoIncrements {
		result += i + " : " + strconv.Itoa(v) + "\n"
	}

	return result
}

// Проверка наличия колонки в схеме
func (s *TableSchema) hasColumnInSchema(columnName string) bool {
	for _, c := range *s.Columns {
		if c.Name == columnName {
			return true
		}
	}

	return false
}

// Колонка
func (s *TableSchema) getColumn(columnName string) (ColumnSchema, error) {
	for _, c := range *s.Columns {
		if c.Name == columnName {
			return c, nil
		}
	}

	return ColumnSchema{}, fmt.Errorf("колонки %s нет в схеме таблицы %s", columnName, s.TableName)
}


// Тип колонки
func (s *TableSchema) getColumnType(columnName string) (int, error) {
	for _, c := range *s.Columns {
		if c.Name == columnName {
			return c.Type, nil
		}
	}

	return -1, fmt.Errorf("колонки %s нет в схеме таблицы %s", columnName, s.TableName)
}

// Индекс колонки
func (s *TableSchema) getColumnIndex(columnName string) (int, error) {
	for i, c := range *s.Columns {
		if c.Name == columnName {
			return i, nil
		}
	}

	return -1, fmt.Errorf("колонки %s нет в схеме таблицы %s", columnName, s.TableName)
}

func (s *TableSchema) writeToFile() error {
	schemaData, err := json.MarshalIndent(s, "", "  ")

	if err != nil {
		return fmt.Errorf("ошибка при сериализации файла схемы в таблице %s: %w", s.TableName, err)
	}

	err = os.WriteFile(getPathTableSchema(s.TableName), schemaData, 0644)

	if err != nil {
		return fmt.Errorf("не удалось записать данные в файл схемы для таблицы: %s", s.TableName)
	}

	return nil
}

func (s *TableSchema) hasAutoIncrementColumn() bool {
	return len(s.AutoIncrements) > 0
}

func (s *TableSchema) getAutoIncrementColumnName() string {
	for _, c := range *s.Columns {
		if c.AutoIncrement {
			return c.Name
		}
	}

	return ""
}

func (s *TableSchema) getAutoIncrementColumnIndex() int {
	for i, c := range *s.Columns {
		if c.AutoIncrement {
			return i
		}
	}

	return -1
}

func (s *TableSchema) getAutoIncrementColumnValue() int {
	column := s.getAutoIncrementColumnName()

	if column == "" {
		return -1
	}

	r, ok := s.AutoIncrements[column]

	if !ok {
		return -1
	}

	return r
}

func (s *TableSchema) incrementColumn(column string) {
	s.AutoIncrements[column]++
}
