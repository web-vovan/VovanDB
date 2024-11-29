package parser

import (
	"strings"
	"unicode"
)

// Типы токенов
const (
	TYPE_UNDEFINED = iota
	TYPE_KEYWORD
	TYPE_OPERATOR
	TYPE_NUMBER
	TYPE_STRING
)

// Ключевые слова
var keywords = map[string]bool{
	"SELECT": true,
	"FROM":   true,
	"WHERE":  true,
}

// Операторы
var operators = map[string]bool{
	"=":  true,
	">":  true,
	"<":  true,
	">=": true,
	"<=": true,
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
func (l *Lexer) Analyze() []Token {
	var tokens []Token

	for {
		ch := l.current()

		if l.isEndString() {
			break
		}

		if l.isSpace() {
			l.next()
			continue
		}

		if unicode.IsLetter(ch) {
			tokens = append(tokens, l.getStringToken())
		} else if unicode.IsDigit(ch) {
			tokens = append(tokens, l.getNumberToken())
		} else {
			tokens = append(tokens, l.getOperatorToken())
		}
	}

	return tokens
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

// Текущий символ является пробелом
func (l *Lexer) isSpace() bool {
	return unicode.IsSpace(l.current())
}

// Конец строки для парсинга
func (l *Lexer) isEndString() bool {
	return l.Ch == 0
}

// Получение строкового токена
func (l *Lexer) getStringToken() Token {
	var builder strings.Builder

	for unicode.IsLetter(l.current()) {
		builder.WriteRune(l.next())
	}

	result := builder.String()

	var tokenType int

	// Проверка на ключевое слово
	if keywords[strings.ToUpper(result)] {
		tokenType = TYPE_KEYWORD
        result = strings.ToUpper(result)
	} else {
		tokenType = TYPE_STRING
	}

	return Token{
		Type:  tokenType,
		Value: result,
	}
}

// Получение числового токена
func (l *Lexer) getNumberToken() Token {
	var builder strings.Builder

	for unicode.IsDigit(l.current()) {
		builder.WriteRune(l.next())
	}

	return Token{
		Type:  TYPE_NUMBER,
		Value: builder.String(),
	}
}

func (l *Lexer) getOperatorToken() Token {
	var builder strings.Builder

	for operators[string(l.current())] {
		builder.WriteRune(l.next())
	}

	return Token{
		Type:  TYPE_OPERATOR,
		Value: builder.String(),
	}
}
