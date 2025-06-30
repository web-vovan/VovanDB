package lexer

import (
	"testing"
	"vovanDB/internal/constants"

    "github.com/google/go-cmp/cmp"
)

func TestAnalyzeSuccess(t *testing.T) {
    type TestData struct {
        sql string
        expectedData []Token
    }

    testData := []TestData{
        {
            sql: "create users",
            expectedData: []Token{
                {
                    Type: constants.TOKEN_KEYWORD,
                    Value: "CREATE",
                },
                {
                    Type: constants.TOKEN_IDENTIFIER,
                    Value: "users",
                },
            },
        },
    }

    for _, test := range testData {
        lexer := NewLexer(test.sql)
        executeResult, _ := lexer.Analyze()

        if diff := cmp.Diff(executeResult, test.expectedData); diff != "" {
            t.Errorf("Mismatch (-want +got):\n%s", diff)
        }
    }
}