package tests

import (
	"encoding/json"
	"testing"
	testHelpers "vovanDB/tests/helpers"
	"vovanDB/internal/database"
)

func TestDeleteExecutor(t *testing.T) {
	type TestData struct {
		testName        string
		sql             string
		expectedSuccess bool
		expectedData    string
		expectedError   string
	}

	defer testHelpers.ClearAllTestDatabase()

	err := testHelpers.CreateTestDataBase()

	if err != nil {
		t.Error(err)
		return
	}

	err = testHelpers.SeedTestDatabase()

	if err != nil {
		t.Error(err)
		return
	}

	testData := []TestData{
		{
			testName: "success delete table",
			sql: `
				DELETE FROM users
				WHERE
				is_admin = false
			`,
			expectedSuccess: true,
			expectedData:    "успешно удалено 2 строк",
			expectedError:   "",
		},
	}

	for _, item := range testData {
		result := database.Execute(item.sql)

		executeResult := database.ExecuteResult{}
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
