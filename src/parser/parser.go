package parser

import (
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

type Parser struct {
	tokens       []lexer.Token
	tokenPointer int
	currentToken lexer.Token
	peekToken    lexer.Token
	errors       []string
}

func New(tokens []lexer.Token) *Parser {
	p := &Parser{errors: []string{}, peekToken: tokens[0], tokens: tokens, tokenPointer: 1}
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currentToken.Kind != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Kind {
	case lexer.VAR:
		return p.parseVarStatement()
	case lexer.CONST:
		return p.parseVarStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.FUN:
		return p.parseFunctionStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.ELSE:
		return p.parseElseStatement()
	case lexer.ELSE_IF:
		return p.parseElseIfStatement()
	default:
		return nil
	}
}

// -----------------------------------------------------------------------------
// Parsing Functions
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionStatement() *ast.Function {
	stmt := &ast.Function{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}

	if !p.expectedPeekToken(lexer.IDENTIFIER) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	stmt.Parameters = p.parseFunctionParameters()

	if p.peekTokenIsOk(lexer.COLON) {
		stmt.ReturnType = p.parseFunctionReturnTypes()
	} else {
		stmt.ReturnType = nil
	}

	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil
	}
	stmt.Body = p.parseFunctionBody()

	return stmt
}

// -----------------------------------------------------------------------------
// Parsing Function Body
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionBody() *ast.FunctionBody {
	block := &ast.FunctionBody{Token: p.currentToken}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for p.currentToken.Kind != lexer.CLOSE_CURLY_BRACKET && p.currentToken.Kind != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// -----------------------------------------------------------------------------
// Parsing Function Return Types
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionReturnTypes() []*ast.FunctionReturnType {
	var listToReturn []*ast.FunctionReturnType

	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil
	}

	for !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
		}
		if !p.expectedPeekToken(lexer.TYPE) {
			return nil
		}
		param := &ast.FunctionReturnType{ReturnType: &ast.Type{Token: p.currentToken, Value: p.currentToken.Value}}
		listToReturn = append(listToReturn, param)
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil
	}

	return listToReturn
}

// -----------------------------------------------------------------------------
// Parsing Function Parameters
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionParameters() []*ast.FunctionParameters {
	var listToReturn []*ast.FunctionParameters

	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil
	}

	for !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
		}
		if !p.expectedPeekToken(lexer.IDENTIFIER) {
			return nil
		}
		param := &ast.FunctionParameters{ParameterName: &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}}

		if !p.expectedPeekToken(lexer.COLON) {
			return nil
		}

		if !p.expectedPeekToken(lexer.TYPE) {
			return nil
		}
		param.ParameterType = &ast.Type{Token: p.currentToken, Value: p.currentToken.Value}

		listToReturn = append(listToReturn, param)
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil
	}

	return listToReturn
}

// -----------------------------------------------------------------------------
// Parsing Return Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	// Skiping for now TODO: Implement this
	for !p.currTokenIsOk(lexer.SEMI_COLON) {
		p.nextToken()
	}
	return stmt
}

// -----------------------------------------------------------------------------
// Parsing Var and Const Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.IDENTIFIER) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}

	if !p.expectedPeekToken(lexer.TYPE) {
		return nil
	}
	stmt.Type = &ast.Type{Token: p.currentToken, Value: p.currentToken.Value}

	if !p.expectedPeekToken(lexer.EQUAL_ASSIGN) {
		return nil
	}

	// Skiping for now TODO: Implement this
	for !p.currTokenIsOk(lexer.SEMI_COLON) {
		p.nextToken()
	}
	return stmt
}

// -----------------------------------------------------------------------------
// Parsing If statements
// -----------------------------------------------------------------------------
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil
	}

	// Skiping for now. TODO: Implement this
	for !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		p.nextToken()
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil
	}
	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil
	}

	stmt.Body = p.parseFunctionBody()
	return stmt
}

// -----------------------------------------------------------------------------
// Parsing Else Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseElseStatement() *ast.ElseStatement {
	stmt := &ast.ElseStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil
	}

	stmt.Body = p.parseFunctionBody()

	return stmt
}

// -----------------------------------------------------------------------------
// Parsing Else If Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseElseIfStatement() *ast.ElseIfStatement {
	stmt := &ast.ElseIfStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil
	}

	// Skiping for now. TODO: Implement this
	for !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		p.nextToken()
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil
	}
	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil
	}

	stmt.Body = p.parseFunctionBody()

	return stmt
}

// -----------------------------------------------------------------------------
// Helper Methods
// -----------------------------------------------------------------------------

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	if p.tokenPointer >= len(p.tokens) {
		p.peekToken = lexer.Token{Kind: lexer.EOF}
	} else {
		p.peekToken = p.tokens[p.tokenPointer]
	}
	p.tokenPointer++
}

func (p *Parser) expectedPeekToken(kind lexer.TokenKind) bool {
	if p.peekTokenIsOk(kind) {
		p.nextToken()
		return true
	} else {
		p.peekError(kind)
		return false
	}
}

func (p *Parser) currTokenIsOk(kind lexer.TokenKind) bool {
	return p.currentToken.Kind == kind
}

func (p *Parser) peekTokenIsOk(kind lexer.TokenKind) bool {
	return p.peekToken.Kind == kind
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(kind lexer.TokenKind) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", lexer.TokenKindString(kind), lexer.TokenKindString(p.peekToken.Kind))
	p.errors = append(p.errors, msg)
}
