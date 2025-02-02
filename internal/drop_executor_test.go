package internal

import (
	"encoding/json"
	"testing"
)

func TestDropExecutor(t *testing.T) {
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

	testData := TestData{
		testName: "success update table",
		sql: `
			DROP TABLE users
		`,
		expectedSuccess: true,
		expectedData:    "таблица users успешно удалена",
		expectedError:   "",
	}

	result := Execute(testData.sql)

	executeResult := ExecuteResult{}
	err = json.Unmarshal([]byte(result), &executeResult)

	if err != nil {
		t.Error(err)
	}

	if executeResult.Success != testData.expectedSuccess {
		t.Errorf("test error: %s, expected: %v, result: %v", testData.testName, testData.expectedSuccess, executeResult.Success)
	}

	if executeResult.Data != testData.expectedData {
		t.Errorf("test error: %s, expected: %s, result: %s", testData.testName, testData.expectedData, executeResult.Data)
	}

	if executeResult.Error != testData.expectedError {
		t.Errorf("test error: %s, expected: %s, result: %s", testData.testName, testData.expectedError, executeResult.Error)
	}

	_, err = getSchema("users")

	if err == nil {
		t.Errorf("test error: файл схемы не удален")
	}

	_, err = getTableData("users")

	if err == nil {
		t.Errorf("test error: файл таблицы не удален")
	}
}
