package parser

import (
	"fmt"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

var (
	FunctionMap = make(map[*ast.Identifier]*ast.Function)
	InForLoop   = false
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
	_            int = iota
	LOWEST           // Lowest precedence, used as default
	ASSIGN           // Assignment operators like =, +=, -=
	LOGICALOR        // Logical OR ||
	LOGICALAND       // Logical AND &&
	DOUBLEEQUALS     // Equality operators ==, !=
	BITWISEORAND     // Bitwise AND/OR operators &, |
	LESSGREATER      // Relational operators <, <=, >, >=
	SUM              // Addition and subtraction +, -
	PRODUCT          // Multiplication, division, modulus *, /, %
	PREFIX           // Unary operators like - (negative), !
	POSTFIX          // Postfix operators like ++, --
	CALL             // Function calls, array indexing ()
	INDEX            // Array indexing []
)

var precedences = map[lexer.TokenKind]int{
	// Assignment operators
	lexer.EQUAL_ASSIGN:  ASSIGN,
	lexer.PLUS_EQUAL:    ASSIGN,
	lexer.MINUS_EQUAL:   ASSIGN,
	lexer.STAR_EQUAL:    ASSIGN,
	lexer.SLASH_EQUAL:   ASSIGN,
	lexer.PERCENT_EQUAL: ASSIGN,

	// Logical operators
	lexer.OR_OR:   LOGICALOR,  // Logical OR ||
	lexer.AND_AND: LOGICALAND, // Logical AND &&

	// Equality operators
	lexer.DOUBLE_EQUAL: DOUBLEEQUALS, // Equality ==, !=
	lexer.NOT_EQUAL:    DOUBLEEQUALS,

	// Bitwise operators
	lexer.AND: BITWISEORAND, // Bitwise AND &
	lexer.OR:  BITWISEORAND, // Bitwise OR |

	// Relational operators
	lexer.LESS_THAN_EQUAL:    LESSGREATER, // Relational <=, >=, <, >
	lexer.GREATER_THAN_EQUAL: LESSGREATER,
	lexer.LESS_THAN:          LESSGREATER,
	lexer.GREATER_THAN:       LESSGREATER,

	// Arithmetic operators
	lexer.PLUS:    SUM, // Addition and subtraction
	lexer.DASH:    SUM,
	lexer.STAR:    PRODUCT, // Multiplication, division, modulus
	lexer.SLASH:   PRODUCT,
	lexer.PERCENT: PRODUCT,

	// Postfix operators
	lexer.PLUS_PLUS:   POSTFIX, // Postfix ++, --
	lexer.MINUS_MINUS: POSTFIX,

	// Function calls and array indexing
	lexer.OPEN_BRACKET: CALL,

	// Array indexing
	lexer.OPEN_SQUARE_BRACKET: INDEX,
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
	p.addPrefix(lexer.OPEN_SQUARE_BRACKET, p.parseArrayValues)
	p.addPrefix(lexer.OPEN_CURLY_BRACKET, p.parseHashMapValues)

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
	p.addInfix(lexer.OPEN_SQUARE_BRACKET, p.parseIndexExpression)

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
	case lexer.CONTINUE:
		return p.parseContinueStatement()
	case lexer.BREAK:
		return p.parseBreakStatement()
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
	return assignmentExpr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.currentToken, Left: left}
	p.nextToken()
	idx := p.parseExpression(LOWEST)
	exp.Index = idx
	if !p.expectedPeekToken(lexer.CLOSE_SQUARE_BRACKET) {
		return nil
	}
	return exp
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

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	if p.peekTokenIsOk(lexer.COMMA) {
		stmt.Expression = &ast.AssignmentExpression{
			Token:    lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
			Left:     &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value},
			Operator: "=",
		}
		list := []ast.Statement{}
		list = append(list, stmt)
		p.nextToken()
		p.nextToken()
		return p.parseMultipleAssignmentStatement(list)
	}

	stmt.Expression = p.parseExpression(LOWEST)
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil
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

	FunctionMap[stmt.Name] = stmt

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
		if !p.expectedPeekToken(lexer.TYPE) {
			return nil
		}

		firstType := p.currentToken
		param := &ast.FunctionReturnType{ReturnType: &ast.Type{Token: firstType, Value: firstType.Value}}
		if p.peekTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
				param.ReturnType.IsArray = true
				param.ReturnType.IsHash = false
				param.ReturnType.SubTypes = nil
				p.nextToken()
			} else if p.peekTokenIsOk(lexer.TYPE) {
				p.nextToken()
				param.ReturnType.IsArray = false
				param.ReturnType.IsHash = true

				subTypes := []*ast.Type{}
				subTypes = append(subTypes, &ast.Type{Token: firstType, Value: firstType.Value, IsArray: false, IsHash: false, SubTypes: nil})
				subTypes = append(subTypes, &ast.Type{Token: p.currentToken, Value: p.currentToken.Value, IsArray: false, IsHash: false, SubTypes: nil})
				param.ReturnType.SubTypes = subTypes

				if !p.expectedPeekToken(lexer.CLOSE_SQUARE_BRACKET) {
					return nil
				}
			}
		} else {
			param.ReturnType.IsArray = false
			param.ReturnType.IsHash = false
			param.ReturnType.SubTypes = nil
		}
		listToReturn = append(listToReturn, param)

		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				msg := fmt.Sprintf("Expected a value after comma, got %s instead", p.peekToken.Value)
				p.errors = append(p.errors, msg)
				return nil
			}
		} else {
			if !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				return nil
			}
		}
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
		firstType := p.currentToken
		param.ParameterType = &ast.Type{Token: firstType, Value: firstType.Value}

		if p.peekTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
				param.ParameterType.IsArray = true
				param.ParameterType.IsHash = false
				param.ParameterType.SubTypes = nil
				p.nextToken()
			} else if p.peekTokenIsOk(lexer.TYPE) {
				p.nextToken()
				param.ParameterType.IsArray = false
				param.ParameterType.IsHash = true

				subTypes := []*ast.Type{}
				subTypes = append(subTypes, &ast.Type{Token: firstType, Value: firstType.Value, IsArray: false, IsHash: false, SubTypes: nil})
				subTypes = append(subTypes, &ast.Type{Token: p.currentToken, Value: p.currentToken.Value, IsArray: false, IsHash: false, SubTypes: nil})
				param.ParameterType.SubTypes = subTypes

				if !p.expectedPeekToken(lexer.CLOSE_SQUARE_BRACKET) {
					return nil
				}
			}
		} else {
			param.ParameterType.IsArray = false
			param.ParameterType.IsHash = false
			param.ParameterType.SubTypes = nil
		}

		listToReturn = append(listToReturn, param)

		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				msg := fmt.Sprintf("Expected a value after comma, got %s instead", p.peekToken.Value)
				p.errors = append(p.errors, msg)
				return nil
			}
		} else {
			if !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				return nil
			}
		}
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
// Parsing Continue Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	if !InForLoop {
		msg := "Continue statement can only be used inside a for loop"
		p.errors = append(p.errors, msg)
		return nil
	}
	stmt := &ast.ContinueStatement{Token: p.currentToken}
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil
	}
	return stmt
}

