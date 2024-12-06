package vovanDB

import (
    "fmt"
)

func createExecutor(s CreateQuery) error {
    err := createValidator(s)

    if err != nil {
        return err
    }

    fmt.Println("create executor")

    return nil
}

func createValidator(s CreateQuery) error {
    // Существование таблицы
    if fileExists(getPathTableMeta(s.Table)) || fileExists(getPathTableData(s.Table)) {
        return fmt.Errorf("невозможно создать таблицу %s, она уже существует", s.Table)
    }

    // Уникальность имен колонок
    nameColumns := make(map[string]bool)
    
    for _, columns := range s.Columns {
        if nameColumns[columns.Name] {
            return fmt.Errorf("дубль колонки %s", columns.Name)
        }

        nameColumns[columns.Name] = true
    }

    return nil
}
