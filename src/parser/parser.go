package parser

import (
	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

type Parser struct {
	tokens    []lexer.Token
	tokenPtr  int
	currToken lexer.Token
	peekToken lexer.Token

	functionMap map[string]*ast.Function
	builtinMap  map[string]bool
	inLoop      bool
	inFunction  bool

	inTesting bool

	prefixParseFns  map[lexer.TokenKind]prefixParseFn
	infixParseFns   map[lexer.TokenKind]infixParseFn
	postfixParseFns map[lexer.TokenKind]postfixParseFn
}

func New(tokens []lexer.Token, inTesting bool) *Parser {
	p := &Parser{
		tokens:      tokens,
		tokenPtr:    1,
		peekToken:   tokens[0],
		inTesting:   inTesting,
		inLoop:      false,
		inFunction:  false,
		functionMap: make(map[string]*ast.Function),
		builtinMap: map[string]bool{
			"print":       true,
			"println":     true,
			"scan":        true,
			"scanln":      true,
			"len":         true,
			"toString":    true,
			"toFloat":     true,
			"toInt":       true,
			"push":        true,
			"pop":         true,
			"insert":      true,
			"remove":      true,
			"getIndex":    true,
			"keys":        true,
			"values":      true,
			"containsKey": true,
			"typeOf":      true,
			"slice":       true,
		},
	}
	p.nextToken()

	p.prefixParseFns = make(map[lexer.TokenKind]prefixParseFn)
	p.addPrefix(lexer.IDENTIFIER, p.parseIdentifier)
	p.addPrefix(lexer.INT, p.parseInteger)
	p.addPrefix(lexer.FLOAT, p.parseFloat)
	p.addPrefix(lexer.BOOL, p.parseBoolean)
	p.addPrefix(lexer.STRING, p.parseString)
	p.addPrefix(lexer.CHAR, p.parseChar)
	p.addPrefix(lexer.NOT, p.parsePrefix)
	p.addPrefix(lexer.DASH, p.parsePrefix)
	p.addPrefix(lexer.OPEN_BRACKET, p.parseGroupedExp)
	p.addPrefix(lexer.OPEN_SQUARE_BRACKET, p.parseArray)
	p.addPrefix(lexer.OPEN_CURLY_BRACKET, p.parseHashMap)

	p.infixParseFns = make(map[lexer.TokenKind]infixParseFn)
	p.addInfix(lexer.PLUS, p.parseInfix)
	p.addInfix(lexer.DASH, p.parseInfix)
	p.addInfix(lexer.SLASH, p.parseInfix)
	p.addInfix(lexer.STAR, p.parseInfix)
	p.addInfix(lexer.PERCENT, p.parseInfix)
	p.addInfix(lexer.DOUBLE_EQUAL, p.parseInfix)
	p.addInfix(lexer.NOT_EQUAL, p.parseInfix)
	p.addInfix(lexer.LESS_THAN_EQUAL, p.parseInfix)
	p.addInfix(lexer.GREATER_THAN_EQUAL, p.parseInfix)
	p.addInfix(lexer.LESS_THAN, p.parseInfix)
	p.addInfix(lexer.GREATER_THAN, p.parseInfix)
	p.addInfix(lexer.AND_AND, p.parseInfix)
	p.addInfix(lexer.OR_OR, p.parseInfix)
	p.addInfix(lexer.AND, p.parseInfix)
	p.addInfix(lexer.OR, p.parseInfix)
	p.addInfix(lexer.OPEN_BRACKET, p.parseCall)
	p.addInfix(lexer.EQUAL_ASSIGN, p.parseAssignment)
	p.addInfix(lexer.PLUS_EQUAL, p.parseAssignment)
	p.addInfix(lexer.MINUS_EQUAL, p.parseAssignment)
	p.addInfix(lexer.STAR_EQUAL, p.parseAssignment)
	p.addInfix(lexer.SLASH_EQUAL, p.parseAssignment)
	p.addInfix(lexer.PERCENT_EQUAL, p.parseAssignment)
	p.addInfix(lexer.OPEN_SQUARE_BRACKET, p.parseIndex)

	p.postfixParseFns = make(map[lexer.TokenKind]postfixParseFn)
	p.addPostfix(lexer.PLUS_PLUS, p.parsePostfix)
	p.addPostfix(lexer.MINUS_MINUS, p.parsePostfix)

	return p
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.currTokenIsOk(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	return program, nil
}
