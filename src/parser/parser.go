package parser

import (
	"errors"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

var (
	FunctionMap = make(map[*ast.Identifier]*ast.Function)
	inForLoop   = false
)

type Parser struct {
	tokens        []lexer.Token
	tokenPointer  int
	previousToken lexer.Token
	currentToken  lexer.Token
	peekToken     lexer.Token

	prefixParseFns  map[lexer.TokenKind]prefixParseFn
	infixParseFns   map[lexer.TokenKind]infixParseFn
	postfixParseFns map[lexer.TokenKind]postfixParseFn
}

type (
	prefixParseFn  func() (ast.Expression, error)
	infixParseFn   func(ast.Expression) (ast.Expression, error)
	postfixParseFn func(ast.Expression, lexer.TokenKind) (ast.Expression, error)
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
	p := &Parser{peekToken: tokens[0], tokens: tokens, tokenPointer: 1}
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

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currentToken.Kind != lexer.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
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

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefixParseFns[p.currentToken.Kind]
	if prefix == nil {
		return nil, errors.New("no prefix parse function for " + lexer.TokenKindString(p.currentToken.Kind))
	}
	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for !p.peekTokenIsOk(lexer.SEMI_COLON) && precedence < p.peekPrecedence() {
		if p.postfixParseFns[p.peekToken.Kind] != nil {
			postfix := p.postfixParseFns[p.peekToken.Kind]
			if postfix == nil {
				return nil, errors.New("no postfix parse function for " + lexer.TokenKindString(p.peekToken.Kind))
			}
			previousOfLeft := p.previousToken.Kind
			p.nextToken()
			leftExp, err = postfix(leftExp, previousOfLeft)
			if err != nil {
				return nil, err
			}
			continue
		}

		infix := p.infixParseFns[p.peekToken.Kind]
		if infix == nil {
			return leftExp, nil
		}
		p.nextToken()
		leftExp, err = infix(leftExp)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

func (p *Parser) parseCallExpression(left ast.Expression) (ast.Expression, error) {
	exp := &ast.CallExpression{Token: p.currentToken, Name: left}
	expArgs, err := p.parseCallArgs()
	if err != nil {
		return nil, err
	}
	exp.Args = expArgs
	return exp, nil
}

func (p *Parser) parseCallArgs() ([]ast.Expression, error) {
	var args []ast.Expression

	if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		p.nextToken()
		return args, nil
	}

	p.nextToken()
	parsedExp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	args = append(args, parsedExp)

	for p.peekTokenIsOk(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		parsedExp, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, parsedExp)
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("expected a closing bracket")
	}

	return args, nil
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) (ast.Expression, error) {
	identifier, ok := left.(*ast.Identifier)
	if !ok {
		return nil, errors.New("Left side of assignment must be an identifier")
	}
	assignmentExpr := &ast.AssignmentExpression{
		Token:    p.currentToken,
		Left:     identifier,
		Operator: p.currentToken.Value,
	}

	p.nextToken()
	parsedExp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	assignmentExpr.Right = parsedExp
	return assignmentExpr, nil
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	exp := &ast.IndexExpression{Token: p.currentToken, Left: left}
	p.nextToken()
	idx, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	exp.Index = idx
	if !p.expectedPeekToken(lexer.CLOSE_SQUARE_BRACKET) {
		return nil, errors.New("Expected a closing square bracket")
	}
	return exp, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.nextToken()
	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("Expected a closing bracket")
	}
	return exp, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Value,
	}
	p.nextToken()
	parsedExp, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}
	expression.Right = parsedExp
	return expression, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Value,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	parsedExp, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	expression.Right = parsedExp
	return expression, nil
}

func (p *Parser) parsePostfixExpression(left ast.Expression, previousOfLeft lexer.TokenKind) (ast.Expression, error) {
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
	return expression, nil
}

func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
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

	parsedExp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Expression = parsedExp
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil, errors.New("Expected a semicolon at the end of the statement")
	}

	return stmt, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

