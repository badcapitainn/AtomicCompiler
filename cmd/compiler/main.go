package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"simplelang/internal/interpreter"
	"simplelang/internal/lexer"
	"simplelang/internal/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: simplelang <source_file>")
		fmt.Println("Example: simplelang examples/hello.sl")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Read source file
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Printf("Compiling and running: %s\n", filename)
	fmt.Println("=" + string(make([]byte, 50, 50)) + "=")

	// Step 1: Lexical Analysis (Tokenization)
	fmt.Println("Step 1: Lexical Analysis...")
	lex := lexer.NewLexer(string(source))
	tokens, err := lex.Tokenize()
	if err != nil {
		fmt.Printf("Lexical error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Generated %d tokens\n", len(tokens)-1) // -1 for EOF token

	// Step 2: Parsing (Syntax Analysis)
	fmt.Println("Step 2: Parsing...")
	parser := parser.NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Parsed %d statements\n", len(ast.Statements))

	// Step 3: Interpretation (Execution)
	fmt.Println("Step 3: Execution...")
	interpreter := interpreter.NewInterpreter()
	err = interpreter.Interpret(ast)
	if err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Program executed successfully!")
}
