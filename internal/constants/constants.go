package constants

// Типы токенов
const (
	TYPE_UNDEFINED = iota
	TYPE_KEYWORD
	TYPE_IDENTIFIER
	TYPE_DIGIT
	TYPE_STRING
	TYPE_BOOL
	TYPE_DATE
	TYPE_DATETIME
	TYPE_OPERATOR
	TYPE_SYMBOL
	TYPE_NULL
)

var TypeNames = map[int]string{
	TYPE_UNDEFINED:  "TYPE_UNDEFINED",
	TYPE_KEYWORD:    "TYPE_KEYWORD",
	TYPE_IDENTIFIER: "TYPE_IDENTIFIER",
	TYPE_DIGIT:      "TYPE_DIGIT",
	TYPE_STRING:     "TYPE_STRING",
	TYPE_BOOL:       "TYPE_BOOL",
	TYPE_DATE:       "TYPE_DATE",
	TYPE_DATETIME:   "TYPE_DATETIME",
	TYPE_OPERATOR:   "TYPE_OPERATOR",
	TYPE_SYMBOL:     "TYPE_SYMBOL",
	TYPE_NULL:       "TYPE_NULL",
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
	"TRUE":  true,
	"FALSE": true,
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

// Типы данных колонок
var ColumnTypes = map[string]int{
	"INT":      TYPE_DIGIT,
	"TEXT":     TYPE_STRING,
	"BOOL":     TYPE_BOOL,
	"DATE":     TYPE_DATE,
	"DATETIME": TYPE_DATETIME,
}
