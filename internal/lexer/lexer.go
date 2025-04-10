package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
	"vovanDB/internal/constants"
	"vovanDB/internal/helpers"
)

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

		// Пропускаем комментарии
		if l.isDashSymbol() {
			err := l.clearComment()

			if err != nil {
				return nil, err
			}

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
			return nil, fmt.Errorf("неизвестный символ %s в позиции %d", string(l.current()), l.Position)
		}
	}

	return tokens, nil
}

// Получение текущего символа и переход к следующему
func (l *Lexer) next() rune {
	currentCh := l.Ch

	var width int
    currentCh, width = utf8.DecodeRuneInString(l.Input[l.Position:]) // Декодируем руну

	l.Position += width

	if l.Position >= len(l.Input) {
		l.Ch = 0
	} else {
		l.Ch, _ = utf8.DecodeRuneInString(l.Input[l.Position:])
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

// Проверка текущего символа на конец строки
func (l *Lexer) isLineFeed() bool {
	return l.current() == '\n'
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
	return constants.Operators[string(l.current())]
}

// Проверка текущего символа на символ
func (l *Lexer) isSymbol() bool {
	return constants.Symbols[string(l.current())]
}

// Проверка текущего символа на одинарную кавычку (в них строковые литералы)
func (l *Lexer) isStringLiteral() bool {
	return l.current() == '\''
}

// Проверка текущего символа дефис (с него начинаются комментарии)
func (l *Lexer) isDashSymbol() bool {
	return l.current() == '-'
}

// Проверка текущего символа на обратный апостроф 
func (l *Lexer) isBacktick() bool {
	return l.current() == '`'
}

// Проверка текущего символа на конец строки
func (l *Lexer) isEnd() bool {
	return l.current() == 0
}

// Получение строкового токена
func (l *Lexer) getStringToken() Token {
	var builder strings.Builder

	for l.isString() || l.current() == '_' {
		builder.WriteRune(l.next())
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

		builder.WriteRune(l.next())
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

		builder.WriteRune(l.next())
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
		builder.WriteRune(l.next())
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
		builder.WriteRune(l.next())
	}

	return Token{
		Type:  constants.TYPE_OPERATOR,
		Value: builder.String(),
	}
}

// Получение токена символа
func (l *Lexer) getSymbolToken() Token {
	return Token{
		Type:  constants.TYPE_SYMBOL,
		Value: string(l.next()),
	}
}

// Очистка от комментариев
func (l *Lexer) clearComment() error {
	l.next()

	if !l.isDashSymbol() {
		return fmt.Errorf("неверный формат комментариев")
	}

	l.next()

	if l.isLineFeed() {
		return nil
	}

	if l.current() != ' ' {
		return fmt.Errorf("неверный формат комментариев, отсутствует пробел")
	}

	for {
		if l.isEnd() || l.isLineFeed() {
			break
		}

		l.next()
	}

	return nil
}
