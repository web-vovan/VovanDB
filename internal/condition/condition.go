package condition

import (
	"fmt"
)

// Условие
type Condition struct {
	Column    string
	Operator  string
	Value     string
	ValueType string
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
