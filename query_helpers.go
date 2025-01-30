package main

// Индексы строк, удовлетворяющие фильтру
func getMatchingRowIndices(tableData *[][]string, tableSchema *TableSchema, conditions *[]Condition) (map[int]bool, error) {
	var result = make(map[int]bool)

	mapConditions, err := transformConditionsToMap(tableSchema, conditions)

    if err != nil {
		return result, err
	}

	for i, line := range *tableData {
		hasFiltered := true
		for j, condition := range mapConditions {
			if condition.Value != line[j] {
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

// Индексы строк, неудовлетворяющие фильтру
func getNotMatchingRowIndices(tableData *[][]string, tableSchema *TableSchema, conditions *[]Condition) (map[int]bool, error) {
	var result = make(map[int]bool)

	mapConditions, err := transformConditionsToMap(tableSchema, conditions)

	if err != nil {
		return result, err
	}

	for i, line := range *tableData {
		hasFiltered := false
		for j, condition := range mapConditions {
			if condition.Value != line[j] {
				hasFiltered = true
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
func transformConditionsToMap(tableSchema *TableSchema, conditions *[]Condition) (map[int]Condition, error) {
	var result = make(map[int]Condition)

	for _, condition := range *conditions {
		index, err := tableSchema.getColumnIndex(condition.Column)

		if err != nil {
			return result, err
		}

		result[index] = condition
	}

	return result, nil
}

// Индексы колонок из схемы, удовлетворяющие списку
func getMatchingColumnIndices(tableSchema *TableSchema, columns []string) (map[int]bool, error) {
	var result = make(map[int]bool)

	if columns[0] == "*" {
		for i := range *tableSchema.Columns {
			result[i] = true
		}
	} else {
		for _, column := range columns {
			index, err := tableSchema.getColumnIndex(column)

			if err != nil {
				return result, err
			}

			result[index] = true
		}
	}

    return result, nil
}
