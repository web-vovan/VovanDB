package internal

import (
	"fmt"
	"strings"
	"vovanDB/internal/constants"
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
    } else if p.isUpdateQuery() {
        sqlQuery, err = updateParser(p)
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
    return p.isKeyword() && p.current().Value == "SELECT"
}

func (p *Parser) isCreateQuery() bool {
    return p.isKeyword() && p.current().Value == "CREATE"
}

func (p *Parser) isDropQuery() bool {
    return p.isKeyword() && p.current().Value == "DROP"
}

func (p *Parser) isInsertQuery() bool {
    return p.isKeyword() && p.current().Value == "INSERT"
}

func (p *Parser) isUpdateQuery() bool {
    return p.isKeyword() && p.current().Value == "UPDATE"
}

func (p *Parser) isIdentifier() bool {
    return p.current().Type == constants.TYPE_IDENTIFIER
}

func (p *Parser) isKeyword() bool {
    return p.current().Type == constants.TYPE_KEYWORD
}

func (p *Parser) isOperator() bool {
    return p.current().Type == constants.TYPE_OPERATOR
}

func (p *Parser) isSymbol() bool {
    return p.current().Type == constants.TYPE_SYMBOL
}

func (p *Parser) isString() bool {
    return p.current().Type == constants.TYPE_STRING
}

func (p *Parser) isDigit() bool {
    return p.current().Type == constants.TYPE_DIGIT
}

func (p *Parser) isBool() bool {
    return p.current().Type == constants.TYPE_BOOL
}

func (p *Parser) isNull() bool {
    return p.current().Type == constants.TYPE_NULL
}

func (p *Parser) isComma() bool {
    return p.isSymbol() && p.current().Value == ","
}

func (p *Parser) isSemicolon() bool {
    return p.isSymbol() && p.current().Value == ";"
}

func (p *Parser) isAsterisk() bool {
    return p.isSymbol() && p.current().Value == "*"
}

func (p *Parser) isAndKeyword() bool {
    return p.isKeyword() && p.current().Value == "AND"
}

func (p *Parser) isAutoIncrementKeyword() bool {
    return p.isKeyword() && p.current().Value == "AUTO_INCREMENT"
}

func (p *Parser) isNotKeyword() bool {
    return p.isKeyword() && p.current().Value == "NOT"
}

func (p *Parser) isOpenParen() bool {
    return p.isSymbol() && p.current().Value == "("
}

func (p *Parser) isCloseParen() bool {
    return p.isSymbol() && p.current().Value == ")"
}

func (p *Parser) isEqualOperator() bool {
    return p.isOperator() && p.current().Value == "="
}

func (p *Parser) getCondition() (Condition, error) {
    nilCondition := Condition{}

    if !p.isIdentifier() {
        return nilCondition, fmt.Errorf("неверная структура в условии where, нет идентификатора")
    }

    column := p.next().Value

    if !p.isOperator() {
        return nilCondition, fmt.Errorf("неверная структура в условии where, нет оператора")
    }

    operator := p.next().Value

    if !p.isString() && !p.isDigit() && !p.isBool() && !p.isNull() {
        return nilCondition, fmt.Errorf("неверная структура в условии where, неподдерживаемый тип в условии")
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
        return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок1")
    }

    name := p.next().Value

    if !p.isIdentifier() {
        return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок2")
    }
    
    typeText := strings.ToUpper(p.next().Value)

    columnType := constants.ColumnTypes[typeText]

    if columnType == 0 {
        return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок3")
    }

	autoIncrement := false

	if p.isAutoIncrementKeyword() {
		autoIncrement = true
		p.next()
	}

	notNull := false

	if p.isNotKeyword() {
		p.next()

		if !p.isNull() {
			return nilCreateColumn, fmt.Errorf("неверная структура в create при указании колонок4")
		}

		p.next()

		notNull = true
	}

	if p.isNull() {
		notNull = false
		p.next()
	}

    return CreateColumn{
        Name: name,
        Type: columnType,
		AutoIncrement: autoIncrement,
		NotNull: notNull,
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

func (p *Parser) getArrayConditions() ([]Condition, error) {
    var conditions []Condition

    // Условия
	for {
		if p.isEnd() || p.isSemicolon() {
			break
		}

		if p.isAndKeyword() {
			p.next()
			continue
		}

		condition, err := p.getCondition()

		if err != nil {
			return nil, err
		}

		conditions = append(conditions, condition)
	}

	if len(conditions) == 0 {
		return nil, fmt.Errorf("нет условий в выражении where")
	}

    return conditions, nil
}

func (p *Parser) getUpdateValues() ([]UpdateValue, error) {
    var values []UpdateValue

    for {
		if p.isKeyword() || p.isSemicolon() || p.isEnd() {
			break
		}

		if p.isComma() {
			p.next()
			continue
		}

		if !p.isIdentifier() {
			return nil, fmt.Errorf("некорректное имя колонки для обновления")
		}

		columnName := p.next().Value

		if !p.isEqualOperator() {
			return nil, fmt.Errorf("некорректный символ для колонки при обновлении")
		}

		p.next()

		if p.isString() || p.isBool() || p.isDigit() {
			columnValue := p.current().Value
			columnType := p.current().Type

			values = append(values, UpdateValue{
				ColumnName: columnName,
				Value:      columnValue,
				Type:       columnType,
			})

			p.next()
		} else {
			return nil, fmt.Errorf("неверный тип для значения в колонки для обновления")
		}
	}

    return values, nil
}