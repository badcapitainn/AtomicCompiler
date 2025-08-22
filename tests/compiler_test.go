package tests

import (
	"simplelang/internal/ast"
	"simplelang/internal/interpreter"
	"simplelang/internal/lexer"
	"simplelang/internal/parser"
	"simplelang/internal/types"
	"testing"
)

func TestLexer(t *testing.T) {
	source := `number x = 42
text message = "Hello World"
boolean flag = true
print x`

	lex := lexer.NewLexer(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}

	// Check that we have the expected number of tokens
	expectedTokens := 15 // including EOF
	if len(tokens) != expectedTokens {
		t.Errorf("Expected %d tokens, got %d", expectedTokens, len(tokens))
	}

	// Check first few tokens
	if tokens[0].Type != lexer.TokenNumberKeyword {
		t.Errorf("Expected TokenNumberKeyword, got %v", tokens[0].Type)
	}
	if tokens[1].Type != lexer.TokenIdentifier {
		t.Errorf("Expected TokenIdentifier, got %v", tokens[1].Type)
	}
	if tokens[2].Type != lexer.TokenAssign {
		t.Errorf("Expected TokenAssign, got %v", tokens[2].Type)
	}
	if tokens[3].Type != lexer.TokenNumber {
		t.Errorf("Expected TokenNumber, got %v", tokens[3].Type)
	}
}

func TestParser(t *testing.T) {
	source := `number x = 42
text message = "Hello World"
print x`

	lex := lexer.NewLexer(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}

	parser := parser.NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	// Check that we have the expected number of statements
	expectedStatements := 3
	if len(program.Statements) != expectedStatements {
		t.Errorf("Expected %d statements, got %d", expectedStatements, len(program.Statements))
	}

	// Check first statement is a variable declaration
	if _, ok := program.Statements[0].(*ast.VariableDeclaration); !ok {
		t.Error("First statement should be a VariableDeclaration")
	}

	// Check second statement is a variable declaration
	if _, ok := program.Statements[1].(*ast.VariableDeclaration); !ok {
		t.Error("Second statement should be a VariableDeclaration")
	}

	// Check third statement is a print statement
	if _, ok := program.Statements[2].(*ast.PrintStatement); !ok {
		t.Error("Third statement should be a PrintStatement")
	}
}

func TestInterpreter(t *testing.T) {
	source := `number x = 10
number y = 5
number result = x + y
print result`

	lex := lexer.NewLexer(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}

	parser := parser.NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	interpreter := interpreter.NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Interpreter failed: %v", err)
	}
}

func TestTypeSystem(t *testing.T) {
	// Test type compatibility
	numberType := types.NumberType{}
	textType := types.TextType{}
	booleanType := types.BooleanType{}

	if !numberType.IsCompatibleWith(types.NumberType{}) {
		t.Error("NumberType should be compatible with NumberType")
	}

	if numberType.IsCompatibleWith(textType) {
		t.Error("NumberType should not be compatible with TextType")
	}

	if !booleanType.IsCompatibleWith(types.BooleanType{}) {
		t.Error("BooleanType should be compatible with BooleanType")
	}

	// Test type from string
	if _, err := types.TypeFromString("number"); err != nil {
		t.Error("Should be able to create NumberType from string")
	}

	if _, err := types.TypeFromString("invalid"); err == nil {
		t.Error("Should not be able to create invalid type from string")
	}
}

func TestArithmetic(t *testing.T) {
	source := `number a = 10
number b = 3
print "Addition: " + (a + b)
print "Subtraction: " + (a - b)
print "Multiplication: " + (a * b)
print "Division: " + (a / b)`

	lex := lexer.NewLexer(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}

	parser := parser.NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	interpreter := interpreter.NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Interpreter failed: %v", err)
	}
}

func TestControlFlow(t *testing.T) {
	source := `number x = 15
if x > 10 then
    print "x is greater than 10"
else
    print "x is less than or equal to 10"
end

loop i from 1 to 3
    print "Loop iteration: " + i
end`

	lex := lexer.NewLexer(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}

	parser := parser.NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	interpreter := interpreter.NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Interpreter failed: %v", err)
	}
}

func TestFunctions(t *testing.T) {
	source := `function add(number a, number b)
    number result = a + b
    print "Result: " + result
end

add(5, 3)
add(10, 20)`

	lex := lexer.NewLexer(source)
	tokens, err := lex.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}

	parser := parser.NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	interpreter := interpreter.NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Interpreter failed: %v", err)
	}
}
