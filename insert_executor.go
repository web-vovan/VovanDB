package vovanDB

import (
	"bytes"
	"fmt"
	"os"
)

func insertExecutor(s InsertQuery) error {
	tableName := s.Table

	err := validateInsertQuery(s)

	if err != nil {
		return err
	}

	var insertData bytes.Buffer

	for _, r := range s.Values {
		insertData.Write(getInsertRowData(r))
	}

	file, err := os.OpenFile(getPathTableData(tableName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	if err != nil {
		return fmt.Errorf("не удалось открыть файл для записи: %w", err)
	}

	_, err = file.Write(insertData.Bytes())

	if err != nil {
		return fmt.Errorf("не удалось записать данные в файл: %w", err)
	}

	return nil
}

// Получение строки с данными
func getInsertRowData(r []InsertValue) []byte {
	var rowBuffer bytes.Buffer

	countValues := len(r)

	for i, v := range r {
		rowBuffer.WriteString(v.Value)

		if i < countValues - 1 {
			rowBuffer.WriteString(";")
		}
	}

	rowBuffer.WriteString("\n")

	return rowBuffer.Bytes()
}
