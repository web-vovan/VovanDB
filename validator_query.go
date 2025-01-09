package vovanDB

import "fmt"

func validateSelectQuery(s SelectQuery) error {
	tableName := s.Table

	// Существование таблицы
	if !fileExists(getPathTableSchema(tableName)) || !fileExists(getPathTableData(tableName)) {
		return fmt.Errorf("невозможно выполнить запрос, таблица %s не существует", tableName)
	}

	// Загружаем схему
	schema, err := getSchema(tableName)

	if err != nil {
		return err
	}

	// Сравниваем названия колонок в select
	if !(len(s.Columns) == 1 && s.Columns[0] == "*") {
		for _, column := range s.Columns {
			if !schema.hasColumnInSchema(column) {
				return fmt.Errorf("запрос невалиден, в таблице %s нет колонки %s", tableName, column)
			}
		}
	}

	err = validateConditions(schema, s.Conditions)

	if err != nil {
		return err
	}

	return nil
}

func validateUpdateQuery(s UpdateQuery) error {
	tableName := s.Table

	// Существование таблицы
	if !fileExists(getPathTableSchema(tableName)) || !fileExists(getPathTableData(tableName)) {
		return fmt.Errorf("невозможно выполнить запрос, таблица %s не существует", tableName)
	}

	// Загружаем схему
	schema, err := getSchema(tableName)

	if err != nil {
		return err
	}

	// Сравниваем названия колонок в update
	for _, value := range s.Values {
		if !schema.hasColumnInSchema(value.ColumnName) {
			return fmt.Errorf("запрос невалиден, в таблице %s нет колонки %s", tableName, value.ColumnName)
		}
	}

	// Сравниваем типы колонок в update
	for _, value := range s.Values {
		columnType, err := schema.getColumnType(value.ColumnName)

		if err != nil {
			return err
		}

		if columnType != value.Type {
			return fmt.Errorf("неверный тип колонки %s в условии where", value.ColumnName)
		}
	}

	err = validateConditions(schema, s.Conditions)

	if err != nil {
		return err
	}

	return nil
}

func validateCreateQuery(s CreateQuery) error {
	// Существование таблицы
	if fileExists(getPathTableSchema(s.Table)) || fileExists(getPathTableData(s.Table)) {
		return fmt.Errorf("невозможно создать таблицу %s, она уже существует", s.Table)
	}

	// Уникальность имен колонок
	nameColumns := make(map[string]bool)

	for _, column := range s.Columns {
		if nameColumns[column.Name] {
			return fmt.Errorf("дубль колонки %s", column.Name)
		}

		nameColumns[column.Name] = true
	}

	// Проверка типа колонки с auto_increment
	for _, column := range s.Columns {
		if column.AutoIncrement && column.Type != TYPE_DIGIT {
			return fmt.Errorf("колонка %s с типом %s не может иметь атрибут auto_increment", column.Name, typeNames[column.Type])
		}
	}

	return nil
}

func validateDropQuery(s DropQuery) error {
	// Существование файлов таблицы
	if !fileExists(getPathTableSchema(s.Table)) || !fileExists(getPathTableData(s.Table)) {
		return fmt.Errorf("невозможно удалить таблицу %s, она не существует", s.Table)
	}

	return nil
}

func validateInsertQuery(s InsertQuery) error {
	tableName := s.Table

	// Существование таблицы
	if !fileExists(getPathTableSchema(tableName)) || !fileExists(getPathTableData(tableName)) {
		return fmt.Errorf("невозможно выполнить запрос, таблица %s не существует", tableName)
	}

	// Загружаем схему
	schema, err := getSchema(tableName)

	if err != nil {
		return err
	}

	// Проверяем что в insert есть все необходимые колонки
	for _, c := range *schema.Columns {
		if c.AutoIncrement {
			continue
		}

		if !hasStringInSlice(c.Name, s.Columns) {
			return fmt.Errorf("запрос невалиден, в запросе нет обязательной колонки %s", c.Name)
		}
	}

	// Сравниваем названия колонок
	for _, c := range s.Columns {
		if !schema.hasColumnInSchema(c) {
			return fmt.Errorf("запрос невалиден, в таблице %s нет колонки %s", tableName, c)
		}
	}

	// Проверяем количество элементов в строках
	for i, rowValues := range s.Values {
		if len(rowValues) != len(s.Columns) {
			return fmt.Errorf("запрос невалиден, в строке %d неверное количество элементов", i+1)
		}
	}

	// Сравниваем типы вставляемых значений
	for i, rowValues := range s.Values {
		for j, value := range rowValues {
			if (*schema.Columns)[j].AutoIncrement {
				continue
			}

			if value.Type != (*schema.Columns)[j].Type {
				return fmt.Errorf("запрос невалиден, в строке %d неверный тип ", i+1)
			}
		}
	}

	return nil
}

func validateConditions(schema TableSchema, conditions []Condition) error {
	// Сравниваем названия колонок в where
	for _, c := range conditions {
		if !schema.hasColumnInSchema(c.Column) {
			return fmt.Errorf("запрос невалиден, в таблице %s нет колонки %s", schema.TableName, c.Column)
		}
	}

	// Сравниваем типы в where
	for _, c := range conditions {
		columnType, err := schema.getColumnType(c.Column)

		if err != nil {
			return err
		}

		if columnType != c.ValueType {
			return fmt.Errorf("неверный тип колонки %s в условии where", c.Column)
		}
	}

	return nil
}
