package interpreter

import (
	"fmt"
	"math"
	"simplelang/internal/ast"
	"simplelang/internal/types"
)

// Environment represents the execution environment
type Environment struct {
	variables map[string]types.Value
	functions map[string]*ast.FunctionDeclaration
	parent    *Environment
}

// NewEnvironment creates a new environment
func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		variables: make(map[string]types.Value),
		functions: make(map[string]*ast.FunctionDeclaration),
		parent:    parent,
	}
}

// SetVariable sets a variable in the current environment
func (e *Environment) SetVariable(name string, value types.Value) {
	e.variables[name] = value
}

// GetVariable gets a variable from the current environment or parent
func (e *Environment) GetVariable(name string) (types.Value, bool) {
	if value, exists := e.variables[name]; exists {
		return value, true
	}
	if e.parent != nil {
		return e.parent.GetVariable(name)
	}
	return nil, false
}

// SetFunction sets a function in the current environment
func (e *Environment) SetFunction(name string, function *ast.FunctionDeclaration) {
	e.functions[name] = function
}

// GetFunction gets a function from the current environment or parent
func (e *Environment) GetFunction(name string) (*ast.FunctionDeclaration, bool) {
	if function, exists := e.functions[name]; exists {
		return function, true
	}
	if e.parent != nil {
		return e.parent.GetFunction(name)
	}
	return nil, false
}

// Interpreter executes the AST
type Interpreter struct {
	environment *Environment
}

// NewInterpreter creates a new interpreter
func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(nil),
	}
}

