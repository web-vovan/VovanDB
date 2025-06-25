package tests

import (
	"testing"
	"vovanDB/internal/helpers"
	schemaHelpers "vovanDB/internal/schema/helpers"
	testHelpers "vovanDB/tests/helpers"
	"vovanDB/internal/database"
)

func TestDropExecutor(t *testing.T) {
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

	testData := TestData{
		testName: "success update table",
		sql: `
			DROP TABLE users
		`,
		expectedSuccess: true,
		expectedData:    "таблица users успешно удалена",
		expectedError:   "",
	}

	executeResult := database.Execute(testData.sql)

	if executeResult.Success != testData.expectedSuccess {
		t.Errorf("test error: %s, expected: %v, result: %v", testData.testName, testData.expectedSuccess, executeResult.Success)
	}

	if executeResult.Data != testData.expectedData {
		t.Errorf("test error: %s, expected: %s, result: %s", testData.testName, testData.expectedData, executeResult.Data)
	}

	if executeResult.Error != testData.expectedError {
		t.Errorf("test error: %s, expected: %s, result: %s", testData.testName, testData.expectedError, executeResult.Error)
	}

	_, err = schemaHelpers.GetSchema("users")

	if err == nil {
		t.Errorf("test error: файл схемы не удален")
	}

	_, err = helpers.GetTableData("users")

	if err == nil {
		t.Errorf("test error: файл таблицы не удален")
	}
}
