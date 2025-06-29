package constants

// Типы токенов
const (
	TOKEN_KEYWORD    = "KEYWORD"
	TOKEN_IDENTIFIER = "IDENTIFIER"
	TOKEN_DIGIT      = "DIGIT"
	TOKEN_STRING     = "STRING"
	TOKEN_BOOL       = "BOOL"
	TOKEN_DATE       = "DATE"
	TOKEN_DATETIME   = "DATETIME"
	TOKEN_OPERATOR   = "OPERATOR"
	TOKEN_SYMBOL     = "SYMBOL"
	TOKEN_NULL       = "NULL"
)

// Типы колонок
const (
	COLUMN_INT      = "INT"
	COLUMN_TEXT     = "TEXT"
	COLUMN_BOOL     = "BOOL"
	COLUMN_DATE     = "DATE"
	COLUMN_DATETIME = "DATETIME"
)

// Маппинг типа из запроса на тип колонки
var ColumnTypes = map[string]string{
	"INT":      COLUMN_INT,
	"TEXT":     COLUMN_TEXT,
	"BOOL":     COLUMN_BOOL,
	"DATE":     COLUMN_DATE,
	"DATETIME": COLUMN_DATETIME,
}

// Маппинг типа токена на тип колонки
var TokenToColumn = map[string]string{
	TOKEN_DIGIT:    COLUMN_INT,
	TOKEN_STRING:   COLUMN_TEXT,
	TOKEN_BOOL:     COLUMN_BOOL,
	TOKEN_DATE:     COLUMN_DATE,
	TOKEN_DATETIME: COLUMN_DATETIME,
}

// Ключевые слова
var Keywords = map[string]bool{
	"SELECT":         true,
	"FROM":           true,
	"WHERE":          true,
	"AND":            true,
	"CREATE":         true,
	"TABLE":          true,
	"DROP":           true,
	"INSERT":         true,
	"INTO":           true,
	"VALUES":         true,
	"UPDATE":         true,
	"DELETE":         true,
	"SET":            true,
	"AUTO_INCREMENT": true,
	"NOT":            true,
	"ORDER":          true,
	"BY":             true,
	"ASC":            true,
	"DESC":           true,
}

// Булевы выражения
var Bools = map[string]bool{
	"true":  true,
	"false": true,
}

// Операторы
var Operators = map[string]bool{
	"=":  true,
	">":  true,
	"<":  true,
	">=": true,
	"<=": true,
}

// Символы
var Symbols = map[string]bool{
	"*": true,
	",": true,
	"(": true,
	")": true,
	";": true,
}

// NULL
var Null = map[string]bool{
	"NULL": true,
}