// Interpret executes a program
func (i *Interpreter) Interpret(program *ast.Program) error {
	for _, statement := range program.Statements {
		_, err := i.executeStatement(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

// executeStatement executes a single statement
func (i *Interpreter) executeStatement(statement ast.Statement) (types.Value, error) {
	switch stmt := statement.(type) {
	case *ast.VariableDeclaration:
		return i.executeVariableDeclaration(stmt)
	case *ast.Assignment:
		return i.executeAssignment(stmt)
	case *ast.IfStatement:
		return i.executeIfStatement(stmt)
	case *ast.LoopStatement:
		return i.executeLoopStatement(stmt)
	case *ast.FunctionDeclaration:
		return i.executeFunctionDeclaration(stmt)
	case *ast.PrintStatement:
		return i.executePrintStatement(stmt)
	default:
		return nil, fmt.Errorf("unknown statement type: %T", statement)
	}
}

// executeVariableDeclaration executes a variable declaration
func (i *Interpreter) executeVariableDeclaration(stmt *ast.VariableDeclaration) (types.Value, error) {
	value, err := i.evaluateExpression(stmt.Value)
	if err != nil {
		return nil, err
	}

	// Type checking
	if !stmt.Type.IsCompatibleWith(value.Type()) {
		return nil, fmt.Errorf("type mismatch: cannot assign %s to variable of type %s", value.Type().String(), stmt.Type.String())
	}

	i.environment.SetVariable(stmt.Name, value)
	return value, nil
}

// executeAssignment executes a variable assignment
func (i *Interpreter) executeAssignment(stmt *ast.Assignment) (types.Value, error) {
	value, err := i.evaluateExpression(stmt.Value)
	if err != nil {
		return nil, err
	}

	// Check if variable exists
	if _, exists := i.environment.GetVariable(stmt.Name); !exists {
		return nil, fmt.Errorf("undefined variable: %s", stmt.Name)
	}

	i.environment.SetVariable(stmt.Name, value)
	return value, nil
}

// executeIfStatement executes an if statement
func (i *Interpreter) executeIfStatement(stmt *ast.IfStatement) (types.Value, error) {
	condition, err := i.evaluateExpression(stmt.Condition)
	if err != nil {
		return nil, err
	}

	// Check if condition is boolean
	if _, ok := condition.Type().(types.BooleanType); !ok {
		return nil, fmt.Errorf("condition must be boolean, got %s", condition.Type().String())
	}

	booleanValue := condition.(types.BooleanValue)
	if booleanValue.Value {
		// Execute then body
		for _, statement := range stmt.ThenBody {
			_, err := i.executeStatement(statement)
			if err != nil {
				return nil, err
			}
		}
	} else {
		// Execute else body
		for _, statement := range stmt.ElseBody {
			_, err := i.executeStatement(statement)
			if err != nil {
				return nil, err
			}
		}
	}

	return types.VoidValue{}, nil
}

// executeLoopStatement executes a loop statement
func (i *Interpreter) executeLoopStatement(stmt *ast.LoopStatement) (types.Value, error) {
	fromValue, err := i.evaluateExpression(stmt.From)
	if err != nil {
		return nil, err
	}

	toValue, err := i.evaluateExpression(stmt.To)
	if err != nil {
		return nil, err
	}

	// Check if both values are numbers
	if _, ok := fromValue.Type().(types.NumberType); !ok {
		return nil, fmt.Errorf("loop bounds must be numbers")
	}
	if _, ok := toValue.Type().(types.NumberType); !ok {
		return nil, fmt.Errorf("loop bounds must be numbers")
	}

	from := fromValue.(types.NumberValue).Value
	to := toValue.(types.NumberValue).Value

	// Create new environment for loop variables
	loopEnv := NewEnvironment(i.environment)
	oldEnv := i.environment
	i.environment = loopEnv

	defer func() {
		i.environment = oldEnv
	}()

	for j := from; j <= to; j++ {
		// Set loop variable
		loopEnv.SetVariable(stmt.Variable, types.NumberValue{Value: j})

		// Execute loop body
		for _, statement := range stmt.Body {
			_, err := i.executeStatement(statement)
			if err != nil {
				return nil, err
			}
		}
	}

	return types.VoidValue{}, nil
}

// executeFunctionDeclaration executes a function declaration
func (i *Interpreter) executeFunctionDeclaration(stmt *ast.FunctionDeclaration) (types.Value, error) {
	i.environment.SetFunction(stmt.Name, stmt)
	return types.VoidValue{}, nil
}

// executePrintStatement executes a print statement
func (i *Interpreter) executePrintStatement(stmt *ast.PrintStatement) (types.Value, error) {
	value, err := i.evaluateExpression(stmt.Value)
	if err != nil {
		return nil, err
	}

	fmt.Println(value.String())
	return types.VoidValue{}, nil
}

// evaluateExpression evaluates an expression
func (i *Interpreter) evaluateExpression(expr ast.Expression) (types.Value, error) {
	switch e := expr.(type) {
	case *ast.Literal:
		return i.evaluateLiteral(e)
	case *ast.Identifier:
		return i.evaluateIdentifier(e)
	case *ast.BinaryExpression:
		return i.evaluateBinaryExpression(e)
	case *ast.UnaryExpression:
		return i.evaluateUnaryExpression(e)
	case *ast.FunctionCall:
		return i.evaluateFunctionCall(e)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

// evaluateLiteral evaluates a literal
func (i *Interpreter) evaluateLiteral(lit *ast.Literal) (types.Value, error) {
	switch lit.Type.(type) {
	case types.NumberType:
		if str, ok := lit.Value.(string); ok {
			var num float64
			_, err := fmt.Sscanf(str, "%f", &num)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", str)
			}
			return types.NumberValue{Value: num}, nil
		}
		return nil, fmt.Errorf("invalid number literal")
	case types.TextType:
		if str, ok := lit.Value.(string); ok {
			return types.TextValue{Value: str}, nil
		}
		return nil, fmt.Errorf("invalid text literal")
	case types.BooleanType:
		if b, ok := lit.Value.(bool); ok {
			return types.BooleanValue{Value: b}, nil
		}
		return nil, fmt.Errorf("invalid boolean literal")
	default:
		return nil, fmt.Errorf("unknown literal type: %s", lit.Type.String())
	}
}

// evaluateIdentifier evaluates an identifier
func (i *Interpreter) evaluateIdentifier(ident *ast.Identifier) (types.Value, error) {
	value, exists := i.environment.GetVariable(ident.Name)
	if !exists {
		return nil, fmt.Errorf("undefined variable: %s", ident.Name)
	}
	return value, nil
}

// evaluateBinaryExpression evaluates a binary expression
func (i *Interpreter) evaluateBinaryExpression(expr *ast.BinaryExpression) (types.Value, error) {
	left, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case "+":
		return i.add(left, right)
	case "-":
		return i.subtract(left, right)
	case "*":
		return i.multiply(left, right)
	case "/":
		return i.divide(left, right)
	case "==":
		return i.equal(left, right)
	case "!=":
		return i.notEqual(left, right)
	case "<":
		return i.lessThan(left, right)
	case "<=":
		return i.lessEqual(left, right)
	case ">":
		return i.greaterThan(left, right)
	case ">=":
		return i.greaterEqual(left, right)
	case "and":
		return i.logicalAnd(left, right)
	case "or":
		return i.logicalOr(left, right)
	default:
		return nil, fmt.Errorf("unknown binary operator: %s", expr.Operator)
	}
}

// evaluateUnaryExpression evaluates a unary expression
func (i *Interpreter) evaluateUnaryExpression(expr *ast.UnaryExpression) (types.Value, error) {
	operand, err := i.evaluateExpression(expr.Operand)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case "-":
		if _, ok := operand.Type().(types.NumberType); !ok {
			return nil, fmt.Errorf("cannot negate non-number value")
		}
		num := operand.(types.NumberValue)
		return types.NumberValue{Value: -num.Value}, nil
	case "!":
		if _, ok := operand.Type().(types.BooleanType); !ok {
			return nil, fmt.Errorf("cannot negate non-boolean value")
		}
		b := operand.(types.BooleanValue)
		return types.BooleanValue{Value: !b.Value}, nil
	default:
		return nil, fmt.Errorf("unknown unary operator: %s", expr.Operator)
	}
}

// evaluateFunctionCall evaluates a function call
func (i *Interpreter) evaluateFunctionCall(call *ast.FunctionCall) (types.Value, error) {
	function, exists := i.environment.GetFunction(call.Name)
	if !exists {
		return nil, fmt.Errorf("undefined function: %s", call.Name)
	}

	// Evaluate arguments
	var args []types.Value
	for _, arg := range call.Arguments {
		value, err := i.evaluateExpression(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, value)
	}

	// Check argument count
	if len(args) != len(function.Parameters) {
		return nil, fmt.Errorf("function %s expects %d arguments, got %d", call.Name, len(function.Parameters), len(args))
	}

	// Create new environment for function execution
	funcEnv := NewEnvironment(i.environment)

	// Set parameters
	for j, param := range function.Parameters {
		// Type checking
		if !param.Type.IsCompatibleWith(args[j].Type()) {
			return nil, fmt.Errorf("type mismatch in function %s: parameter %s expects %s, got %s",
				call.Name, param.Name, param.Type.String(), args[j].Type().String())
		}
		funcEnv.SetVariable(param.Name, args[j])
	}

	// Execute function body
	oldEnv := i.environment
	i.environment = funcEnv

	defer func() {
		i.environment = oldEnv
	}()

	for _, statement := range function.Body {
		_, err := i.executeStatement(statement)
		if err != nil {
			return nil, err
		}
	}

	return types.VoidValue{}, nil
}

// Arithmetic operations
func (i *Interpreter) add(left, right types.Value) (types.Value, error) {
	// Number + Number = Number
	if _, ok := left.Type().(types.NumberType); ok {
		if _, ok := right.Type().(types.NumberType); ok {
			l := left.(types.NumberValue).Value
			r := right.(types.NumberValue).Value
			return types.NumberValue{Value: l + r}, nil
		}
	}

	// Text + Text = Text (concatenation)
	if _, ok := left.Type().(types.TextType); ok {
		if _, ok := right.Type().(types.TextType); ok {
			l := left.(types.TextValue).Value
			r := right.(types.TextValue).Value
			return types.TextValue{Value: l + r}, nil
		}
	}

	// Text + Number = Text (concatenation with number converted to string)
	if _, ok := left.Type().(types.TextType); ok {
		if _, ok := right.Type().(types.NumberType); ok {
			l := left.(types.TextValue).Value
			r := right.(types.NumberValue).Value
			return types.TextValue{Value: l + fmt.Sprintf("%g", r)}, nil
		}
	}

	// Number + Text = Text (concatenation with number converted to string)
	if _, ok := left.Type().(types.NumberType); ok {
		if _, ok := right.Type().(types.TextType); ok {
			l := left.(types.NumberValue).Value
			r := right.(types.TextValue).Value
			return types.TextValue{Value: fmt.Sprintf("%g", l) + r}, nil
		}
	}

	return nil, fmt.Errorf("cannot add %s and %s", left.Type().String(), right.Type().String())
}

func (i *Interpreter) subtract(left, right types.Value) (types.Value, error) {
	if _, ok := left.Type().(types.NumberType); ok {
		if _, ok := right.Type().(types.NumberType); ok {
			l := left.(types.NumberValue).Value
			r := right.(types.NumberValue).Value
			return types.NumberValue{Value: l - r}, nil
		}
	}
	return nil, fmt.Errorf("cannot subtract %s from %s", right.Type().String(), left.Type().String())
}

func (i *Interpreter) multiply(left, right types.Value) (types.Value, error) {
	if _, ok := left.Type().(types.NumberType); ok {
		if _, ok := right.Type().(types.NumberType); ok {
			l := left.(types.NumberValue).Value
			r := right.(types.NumberValue).Value
			return types.NumberValue{Value: l * r}, nil
		}
	}
	return nil, fmt.Errorf("cannot multiply %s and %s", left.Type().String(), right.Type().String())
}

func (i *Interpreter) divide(left, right types.Value) (types.Value, error) {
	if _, ok := left.Type().(types.NumberType); ok {
		if _, ok := right.Type().(types.NumberType); ok {
			l := left.(types.NumberValue).Value
			r := right.(types.NumberValue).Value
			if r == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return types.NumberValue{Value: l / r}, nil
		}
	}
	return nil, fmt.Errorf("cannot divide %s by %s", left.Type().String(), right.Type().String())
}

// Comparison operations
func (i *Interpreter) equal(left, right types.Value) (types.Value, error) {
	if left.Type() != right.Type() {
		return types.BooleanValue{Value: false}, nil
	}

	switch l := left.(type) {
	case types.NumberValue:
		r := right.(types.NumberValue)
		return types.BooleanValue{Value: math.Abs(l.Value-r.Value) < 1e-9}, nil
	case types.TextValue:
		r := right.(types.TextValue)
		return types.BooleanValue{Value: l.Value == r.Value}, nil
	case types.BooleanValue:
		r := right.(types.BooleanValue)
		return types.BooleanValue{Value: l.Value == r.Value}, nil
	default:
		return types.BooleanValue{Value: false}, nil
	}
}

func (i *Interpreter) notEqual(left, right types.Value) (types.Value, error) {
	result, err := i.equal(left, right)
	if err != nil {
		return nil, err
	}
	return types.BooleanValue{Value: !result.(types.BooleanValue).Value}, nil
}

func (i *Interpreter) lessThan(left, right types.Value) (types.Value, error) {
	if _, ok := left.Type().(types.NumberType); ok {
		if _, ok := right.Type().(types.NumberType); ok {
			l := left.(types.NumberValue).Value
			r := right.(types.NumberValue).Value
			return types.BooleanValue{Value: l < r}, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %s and %s", left.Type().String(), right.Type().String())
}

func (i *Interpreter) lessEqual(left, right types.Value) (types.Value, error) {
	if _, ok := left.Type().(types.NumberType); ok {
		if _, ok := right.Type().(types.NumberType); ok {
			l := left.(types.NumberValue).Value
			r := right.(types.NumberValue).Value
			return types.BooleanValue{Value: l <= r}, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %s and %s", left.Type().String(), right.Type().String())
}

func (i *Interpreter) greaterThan(left, right types.Value) (types.Value, error) {
	if _, ok := left.Type().(types.NumberType); ok {
		if _, ok := right.Type().(types.NumberType); ok {
			l := left.(types.NumberValue).Value
			r := right.(types.NumberValue).Value
			return types.BooleanValue{Value: l > r}, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %s and %s", left.Type().String(), right.Type().String())
}

func (i *Interpreter) greaterEqual(left, right types.Value) (types.Value, error) {
	if _, ok := left.Type().(types.NumberType); ok {
		if _, ok := right.Type().(types.NumberType); ok {
			l := left.(types.NumberValue).Value
			r := right.(types.NumberValue).Value
			return types.BooleanValue{Value: l >= r}, nil
		}
	}
	return nil, fmt.Errorf("cannot compare %s and %s", left.Type().String(), right.Type().String())
}

// Logical operations
func (i *Interpreter) logicalAnd(left, right types.Value) (types.Value, error) {
	if _, ok := left.Type().(types.BooleanType); ok {
		if _, ok := right.Type().(types.BooleanType); ok {
			l := left.(types.BooleanValue).Value
			r := right.(types.BooleanValue).Value
			return types.BooleanValue{Value: l && r}, nil
		}
	}
	return nil, fmt.Errorf("cannot perform logical AND on %s and %s", left.Type().String(), right.Type().String())
}

func (i *Interpreter) logicalOr(left, right types.Value) (types.Value, error) {
	if _, ok := left.Type().(types.BooleanType); ok {
		if _, ok := right.Type().(types.BooleanType); ok {
			l := left.(types.BooleanValue).Value
			r := right.(types.BooleanValue).Value
			return types.BooleanValue{Value: l || r}, nil
		}
	}
	return nil, fmt.Errorf("cannot perform logical OR on %s and %s", left.Type().String(), right.Type().String())
}
