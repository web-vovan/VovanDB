package parser

import "fmt"

// Условие
type Condition struct {
	Column   string
	Operator string
	Value    string
}

func (c Condition) String() string {
	return fmt.Sprintf("Column: %s, Operator: %s, Value: %s", c.Column, c.Operator, c.Value)
}