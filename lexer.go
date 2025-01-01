package vovanDB

import (
	"fmt"
	"strings"
	"unicode"
)

// Типы токенов
const (
	TYPE_UNDEFINED = iota
	TYPE_KEYWORD
	TYPE_IDENTIFIER
	TYPE_DIGIT
	TYPE_STRING
	TYPE_BOOL
	TYPE_OPERATOR
	TYPE_SYMBOL
)

var typeNames = map[int]string{
	TYPE_UNDEFINED:  "TYPE_UNDEFINED",
	TYPE_KEYWORD:    "TYPE_KEYWORD",
	TYPE_IDENTIFIER: "TYPE_IDENTIFIER",
	TYPE_DIGIT:      "TYPE_DIGIT",
	TYPE_STRING:     "TYPE_STRING",
	TYPE_BOOL:       "TYPE_BOOL",
	TYPE_OPERATOR:   "TYPE_OPERATOR",
	TYPE_SYMBOL:     "TYPE_SYMBOL",
}

// Ключевые слова
var keywords = map[string]bool{
	"SELECT": true,
	"FROM":   true,
	"WHERE":  true,
	"AND":    true,
	"CREATE":    true,
	"TABLE":    true,
	"DROP": true,
	"INSERT": true,
	"INTO": true,
	"VALUES": true,
	"UPDATE": true,
	"SET": true,
}

// Булевы выражения
var bools = map[string]bool{
	"TRUE":  true,
	"FALSE": true,
}

// Операторы
var operators = map[string]bool{
	"=":  true,
	">":  true,
	"<":  true,
	">=": true,
	"<=": true,
}

// Символы
var symbols = map[string]bool{
	"*": true,
	",": true,
	"(": true,
	")": true,
	";": true,
}

type Lexer struct {
	Input    string // запрос для лексического анализа
	Ch       rune   // текущий символ
	Position int    // текущая позиция
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		Input:    input,
		Ch:       rune(input[0]),
		Position: 0,
	}
}

// Лексический анализ
func (l *Lexer) Analyze() ([]Token, error) {
	var tokens []Token

	for {
		if l.isEnd() {
			break
		}

		if l.isSpace() {
			l.next()
			continue
		}

		if l.isString() {
			tokens = append(tokens, l.getStringToken())
		} else if l.isStringLiteral() {
			token, err := l.getStringLiteralToken()

			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
		} else if l.isDigit() {
			tokens = append(tokens, l.getDigitToken())
		} else if l.isOperator() {
			tokens = append(tokens, l.getOperatorToken())
		} else if l.isSymbol() {
			tokens = append(tokens, l.getSymbolToken())
		} else {
			return nil, fmt.Errorf("неизвестный символ %s в позиции %d", string(l.current()), l.Position)
		}
	}

	return tokens, nil
}

// Получение текущего символа и переход к следующему
func (l *Lexer) next() rune {
	currentCh := l.Ch

	l.Position++

	if l.Position >= len(l.Input) {
		l.Ch = 0
	} else {
		l.Ch = rune(l.Input[l.Position])
	}

	return currentCh
}

// Получение текущего символа
func (l *Lexer) current() rune {
	return l.Ch
}

// Проверка текущего символа на пробел
func (l *Lexer) isSpace() bool {
	return unicode.IsSpace(l.current())
}

// Проверка текущего символа на строку
func (l *Lexer) isString() bool {
	return unicode.IsLetter(l.current())
}

// Проверка текущего символа на число
func (l *Lexer) isDigit() bool {
	return unicode.IsDigit(l.current())
}

// Проверка текущего символа на оператор
func (l *Lexer) isOperator() bool {
	return operators[string(l.current())]
}

// Проверка текущего символа на символ
func (l *Lexer) isSymbol() bool {
	return symbols[string(l.current())]
}

// Проверка текущего символа на одинарную кавычку (в них строковые литералы)
func (l *Lexer) isStringLiteral() bool {
	return string(l.current()) == "'"
}

// Проверка текущего символа на конец строки
func (l *Lexer) isEnd() bool {
	return l.current() == 0
}

// Получение строкового токена
func (l *Lexer) getStringToken() Token {
	var builder strings.Builder

	for l.isString() {
		builder.WriteRune(l.next())
	}

	result := builder.String()

	var tokenType int

	if keywords[strings.ToUpper(result)] {
		tokenType = TYPE_KEYWORD
		result = strings.ToUpper(result)
	} else if bools[strings.ToUpper(result)] {
		tokenType = TYPE_BOOL
		result = strings.ToUpper(result)
	} else {
		tokenType = TYPE_IDENTIFIER
	}

	return Token{
		Type:  tokenType,
		Value: result,
	}
}

// Получение токена строкового литерала
func (l *Lexer) getStringLiteralToken() (Token, error) {
	var builder strings.Builder

	if l.isStringLiteral() {
		l.next()
	}

	hasEndStringLiteral := false

	for {
		if l.isEnd() {
			break
		}

		if l.isStringLiteral() {
			hasEndStringLiteral = true
			l.next()
			break
		}

		builder.WriteRune(l.next())
	}

	if !hasEndStringLiteral {
		return Token{}, fmt.Errorf("отсутствует закрывающая кавычка для строки")
	}

	return Token{
		Type:  TYPE_STRING,
		Value: builder.String(),
	}, nil
}

// Получение числового токена
func (l *Lexer) getDigitToken() Token {
	var builder strings.Builder

	for l.isDigit() {
		builder.WriteRune(l.next())
	}

	return Token{
		Type:  TYPE_DIGIT,
		Value: builder.String(),
	}
}

// Получение токена оператора
func (l *Lexer) getOperatorToken() Token {
	var builder strings.Builder

	for l.isOperator() {
		builder.WriteRune(l.next())
	}

	return Token{
		Type:  TYPE_OPERATOR,
		Value: builder.String(),
	}
}

// Получение токена символа
func (l *Lexer) getSymbolToken() Token {
	return Token{
		Type:  TYPE_SYMBOL,
		Value: string(l.next()),
	}
}
