package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

type ExecuteResult struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
	Error   string `json:"error"`
	Time    string `json:"time"`
	Memory  string `json:"memory"`
}

func main() {
	sql := `
	CREATE TABLE Customers (
		id int AUTO_INCREMENT,
		CustomerName text NULL,
		ContactName text,
		Address text not null,
		City text,
		Country text,
		date date
	);
	`

	start := time.Now()

	// Чтение .env файла
	godotenv.Load(".env")

	// Лексический анализатор
	lexer := NewLexer(sql)

	// Токены лексического анализа
	tokens, err := lexer.Analyze()

	if err != nil {
		result, _ := json.Marshal(
			ExecuteResult{
				Success: false,
				Error:   err.Error(),
			})
		fmt.Println(string(result))

		return
	}

	// Парсер
	parser := NewParser(tokens)

	// Подготовленный запрос
	sqlQuery, err := parser.parse()

	if err != nil {
		result, _ := json.Marshal(
			ExecuteResult{
				Success: false,
				Error:   err.Error(),
			})
		fmt.Println(string(result))

		return
	}

	// Executor
	executor := NewExecutor(sqlQuery)

	// Выполняем запрос
	data, err := executor.executeQuery()

	if err != nil {
		result, _ := json.Marshal(
			ExecuteResult{
				Success: false,
				Error:   err.Error(),
			})
		fmt.Println(string(result))

		return
	}

	duration := time.Since(start)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	result, _ := json.Marshal(
		ExecuteResult{
			Success: true,
			Data:    data,
			Error:   "",
			Time:    duration.String(),
			Memory:  humanReadableBytes(memStats.Alloc),
		})

	fmt.Println(string(result))
}
