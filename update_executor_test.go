package vovanDB

import (
	"encoding/json"
	"testing"
)

func TestUpdateExecutor(t *testing.T) {
	type TestData struct {
		testName        string
		sql             string
		expectedSuccess bool
		expectedData    string
		expectedError   string
	}

	defer clearAllDatabase()

	err := createTestDataBase()

	if err != nil {
		t.Error(err)
		return
	}

	err = seedTestDatabase()

	if err != nil {
		t.Error(err)
		return
	}

	testData := []TestData{
		{
			testName: "success update table",
			sql: `
				UPDATE users
				SET is_admin = true
				WHERE
				is_admin = false
			`,
			expectedSuccess: true,
			expectedData:    "успешно обновлено 2 строк",
			expectedError:   "",
		},
	}

	for _, item := range testData {
		result := Execute(item.sql)

		executeResult := ExecuteResult{}
		err := json.Unmarshal([]byte(result), &executeResult)

		if err != nil {
			t.Error(err)
		}

		if executeResult.Success != item.expectedSuccess {
			t.Errorf("test error: %s, expected: %v, result: %v", item.testName, item.expectedSuccess, executeResult.Success)
		}

		if executeResult.Data != item.expectedData {
			t.Errorf("test error: %s, expected: %s, result: %s", item.testName, item.expectedData, executeResult.Data)
		}

		if executeResult.Error != item.expectedError {
			t.Errorf("test error: %s, expected: %s, result: %s", item.testName, item.expectedError, executeResult.Error)
		}
	}
}
