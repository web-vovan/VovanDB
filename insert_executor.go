package vovanDB

import (
	"bytes"
	"strconv"
)

func insertExecutor(s InsertQuery) error {
	tableName := s.Table
	schema, err := getSchema(tableName)

	if err != nil {
		return err
	}

	err = validateInsertQuery(s)

	if err != nil {
		return err
	}

	var insertData bytes.Buffer

	autoIncrementColumnName := schema.getAutoIncrementColumnName()

	for _, r := range s.Values {
		insertData.Write(getInsertRowData(r, &schema, autoIncrementColumnName))
	}

	err = writeDataInTable(insertData.Bytes(), tableName)

	if err != nil {
		return err
	}

	err = schema.writeToFile()

	if err != nil {
		return err
	}

	return nil
}

// Получение строки с данными
func getInsertRowData(r []InsertValue, s *TableSchema, autoIncrementColumnName string) []byte {
	var rowBuffer bytes.Buffer

	countValues := len(r)

	if autoIncrementColumnName != "" {
		s.incrementColumn(autoIncrementColumnName)
		rowBuffer.WriteString(strconv.Itoa(s.AutoIncrements[autoIncrementColumnName]) + ";")
	}

	for i, v := range r {
		rowBuffer.WriteString(v.Value)

		if i < countValues-1 {
			rowBuffer.WriteString(";")
		}
	}

	rowBuffer.WriteString("\n")

	return rowBuffer.Bytes()
}
