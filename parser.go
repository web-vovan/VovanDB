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

// Парсинг токенов
func (p *Parser) parse() (SQLQuery, error) {
    var sqlQuery SQLQuery
    var err error

	if p.isCreateQuery() {
		sqlQuery, err = createParser(p)		
	} else if p.isSelectQuery() {
		sqlQuery, err = selectParser(p)
	} else if p.isDropQuery() {
        sqlQuery, err = dropParser(p)
    } else if p.isInsertQuery() {
        sqlQuery, err = insertParser(p)
    } else {
		return sqlQuery, fmt.Errorf("данный тип запроса пока не поддерживается %s", p.Tokens[0].Value)
	}

	if err != nil {
		return sqlQuery, err
	}

    return sqlQuery, nil
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

func (p *Parser) isDropQuery() bool {
    t := p.current()

    return t.Type == TYPE_KEYWORD && t.Value == "DROP"
}

func (p *Parser) isInsertQuery() bool {
    t := p.current()

    return t.Type == TYPE_KEYWORD && t.Value == "INSERT"
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

    if !p.isIdentifier() {
        fmt.Println(p.current())
        return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок")
    }

    name := p.next().Value

    if !p.isIdentifier() {
        return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок")
    }
    
    typeText := strings.ToUpper(p.next().Value)

    columnType := columnTypes[typeText]

    if columnType == 0 {
        return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок")
    } 

    return CreateColumn{
        Name: name,
        Type: columnType,
    }, nil
}

func (p *Parser) getInsertRow() ([]InsertValue, error) {
    var rowValues []InsertValue

    if !p.isOpenParen() {
		return rowValues, fmt.Errorf("неверная структура запроса, отсутствует символ (")
	}

    p.next()

    hasCloseParen := false

    for {
        if p.isComma() {
			p.next()
			continue
		}

		if p.isCloseParen() {
            hasCloseParen = true
			p.next()
			break
		}

		rowValues = append(rowValues, InsertValue{
			Type: p.current().Type,
            Value: p.current().Value,
		})

        p.next()
    }

    if !hasCloseParen {
        return rowValues, fmt.Errorf("неверная структура запроса, отсутствует символ )")
    }

    return rowValues, nil
}