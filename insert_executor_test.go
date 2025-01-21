package vovanDB

import (
	"encoding/json"
	"testing"
)

func TestInsertExecutor(t *testing.T) {
	type TestData struct {
		testName              string
		sql                   string
		expectedSuccess       bool
		expectedData          string
		expectedError         string
		expectedAutoIncrement int
	}

	defer clearAllDatabase()

	err := createTestDataBase()

	if err != nil {
		t.Error(err)
		return
	}

	testData := []TestData{
		{
			testName: "success insert table",
			sql: `
				INSERT INTO users (id, name, age, is_admin, date)
				VALUES
				(1, 'vova', 38, true, '2025-01-28'),
				(2, 'katay', 33, false, '2025-01-28'),
				(3, 'sacha', 38, false, '2025-01-28');
			`,
			expectedSuccess:       true,
			expectedData:          "данные в таблицу users успешно добавлены",
			expectedError:         "",
			expectedAutoIncrement: 3,
		},
		{
			testName: "success insert table",
			sql: `
				INSERT INTO users (name, age, is_admin, date)
				VALUES
				('vova2', 38, true, '2025-01-28'),
				('katay2', 33, false, '2025-01-28'),
				('sacha2', 38, false, '2025-01-28');
			`,
			expectedSuccess:       true,
			expectedData:          "данные в таблицу users успешно добавлены",
			expectedError:         "",
			expectedAutoIncrement: 6,
		},
	}

	for _, item := range testData {
		result := Execute(item.sql)

		executeResult := ExecuteResult{}
		err := json.Unmarshal([]byte(result), &executeResult)

		schema, _ := getSchema("users")
		schemaAutoIncrement := schema.getAutoIncrementColumnValue()

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

		if schemaAutoIncrement != item.expectedAutoIncrement {
			t.Errorf("test error: %s, expected: %d, result: %d", item.testName, item.expectedAutoIncrement, schemaAutoIncrement)
		}
	}
}
