package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"vovanDB/internal/database"
)

func Run() {
	var executeResult database.ExecuteResult

    if len(os.Args) != 2 {
		executeResult = database.ExecuteResult{
            Success: false,
            Error:   "Ожидается один параметр в качестве sql запроса",
        }
	} else {
		sql := os.Args[1]
		executeResult = database.Execute(sql)
	}

	result, err := json.Marshal(executeResult)
	if err != nil {
		fmt.Println("Ошибка после выполнения запроса: %w", err)
	}

	fmt.Println(string(result))
}