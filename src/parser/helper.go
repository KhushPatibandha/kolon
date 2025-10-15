package parser

import (
	"errors"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Helper Methods
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	if p.tokenPtr >= len(p.tokens) {
		p.peekToken = lexer.Token{Kind: lexer.EOF}
	} else {
		p.peekToken = p.tokens[p.tokenPtr]
	}
	p.tokenPtr++
}

func (p *Parser) expectedPeekToken(kind lexer.TokenKind) bool {
	if p.peekTokenIsOk(kind) {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) currTokenIsOk(kind lexer.TokenKind) bool {
	return p.currToken.Kind == kind
}

func (p *Parser) peekTokenIsOk(kind lexer.TokenKind) bool {
	return p.peekToken.Kind == kind
}

func (p *Parser) addPrefix(tokenKind lexer.TokenKind, fn prefixParseFn) {
	p.prefixParseFns[tokenKind] = fn
}

func (p *Parser) addInfix(tokenKind lexer.TokenKind, fn infixParseFn) {
	p.infixParseFns[tokenKind] = fn
}

func (p *Parser) addPostfix(tokenKind lexer.TokenKind, fn postfixParseFn) {
	p.postfixParseFns[tokenKind] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Kind]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currToken.Kind]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) compareFunctionSig(f1, f2 *ast.Function) bool {
	if !f1.Name.Equals(f2.Name) ||
		f1.Parameters == nil && f2.Parameters != nil ||
		f1.Parameters != nil && f2.Parameters == nil ||
		len(f1.Parameters) != len(f2.Parameters) ||
		f1.ReturnTypes == nil && f2.ReturnTypes != nil ||
		f1.ReturnTypes != nil && f2.ReturnTypes == nil ||
		len(f1.ReturnTypes) != len(f2.ReturnTypes) {
		return false
	}
	for i := range f1.Parameters {
		if !f1.Parameters[i].Equals(f2.Parameters[i]) {
			return false
		}
	}
	for i := range f1.ReturnTypes {
		if !f1.ReturnTypes[i].Equals(f2.ReturnTypes[i]) {
			return false
		}
	}
	return true
}

func (p *Parser) assignDefaultValue(t *ast.Type) (ast.Expression, error) {
	switch t.Name {
	case "int":
		return &ast.Integer{Token: lexer.Token{Kind: lexer.INT, Value: "0"}, Value: 0}, nil
	case "float":
		return &ast.Float{Token: lexer.Token{Kind: lexer.FLOAT, Value: "0.0"}, Value: 0.0}, nil
	case "bool":
		return &ast.Bool{Token: lexer.Token{Kind: lexer.BOOL, Value: "false"}, Value: false}, nil
	case "string":
		return &ast.String{Token: lexer.Token{Kind: lexer.STRING, Value: "\"\""}, Value: "\"\""}, nil
	case "char":
		return &ast.Char{Token: lexer.Token{Kind: lexer.CHAR, Value: "''"}, Value: "''"}, nil
	default:
		return nil, errors.New("can only assign default values to `int`, `float`, `bool`, `string` and `char` types")
	}
}

func (p *Parser) assignTypeToValue(exp ast.Expression, t *ast.Type) {
	switch t.Kind {
	case ast.TypeArray:
		if arr, ok := exp.(*ast.Array); ok {
			arr.Type = t.ElementType
		}
	case ast.TypeHashMap:
		if hash, ok := exp.(*ast.HashMap); ok {
			hash.KeyType = t.KeyType
			hash.ValueType = t.ValueType
		}
	default:
		return
	}
}
