package ast

import "simplelang/internal/types"

// Node represents any AST node
type Node interface {
	Accept(visitor Visitor) interface{}
}

// Visitor pattern for AST traversal
type Visitor interface {
	VisitProgram(node *Program) interface{}
	VisitStatement(node Statement) interface{}
	VisitExpression(node Expression) interface{}
	VisitVariableDeclaration(node *VariableDeclaration) interface{}
	VisitAssignment(node *Assignment) interface{}
	VisitIfStatement(node *IfStatement) interface{}
	VisitLoopStatement(node *LoopStatement) interface{}
	VisitFunctionDeclaration(node *FunctionDeclaration) interface{}
	VisitFunctionCall(node *FunctionCall) interface{}
	VisitPrintStatement(node *PrintStatement) interface{}
	VisitBinaryExpression(node *BinaryExpression) interface{}
	VisitUnaryExpression(node *UnaryExpression) interface{}
	VisitLiteral(node *Literal) interface{}
	VisitIdentifier(node *Identifier) interface{}
}

// Program represents the root of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) Accept(visitor Visitor) interface{} {
	return visitor.VisitProgram(p)
}

// Statement represents any statement in the language
type Statement interface {
	Node
	IsStatement()
}

// Expression represents any expression in the language
type Expression interface {
	Node
	IsExpression()
}

// VariableDeclaration represents a variable declaration
type VariableDeclaration struct {
	Type  types.Type
	Name  string
	Value Expression
}

func (v *VariableDeclaration) Accept(visitor Visitor) interface{} {
	return visitor.VisitVariableDeclaration(v)
}

func (v *VariableDeclaration) IsStatement() {}

// Assignment represents a variable assignment
type Assignment struct {
	Name  string
	Value Expression
}

func (a *Assignment) Accept(visitor Visitor) interface{} {
	return visitor.VisitAssignment(a)
}

func (a *Assignment) IsStatement() {}

// IfStatement represents an if-else statement
type IfStatement struct {
	Condition Expression
	ThenBody  []Statement
	ElseBody  []Statement
}

func (i *IfStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitIfStatement(i)
}

func (i *IfStatement) IsStatement() {}

// LoopStatement represents a loop
type LoopStatement struct {
	Variable string
	From     Expression
	To       Expression
	Body     []Statement
}

func (l *LoopStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitLoopStatement(l)
}

func (l *LoopStatement) IsStatement() {}

// FunctionDeclaration represents a function definition
type FunctionDeclaration struct {
	Name       string
	Parameters []Parameter
	ReturnType types.Type
	Body       []Statement
}

type Parameter struct {
	Name string
	Type types.Type
}

func (f *FunctionDeclaration) Accept(visitor Visitor) interface{} {
	return visitor.VisitFunctionDeclaration(f)
}

func (f *FunctionDeclaration) IsStatement() {}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string
	Arguments []Expression
}

func (f *FunctionCall) Accept(visitor Visitor) interface{} {
	return visitor.VisitFunctionCall(f)
}

func (f *FunctionCall) IsExpression() {}

// PrintStatement represents a print statement
type PrintStatement struct {
	Value Expression
}

func (p *PrintStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitPrintStatement(p)
}

func (p *PrintStatement) IsStatement() {}

// BinaryExpression represents a binary operation
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryExpression(b)
}

func (b *BinaryExpression) IsExpression() {}

// UnaryExpression represents a unary operation
type UnaryExpression struct {
	Operator string
	Operand  Expression
}

func (u *UnaryExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpression(u)
}

func (u *UnaryExpression) IsExpression() {}

// Literal represents a literal value
type Literal struct {
	Value interface{}
	Type  types.Type
}

func (l *Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteral(l)
}

func (l *Literal) IsExpression() {}

// Identifier represents a variable reference
type Identifier struct {
	Name string
}

func (i *Identifier) Accept(visitor Visitor) interface{} {
	return visitor.VisitIdentifier(i)
}

func (i *Identifier) IsExpression() {}