// -----------------------------------------------------------------------------
// Parsing Break Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	if !InForLoop {
		msg := "Break statement can only be used inside a for loop"
		p.errors = append(p.errors, msg)
		return nil
	}
	stmt := &ast.BreakStatement{Token: p.currentToken}
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil
	}
	return stmt
}

// -----------------------------------------------------------------------------
// Parsing Multiple assignment statement
// -----------------------------------------------------------------------------
func (p *Parser) parseMultipleAssignmentStatement(list []ast.Statement) *ast.MultiValueAssignStmt {
	stmt := &ast.MultiValueAssignStmt{Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="}}

	if p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
		msg := fmt.Sprintf("Expected an identifier or a var/const keyword, got %s instead", p.currentToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	for !p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
		if p.currTokenIsOk(lexer.VAR) || p.currTokenIsOk(lexer.CONST) {
			varSig := p.parseVarStmtSig()
			list = append(list, varSig)
			p.nextToken()
		} else if p.currTokenIsOk(lexer.IDENTIFIER) {
			expStmt := &ast.ExpressionStatement{
				Token: p.currentToken,
				Expression: &ast.AssignmentExpression{
					Token:    lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
					Left:     &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value},
					Operator: "=",
				},
			}
			list = append(list, expStmt)
			p.nextToken()
		} else {
			errMsg := fmt.Sprintf("Expected an identifier or a var/const keyword, got %s instead", p.currentToken.Value)
			p.errors = append(p.errors, errMsg)
			return nil
		}

		if p.currTokenIsOk(lexer.COMMA) {
			p.nextToken()
		}
	}

	p.nextToken()

	valueList := []ast.Expression{}
	for !p.currTokenIsOk(lexer.SEMI_COLON) {
		valueList = append(valueList, p.parseExpression(LOWEST))
		p.nextToken()
		if p.currTokenIsOk(lexer.COMMA) {
			p.nextToken()
		}
	}

	// ValueList could have a CallExpression, or it would have equal number of expressions as the number of objects in the list
	// if we have a call expression, traverse the object list and assign the call expression to each object
	// if we have equal number of expressions, assign each expression to the corresponding object

	if len(valueList) == 1 {
		if callExp, ok := valueList[0].(*ast.CallExpression); ok {
			for i, obj := range list {
				varObj, ok := obj.(*ast.VarStatement)
				if ok {
					varObj.Value = callExp
					list[i] = varObj
				} else {
					expStmtObj, ok := obj.(*ast.ExpressionStatement)
					if ok {
						expStmtObjExp, ok1 := expStmtObj.Expression.(*ast.AssignmentExpression)
						if ok1 {
							expStmtObjExp.Right = callExp
							expStmtObj.Expression = expStmtObjExp
							list[i] = expStmtObj
						}
					}
				}
			}
		}
	} else {
		// there are multiple expressions
		if len(valueList) != len(list) {
			errMsg := "Number of expressions on the right do not match the number of decelaraions on the left"
			p.errors = append(p.errors, errMsg)
			return nil
		}

		for i, obj := range list {
			varObj, ok := obj.(*ast.VarStatement)
			if ok {
				varObj.Value = valueList[i]
				if _, ok := valueList[i].(*ast.ArrayValue); ok {
					varObj.Value.(*ast.ArrayValue).Type = varObj.Type
				}
				if _, ok := valueList[i].(*ast.HashMap); ok {
					varObj.Value.(*ast.HashMap).KeyType = varObj.Type.SubTypes[0]
					varObj.Value.(*ast.HashMap).ValueType = varObj.Type.SubTypes[1]
				}
				list[i] = varObj
			} else {
				expStmtObj, ok := obj.(*ast.ExpressionStatement)
				if ok {
					expStmtObjExp, ok1 := expStmtObj.Expression.(*ast.AssignmentExpression)
					if ok1 {
						expStmtObjExp.Right = valueList[i]
						expStmtObj.Expression = expStmtObjExp
						list[i] = expStmtObj
					}
				}
			}
		}
	}

	stmt.Objects = list
	return stmt
}

