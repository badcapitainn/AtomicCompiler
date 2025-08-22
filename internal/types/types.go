package types

import "fmt"

// Type represents a SimpleLang data type
type Type interface {
	String() string
	IsCompatibleWith(other Type) bool
}

// Basic types
type NumberType struct{}
type TextType struct{}
type BooleanType struct{}
type VoidType struct{}

func (n NumberType) String() string  { return "number" }
func (t TextType) String() string    { return "text" }
func (b BooleanType) String() string { return "boolean" }
func (v VoidType) String() string    { return "void" }

func (n NumberType) IsCompatibleWith(other Type) bool {
	switch other.(type) {
	case NumberType:
		return true
	default:
		return false
	}
}

func (t TextType) IsCompatibleWith(other Type) bool {
	switch other.(type) {
	case TextType:
		return true
	default:
		return false
	}
}

func (b BooleanType) IsCompatibleWith(other Type) bool {
	switch other.(type) {
	case BooleanType:
		return true
	default:
		return false
	}
}

func (v VoidType) IsCompatibleWith(other Type) bool {
	return true
}

// TypeFromString converts a string representation to a Type
func TypeFromString(typeStr string) (Type, error) {
	switch typeStr {
	case "number":
		return NumberType{}, nil
	case "text":
		return TextType{}, nil
	case "boolean":
		return BooleanType{}, nil
	case "void":
		return VoidType{}, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", typeStr)
	}
}

// Value represents a runtime value
type Value interface {
	Type() Type
	String() string
}

type NumberValue struct {
	Value float64
}

func (n NumberValue) Type() Type     { return NumberType{} }
func (n NumberValue) String() string { return fmt.Sprintf("%g", n.Value) }

type TextValue struct {
	Value string
}

func (t TextValue) Type() Type     { return TextType{} }
func (t TextValue) String() string { return t.Value }

type BooleanValue struct {
	Value bool
}

func (b BooleanValue) Type() Type     { return BooleanType{} }
func (b BooleanValue) String() string { return fmt.Sprintf("%t", b.Value) }

type VoidValue struct{}

func (v VoidValue) Type() Type     { return VoidType{} }
func (v VoidValue) String() string { return "void" }
