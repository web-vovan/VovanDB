package vovanDB

import (
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

func getPathTableMeta(table string) string {
    return filepath.Join(databaseDir, table +".meta")
}

func getPathTableData(table string) string {
    return filepath.Join(databaseDir, table + ".data")
}