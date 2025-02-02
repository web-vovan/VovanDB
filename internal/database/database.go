package database

import (
	"encoding/json"
	"runtime"
	"time"
	"fmt"
	"vovanDB/internal/lexer"
	"vovanDB/internal/parser"
	"vovanDB/internal"

	"github.com/joho/godotenv"
)

type ExecuteResult struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
	Error   string `json:"error"`
	Time    string `json:"time"`
	Memory  string `json:"memory"`
}

func Execute(sql string) string {
    start := time.Now()

	// Чтение .env файла
	godotenv.Load(".env")

	// Лексический анализатор
	lexer := lexer.NewLexer(sql)

	// Токены лексического анализа
	tokens, err := lexer.Analyze()

	if err != nil {
		result, _ := json.Marshal(
			ExecuteResult{
				Success: false,
				Error:   err.Error(),
			})
		return string(result)
	}

	// Парсер
	parser := parser.NewParser(tokens)

	// Подготовленный запрос
	sqlQuery, err := parser.Parse()

	if err != nil {
		result, _ := json.Marshal(
			ExecuteResult{
				Success: false,
				Error:   err.Error(),
			})
		return string(result)
	}

	// Executor
	executor := internal.NewExecutor(sqlQuery)

	// Выполняем запрос
	data, err := executor.ExecuteQuery()

	if err != nil {
		result, _ := json.Marshal(
			ExecuteResult{
				Success: false,
				Error:   err.Error(),
			})

		return string(result)
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

	return string(result)
}

func ErrorArgs() string {
    result, _ := json.Marshal(
        ExecuteResult{
            Success: false,
            Error:   "Укажите один параметр в качестве sql запроса",
        })

    return string(result)
}

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