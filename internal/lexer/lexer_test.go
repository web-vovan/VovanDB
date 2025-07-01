package lexer

import (
	"testing"
	"vovanDB/internal/constants"

	"github.com/google/go-cmp/cmp"
)

func TestAnalyzeSuccess(t *testing.T) {
	type TestData struct {
		sql          string
		expectedData []Token
	}

	testData := []TestData{
		{
			sql: "select id, name2, `short_name` from users where is_admin=false",
			expectedData: []Token{
				{Type: constants.TOKEN_KEYWORD, Value: "SELECT"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "id"},
				{Type: constants.TOKEN_SYMBOL, Value: ","},
				{Type: constants.TOKEN_IDENTIFIER, Value: "name2"},
				{Type: constants.TOKEN_SYMBOL, Value: ","},
				{Type: constants.TOKEN_IDENTIFIER, Value: "short_name"},
				{Type: constants.TOKEN_KEYWORD, Value: "FROM"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "users"},
				{Type: constants.TOKEN_KEYWORD, Value: "WHERE"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "is_admin"},
				{Type: constants.TOKEN_OPERATOR, Value: "="},
				{Type: constants.TOKEN_BOOL, Value: "false"},
			},
		},
		{
			sql: "select `short_name23^&%`",
			expectedData: []Token{
				{Type: constants.TOKEN_KEYWORD, Value: "SELECT"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "short_name23^&%"},
			},
		},
		{
			sql: "select * from users where id >= 123",
			expectedData: []Token{
				{Type: constants.TOKEN_KEYWORD, Value: "SELECT"},
				{Type: constants.TOKEN_SYMBOL, Value: "*"},
				{Type: constants.TOKEN_KEYWORD, Value: "FROM"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "users"},
				{Type: constants.TOKEN_KEYWORD, Value: "WHERE"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "id"},
				{Type: constants.TOKEN_OPERATOR, Value: ">="},
				{Type: constants.TOKEN_DIGIT, Value: "123"},
			},
		},
		{
			sql: `
            -- comment1
            select * 
            from users -- comment2
            `,
			expectedData: []Token{
				{Type: constants.TOKEN_KEYWORD, Value: "SELECT"},
				{Type: constants.TOKEN_SYMBOL, Value: "*"},
				{Type: constants.TOKEN_KEYWORD, Value: "FROM"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "users"},
			},
		},
		{
			sql: "select `select`",
			expectedData: []Token{
				{Type: constants.TOKEN_KEYWORD, Value: "SELECT"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "select"},
			},
		},
        {
			sql: "create id auto_increment",
			expectedData: []Token{
				{Type: constants.TOKEN_KEYWORD, Value: "CREATE"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "id"},
				{Type: constants.TOKEN_KEYWORD, Value: "AUTO_INCREMENT"},
			},
		},
        {
			sql: "where date = '2025-08-01'",
			expectedData: []Token{
				{Type: constants.TOKEN_KEYWORD, Value: "WHERE"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "date"},
                {Type: constants.TOKEN_OPERATOR, Value: "="},
                {Type: constants.TOKEN_DATE, Value: "2025-08-01"},
			},
		},
        {
			sql: "where date = '2025-08-01 23:59:59'",
			expectedData: []Token{
				{Type: constants.TOKEN_KEYWORD, Value: "WHERE"},
				{Type: constants.TOKEN_IDENTIFIER, Value: "date"},
                {Type: constants.TOKEN_OPERATOR, Value: "="},
                {Type: constants.TOKEN_DATETIME, Value: "2025-08-01 23:59:59"},
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