func (p *Parser) parseIntegerValue() (ast.Expression, error) {
	intVal := &ast.IntegerValue{Token: p.currentToken}
	value, err := strconv.ParseInt(p.currentToken.Value, 0, 64)
	if err != nil {
		return nil, errors.New("could not parse " + p.currentToken.Value + " as integer")
	}
	intVal.Value = value
	return intVal, nil
}

func (p *Parser) parseFloatValue() (ast.Expression, error) {
	floatVal := &ast.FloatValue{Token: p.currentToken}
	value, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		return nil, errors.New("could not parse " + p.currentToken.Value + " as float")
	}
	floatVal.Value = value
	return floatVal, nil
}

func (p *Parser) parseBooleanValue() (ast.Expression, error) {
	val := p.currentToken.Value
	if val == "true" {
		return &ast.BooleanValue{Token: p.currentToken, Value: true}, nil
	}
	return &ast.BooleanValue{Token: p.currentToken, Value: false}, nil
}

func (p *Parser) parseStringValue() (ast.Expression, error) {
	return &ast.StringValue{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

func (p *Parser) parseCharValue() (ast.Expression, error) {
	return &ast.CharValue{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

// -----------------------------------------------------------------------------
// Parsing Functions
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionStatement() (*ast.Function, error) {
	stmt := &ast.Function{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("Expected a colon after the fun keyword")
	}

	if !p.expectedPeekToken(lexer.IDENTIFIER) {
		return nil, errors.New("Expected a function name after the colon")
	}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	parameters, err := p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}
	stmt.Parameters = parameters

	if p.peekTokenIsOk(lexer.COLON) {
		returnType, err := p.parseFunctionReturnTypes()
		if err != nil {
			return nil, err
		}
		stmt.ReturnType = returnType
	} else {
		stmt.ReturnType = nil
	}

	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil, errors.New("Expected an open curly bracket after the function signature")
	}
	funBody, err := p.parseFunctionBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = funBody

	FunctionMap[stmt.Name] = stmt

	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Function Body
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionBody() (*ast.FunctionBody, error) {
	block := &ast.FunctionBody{Token: p.currentToken}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for p.currentToken.Kind != lexer.CLOSE_CURLY_BRACKET && p.currentToken.Kind != lexer.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}
	return block, nil
}

// -----------------------------------------------------------------------------
// Parsing Function Return Types
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionReturnTypes() ([]*ast.FunctionReturnType, error) {
	var listToReturn []*ast.FunctionReturnType

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("Expected a colon after the function name and parameters")
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil, errors.New("Expected an open bracket after the colon")
	}

	for !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		if !p.expectedPeekToken(lexer.TYPE) {
			return nil, errors.New("Expected a type after the open bracket")
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
					return nil, errors.New("Expected a closing square bracket after value type")
				}
			} else {
				return nil, errors.New("Expected a closing square bracket for array or a type for hashmap. but got=" + p.peekToken.Value)
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
				return nil, errors.New("Expected a value after comma, got " + p.peekToken.Value + " instead")
			}
		} else {
			if !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				return nil, errors.New("Expected a closing bracket or a comma after the parameter type but got=" + p.peekToken.Value)
			}
		}
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("Expected a closing bracket after the return types")
	}

	return listToReturn, nil
}

