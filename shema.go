package vovanDB

type Table struct {
	Schema *TableSchema
}

type TableSchema struct {
	TableName string         `json:"tableName"`
	Columns   *[]ColumnSchema `json:"columns"`
}

type ColumnSchema struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

func (s TableSchema) String() string {
	var result string

	result = "Table: " + s.TableName + "\n"
	result += "Columns: \n"

	for _, c := range *s.Columns {
		result += c.Name + " : " + typeNames[c.Type] + "\n"
	}

	return result
}

func (s *TableSchema) hasColumnInSchema(columnName string) bool {
	for _, c := range *s.Columns {
		if c.Name == columnName {
			return true
		}
	}

	return false
}
