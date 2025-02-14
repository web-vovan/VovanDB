package condition

import (
	"fmt"
	"vovanDB/internal/constants"
)

// Условие
type Condition struct {
	Column    string
	Operator  string
	Value     string
	ValueType int
}

func (c Condition) String() string {
	return fmt.Sprintf("\nColumn: %s\nOperator: %s\nValue: %s\nValueType: %s\n===", c.Column, c.Operator, c.Value, constants.TypeNames[c.ValueType])
}

func (c *Condition) Compare(i string) (bool, error) {
	switch c.Operator {
	case "=":
		return i == c.Value, nil
	case "!=":
		return i != c.Value, nil
	case ">":
		return i > c.Value, nil
	case ">=":
		return i >= c.Value, nil
	case "<":
		return i < c.Value, nil
	case "<=":
		return i <= c.Value,nil
	default:
		return false, fmt.Errorf("оператор %s не поддерживается", c.Operator)
	}
}
