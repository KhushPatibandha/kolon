package lexer

import "fmt"

type TokenKind int

const (
	EOF        TokenKind = iota // PERF:
	STRING                      // PERF:
	CHAR                        // PERF:
	INT                         // PERF:
	LONG                        // PERF:
	FLOAT                       // FIX: look into precision
	DOUBLE                      // FIX: look into precision
	BOOL                        // PERF:
	IDENTIFIER                  // PERF:
	// ENUM

	TYPE // PERF:

	PRINT   // PERF:
	PRINTLN // PERF:

	// EXPONENT // TODO: Include in test

	OPEN_BRACKET         // PERF:
	CLOSE_BRACKET        // PERF:
	OPEN_SQUARE_BRACKET  // PERF:
	CLOSE_SQUARE_BRACKET // PERF:
	OPEN_CURLY_BRACKET   // PERF:
	CLOSE_CURLY_BRACKET  // PERF:

	LESS_THAN          // PERF:
	GREATER_THAN       // PERF:
	LESS_THAN_EQUAL    // PERF:
	GREATER_THAN_EQUAL // PERF:

	EQUAL_ASSIGN // PERF:
	DOUBLE_EQUAL // PERF:
	NOT          // PERF:
	NOT_EQUAL    // PERF:

	PLUS    // PERF:
	DASH    // PERF:
	STAR    // PERF:
	SLASH   // PERF:
	PERCENT // PERF:

	PLUS_PLUS     // PERF:
	PLUS_EQUAL    // PERF:
	MINUS_MINUS   // PERF:
	MINUS_EQUAL   // PERF:
	STAR_EQUAL    // PERF:
	SLASH_EQUAL   // PERF:
	PERCENT_EQUAL // PERF:

	COLON         // PERF:
	SEMI_COLON    // PERF:
	DOT           // PERF:
	COMMA         // PERF:
	QUESTION_MARK // PERF:
	DOT_DOT       // PERF:

	AND     // PERF:
	OR      // PERF:
	AND_AND // PERF:
	OR_OR   // PERF:

	IN      // PERF:
	VAR     // PERF:
	CONST   // PERF:
	FUN     // PERF:
	IF      // PERF:
	ELSE    // PERF:
	ELSE_IF // PERF:
	FOR     // PERF:
	RETURN  // PERF:

	// TODO: Can't parse these yet.
	IMPORT
	FROM
	PACKAGE
	STRUCT
	// EXPORT
)

var reservedWords = map[string]TokenKind{
	"var":     VAR,
	"const":   CONST,
	"fun":     FUN,
	"if":      IF,
	"else":    ELSE,
	"for":     FOR,
	"else if": ELSE_IF,
	"println": PRINTLN,
	"print":   PRINT,
	"true":    BOOL,
	"false":   BOOL,
	"return":  RETURN,
	"in":      IN,
	"auto":    TYPE,
	"string":  TYPE,
	"char":    TYPE,
	"int":     TYPE,
	"long":    TYPE,
	"float":   TYPE,
	"double":  TYPE,
	"bool":    TYPE,
	"package": PACKAGE,
	"struct":  STRUCT,
	"import":  IMPORT,
	"from":    FROM,

	// "enum":  ENUM,
	// "export": EXPORT,
}

type Token struct {
	Kind  TokenKind
	Value string
}

func (token Token) Help() {
	if token.Kind == STRING || token.Kind == INT || token.Kind == BOOL || token.Kind == CHAR || token.Kind == LONG || token.Kind == FLOAT || token.Kind == DOUBLE || token.Kind == IDENTIFIER || token.Kind == TYPE || token.Kind == PRINT || token.Kind == PRINTLN || token.Kind == PACKAGE || token.Kind == STRUCT || token.Kind == IMPORT || token.Kind == FROM {
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
	case LONG:
		return "LONG"
	case FLOAT:
		return "FLOAT"
	case DOUBLE:
		return "DOUBLE"
	case BOOL:
		return "BOOL"
		// case ENUM:
		// return "ENUM"
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
	case DOT:
		return "DOT"
	case COMMA:
		return "COMMA"
	case QUESTION_MARK:
		return "QUESTION_MARK"
	case DOT_DOT:
		return "DOT_DOT"
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
	// case EXPONENT:
	// return "EXPONENT"
	case PRINTLN:
		return "PRINTLN"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case IN:
		return "IN"
	case IMPORT:
		return "IMPORT"
	case FROM:
		return "FROM"
	case PACKAGE:
		return "PACKAGE"
	case STRUCT:
		return "STRUCT"
	default:
		return fmt.Sprintf("unknown(%d)", tKind)
	}
}
