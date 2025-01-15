package vovanDB

import (
	"encoding/json"
	"runtime"
	"time"
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
		return string(result)
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
		return string(result)
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
