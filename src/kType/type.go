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

func (t *Type) Equals(other *Type) bool {
	if other == nil || t.Kind != other.Kind {
		return false
	}
	switch t.Kind {
	case TypeBase:
		if other.Kind != TypeBase {
			return false
		}
		return t.Name == other.Name
	case TypeArray:
		if other.Kind != TypeArray {
			return false
		}
		if other.ElementType == nil {
			other.ElementType = t.ElementType
			return true
		}
		return t.ElementType.Equals(other.ElementType)
	default:
		if other.Kind != TypeHashMap {
			return false
		}
		if other.KeyType == nil && other.ValueType == nil {
			other.KeyType = t.KeyType
			other.ValueType = t.ValueType
			return true
		}
		return t.KeyType.Equals(other.KeyType) && t.ValueType.Equals(other.ValueType)
	}
}

func (t *Type) PrintTypeHelper() {
	switch t.Kind {
	case TypeBase:
		fmt.Printf("TypeBase: %s\n", t.Name)
	case TypeArray:
		fmt.Printf("TypeArray:\n")
		fmt.Printf("  ElementType:\n")
		t.ElementType.PrintTypeHelper()
	case TypeHashMap:
		fmt.Printf("TypeHashMap:\n")
		fmt.Printf("  KeyType:\n")
		t.KeyType.PrintTypeHelper()
		fmt.Printf("  ValueType:\n")
		t.ValueType.PrintTypeHelper()
	default:
		fmt.Println("Unknown TypeKind")
	}
}

func (t *Type) TypeKindToString() string {
	switch t.Kind {
	case TypeBase:
		return "TypeBase"
	case TypeArray:
		return "TypeArray"
	case TypeHashMap:
		return "TypeHashMap"
	default:
		return "UnknownTypeKind"
	}
}
