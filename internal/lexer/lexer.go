package lexer

import (
	"fmt"
	"unicode"
)

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenError

	// Literals
	TokenNumber
	TokenText
	TokenBoolean

	// Identifiers
	TokenIdentifier

	// Keywords
	TokenNumberKeyword
	TokenTextKeyword
	TokenBooleanKeyword
	TokenFunction
	TokenIf
	TokenThen
	TokenElse
	TokenEnd
	TokenLoop
	TokenFrom
	TokenTo
	TokenPrint

	// Operators
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenAssign
	TokenEqual
	TokenNotEqual
	TokenLessThan
	TokenLessEqual
	TokenGreaterThan
	TokenGreaterEqual
	TokenAnd
	TokenOr
	TokenNot

	// Delimiters
	TokenLeftParen
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenSemicolon
	TokenColon
)

// Token represents a single token from the source code
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Column  int
	Literal interface{}
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %d, Value: '%s', Line: %d, Column: %d}", t.Type, t.Value, t.Line, t.Column)
}

// Lexer breaks source code into tokens
type Lexer struct {
	input    string
	position int
	line     int
	column   int
	tokens   []Token
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:    input,
		position: 0,
		line:     1,
		column:   1,
		tokens:   []Token{},
	}
}

// Tokenize breaks the input into tokens
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.position < len(l.input) {
		l.skipWhitespace()

		if l.position >= len(l.input) {
			break
		}

		token, err := l.nextToken()
		if err != nil {
			return nil, err
		}

		if token.Type == TokenError {
			return nil, fmt.Errorf("lexical error at line %d, column %d: %s", token.Line, token.Column, token.Value)
		}

		l.tokens = append(l.tokens, token)
	}

	l.tokens = append(l.tokens, Token{Type: TokenEOF, Line: l.line, Column: l.column})
	return l.tokens, nil
}

func (l *Lexer) nextToken() (Token, error) {
	char := l.currentChar()

	switch {
	case unicode.IsDigit(char):
		return l.readNumber(), nil
	case char == '"':
		return l.readText(), nil
	case unicode.IsLetter(char):
		return l.readIdentifierOrKeyword(), nil
	case char == '+':
		l.advance()
		return Token{Type: TokenPlus, Value: "+", Line: l.line, Column: l.column - 1}, nil
	case char == '-':
		l.advance()
		return Token{Type: TokenMinus, Value: "-", Line: l.line, Column: l.column - 1}, nil
	case char == '*':
		l.advance()
		return Token{Type: TokenMultiply, Value: "*", Line: l.line, Column: l.column - 1}, nil
	case char == '/':
		l.advance()
		return Token{Type: TokenDivide, Value: "/", Line: l.line, Column: l.column - 1}, nil
	case char == '=':
		l.advance()
		if l.currentChar() == '=' {
			l.advance()
			return Token{Type: TokenEqual, Value: "==", Line: l.line, Column: l.column - 2}, nil
		}
		return Token{Type: TokenAssign, Value: "=", Line: l.line, Column: l.column - 1}, nil
	case char == '<':
		l.advance()
		if l.currentChar() == '=' {
			l.advance()
			return Token{Type: TokenLessEqual, Value: "<=", Line: l.line, Column: l.column - 2}, nil
		}
		return Token{Type: TokenLessThan, Value: "<", Line: l.line, Column: l.column - 1}, nil
	case char == '>':
		l.advance()
		if l.currentChar() == '=' {
			l.advance()
			return Token{Type: TokenGreaterEqual, Value: ">=", Line: l.line, Column: l.column - 2}, nil
		}
		return Token{Type: TokenGreaterThan, Value: ">", Line: l.line, Column: l.column - 1}, nil
	case char == '!':
		l.advance()
		if l.currentChar() == '=' {
			l.advance()
			return Token{Type: TokenNotEqual, Value: "!=", Line: l.line, Column: l.column - 2}, nil
		}
		return Token{Type: TokenNot, Value: "!", Line: l.line, Column: l.column - 1}, nil
	case char == '(':
		l.advance()
		return Token{Type: TokenLeftParen, Value: "(", Line: l.line, Column: l.column - 1}, nil
	case char == ')':
		l.advance()
		return Token{Type: TokenRightParen, Value: ")", Line: l.line, Column: l.column - 1}, nil
	case char == '{':
		l.advance()
		return Token{Type: TokenLeftBrace, Value: "{", Line: l.line, Column: l.column - 1}, nil
	case char == '}':
		l.advance()
		return Token{Type: TokenRightBrace, Value: "}", Line: l.line, Column: l.column - 1}, nil
	case char == ',':
		l.advance()
		return Token{Type: TokenComma, Value: ",", Line: l.line, Column: l.column - 1}, nil
	case char == ';':
		l.advance()
		return Token{Type: TokenSemicolon, Value: ";", Line: l.line, Column: l.column - 1}, nil
	case char == ':':
		l.advance()
		return Token{Type: TokenColon, Value: ":", Line: l.line, Column: l.column - 1}, nil
	default:
		return Token{Type: TokenError, Value: fmt.Sprintf("unexpected character: %c", char), Line: l.line, Column: l.column}, nil
	}
}

