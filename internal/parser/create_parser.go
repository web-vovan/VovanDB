package parser

import (
	"fmt"
	"strings"
	"vovanDB/internal/constants"
)

type CreateColumn struct {
	Name string
	Type int
	AutoIncrement bool
	NotNull bool
}

type CreateQuery struct {
	Table   string
	Columns []CreateColumn
}

func (q CreateQuery) String() string {
	return fmt.Sprintf("Table: %s\nColumns: %s\n", q.Table, q.Columns)
}

func (c CreateColumn) String() string {
	return fmt.Sprintf("\nName: %s\nType: %s\nAutoIncrement: %t\nNotNull: %t", c.Name, constants.TypeNames[c.Type], c.AutoIncrement, c.NotNull)
}

func (q CreateQuery) QueryType() string {
	return "CREATE"
}

func createParser(p *Parser) (SQLQuery, error) {
	var table string
	var columns []CreateColumn

	if p.isCreateQuery() {
		p.next()
	} else {
		return nil, fmt.Errorf("неверная структура запроса для create")
	}

	if !p.isAndKeyword() && p.current().Value != "TABLE" {
		return nil, fmt.Errorf("неверная структура запроса для create")
	}

	p.next()

	if !p.isIdentifier() {
		return nil, fmt.Errorf("неверная структура запроса для create")
	}

	table = p.next().Value

	if !p.isOpenParen() {
		return nil, fmt.Errorf("неверная структура запроса, отсутствует символ (")
	}

	p.next()

	hasCloseParen := false

	for {
		if p.isEnd() {
			break
		}

		if p.isComma() {
			p.next()
			continue
		}

		if p.isCloseParen() {
			hasCloseParen = true
			p.next()
			break
		}

		column, err := p.getCreateColumn()

		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	if p.isSemicolon() {
		p.next()
	}

	if !hasCloseParen || !p.isEnd() {
		return nil, fmt.Errorf("неверная структура запроса, отсутствует символ )")
	}

	if !p.isEnd() {
		return nil, fmt.Errorf("неверная структура запроса")
	}

	return CreateQuery{
		Table:   table,
		Columns: columns,
	}, nil
}

func (p *Parser) getCreateColumn() (CreateColumn, error) {
    nilCreateColumn := CreateColumn{}

    if !p.isIdentifier() {
        return nilCreateColumn, fmt.Errorf("ошибка в структуре create, неверное имя для колонки %s", p.current().Value)
    }

    name := p.next().Value

    if !p.isIdentifier() {
        return nilCreateColumn, fmt.Errorf("ошибка в структуре create, неверный тип для колонки %s", p.current().Value)
    }
    
    typeText := strings.ToUpper(p.next().Value)

    columnType := constants.ColumnTypes[typeText]

    if columnType == 0 {
        return nilCreateColumn, fmt.Errorf("ошибка в структуре create, не удалось определить тип для колонки %s", p.current().Value)
    }

	notNull := true

	if p.isNotKeyword() {
		p.next()

		if !p.isNull() {
			return nilCreateColumn, fmt.Errorf("ошибка в структуре create, после NOT ожидается выражение NULL, а указано %s", p.current().Value)
		}

		p.next()

		notNull = true
	}

	if p.isNull() {
		notNull = false
		p.next()
	}

	autoIncrement := false

	if p.isAutoIncrementKeyword() {
		autoIncrement = true
		p.next()
	}

    return CreateColumn{
        Name: name,
        Type: columnType,
		AutoIncrement: autoIncrement,
		NotNull: notNull,
    }, nil
}
