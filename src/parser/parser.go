package parser

import (
	"fmt"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

type Parser struct {
	tokens        []lexer.Token
	tokenPointer  int
	previousToken lexer.Token
	currentToken  lexer.Token
	peekToken     lexer.Token
	errors        []string

	prefixParseFns  map[lexer.TokenKind]prefixParseFn
	infixParseFns   map[lexer.TokenKind]infixParseFn
	postfixParseFns map[lexer.TokenKind]postfixParseFn
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression, lexer.TokenKind) ast.Expression
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	DOUBLEEQUALS
	LOGICALORAND
	BITWISEORAND
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	POSTFIX
	CALL
)

var precedences = map[lexer.TokenKind]int{
	lexer.DOUBLE_EQUAL:       DOUBLEEQUALS,
	lexer.NOT_EQUAL:          DOUBLEEQUALS,
	lexer.LESS_THAN_EQUAL:    LESSGREATER,
	lexer.GREATER_THAN_EQUAL: LESSGREATER,
	lexer.LESS_THAN:          LESSGREATER,
	lexer.GREATER_THAN:       LESSGREATER,
	lexer.PLUS:               SUM,
	lexer.DASH:               SUM,
	lexer.SLASH:              PRODUCT,
	lexer.STAR:               PRODUCT,
	lexer.PERCENT:            PRODUCT,
	lexer.AND_AND:            LOGICALORAND,
	lexer.OR_OR:              LOGICALORAND,
	lexer.AND:                BITWISEORAND,
	lexer.OR:                 BITWISEORAND,
	lexer.PLUS_PLUS:          POSTFIX,
	lexer.MINUS_MINUS:        POSTFIX,
	lexer.OPEN_BRACKET:       CALL,
	lexer.EQUAL_ASSIGN:       ASSIGN,
	lexer.PLUS_EQUAL:         ASSIGN,
	lexer.MINUS_EQUAL:        ASSIGN,
	lexer.STAR_EQUAL:         ASSIGN,
	lexer.SLASH_EQUAL:        ASSIGN,
	lexer.PERCENT_EQUAL:      ASSIGN,
}

func New(tokens []lexer.Token) *Parser {
	p := &Parser{errors: []string{}, peekToken: tokens[0], tokens: tokens, tokenPointer: 1}
	p.nextToken()

	p.prefixParseFns = make(map[lexer.TokenKind]prefixParseFn)
	p.addPrefix(lexer.IDENTIFIER, p.parseIdentifier)
	p.addPrefix(lexer.INT, p.parseIntegerValue)
	p.addPrefix(lexer.FLOAT, p.parseFloatValue)
	p.addPrefix(lexer.BOOL, p.parseBooleanValue)
	p.addPrefix(lexer.STRING, p.parseStringValue)
	p.addPrefix(lexer.CHAR, p.parseCharValue)
	p.addPrefix(lexer.NOT, p.parsePrefixExpression)
	p.addPrefix(lexer.DASH, p.parsePrefixExpression)
	p.addPrefix(lexer.OPEN_BRACKET, p.parseGroupedExpression)

	p.infixParseFns = make(map[lexer.TokenKind]infixParseFn)
	p.addInfix(lexer.PLUS, p.parseInfixExpression)
	p.addInfix(lexer.DASH, p.parseInfixExpression)
	p.addInfix(lexer.SLASH, p.parseInfixExpression)
	p.addInfix(lexer.STAR, p.parseInfixExpression)
	p.addInfix(lexer.PERCENT, p.parseInfixExpression)
	p.addInfix(lexer.DOUBLE_EQUAL, p.parseInfixExpression)
	p.addInfix(lexer.NOT_EQUAL, p.parseInfixExpression)
	p.addInfix(lexer.LESS_THAN_EQUAL, p.parseInfixExpression)
	p.addInfix(lexer.GREATER_THAN_EQUAL, p.parseInfixExpression)
	p.addInfix(lexer.LESS_THAN, p.parseInfixExpression)
	p.addInfix(lexer.GREATER_THAN, p.parseInfixExpression)
	p.addInfix(lexer.AND_AND, p.parseInfixExpression)
	p.addInfix(lexer.OR_OR, p.parseInfixExpression)
	p.addInfix(lexer.AND, p.parseInfixExpression)
	p.addInfix(lexer.OR, p.parseInfixExpression)
	p.addInfix(lexer.OPEN_BRACKET, p.parseCallExpression)
	p.addInfix(lexer.EQUAL_ASSIGN, p.parseAssignmentExpression)
	p.addInfix(lexer.PLUS_EQUAL, p.parseAssignmentExpression)
	p.addInfix(lexer.MINUS_EQUAL, p.parseAssignmentExpression)
	p.addInfix(lexer.STAR_EQUAL, p.parseAssignmentExpression)
	p.addInfix(lexer.SLASH_EQUAL, p.parseAssignmentExpression)
	p.addInfix(lexer.PERCENT_EQUAL, p.parseAssignmentExpression)

	p.postfixParseFns = make(map[lexer.TokenKind]postfixParseFn)
	p.addPostfix(lexer.PLUS_PLUS, p.parsePostfixExpression)
	p.addPostfix(lexer.MINUS_MINUS, p.parsePostfixExpression)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currentToken.Kind != lexer.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
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
	case lexer.FOR:
		return p.parseForLoop()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Kind]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Kind)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIsOk(lexer.SEMI_COLON) && precedence < p.peekPrecedence() {

		if p.postfixParseFns[p.peekToken.Kind] != nil {
			postfix := p.postfixParseFns[p.peekToken.Kind]
			if postfix == nil {
				p.noPostfixParseFnError(p.peekToken.Kind)
				return nil
			}
			previousOfLeft := p.previousToken.Kind
			p.nextToken()
			leftExp = postfix(leftExp, previousOfLeft)
			continue
		}

		infix := p.infixParseFns[p.peekToken.Kind]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken, Name: left}
	exp.Args = p.parseCallArgs()
	return exp
}

