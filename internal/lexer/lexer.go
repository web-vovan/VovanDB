package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"vovanDB/internal/constants"
	"vovanDB/internal/helpers"
)

type Lexer struct {
	Sql      []rune // запрос для лексического анализа
	Ch       rune   // текущий символ
	Position int    // текущая позиция
}

func NewLexer(input string) *Lexer {
	sql := []rune(input)

	return &Lexer{
		Sql:      sql,
		Ch:       sql[0],
		Position: 0,
	}
}

// Лексический анализ
func (l *Lexer) Analyze() ([]Token, error) {
	var tokens []Token

	for {
		// Выходим если строка закончилась
		if l.isEnd() {
			break
		}

		// Пропускаем пробелы и переносы строк
		if l.isSpace() {
			l.next()
			continue
		}

		// Пропускаем комментарии
		isClearComment, err := l.clearComment()
		if err != nil {
			return nil, err
		}
		if isClearComment {
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
		} else if l.isBacktick() {
			token, err := l.getBacktickToken()
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
			return nil, fmt.Errorf("неизвестный символ %s в позиции %d", string(l.Ch), l.Position)
		}
	}

	return tokens, nil
}

// Переход к следующей руне
func (l *Lexer) next() {
	l.Position++

	if l.Position >= len(l.Sql) {
		l.Ch = 0
	} else {
		l.Ch = l.Sql[l.Position]
	}
}

// Проверка текущего символа на пробел
func (l *Lexer) isSpace() bool {
	return unicode.IsSpace(l.Ch)
}

// Проверка текущего символа на конец строки
func (l *Lexer) isLineFeed() bool {
	return l.Ch == '\n'
}

// Проверка текущей руны на строку
func (l *Lexer) isString() bool {
	return unicode.IsLetter(l.Ch)
}

// Проверка текущей руны на число
func (l *Lexer) isDigit() bool {
	return unicode.IsDigit(l.Ch)
}

// Проверка текущей руны на оператор
func (l *Lexer) isOperator() bool {
	return constants.Operators[string(l.Ch)]
}

// Проверка текущей руны на символ
func (l *Lexer) isSymbol() bool {
	return constants.Symbols[string(l.Ch)]
}

// Проверка текущей руны на одинарную кавычку (в них строковые литералы)
func (l *Lexer) isStringLiteral() bool {
	return l.Ch == '\''
}

// Проверка текущей руны дефис (с него начинаются комментарии)
func (l *Lexer) isDashSymbol() bool {
	return l.Ch == '-'
}

// Проверка текущей руны на обратный апостроф
func (l *Lexer) isBacktick() bool {
	return l.Ch == '`'
}

// Проверка текущей руны на конец строки
func (l *Lexer) isEnd() bool {
	return l.Ch == 0
}

// Получение строкового токена
func (l *Lexer) getStringToken() Token {
	var builder strings.Builder

	for l.isString() || l.Ch == '_' {
		builder.WriteRune(l.Ch)
		l.next()
	}

	result := builder.String()

	var tokenType int

	if constants.Keywords[strings.ToUpper(result)] {
		tokenType = constants.TYPE_KEYWORD
		result = strings.ToUpper(result)
	} else if constants.Bools[strings.ToLower(result)] {
		tokenType = constants.TYPE_BOOL
		result = strings.ToLower(result)
	} else if constants.Null[strings.ToUpper(result)] {
		tokenType = constants.TYPE_NULL
		result = strings.ToUpper(result)
	} else {
		tokenType = constants.TYPE_IDENTIFIER
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

		builder.WriteRune(l.Ch)
		l.next()
	}

	if !hasEndStringLiteral {
		return Token{}, fmt.Errorf("отсутствует закрывающая кавычка для строки")
	}

	if helpers.IsValidDate(builder.String()) {
		return Token{
			Type:  constants.TYPE_DATE,
			Value: builder.String(),
		}, nil
	}

	if helpers.IsValidDatetime(builder.String()) {
		return Token{
			Type:  constants.TYPE_DATETIME,
			Value: builder.String(),
		}, nil
	}

	return Token{
		Type:  constants.TYPE_STRING,
		Value: builder.String(),
	}, nil
}

// Получение токена обратного апострофа
func (l *Lexer) getBacktickToken() (Token, error) {
	var builder strings.Builder

	if l.isBacktick() {
		l.next()
	}

	hasEndBacktick := false

	for {
		if l.isEnd() {
			break
		}

		if l.isBacktick() {
			hasEndBacktick = true
			l.next()
			break
		}

		builder.WriteRune(l.Ch)
		l.next()
	}

	if !hasEndBacktick {
		return Token{}, fmt.Errorf("отсутствует закрывающий обратный апостроф `")
	}

	return Token{
		Type:  constants.TYPE_IDENTIFIER,
		Value: builder.String(),
	}, nil
}

// Получение числового токена
func (l *Lexer) getDigitToken() Token {
	var builder strings.Builder

	for l.isDigit() {
		builder.WriteRune(l.Ch)
		l.next()
	}

	return Token{
		Type:  constants.TYPE_DIGIT,
		Value: builder.String(),
	}
}

// Получение токена оператора
func (l *Lexer) getOperatorToken() Token {
	var builder strings.Builder

	for l.isOperator() {
		builder.WriteRune(l.Ch)
		l.next()
	}

	return Token{
		Type:  constants.TYPE_OPERATOR,
		Value: builder.String(),
	}
}

// Получение токена символа
func (l *Lexer) getSymbolToken() Token {
	token := Token{
		Type:  constants.TYPE_SYMBOL,
		Value: string(l.Ch),
	}

	l.next()

	return token
}

// Очистка от комментариев
func (l *Lexer) clearComment() (bool, error) {
	if !l.isDashSymbol() {
		return false, nil
	}

	l.next()

	if !l.isDashSymbol() {
		return false, fmt.Errorf("неверный формат комментариев")
	}

	l.next()

	if l.isLineFeed() {
		return true, nil
	}

	if l.Ch != ' ' {
		return false, fmt.Errorf("неверный формат комментариев, отсутствует пробел")
	}

	for {
		if l.isEnd() || l.isLineFeed() {
			break
		}

		l.next()
	}

	return true, nil
}