func (l *Lexer) readNumber() Token {
	start := l.position
	startColumn := l.column

	for l.position < len(l.input) && (unicode.IsDigit(l.currentChar()) || l.currentChar() == '.') {
		l.advance()
	}

	value := l.input[start:l.position]
	return Token{
		Type:    TokenNumber,
		Value:   value,
		Line:    l.line,
		Column:  startColumn,
		Literal: value,
	}
}

func (l *Lexer) readText() Token {
	startColumn := l.column
	l.advance() // skip opening quote

	start := l.position
	for l.position < len(l.input) && l.currentChar() != '"' {
		if l.currentChar() == '\n' {
			l.line++
			l.column = 1
		}
		l.advance()
	}

	if l.position >= len(l.input) {
		return Token{
			Type:   TokenError,
			Value:  "unterminated string",
			Line:   l.line,
			Column: startColumn,
		}
	}

	value := l.input[start:l.position]
	l.advance() // skip closing quote

	return Token{
		Type:    TokenText,
		Value:   value,
		Line:    l.line,
		Column:  startColumn,
		Literal: value,
	}
}

func (l *Lexer) readIdentifierOrKeyword() Token {
	start := l.position
	startColumn := l.column

	for l.position < len(l.input) && (unicode.IsLetter(l.currentChar()) || unicode.IsDigit(l.currentChar()) || l.currentChar() == '_') {
		l.advance()
	}

	value := l.input[start:l.position]
	tokenType := l.getKeywordType(value)

	if tokenType == TokenBoolean && (value == "true" || value == "false") {
		return Token{
			Type:    TokenBoolean,
			Value:   value,
			Line:    l.line,
			Column:  startColumn,
			Literal: value == "true",
		}
	}

	return Token{
		Type:    tokenType,
		Value:   value,
		Line:    l.line,
		Column:  startColumn,
		Literal: value,
	}
}

func (l *Lexer) getKeywordType(value string) TokenType {
	switch value {
	case "number":
		return TokenNumberKeyword
	case "text":
		return TokenTextKeyword
	case "boolean":
		return TokenBooleanKeyword
	case "function":
		return TokenFunction
	case "if":
		return TokenIf
	case "then":
		return TokenThen
	case "else":
		return TokenElse
	case "end":
		return TokenEnd
	case "loop":
		return TokenLoop
	case "from":
		return TokenFrom
	case "to":
		return TokenTo
	case "print":
		return TokenPrint
	default:
		return TokenIdentifier
	}
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) && unicode.IsSpace(l.currentChar()) {
		if l.currentChar() == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.advance()
	}
}

func (l *Lexer) currentChar() rune {
	if l.position >= len(l.input) {
		return 0
	}
	return rune(l.input[l.position])
}

func (l *Lexer) advance() {
	l.position++
	l.column++
}
