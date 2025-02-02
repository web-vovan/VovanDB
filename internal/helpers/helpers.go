package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func GetPathTableSchema(table string) string {
	return filepath.Join(GetDataBaseDir(), table+".schema")
}

func GetPathTableData(table string) string {
	return filepath.Join(GetDataBaseDir(), table+".data")
}

func CreateTableFiles(tableName string) error {
	err := CreateDatabaseDirIfNotExists()

	if err != nil {
		return err
	}

	// Создаем файл схемы
	schemaFile, err := os.Create(GetPathTableSchema(tableName))

	if err != nil {
		return fmt.Errorf("не удалось создать файл схемы : %s", GetPathTableSchema(tableName))
	}

	// Создаем файл с данными
	dataFile, err := os.Create(GetPathTableData(tableName))

	if err != nil {
		return fmt.Errorf("не удалось создать файл схемы : %s", GetPathTableData(tableName))
	}

	defer schemaFile.Close()
	defer dataFile.Close()

	return nil
}

func CreateDatabaseDirIfNotExists() error {
	info, err := os.Stat(GetDataBaseDir())

	if os.IsNotExist(err) || !info.IsDir() {
		err := os.MkdirAll(GetDataBaseDir(), os.ModePerm)

		if err != nil {
			return fmt.Errorf("не удалось создать директорию %s для хранения данных", GetDataBaseDir())
		}
	}

	return nil
}

// Получение данных таблицы
func GetTableData(tableName string) ([][]string, error) {
	var result [][]string
	rawData, err := os.ReadFile(GetPathTableData(tableName))

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

func HasStringInSlice(value string, slice []string) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}

	return false
}

func WriteDataInTable(data []byte, tableName string) error {
	file, err := os.OpenFile(GetPathTableData(tableName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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

func IsValidDate(str string) bool {
	_, err := time.Parse("2006-01-02", str)

	if err == nil {
		return true
	}

	return false
}

func IsValidDatetime(str string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", str)

	if err == nil {
		return true
	}

	return false
}

func GetDataBaseDir() string {
	if os.Getenv("VOVAN_DB_DATABASE_DIR") == "" {
		return "database"
	}

	return os.Getenv("VOVAN_DB_DATABASE_DIR")
}

func ClearAllDatabase() {
	os.RemoveAll(GetDataBaseDir())
}
