package parser

import (
	"fmt"
	"simplelang/internal/ast"
	"simplelang/internal/lexer"
	"simplelang/internal/types"
)

// Parser converts tokens into an AST
type Parser struct {
	tokens []lexer.Token
	pos    int
}

// NewParser creates a new parser
func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Parse parses the tokens and returns an AST
func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}

	for p.current().Type != lexer.TokenEOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	token := p.current()

	switch token.Type {
	case lexer.TokenNumberKeyword, lexer.TokenTextKeyword, lexer.TokenBooleanKeyword:
		return p.parseVariableDeclaration()
	case lexer.TokenIdentifier:
		// Look ahead to see if this is an assignment
		if p.peek().Type == lexer.TokenAssign {
			return p.parseAssignment()
		}
		return p.parseExpressionStatement()
	case lexer.TokenIf:
		return p.parseIfStatement()
	case lexer.TokenLoop:
		return p.parseLoopStatement()
	case lexer.TokenFunction:
		return p.parseFunctionDeclaration()
	case lexer.TokenPrint:
		return p.parsePrintStatement()
	default:
		return nil, fmt.Errorf("unexpected token at line %d, column %d: %s", token.Line, token.Column, token.Value)
	}
}

func (p *Parser) parseVariableDeclaration() (*ast.VariableDeclaration, error) {
	typeToken := p.current()
	p.advance()

	if p.current().Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("expected identifier after type, got %s", p.current().Value)
	}

	name := p.current().Value
	p.advance()

	if p.current().Type != lexer.TokenAssign {
		return nil, fmt.Errorf("expected '=' after variable name, got %s", p.current().Value)
	}
	p.advance()

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	varType, err := types.TypeFromString(typeToken.Value)
	if err != nil {
		return nil, err
	}

	return &ast.VariableDeclaration{
		Type:  varType,
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) parseAssignment() (*ast.Assignment, error) {
	name := p.current().Value
	p.advance() // consume identifier

	if p.current().Type != lexer.TokenAssign {
		return nil, fmt.Errorf("expected '=' after variable name, got %s", p.current().Value)
	}
	p.advance()

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.Assignment{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) parseIfStatement() (*ast.IfStatement, error) {
	p.advance() // consume 'if'

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.current().Type != lexer.TokenThen {
		return nil, fmt.Errorf("expected 'then' after condition, got %s", p.current().Value)
	}
	p.advance()

	var thenBody []ast.Statement
	for p.current().Type != lexer.TokenElse && p.current().Type != lexer.TokenEnd && p.current().Type != lexer.TokenEOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		thenBody = append(thenBody, stmt)
	}

	var elseBody []ast.Statement
	if p.current().Type == lexer.TokenElse {
		p.advance()
		for p.current().Type != lexer.TokenEnd && p.current().Type != lexer.TokenEOF {
			stmt, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			elseBody = append(elseBody, stmt)
		}
	}

	if p.current().Type != lexer.TokenEnd {
		return nil, fmt.Errorf("expected 'end' after if statement, got %s", p.current().Value)
	}
	p.advance()

	return &ast.IfStatement{
		Condition: condition,
		ThenBody:  thenBody,
		ElseBody:  elseBody,
	}, nil
}

func (p *Parser) parseLoopStatement() (*ast.LoopStatement, error) {
	p.advance() // consume 'loop'

	if p.current().Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("expected identifier after 'loop', got %s", p.current().Value)
	}

	variable := p.current().Value
	p.advance()

	if p.current().Type != lexer.TokenFrom {
		return nil, fmt.Errorf("expected 'from' after loop variable, got %s", p.current().Value)
	}
	p.advance()

	fromExpr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.current().Type != lexer.TokenTo {
		return nil, fmt.Errorf("expected 'to' after 'from' expression, got %s", p.current().Value)
	}
	p.advance()

	toExpr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	var body []ast.Statement
	for p.current().Type != lexer.TokenEnd && p.current().Type != lexer.TokenEOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	if p.current().Type != lexer.TokenEnd {
		return nil, fmt.Errorf("expected 'end' after loop body, got %s", p.current().Value)
	}
	p.advance()

	return &ast.LoopStatement{
		Variable: variable,
		From:     fromExpr,
		To:       toExpr,
		Body:     body,
	}, nil
}