func (p *Parser) parseCallArgs() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIsOk(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil
	}

	return args
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	identifier, ok := left.(*ast.Identifier)
	if !ok {
		p.errors = append(p.errors, "Left side of assignment must be an identifier")
		return nil
	}
	assignmentExpr := &ast.AssignmentExpression{
		Token:    p.currentToken,
		Left:     identifier,
		Operator: p.currentToken.Value,
	}

	p.nextToken()
	assignmentExpr.Right = p.parseExpression(LOWEST)
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil
	}
	return assignmentExpr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Value,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Value,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression, previousOfLeft lexer.TokenKind) ast.Expression {
	expression := &ast.PostfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Value,
		Left:     left,
	}

	if p.peekTokenIsOk(lexer.SEMI_COLON) && previousOfLeft == lexer.SEMI_COLON || previousOfLeft == lexer.OPEN_CURLY_BRACKET {
		expression.IsStmt = true
	} else {
		expression.IsStmt = false
	}

	// p.nextToken()
	return expression
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIsOk(lexer.SEMI_COLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
}

func (p *Parser) parseIntegerValue() ast.Expression {
	intVal := &ast.IntegerValue{Token: p.currentToken}
	value, err := strconv.ParseInt(p.currentToken.Value, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}
	intVal.Value = value
	return intVal
}

func (p *Parser) parseFloatValue() ast.Expression {
	floatVal := &ast.FloatValue{Token: p.currentToken}
	value, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.currentToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}
	floatVal.Value = value
	return floatVal
}

func (p *Parser) parseBooleanValue() ast.Expression {
	val := p.currentToken.Value
	if val == "true" {
		return &ast.BooleanValue{Token: p.currentToken, Value: true}
	}
	return &ast.BooleanValue{Token: p.currentToken, Value: false}
}

func (p *Parser) parseStringValue() ast.Expression {
	return &ast.StringValue{Token: p.currentToken, Value: p.currentToken.Value}
}

