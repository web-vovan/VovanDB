package tests

import (
	"testing"
	"vovanDB/internal/database"
	testHelpers "vovanDB/tests/helpers"
)

func TestCreateExecutor(t *testing.T) {
	type TestData struct {
		testName        string
		sql             string
		expectedSuccess bool
		expectedData    string
		expectedError   string
	}

	defer testHelpers.ClearAllTestDatabase()

	testData := []TestData{
		{
			testName: "success create table",
			sql: `
			-- какой-то комментарий
			CREATE TABLE users (
				id int AUTO_INCREMENT,
				name text NULL,
				age int,
				is_admin bool,
				date date
			);
			`,
			expectedSuccess: true,
			expectedData:    "таблица users успешно создана",
			expectedError:   "",
		},
		{
			testName: "success create table",
			sql: "CREATE TABLE users (id int AUTO_INCREMENT,`name` text NULL,age int,is_admin bool,`date` date)",
			expectedSuccess: true,
			expectedData:    "таблица users успешно создана",
			expectedError:   "",
		},
	}

	for _, item := range testData {
		executeResult := database.Execute(item.sql)

		if executeResult.Success != item.expectedSuccess {
			t.Errorf("test error: %s, expected: %v, result: %v", item.testName, item.expectedSuccess, executeResult.Success)
		}

		if executeResult.Data != item.expectedData {
			t.Errorf("test error: %s, expected: %s, result: %s", item.testName, item.expectedData, executeResult.Data)
		}

		if executeResult.Error != item.expectedError {
			t.Errorf("test error: %s, expected: %s, result: %s", item.testName, item.expectedError, executeResult.Error)
		}

		testHelpers.ClearAllTestDatabase()
	}
}
