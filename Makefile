test:
	go test ./tests ./internal/lexer -v -count=1

test-lexer:
	go test ./internal/lexer -v -count=1