package vovanDB

import (
	"fmt"
	"os"
	"path/filepath"
)

const databaseDir = "database"

func fileExists(path string) bool {
    _, err := os.Stat(path)

    if os.IsNotExist(err) {
        return false
    }

    return err == nil
}

func getPathTableSchema(table string) string {
    return filepath.Join(databaseDir, table +".schema")
}

func getPathTableData(table string) string {
    return filepath.Join(databaseDir, table + ".data")
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
    info, err := os.Stat(databaseDir)

    if os.IsNotExist(err) || !info.IsDir() {
        err := os.MkdirAll(databaseDir, os.ModePerm)

        if err != nil {
            return fmt.Errorf("не удалось создать директорию %s для хранения данных", databaseDir)
        }
    }

    return nil
}