package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func humanReadableBytes(bytes uint64) string {
	const uint = 1024

	if bytes < uint {
		return fmt.Sprintf("%d B", bytes)
	}

	if bytes < uint*uint {
		return fmt.Sprintf("%d KB", bytes/uint)
	}

	return fmt.Sprintf("%d MB", bytes/(uint*uint))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func getPathTableSchema(table string) string {
	return filepath.Join(getDataBaseDir(), table+".schema")
}

func getPathTableData(table string) string {
	return filepath.Join(getDataBaseDir(), table+".data")
}

func createTableFiles(tableName string) error {
	err := createDatabaseDirIfNotExists()

	if err != nil {
		return err
	}

	// Создаем файл схемы
	schemaFile, err := os.Create(getPathTableSchema(tableName))

	if err != nil {
		return fmt.Errorf("не удалось создать файл схемы : %s", getPathTableSchema(tableName))
	}

	// Создаем файл с данными
	dataFile, err := os.Create(getPathTableData(tableName))

	if err != nil {
		return fmt.Errorf("не удалось создать файл схемы : %s", getPathTableData(tableName))
	}

	defer schemaFile.Close()
	defer dataFile.Close()

	return nil
}

func createDatabaseDirIfNotExists() error {
	info, err := os.Stat(getDataBaseDir())

	if os.IsNotExist(err) || !info.IsDir() {
		err := os.MkdirAll(getDataBaseDir(), os.ModePerm)

		if err != nil {
			return fmt.Errorf("не удалось создать директорию %s для хранения данных", getDataBaseDir())
		}
	}

	return nil
}

// Получение схемы таблицы
func getSchema(tableName string) (TableSchema, error) {
	var schema TableSchema
	schemaData, err := os.ReadFile(getPathTableSchema(tableName))

	if err != nil {
		return schema, fmt.Errorf("не удалось загрузить файл схемы для таблицы %s", tableName)
	}

	err = json.Unmarshal(schemaData, &schema)

	if err != nil {
		return schema, fmt.Errorf("ошибка при декодировании файла схемы для таблицы %s", tableName)
	}

	return schema, nil
}

// Получение данных таблицы
func getTableData(tableName string) ([][]string, error) {
	var result [][]string
	rawData, err := os.ReadFile(getPathTableData(tableName))

	if err != nil {
		return result, fmt.Errorf("не удалось прочитать файл с данными таблицы %s", tableName)
	}

	for _, line := range strings.Split(string(rawData), "\n") {
		if line != "" {
			result = append(result, strings.Split(line, ";"))
		}
	}

	return result, nil
}

func hasStringInSlice(value string, slice []string) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}

	return false
}

func writeDataInTable(data []byte, tableName string) error {
	file, err := os.OpenFile(getPathTableData(tableName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	if err != nil {
		return fmt.Errorf("не удалось открыть файл для записи: %w", err)
	}

	_, err = file.Write(data)

	if err != nil {
		return fmt.Errorf("не удалось записать данные в файл: %w", err)
	}

	return nil
}

func isValidDate(str string) bool {
	_, err := time.Parse("2006-01-02", str)

	if err == nil {
		return true
	}

	return false
}

func isValidDatetime(str string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", str)

	if err == nil {
		return true
	}

	return false
}

func getDataBaseDir() string {
	if os.Getenv("VOVAN_DB_DATABASE_DIR") == "" {
		return "database"
	}

	return os.Getenv("VOVAN_DB_DATABASE_DIR")
}

func clearAllDatabase() {
	os.RemoveAll(getDataBaseDir())
}

func createTestDataBase() error {
	sql := `
        CREATE TABLE users (
            id int AUTO_INCREMENT,
            name text NULL,
            age int,
            is_admin bool,
            date date
        );`

	result := Execute(sql)

	executeResult := ExecuteResult{}
	err := json.Unmarshal([]byte(result), &executeResult)

	if err != nil {
		return fmt.Errorf("не удалось создать тестовую таблицу: ошибка при декодировании результата")
	}

	if executeResult.Success != true {
		return fmt.Errorf("не удалось создать тестовую таблицу: %s", executeResult.Error)
	}

	return nil
}

func seedTestDatabase() error {
	sql := `
        INSERT INTO users (id, name, age, is_admin, date)
        VALUES
        (1, 'vova', 38, true, '2025-01-28'),
        (2, 'katya', 33, false, '2025-01-28'),
        (3, 'sacha', 38, false, '2025-01-28');
    `

	result := Execute(sql)

	executeResult := ExecuteResult{}
	err := json.Unmarshal([]byte(result), &executeResult)

	if err != nil {
		return fmt.Errorf("не удалось наполнить таблицу тестовыми данными: ошибка при декодировании результата")
	}

	if executeResult.Success != true {
		return fmt.Errorf("не удалось наполнить таблицу тестовыми данными: %s", executeResult.Error)
	}

	return nil
}
