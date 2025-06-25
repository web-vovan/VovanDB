package tests

import (
	"testing"
	testHelpers "vovanDB/tests/helpers"
	"vovanDB/internal/database"
)

func TestSelectExecutor(t *testing.T) {
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
			testName: "success select table",
			sql: `
				SELECT *
				FROM users
				WHERE
				is_admin = false
			`,
			expectedSuccess: true,
			expectedData:    "[{\"id\":2,\"name\":\"katya\",\"age\":33,\"is_admin\":false,\"date\":\"2025-01-28\"},{\"id\":3,\"name\":\"sacha\",\"age\":38,\"is_admin\":false,\"date\":\"2025-01-28\"}]",
			expectedError:   "",
		},
		{
			testName: "success select table",
			sql: `
				SELECT *
				FROM users
			`,
			expectedSuccess: true,
			expectedData:    "[{\"id\":1,\"name\":\"vova\",\"age\":38,\"is_admin\":true,\"date\":\"2025-01-28\"},{\"id\":2,\"name\":\"katya\",\"age\":33,\"is_admin\":false,\"date\":\"2025-01-28\"},{\"id\":3,\"name\":\"sacha\",\"age\":38,\"is_admin\":false,\"date\":\"2025-01-28\"}]",
			expectedError:   "",
		},
		{
			testName: "success select table",
			sql: `
				SELECT id, name
				FROM users
				WHERE
				is_admin = false
			`,
			expectedSuccess: true,
			expectedData:    "[{\"id\":2,\"name\":\"katya\"},{\"id\":3,\"name\":\"sacha\"}]",
			expectedError:   "",
		},
		{
			testName: "success select table order desc",
			sql: `
				SELECT age
				FROM users
				ORDER BY age DESC
			`,
			expectedSuccess: true,
			expectedData:    "[{\"age\":38},{\"age\":38},{\"age\":33}]",
			expectedError:   "",
		},
		{
			testName: "success select table asc",
			sql: `
				SELECT age
				FROM users
				ORDER BY age ASC
			`,
			expectedSuccess: true,
			expectedData:    "[{\"age\":33},{\"age\":38},{\"age\":38}]",
			expectedError:   "",
		},
		{
			testName: "success select table",
			sql: `
				SELECT id, name, age
				FROM users
				WHERE
				age > 33
			`,
			expectedSuccess: true,
			expectedData:    "[{\"id\":1,\"name\":\"vova\",\"age\":38},{\"id\":3,\"name\":\"sacha\",\"age\":38}]",
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
	}
}
