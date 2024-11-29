package parser

import "fmt"

type Token struct {
	Type  int
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %d, Value: %s", t.Type, t.Value)
}

func printTokens(tokens []Token) {
    for _, token := range tokens {
        fmt.Println(token)
    }
}