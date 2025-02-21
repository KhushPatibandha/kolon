package lexer

import "fmt"

type TokenKind int

const (
	EOF TokenKind = iota
	STRING
	CHAR
	INT
	FLOAT
	BOOL
	IDENTIFIER
	TYPE
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_SQUARE_BRACKET
	CLOSE_SQUARE_BRACKET
	OPEN_CURLY_BRACKET
	CLOSE_CURLY_BRACKET
	LESS_THAN
	GREATER_THAN
	LESS_THAN_EQUAL
	GREATER_THAN_EQUAL
	EQUAL_ASSIGN
	DOUBLE_EQUAL
	NOT
	NOT_EQUAL
	PLUS
	DASH
	STAR
	SLASH
	PERCENT
	PLUS_PLUS
	PLUS_EQUAL
	MINUS_MINUS
	MINUS_EQUAL
	STAR_EQUAL
	SLASH_EQUAL
	PERCENT_EQUAL
	COLON
	SEMI_COLON
	COMMA
	AND
	OR
	AND_AND
	OR_OR
	VAR
	CONST
	FUN
	IF
	ELSE
	ELSE_IF
	FOR
	WHILE
	RETURN
	CONTINUE
	BREAK
)

var reservedWords = map[string]TokenKind{
	"var":      VAR,
	"const":    CONST,
	"fun":      FUN,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"while":    WHILE,
	"else if":  ELSE_IF,
	"true":     BOOL,
	"false":    BOOL,
	"return":   RETURN,
	"string":   TYPE,
	"char":     TYPE,
	"int":      TYPE,
	"float":    TYPE,
	"bool":     TYPE,
	"continue": CONTINUE,
	"break":    BREAK,
}

type Token struct {
	Kind  TokenKind
	Value string
}

func (token Token) Help() {
	if token.Kind == STRING || token.Kind == INT || token.Kind == BOOL || token.Kind == CHAR || token.Kind == FLOAT || token.Kind == IDENTIFIER || token.Kind == TYPE {
		fmt.Printf("%s(%s)\n", TokenKindString(token.Kind), token.Value)
	} else {
		fmt.Printf("%s()\n", TokenKindString(token.Kind))
	}
}

func GetNewToken(k TokenKind, v string) Token {
	return Token{k, v}
}

func TokenKindString(tKind TokenKind) string {
	switch tKind {
	case EOF:
		return "EOF"
	case STRING:
		return "STRING"
	case CHAR:
		return "CHAR"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case BOOL:
		return "BOOL"
	case TYPE:
		return "TYPE"
	case IDENTIFIER:
		return "IDENTIFIER"
	case OPEN_BRACKET:
		return "OPEN_BRACKET"
	case CLOSE_BRACKET:
		return "CLOSE_BRACKET"
	case OPEN_SQUARE_BRACKET:
		return "OPEN_SQUARE_BRACKET"
	case CLOSE_SQUARE_BRACKET:
		return "CLOSE_SQUARE_BRACKET"
	case OPEN_CURLY_BRACKET:
		return "OPEN_CURLY_BRACKET"
	case CLOSE_CURLY_BRACKET:
		return "CLOSE_CURLY_BRACKET"
	case LESS_THAN:
		return "LESS_THAN"
	case GREATER_THAN:
		return "GREATER_THAN"
	case LESS_THAN_EQUAL:
		return "LESS_THAN_EQUAL"
	case GREATER_THAN_EQUAL:
		return "GREATER_THAN_EQUAL"
	case EQUAL_ASSIGN:
		return "EQUAL_ASSIGN"
	case DOUBLE_EQUAL:
		return "DOUBLE_EQUAL"
	case NOT:
		return "NOT"
	case NOT_EQUAL:
		return "NOT_EQUAL"
	case PLUS:
		return "PLUS"
	case DASH:
		return "DASH"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case PERCENT:
		return "PERCENT"
	case PLUS_PLUS:
		return "PLUS_PLUS"
	case PLUS_EQUAL:
		return "PLUS_EQUAL"
	case MINUS_MINUS:
		return "MINUS_MINUS"
	case MINUS_EQUAL:
		return "MINUS_EQUAL"
	case STAR_EQUAL:
		return "STAR_EQUAL"
	case SLASH_EQUAL:
		return "SLASH_EQUAL"
	case PERCENT_EQUAL:
		return "PERCENT_EQUAL"
	case COLON:
		return "COLON"
	case SEMI_COLON:
		return "SEMI_COLON"
	case COMMA:
		return "COMMA"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case AND_AND:
		return "AND_AND"
	case OR_OR:
		return "OR_OR"
	case VAR:
		return "VAR"
	case CONST:
		return "CONST"
	case FUN:
		return "FUN"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case ELSE_IF:
		return "ELSE_IF"
	case FOR:
		return "FOR"
	case WHILE:
		return "WHILE"
	case RETURN:
		return "RETURN"
	case CONTINUE:
		return "CONTINUE"
	case BREAK:
		return "BREAK"
	default:
		return fmt.Sprintf("unknown(%d)", tKind)
	}
}
