package vovanDB

import (
	"fmt"
	"runtime"
	"time"
)

func Execute(sql string) (string, error) {
	start := time.Now()

	// Лексический анализатор
	lexer := NewLexer(sql)

	// Токены лексического анализа
	tokens, err := lexer.Analyze()

	if err != nil {
		return "", err
	}

	// Парсер
	parser := NewParser(tokens)

	// Подготовленный запрос
	sqlQuery, err := parser.parse()

	if err != nil {
		return "", err
	}

	// Executor
	executor := NewExecutor(sqlQuery)

	// Выполняем запрос
	err = executor.executeQuery()

	if err != nil {
		return "", err
	}

	duration := time.Since(start)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	statisticData := fmt.Sprintf("время выполнения: %s \nвыделено памяти: %s", duration.String(), humanReadableBytes(memStats.Alloc))

	return statisticData, nil
}

func humanReadableBytes(bytes uint64) string {
	const uint = 1024

	if bytes < uint {
		return fmt.Sprintf("%d B", bytes)
	}

	if bytes < uint * uint {
		return fmt.Sprintf("%d KB", bytes/uint)
	}

	return fmt.Sprintf("%d MB", bytes/(uint * uint))
}