// -----------------------------------------------------------------------------
// Parsing Function Parameters
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionParameters() ([]*ast.FunctionParameters, error) {
	var listToReturn []*ast.FunctionParameters

	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil, errors.New("Expected an open bracket after the function name")
	}

	for !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		if !p.expectedPeekToken(lexer.IDENTIFIER) {
			return nil, errors.New("Expected an identifier or close bracket after the open bracket")
		}
		param := &ast.FunctionParameters{ParameterName: &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}}

		if !p.expectedPeekToken(lexer.COLON) {
			return nil, errors.New("Expected a colon after the parameter name")
		}

		if !p.expectedPeekToken(lexer.TYPE) {
			return nil, errors.New("Expected a type after the colon")
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
					return nil, errors.New("Expected a closing square bracket after value type")
				}
			} else {
				return nil, errors.New("Expected a closing square bracket for array or a type for hashmap. but got=" + p.peekToken.Value)
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
				return nil, errors.New("Expected a value after comma got " + p.peekToken.Value + " instead")
			}
		} else {
			if !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				return nil, errors.New("Expected a closing bracket or a comma after the parameter type but got=" + p.peekToken.Value)
			}
		}
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("Expected a closing bracket after the parameters")
	}

	return listToReturn, nil
}

// -----------------------------------------------------------------------------
// Parsing Return Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	if p.peekTokenIsOk(lexer.SEMI_COLON) {
		stmt.Value = nil
		p.nextToken()
		return stmt, nil
	}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("Expected a colon after the return keyword")
	}

	// We have multiple return values
	if p.peekTokenIsOk(lexer.OPEN_BRACKET) {
		p.nextToken()
		if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
			return nil, errors.New("Expected values after open bracket")
		}
		p.nextToken()

		for !p.currTokenIsOk(lexer.CLOSE_BRACKET) {
			parsedExp, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			stmt.Value = append(stmt.Value, parsedExp)
			if p.peekTokenIsOk(lexer.COMMA) {
				p.nextToken()
				if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
					return nil, errors.New("Expected a value after comma")
				}
				p.nextToken()
			} else if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				p.nextToken()
				break
			}
		}
	} else {
		// only 1 return value
		p.nextToken()
		parsedExp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		stmt.Value = append(stmt.Value, parsedExp)
		if p.peekTokenIsOk(lexer.COMMA) {
			return nil, errors.New("Expected a semicolon, got a COMMA insted. To return multiple values, use `return: (val1, val2, ...)`")
		}
	}

	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil, errors.New("Expected a semicolon at the end of the statement")
	}

	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Continue Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseContinueStatement() (*ast.ContinueStatement, error) {
	if !inForLoop {
		return nil, errors.New("Continue statement can only be used inside a for loop")
	}
	stmt := &ast.ContinueStatement{Token: p.currentToken}
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil, errors.New("Expected a semicolon at the end of the statement")
	}
	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Break Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseBreakStatement() (*ast.BreakStatement, error) {
	if !inForLoop {
		return nil, errors.New("Break statement can only be used inside a for loop")
	}
	stmt := &ast.BreakStatement{Token: p.currentToken}
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil, errors.New("Expected a semicolon at the end of the statement")
	}
	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Multiple assignment statement
// -----------------------------------------------------------------------------
func (p *Parser) parseMultipleAssignmentStatement(list []ast.Statement) (*ast.MultiValueAssignStmt, error) {
	stmt := &ast.MultiValueAssignStmt{Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="}}

	if p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
		return nil, errors.New("Expected an identifier or a var/const keyword, got " + p.currentToken.Value + " instead")
	}

	for !p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
		if p.currTokenIsOk(lexer.VAR) || p.currTokenIsOk(lexer.CONST) {
			varSig, err := p.parseVarStmtSig()
			if err != nil {
				return nil, err
			}
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
			return nil, errors.New("Expected an identifier or a var/const keyword, got " + p.currentToken.Value + " instead")
		}

		if p.currTokenIsOk(lexer.COMMA) {
			if p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
				return nil, errors.New("Expected an identifier or a var/const keyword, got " + p.peekToken.Value + " instead")
			}
			p.nextToken()
		} else if p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
			break
		}
	}

	p.nextToken()

	valueList := []ast.Expression{}
	for !p.currTokenIsOk(lexer.SEMI_COLON) {
		parsedExp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		valueList = append(valueList, parsedExp)
		p.nextToken()
		if p.currTokenIsOk(lexer.COMMA) {
			p.nextToken()
		}
	}

	// ValueList could have a CallExpression, or it would have equal number of expressions as the number of objects in the list
	// if we have a call expression, traverse the object list and assign the call expression to each object
	// if we have equal number of expressions, assign each expression to the corresponding object

	if len(valueList) == 1 {
		stmt.SingleCallExp = true
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
		} else {
			return nil, errors.New("Number of expressions on the right do not match the number of decelaraions on the left")
		}
	} else {
		// there are multiple expressions
		if len(valueList) != len(list) {
			return nil, errors.New("Number of expressions on the right do not match the number of decelaraions on the left")
		}
		stmt.SingleCallExp = false

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
				} else {
					return nil, errors.New("Expected an variable identifier or a var/const statement, got " + p.currentToken.Value + " instead")
				}
			}
		}
	}

	stmt.Objects = list
	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Arrays {...}