func (p *Parser) parseFunctionDeclaration() (*ast.FunctionDeclaration, error) {
	p.advance() // consume 'function'

	if p.current().Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("expected function name after 'function', got %s", p.current().Value)
	}

	name := p.current().Value
	p.advance()

	if p.current().Type != lexer.TokenLeftParen {
		return nil, fmt.Errorf("expected '(' after function name, got %s", p.current().Value)
	}
	p.advance()

	var parameters []ast.Parameter
	for p.current().Type != lexer.TokenRightParen {
		if len(parameters) > 0 {
			if p.current().Type != lexer.TokenComma {
				return nil, fmt.Errorf("expected ',' between parameters, got %s", p.current().Value)
			}
			p.advance()
		}

		if p.current().Type != lexer.TokenNumberKeyword && p.current().Type != lexer.TokenTextKeyword && p.current().Type != lexer.TokenBooleanKeyword {
			return nil, fmt.Errorf("expected parameter type, got %s", p.current().Value)
		}

		paramType, err := types.TypeFromString(p.current().Value)
		if err != nil {
			return nil, err
		}
		p.advance()

		if p.current().Type != lexer.TokenIdentifier {
			return nil, fmt.Errorf("expected parameter name, got %s", p.current().Value)
		}

		parameters = append(parameters, ast.Parameter{
			Name: p.current().Value,
			Type: paramType,
		})
		p.advance()
	}
	p.advance() // consume ')'

	var body []ast.Statement
	for p.current().Type != lexer.TokenEnd && p.current().Type != lexer.TokenEOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	if p.current().Type != lexer.TokenEnd {
		return nil, fmt.Errorf("expected 'end' after function body, got %s", p.current().Value)
	}
	p.advance()

	return &ast.FunctionDeclaration{
		Name:       name,
		Parameters: parameters,
		ReturnType: types.VoidType{},
		Body:       body,
	}, nil
}

func (p *Parser) parsePrintStatement() (*ast.PrintStatement, error) {
	p.advance() // consume 'print'

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.PrintStatement{
		Value: value,
	}, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseLogicalOr()
}

func (p *Parser) parseLogicalOr() (ast.Expression, error) {
	left, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.TokenOr {
		operator := p.current().Value
		p.advance()

		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseLogicalAnd() (ast.Expression, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.TokenAnd {
		operator := p.current().Value
		p.advance()

		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseEquality() (ast.Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.TokenEqual || p.current().Type == lexer.TokenNotEqual {
		operator := p.current().Value
		p.advance()

		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseComparison() (ast.Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.TokenLessThan || p.current().Type == lexer.TokenLessEqual ||
		p.current().Type == lexer.TokenGreaterThan || p.current().Type == lexer.TokenGreaterEqual {
		operator := p.current().Value
		p.advance()

		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseTerm() (ast.Expression, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.TokenPlus || p.current().Type == lexer.TokenMinus {
		operator := p.current().Value
		p.advance()

		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseFactor() (ast.Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.TokenMultiply || p.current().Type == lexer.TokenDivide {
		operator := p.current().Value
		p.advance()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseUnary() (ast.Expression, error) {
	if p.current().Type == lexer.TokenMinus || p.current().Type == lexer.TokenNot {
		operator := p.current().Value
		p.advance()

		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpression{
			Operator: operator,
			Operand:  operand,
		}, nil
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (ast.Expression, error) {
	token := p.current()

	switch token.Type {
	case lexer.TokenNumber:
		p.advance()
		return &ast.Literal{
			Value: token.Literal,
			Type:  types.NumberType{},
		}, nil

	case lexer.TokenText:
		p.advance()
		return &ast.Literal{
			Value: token.Literal,
			Type:  types.TextType{},
		}, nil

	case lexer.TokenBoolean:
		p.advance()
		return &ast.Literal{
			Value: token.Literal,
			Type:  types.BooleanType{},
		}, nil

	case lexer.TokenIdentifier:
		name := token.Value
		p.advance()

		// Check if this is a function call
		if p.current().Type == lexer.TokenLeftParen {
			return p.parseFunctionCall(name)
		}

		return &ast.Identifier{Name: name}, nil

	case lexer.TokenLeftParen:
		p.advance()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if p.current().Type != lexer.TokenRightParen {
			return nil, fmt.Errorf("expected ')', got %s", p.current().Value)
		}
		p.advance()

		return expr, nil

	default:
		return nil, fmt.Errorf("unexpected token: %s", token.Value)
	}
}

func (p *Parser) parseFunctionCall(name string) (*ast.FunctionCall, error) {
	p.advance() // consume '('

	var arguments []ast.Expression
	for p.current().Type != lexer.TokenRightParen {
		if len(arguments) > 0 {
			if p.current().Type != lexer.TokenComma {
				return nil, fmt.Errorf("expected ',' between arguments, got %s", p.current().Value)
			}
			p.advance()
		}

		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	if p.current().Type != lexer.TokenRightParen {
		return nil, fmt.Errorf("expected ')', got %s", p.current().Value)
	}
	p.advance()

	return &ast.FunctionCall{
		Name:      name,
		Arguments: arguments,
	}, nil
}

func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// For now, we'll just return the expression as a statement
	// In a more sophisticated parser, you might want to handle this differently
	return &ast.PrintStatement{Value: expr}, nil
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek() lexer.Token {
	if p.pos+1 >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) advance() {
	p.pos++
}