// -----------------------------------------------------------------------------
// Parsing Arrays {...}
// -----------------------------------------------------------------------------
func (p *Parser) parseArrayValues() ast.Expression {
	if !p.currTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		msg := fmt.Sprintf("Expected an open square bracket, got %s instead", p.currentToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}
	array := &ast.ArrayValue{Token: p.currentToken}
	p.nextToken()

	list := []ast.Expression{}
	for !p.currTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
		list = append(list, p.parseExpression(LOWEST))
		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
				msg := fmt.Sprintf("Expected a value after comma, got %s instead", p.peekToken.Value)
				p.errors = append(p.errors, msg)
				return nil
			}
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			p.nextToken()
			break
		}
	}
	array.Values = list
	return array
}

// -----------------------------------------------------------------------------
// Parsing HashMap {...}
// -----------------------------------------------------------------------------
func (p *Parser) parseHashMapValues() ast.Expression {
	if !p.currTokenIsOk(lexer.OPEN_CURLY_BRACKET) {
		msg := fmt.Sprintf("Expected an open curly bracket, got %s instead", p.currentToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}
	hashMap := &ast.HashMap{Token: p.currentToken}
	pairs := make(map[ast.Expression]ast.Expression)
	p.nextToken()

	for !p.currTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
		key := p.parseExpression(LOWEST)
		if !p.expectedPeekToken(lexer.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)
		pairs[key] = value

		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
				msg := fmt.Sprintf("Expected a value after comma, got %s instead", p.peekToken.Value)
				p.errors = append(p.errors, msg)
				return nil
			}
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
			p.nextToken()
			break
		}
	}
	hashMap.Pairs = pairs
	return hashMap
}

