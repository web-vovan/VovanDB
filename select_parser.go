package parser


import (
	"fmt"
)

type SelectQuery struct {
	Table     string
	Columns   []string
	Condition *Condition
}

func (q SelectQuery) String() string {
    return fmt.Sprintf("Table: %s, Columns: %s, Condition: (%s)", q.Table, q.Columns, *q.Condition)
}

func (q SelectQuery) QueryType() string {
	return "SELECT"
}

func selectParser(tokens []Token) (SQLQuery, error) {
	var columns []string
	var table string
	var condition Condition

	index := 1

    // Колонки
    if tokens[index].Type == TYPE_IDENTIFIER {
        columns = append(columns, tokens[index].Value)
        index++
    } else {
        return nil, fmt.Errorf("ожидается имя колонки, вместо этого %s", tokens[index].Value)
    }

    // Таблица
    if tokens[index].Value != "FROM" {
        return nil, fmt.Errorf("ожидается поле FROM, вместо этого %s", tokens[index].Value)
    }

    index++
    table = tokens[index].Value

    // Условия
    index++

    // Есть условие WHERE
    if index <= len(tokens) {
        if tokens[index].Value != "WHERE" {
            return nil, fmt.Errorf("ожидается поле WHERE, вместо этого %s", tokens[index].Value)
        }

        index++

        if tokens[index].Type != TYPE_IDENTIFIER {
            return nil, fmt.Errorf("в условии ожидается строка, вместо этого %d", tokens[index].Type)
        }

        conditionColumn := tokens[index].Value

        index++

        if tokens[index].Type != TYPE_OPERATOR {
            return nil, fmt.Errorf("в условии ожидается оператор, вместо этого %d", tokens[index].Type)
        }

        conditionOperator := tokens[index].Value

        index++

        conditionValue := tokens[index].Value

        condition = Condition{
            Column:   conditionColumn,
            Operator: conditionOperator,
            Value:    conditionValue,
        }
    }

    return SelectQuery{
        Table: table,
        Columns: columns,
        Condition: &condition,
    }, nil
}