// -----------------------------------------------------------------------------
func (p *Parser) parseArrayValues() (ast.Expression, error) {
	if !p.currTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		return nil, errors.New("Expected an open square bracket, got " + p.currentToken.Value + " instead")
	}
	array := &ast.ArrayValue{Token: p.currentToken}
	p.nextToken()

	list := []ast.Expression{}
	for !p.currTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
		parsedExp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		list = append(list, parsedExp)
		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
				return nil, errors.New("Expected a value after comma, got " + p.peekToken.Value + " instead")
			}
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			p.nextToken()
			break
		}
	}
	array.Values = list
	return array, nil
}

// -----------------------------------------------------------------------------
// Parsing HashMap {...}
// -----------------------------------------------------------------------------
func (p *Parser) parseHashMapValues() (ast.Expression, error) {
	if !p.currTokenIsOk(lexer.OPEN_CURLY_BRACKET) {
		return nil, errors.New("Expected an open curly bracket, got " + p.currentToken.Value + " instead")
	}
	hashMap := &ast.HashMap{Token: p.currentToken}
	pairs := make(map[ast.Expression]ast.Expression)
	p.nextToken()

	for !p.currTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
		parsedExp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		key := parsedExp
		if !p.expectedPeekToken(lexer.COLON) {
			return nil, errors.New("Expected a colon after the key")
		}
		p.nextToken()
		parsedExp, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		value := parsedExp
		pairs[key] = value

		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
				return nil, errors.New("Expected a value after comma, got " + p.peekToken.Value + " instead")
			}
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
			p.nextToken()
			break
		}
	}
	hashMap.Pairs = pairs
	return hashMap, nil
}

// -----------------------------------------------------------------------------
// Parsing Var and Const Statements
// -----------------------------------------------------------------------------
func (p *Parser) parseVarStmtSig() (*ast.VarStatement, error) {
	stmt := &ast.VarStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.IDENTIFIER) {
		if stmt.Token.Kind == lexer.CONST {
			return nil, errors.New("Expected an identifier after const keyword")
		}
		return nil, errors.New("Expected an identifier after var keyword")
	}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("Expected a colon after the identifier")
	}

	if !p.expectedPeekToken(lexer.TYPE) {
		return nil, errors.New("Expected a type after the colon")
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
				return nil, errors.New("Expected a closing square bracket")
			}
		}
	} else {
		// no array or hashmap, just simple type
		stmt.Type.IsArray = false
		stmt.Type.IsHash = false
		stmt.Type.SubTypes = nil
	}

	return stmt, nil
}

func (p *Parser) parseVarStatement() (ast.Statement, error) {
	stmt, err := p.parseVarStmtSig()
	if err != nil {
		return nil, err
	}

	// we have multiple assignments
	if p.peekTokenIsOk(lexer.COMMA) {
		list := []ast.Statement{}
		list = append(list, stmt)
		p.nextToken()
		p.nextToken()
		return p.parseMultipleAssignmentStatement(list)
	}

	if stmt.Type.IsArray && !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		return nil, errors.New("Array type must be initialized with values, for empty array use `[]")
	} else if stmt.Type.IsHash && !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		return nil, errors.New("Hashmap type must be initialized with values, for empty hashmap use `{}`")
	} else if stmt.Token.Kind == lexer.CONST && !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		return nil, errors.New("Const " + stmt.Name.Value + " must be initialized")
	}

	if p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		p.nextToken()
		p.nextToken()
		parsedExp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		stmt.Value = parsedExp

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
		return nil, errors.New("Expected a semicolon at the end of the statement")
	}
	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing If statements