// -----------------------------------------------------------------------------
// Parsing Var and Const Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseVarStmtSig() *ast.VarStatement {
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

	firstType := p.currentToken
	stmt.Type = &ast.Type{Token: firstType, Value: firstType.Value}

	// check for array
	if p.peekTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		p.nextToken()
		if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			// is an array
			stmt.Type.IsArray = true
			stmt.Type.IsHash = false
			stmt.Type.SubTypes = nil
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.TYPE) {
			// is a hashmap
			p.nextToken()
			stmt.Type.IsArray = false
			stmt.Type.IsHash = true

			subTypes := []*ast.Type{}
			subTypes = append(subTypes, &ast.Type{Token: firstType, Value: firstType.Value, IsArray: false, IsHash: false, SubTypes: nil})
			subTypes = append(subTypes, &ast.Type{Token: p.currentToken, Value: p.currentToken.Value, IsArray: false, IsHash: false, SubTypes: nil})
			stmt.Type.SubTypes = subTypes

			if !p.expectedPeekToken(lexer.CLOSE_SQUARE_BRACKET) {
				return nil
			}
		}
	} else {
		// no array or hashmap, just simple type
		stmt.Type.IsArray = false
		stmt.Type.IsHash = false
		stmt.Type.SubTypes = nil
	}

	return stmt
}

func (p *Parser) parseVarStatement() ast.Statement {
	stmt := p.parseVarStmtSig()

	// we have multiple assignments
	if p.peekTokenIsOk(lexer.COMMA) {
		list := []ast.Statement{}
		list = append(list, stmt)
		p.nextToken()
		p.nextToken()
		return p.parseMultipleAssignmentStatement(list)
	}

	if stmt.Type.IsArray && !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		msg := fmt.Sprintf("Array type %s must be initialized with values, even `[]` empty", stmt.Type.Value)
		p.errors = append(p.errors, msg)
		return nil
	} else if stmt.Type.IsHash && !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		msg := fmt.Sprintf("Hashmap type %s must be initialized with values, even `{}` empty", stmt.Type.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	if p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)

		// if the type has isArray = true, exec this:
		if stmt.Type.IsArray {
			if array, ok := stmt.Value.(*ast.ArrayValue); ok {
				array.Type = stmt.Type
				stmt.Value = array
			}
		} else if stmt.Type.IsHash {
			if hash, ok := stmt.Value.(*ast.HashMap); ok {
				hash.KeyType = stmt.Type.SubTypes[0]
				hash.ValueType = stmt.Type.SubTypes[1]
				stmt.Value = hash
			}
		}
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
	InForLoop = true
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
	stmt.Left = p.parseVarStatement().(*ast.VarStatement)
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

	InForLoop = false
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
