package internal

import (
	"bytes"
	"fmt"
	"strconv"
	"vovanDB/internal/helpers"
	"vovanDB/internal/parser"
	schemaHelpers "vovanDB/internal/schema/helpers"
)

func insertExecutor(s parser.InsertQuery) (string, error) {
	tableName := s.Table
	schema, err := schemaHelpers.GetSchema(tableName)

	if err != nil {
		return "", err
	}

	err = validateInsertQuery(s)

	if err != nil {
		return "", err
	}

	var insertData bytes.Buffer

	hasAutoIncrementColumn := schema.HasAutoIncrementColumn()
	addAutoIncrementData := false
	autoIncrementColumnName := schema.GetAutoIncrementColumnName()

	if hasAutoIncrementColumn {
		if !helpers.HasStringInSlice(autoIncrementColumnName, s.Columns) {
			addAutoIncrementData = true
		}
	}

	for _, r := range s.Values {
		if addAutoIncrementData {
			schema.IncrementColumn(autoIncrementColumnName)
			insertData.WriteString(strconv.Itoa(schema.AutoIncrements[autoIncrementColumnName]) + ";")
		}

		insertData.Write(getInsertRowData(r))
	}

	// Обновляем схему значением последнего элемента из колонки auto_increment
	if hasAutoIncrementColumn && !addAutoIncrementData {
		newAutoIncrement, _ := strconv.Atoi(s.Values[len(s.Values)-1][0].Value)
		schema.AutoIncrements[autoIncrementColumnName] = newAutoIncrement
	}

	err = helpers.WriteDataInTable(insertData.Bytes(), tableName)

	if err != nil {
		return "", err
	}

	err = schema.WriteToFile()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("данные в таблицу %s успешно добавлены", tableName), nil
}

// Получение строки с данными
func getInsertRowData(r []parser.InsertValue) []byte {
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
