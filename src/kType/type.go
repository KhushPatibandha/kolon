package ktype

import (
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Type
// ------------------------------------------------------------------------------------------------------------------
type TypeKind int

const (
	TypeBase    TypeKind = iota // For BaseTypes -- Refer to BaseType interface
	TypeArray                   // For Array types
	TypeHashMap                 // For HashMap types
)

type TypeCheckResult struct {
	Types   []*Type
	TypeLen int
}

type Type struct {
	Kind  TypeKind
	Token lexer.Token

	// For BaseTypes -- Kind == TypeBase
	// eg: int, float, etc...
	Name string

	// For Array types -- Kind == TypeArray
	ElementType *Type

	// For HashMap types -- Kind == TypeHashMap
	KeyType   *Type
	ValueType *Type
}

func (t *Type) TokenValue() string { return t.Token.Value }
func (t *Type) String() string {
	switch t.Kind {
	case TypeBase:
		if t.Name == "" {
			return "unknown"
		}
		return t.Name
	case TypeArray:
		if t.ElementType == nil {
			return "unknown[]"
		}
		return fmt.Sprintf("%s[]", t.ElementType.String())
	default:
		if t.KeyType == nil && t.ValueType == nil {
			return "unknown[unknown]"
		}
		return fmt.Sprintf("%s[%s]", t.KeyType.String(), t.ValueType.String())
	}
}
