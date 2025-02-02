package internal

import (
	"fmt"
	"vovanDB/internal/constants"
)

type Token struct {
	Type  int
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %s;Value: %s ||", constants.TypeNames[t.Type], t.Value)
}

func printTokens(tokens []Token) {
	for _, token := range tokens {
		fmt.Println(token)
	}
}
