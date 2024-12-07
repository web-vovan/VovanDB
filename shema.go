package vovanDB

type Table struct {
	Schema *TableSchema
	Data   *Data
}

type TableSchema struct {
	TableName string         `json:"tableName"`
	Columns   *[]ColumnSchema `json:"columns"`
}

type ColumnSchema struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type Data struct {
}