func (p *Parser) parseCharValue() ast.Expression {
	return &ast.CharValue{Token: p.currentToken, Value: p.currentToken.Value}
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
		block.Statements = append(block.Statements, stmt)
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

	if p.peekTokenIsOk(lexer.SEMI_COLON) {
		stmt.Value = nil
		p.nextToken()
		return stmt
	}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}

	// We have multiple return values
	if p.peekTokenIsOk(lexer.OPEN_BRACKET) {
		if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
			return nil
		}

		for !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
			if p.peekTokenIsOk(lexer.COMMA) {
				p.nextToken()
			}
			p.nextToken()
			stmt.Value = append(stmt.Value, p.parseExpression(LOWEST))
		}

		if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
			return nil
		}
	} else {
		// only 1 return value
		p.nextToken()
		stmt.Value = append(stmt.Value, p.parseExpression(LOWEST))
	}

	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil
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

	if p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		p.nextToken()

		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	} else {
		// If the value is not assigned, then the default value will be assigned
		switch stmt.Type.Value {
		case "int":
			stmt.Value = &ast.IntegerValue{Token: lexer.Token{Kind: lexer.INT, Value: "0"}, Value: 0}
		case "float":
			stmt.Value = &ast.FloatValue{Token: lexer.Token{Kind: lexer.FLOAT, Value: "0.0"}, Value: 0.0}
		case "bool":
			stmt.Value = &ast.BooleanValue{Token: lexer.Token{Kind: lexer.BOOL, Value: "false"}, Value: false}
		case "string":
			stmt.Value = &ast.StringValue{Token: lexer.Token{Kind: lexer.STRING, Value: ""}, Value: ""}
		case "char":
			stmt.Value = &ast.CharValue{Token: lexer.Token{Kind: lexer.CHAR, Value: ""}, Value: ""}
		}
	}

	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil
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

	stmt.Value = p.parseGroupedExpression()

	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil
	}

	stmt.Body = p.parseFunctionBody()

	// -----------------------------------------------------------------------------
	// Parsing Else If Statements
	// -----------------------------------------------------------------------------
	var elseIfList []*ast.ElseIfStatement
	if p.peekTokenIsOk(lexer.ELSE_IF) {
		for p.peekTokenIsOk(lexer.ELSE_IF) {
			p.nextToken()
			elseIfStmt := &ast.ElseIfStatement{Token: p.currentToken}

			if !p.expectedPeekToken(lexer.COLON) {
				return nil
			}
			if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
				return nil
			}

			elseIfStmt.Value = p.parseGroupedExpression()

			if !p.expectedPeekToken(lexer.COLON) {
				return nil
			}
			if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
				return nil
			}

			elseIfStmt.Body = p.parseFunctionBody()

			elseIfList = append(elseIfList, elseIfStmt)
		}
		stmt.MultiConseq = elseIfList
	} else {
		stmt.MultiConseq = nil
	}

	// -----------------------------------------------------------------------------
	// Parsing Else Statements
	// -----------------------------------------------------------------------------
	if p.peekTokenIsOk(lexer.ELSE) {
		p.nextToken()
		elseStmt := &ast.ElseStatement{Token: p.currentToken}

		if !p.expectedPeekToken(lexer.COLON) {
			return nil
		}
		if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
			return nil
		}
		elseStmt.Body = p.parseFunctionBody()
		stmt.Consequence = elseStmt
	} else {
		stmt.Consequence = nil
	}

	return stmt
}

// -----------------------------------------------------------------------------
// Parsing For Loop
// -----------------------------------------------------------------------------
func (p *Parser) parseForLoop() *ast.ForLoopStatement {
	stmt := &ast.ForLoopStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil
	}

	if !p.expectedPeekToken(lexer.VAR) {
		return nil
	}
	stmt.Left = p.parseVarStatement()
	p.nextToken()

	stmt.Middle = p.parseExpression(LOWEST).(*ast.InfixExpression)

	p.nextToken()
	p.nextToken()

	stmt.Right = p.parseExpression(LOWEST).(*ast.PostfixExpression)

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
	p.previousToken = p.currentToken
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

func (p *Parser) addPrefix(tokenKind lexer.TokenKind, fn prefixParseFn) {
	p.prefixParseFns[tokenKind] = fn
}

func (p *Parser) addInfix(tokenKind lexer.TokenKind, fn infixParseFn) {
	p.infixParseFns[tokenKind] = fn
}

func (p *Parser) addPostfix(tokenKind lexer.TokenKind, fn postfixParseFn) {
	p.postfixParseFns[tokenKind] = fn
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenKind) {
	msg := fmt.Sprintf("No prefix parse function for %s found", lexer.TokenKindString(t))
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPostfixParseFnError(t lexer.TokenKind) {
	msg := fmt.Sprintf("No postfix parse function for %s found", lexer.TokenKindString(t))
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Kind]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Kind]; ok {
		return p
	}
	return LOWEST
}
