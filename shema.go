package vovanDB

type TableSchema struct {
    TableName string
    Columns *[]ColumnSchema
}

type ColumnSchema struct {
    Name string
    Type int
}