package internal

import (
	"bytes"
	"vovanDB/internal/schema"
	"vovanDB/internal/condition"
)

// Индексы строк, удовлетворяющие фильтру
func GetMatchingRowIndices(tableData *[][]string, tableSchema *schema.TableSchema, conditions *[]condition.Condition) (map[int]bool, error) {
	var result = make(map[int]bool)

	mapConditions, err := TransformConditionsToMap(tableSchema, conditions)

    if err != nil {
		return result, err
	}

	for i, line := range *tableData {
		hasFiltered := true
		for j, condition := range mapConditions {
			compareResult, err := condition.Compare(line[j])

			if err != nil {
				return result, err
			}

			if !compareResult {
				hasFiltered = false
				break
			}
		}

		if hasFiltered {
			result[i] = true
		}
	}

    return result, nil
}

// Преобразование массива с условиями в мапу
func TransformConditionsToMap(tableSchema *schema.TableSchema, conditions *[]condition.Condition) (map[int]condition.Condition, error) {
	var result = make(map[int]condition.Condition)

	for _, condition := range *conditions {
		index, err := tableSchema.GetColumnIndex(condition.Column)

		if err != nil {
			return result, err
		}

		result[index] = condition
	}

	return result, nil
}

// Индексы колонок из схемы, удовлетворяющие списку
func GetMatchingColumnIndices(tableSchema *schema.TableSchema, columns []string) (map[int]bool, error) {
	var result = make(map[int]bool)

	if columns[0] == "*" {
		for i := range *tableSchema.Columns {
			result[i] = true
		}
	} else {
		for _, column := range columns {
			index, err := tableSchema.GetColumnIndex(column)

			if err != nil {
				return result, err
			}

			result[index] = true
		}
	}

    return result, nil
}

// Подготовка строки для записи в таблицу с данными
func TransformArrStringToRowData(r []string) []byte {
	var rowBuffer bytes.Buffer

	countValues := len(r)

	for i, v := range r {
		rowBuffer.WriteString(v)

		if i < countValues-1 {
			rowBuffer.WriteString(";")
		}
	}

	rowBuffer.WriteString("\n")

	return rowBuffer.Bytes()
}
