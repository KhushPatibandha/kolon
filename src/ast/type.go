package ast

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

type Type struct {
	Kind  TypeKind
	Token *lexer.Token

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
		return t.Name
	case TypeArray:
		return fmt.Sprintf("%s[]", t.ElementType.String())
	default:
		return fmt.Sprintf("%s[%s]", t.KeyType.String(), t.ValueType.String())
	}
}

func (t *Type) Equals(other *Type) bool {
	if other == nil || t.Kind != other.Kind {
		return false
	}
	switch t.Kind {
	case TypeBase:
		return t.Name == other.Name
	case TypeArray:
		return t.ElementType.Equals(other.ElementType)
	default:
		return t.KeyType.Equals(other.KeyType) && t.ValueType.Equals(other.ValueType)
	}
}