// -----------------------------------------------------------------------------
func (p *Parser) parseIfStatement() (*ast.IfStatement, error) {
	stmt := &ast.IfStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("Expected a colon after `if` keyword")
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil, errors.New("Expected an open bracket after the colon")
	}

	groupedExpValue, err := p.parseGroupedExpression()
	if err != nil {
		return nil, err
	}
	stmt.Value = groupedExpValue

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("Expected a colon after the expression")
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil, errors.New("Expected an open curly bracket after the colon")
	}

	body, err := p.parseFunctionBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	// -----------------------------------------------------------------------------
	// Parsing Else If Statements
	// -----------------------------------------------------------------------------
	var elseIfList []*ast.ElseIfStatement
	if p.peekTokenIsOk(lexer.ELSE_IF) {
		for p.peekTokenIsOk(lexer.ELSE_IF) {
			p.nextToken()
			elseIfStmt := &ast.ElseIfStatement{Token: p.currentToken}

			if !p.expectedPeekToken(lexer.COLON) {
				return nil, errors.New("Expected a colon after the else if keyword")
			}
			if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
				return nil, errors.New("Expected an open bracket after the colon")
			}

			elseIfValue, err := p.parseGroupedExpression()
			if err != nil {
				return nil, err
			}
			elseIfStmt.Value = elseIfValue

			if !p.expectedPeekToken(lexer.COLON) {
				return nil, errors.New("Expected a colon after the grouped expression")
			}
			if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
				return nil, errors.New("Expected an open curly bracket after the colon")
			}

			elseIFBody, err := p.parseFunctionBody()
			if err != nil {
				return nil, err
			}
			elseIfStmt.Body = elseIFBody

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
			return nil, errors.New("Expected a colon after the else keyword")
		}
		if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
			return nil, errors.New("Expected an open curly bracket after the colon")
		}
		elseBody, err := p.parseFunctionBody()
		if err != nil {
			return nil, err
		}
		elseStmt.Body = elseBody
		stmt.Consequence = elseStmt
	} else {
		stmt.Consequence = nil
	}

	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing For Loop
// -----------------------------------------------------------------------------
func (p *Parser) parseForLoop() (*ast.ForLoopStatement, error) {
	inForLoop = true
	stmt := &ast.ForLoopStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("Expected a colon after the for keyword")
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil, errors.New("Expected an open bracket after the colon")
	}

	if !p.expectedPeekToken(lexer.VAR) {
		return nil, errors.New("Expected a var keyword after the open bracket")
	}
	varStmt, err := p.parseVarStatement()
	if err != nil {
		return nil, err
	}
	stmt.Left = varStmt.(*ast.VarStatement)
	p.nextToken()

	parsedExp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Middle = parsedExp.(*ast.InfixExpression)

	p.nextToken()
	p.nextToken()

	parsedExp, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Right = parsedExp.(*ast.PostfixExpression)

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("Expected a closing bracket after the postfix expression")
	}
	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("Expected a colon after the closing bracket")
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil, errors.New("Expected an open curly bracket after the colon")
	}
	stmtBody, err := p.parseFunctionBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = stmtBody

	inForLoop = false
	return stmt, nil
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
		return false
	}
}

func (p *Parser) currTokenIsOk(kind lexer.TokenKind) bool {
	return p.currentToken.Kind == kind
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
	if p, ok := precedences[p.currentToken.Kind]; ok {
		return p
	}
	return LOWEST
}
