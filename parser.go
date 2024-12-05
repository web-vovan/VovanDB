package vovanDB

import (
	"fmt"
	"strings"
)

type SQLQuery interface {
	QueryType() string
}

// Парсер
type Parser struct {
	Tokens []Token
	CurrentToken Token
	Position int
}

// Новый парсер
func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens: tokens,
		CurrentToken: tokens[0],
		Position: 0,
	}
}

// Получение текущего токена
func (p *Parser) current() Token {
	return p.CurrentToken
}

// Получение текущего токена и переход к следующему
func (p *Parser) next() Token {
	currentToken := p.current()

	p.Position++

	if p.Position >= len(p.Tokens) {
		p.CurrentToken = Token{}
	} else {
		p.CurrentToken = p.Tokens[p.Position]	
	}

	return currentToken
}

// Конец списка токенов
func (p *Parser) isEnd() bool {
	return p.CurrentToken.Value == ""
}

func (p *Parser) isSelectQuery() bool {
    t := p.current()

    return t.Type == TYPE_KEYWORD && t.Value == "SELECT"
}

func (p *Parser) isCreateQuery() bool {
    t := p.current()

    return t.Type == TYPE_KEYWORD && t.Value == "CREATE"
}

func (p *Parser) isIdentifier() bool {
    return p.current().Type == TYPE_IDENTIFIER
}

func (p *Parser) isOperator() bool {
    return p.current().Type == TYPE_OPERATOR
}

func (p *Parser) isSymbol() bool {
    return p.current().Type == TYPE_OPERATOR
}

func (p *Parser) isString() bool {
    return p.current().Type == TYPE_STRING
}

func (p *Parser) isDigit() bool {
    return p.current().Type == TYPE_DIGIT
}

func (p *Parser) isBool() bool {
    return p.current().Type == TYPE_BOOL
}

func (p *Parser) isComma() bool {
    return p.current().Type == TYPE_SYMBOL && p.current().Value == ","
}

func (p *Parser) isSemicolon() bool {
    return p.current().Type == TYPE_SYMBOL && p.current().Value == ";"
}

func (p *Parser) isAsterisk() bool {
    return p.current().Type == TYPE_SYMBOL && p.current().Value == "*"
}

func (p *Parser) isAndKeyword() bool {
    return p.current().Type == TYPE_KEYWORD && p.current().Value == "AND"
}

func (p *Parser) isOpenParen() bool {
    return p.current().Type == TYPE_SYMBOL && p.current().Value == "("
}

func (p *Parser) isCloseParen() bool {
    return p.current().Type == TYPE_SYMBOL && p.current().Value == ")"
}

func (p *Parser) getCondition() (Condition, error) {
    nilCondition := Condition{}

    if !p.isIdentifier() {
        return nilCondition, fmt.Errorf("неверная структура в условии where")
    }

    column := p.next().Value

    if !p.isOperator() {
        return nilCondition, fmt.Errorf("неверная структура в условии where")
    }

    operator := p.next().Value

    if !p.isString() && !p.isDigit() && !p.isBool() {
        return nilCondition, fmt.Errorf("неверная структура в условии where")
    }

    value := p.current().Value
    valueType := p.current().Type

    p.next()

    return Condition{
        Column: column,
        Operator: operator,
        Value: value,
        ValueType: valueType,
    }, nil
}

func (p *Parser) getCreateColumn() (CreateColumn, error) {
    nilCreateColumn := CreateColumn{}

    // fmt.Println(p.current())
    if !p.isIdentifier() {
        fmt.Println(p.current())
        return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок 1")
    }

    name := p.next().Value

    if !p.isIdentifier() {
        return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок 2")
    }

    
    typeText := strings.ToUpper(p.next().Value)

    columnType := columnTypes[typeText]

    if columnType == 0 {
        return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок 3")
    } 

    return CreateColumn{
        Name: name,
        Type: columnType,
    }, nil
}