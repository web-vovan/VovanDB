package parser

// Парсер
type Parser struct {
	Tokens []Token
	CurrentToken Token
	Position int
}

// Новый парсер
func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens: tokens,
		CurrentToken: tokens[0],
		Position: 0,
	}
}

// Получение текущего токена
func (p *Parser) current() Token {
	return p.CurrentToken
}

// Получение текущего токена и переход к следующему
func (p *Parser) next() Token {
	currentToken := p.current()

	p.Position++

	if p.Position >= len(p.Tokens) {
		p.CurrentToken = Token{}
	} else {
		p.CurrentToken = p.Tokens[p.Position]	
	}

	return currentToken
}

// Конец списка токенов
func (p *Parser) isEnd() bool {
	return p.CurrentToken.Value == ""
}

func (p *Parser) isSelectQuery() bool {
    t := p.current()
    
    if (t.Type == TYPE_KEYWORD && t.Value == "SELECT") {
        p.next()
        return true
    }

    return false
}