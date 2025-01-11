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

	hasAutoIncrementColumn := schema.hasAutoIncrementColumn()
	addAutoIncrementData := false
	autoIncrementColumnName := schema.getAutoIncrementColumnName()

	if hasAutoIncrementColumn {
		if !hasStringInSlice(autoIncrementColumnName, s.Columns) {
			addAutoIncrementData = true
		}
	}

	for _, r := range s.Values {
		if addAutoIncrementData {
			schema.incrementColumn(autoIncrementColumnName)
			insertData.WriteString(strconv.Itoa(schema.AutoIncrements[autoIncrementColumnName]) + ";")
		}

		insertData.Write(getInsertRowData(r))
	}

	// Обновляем схему значением последнего элемента из колонки auto_increment
	if hasAutoIncrementColumn && !addAutoIncrementData {
		newAutoIncrement, _ := strconv.Atoi(s.Values[len(s.Values)-1][0].Value)
		schema.AutoIncrements[autoIncrementColumnName] = newAutoIncrement
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
func getInsertRowData(r []InsertValue) []byte {
	var rowBuffer bytes.Buffer

	countValues := len(r)

	for i, v := range r {
		rowBuffer.WriteString(v.Value)

		if i < countValues-1 {
			rowBuffer.WriteString(";")
		}
	}

	rowBuffer.WriteString("\n")

	return rowBuffer.Bytes()
}
