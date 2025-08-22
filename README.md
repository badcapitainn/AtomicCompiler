# SimpleLang Compiler

A simple educational programming language compiler written in Go, designed to teach fundamental programming concepts.

## What is SimpleLang?

SimpleLang is a beginner-friendly programming language that covers essential programming fundamentals:

- **Data Types**: integers, strings, booleans
- **Variables**: declaration and assignment
- **Control Flow**: if statements, loops
- **Functions**: basic function definitions and calls
- **Input/Output**: print statements and user input

## Language Features

### Data Types
```
number x = 42
text name = "Hello World"
boolean isTrue = true
```

### Variables
```
number age = 25
text message = "Welcome to SimpleLang!"
```

### Control Flow
```
if age > 18 then
    print "You are an adult"
else
    print "You are a minor"
end

loop i from 1 to 5
    print i
end
```

### Functions
```
function greet(text name)
    print "Hello " + name
end

greet("Alice")
```

## Project Structure

```
├── cmd/
│   └── compiler/          # Main compiler executable
├── internal/
│   ├── lexer/            # Lexical analysis
│   ├── parser/           # Syntax parsing
│   ├── ast/              # Abstract Syntax Tree
│   ├── interpreter/      # Code execution
│   └── types/            # Type system
├── examples/              # Sample SimpleLang programs
└── tests/                # Test files
```

## Getting Started

### Prerequisites
- Go 1.21 or later

### Installation
```bash
git clone <repository>
cd SimpleLang
go mod tidy
```

### Running the Compiler
```bash
go run cmd/compiler/main.go examples/hello.sl
```

### Building
```bash
go build -o simplelang cmd/compiler/main.go
```

## Example Programs

Check the `examples/` directory for sample SimpleLang programs that demonstrate various language features.

## Learning Objectives

This compiler is designed to help students understand:

1. **Lexical Analysis**: How source code is broken into tokens
2. **Parsing**: How tokens are organized into a syntax tree
3. **Type Checking**: How programming languages ensure type safety
4. **Code Generation**: How high-level code becomes executable instructions
5. **Runtime**: How programs execute and manage memory

## Contributing

Feel free to extend the language with additional features like:
- Arrays and lists
- More complex data structures
- Object-oriented features
- Standard library functions

## License

MIT License - feel free to use this for educational purposes!
