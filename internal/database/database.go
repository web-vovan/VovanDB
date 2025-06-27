package database

import (
	"runtime"
	"time"
	"fmt"
	"vovanDB/internal/lexer"
	"vovanDB/internal/parser"
	"vovanDB/internal/executor"

	"github.com/joho/godotenv"
)

type ExecuteResult struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
	Error   string `json:"error"`
	Time    string `json:"time"`
	Memory  string `json:"memory"`
}

func Execute(sql string) *ExecuteResult {
    start := time.Now()

	// Чтение .env файла
	godotenv.Load(".env")

	// Лексический анализатор
	lexer := lexer.NewLexer(sql)

	// Токены лексического анализа
	tokens, err := lexer.Analyze()
	if err != nil {
		return getExecuteErrorResult(err)
	}

	// Парсер
	parser := parser.NewParser(tokens)

	// Подготовленный запрос
	sqlQuery, err := parser.Parse()
	if err != nil {
		return getExecuteErrorResult(err)
	}

	// Executor
	executor := executor.NewExecutor(sqlQuery)

	// Выполняем запрос
	data, err := executor.ExecuteQuery()
	if err != nil {
		return getExecuteErrorResult(err)
	}

	duration := time.Since(start)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &ExecuteResult{
		Success: true,
		Data:    data,
		Error:   "",
		Time:    duration.String(),
		Memory:  humanReadableBytes(memStats.Alloc),
	}
}

func getExecuteErrorResult(err error) *ExecuteResult {
	return &ExecuteResult{
		Success: false,
		Error:   err.Error(),
	}
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