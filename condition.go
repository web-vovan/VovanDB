package main

import "fmt"

// Условие
type Condition struct {
	Column    string
	Operator  string
	Value     string
	ValueType int
}

func (c Condition) String() string {
	return fmt.Sprintf("\nColumn: %s\nOperator: %s\nValue: %s\nValueType: %s\n===", c.Column, c.Operator, c.Value, typeNames[c.ValueType])
}
