package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

var (
	FunctionMap      = make(map[string]*ast.Function)
	localFunctionMap = make(map[string]bool)
	inForLoop        = false
	inFunction       = false
)

type Parser struct {
	tokens        []lexer.Token
	tokenPointer  int
	previousToken lexer.Token
	currentToken  lexer.Token
	peekToken     lexer.Token
	inTesting     bool

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
	_ int = iota
	LOWEST
	ASSIGN
	LOGICALOR
	LOGICALAND
	DOUBLEEQUALS
	BITWISEORAND
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	POSTFIX
	CALL
	INDEX
)

var precedences = map[lexer.TokenKind]int{
	lexer.EQUAL_ASSIGN:  ASSIGN,
	lexer.PLUS_EQUAL:    ASSIGN,
	lexer.MINUS_EQUAL:   ASSIGN,
	lexer.STAR_EQUAL:    ASSIGN,
	lexer.SLASH_EQUAL:   ASSIGN,
	lexer.PERCENT_EQUAL: ASSIGN,

	lexer.OR_OR:   LOGICALOR,
	lexer.AND_AND: LOGICALAND,

	lexer.DOUBLE_EQUAL: DOUBLEEQUALS,
	lexer.NOT_EQUAL:    DOUBLEEQUALS,

	lexer.AND: BITWISEORAND,
	lexer.OR:  BITWISEORAND,

	lexer.LESS_THAN_EQUAL:    LESSGREATER,
	lexer.GREATER_THAN_EQUAL: LESSGREATER,
	lexer.LESS_THAN:          LESSGREATER,
	lexer.GREATER_THAN:       LESSGREATER,

	lexer.PLUS:    SUM,
	lexer.DASH:    SUM,
	lexer.STAR:    PRODUCT,
	lexer.SLASH:   PRODUCT,
	lexer.PERCENT: PRODUCT,

	lexer.PLUS_PLUS:   POSTFIX,
	lexer.MINUS_MINUS: POSTFIX,

	lexer.OPEN_BRACKET: CALL,

	lexer.OPEN_SQUARE_BRACKET: INDEX,
}

func New(tokens []lexer.Token, inTesting bool) *Parser {
	p := &Parser{inTesting: inTesting, peekToken: tokens[0], tokens: tokens, tokenPointer: 1}
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
	// check if all functions called are defined
	for funcName := range localFunctionMap {
		if _, ok := FunctionMap[funcName]; !ok {
			if _, ok := builtinMap[funcName]; !ok {
				return nil, errors.New("function `" + funcName + "` is called but not defined")
			}
		}
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
	case lexer.WHILE:
		return p.parseWhileLoop()
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
		return nil, errors.New("no prefix parse function for: " + lexer.TokenKindString(p.currentToken.Kind))
	}
	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for !p.peekTokenIsOk(lexer.SEMI_COLON) && precedence < p.peekPrecedence() {
		if p.postfixParseFns[p.peekToken.Kind] != nil {
			postfix := p.postfixParseFns[p.peekToken.Kind]
			if postfix == nil {
				return nil, errors.New("no postfix parse function for: " + lexer.TokenKindString(p.peekToken.Kind))
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

	// add to localFunctionMap to check at the end of parsing if function is defined or not
	if _, ok := leftExp.(*ast.CallExpression); ok {
		localFunctionMap[leftExp.(*ast.CallExpression).Name.(*ast.Identifier).Value] = true
	}

	return leftExp, nil
}

// -----------------------------------------------------------------------------
// Parsing Expressions
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// Parsing Identifier Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

// -----------------------------------------------------------------------------
// Parsing Integer Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseIntegerValue() (ast.Expression, error) {
	intVal := &ast.IntegerValue{Token: p.currentToken}
	value, err := strconv.ParseInt(p.currentToken.Value, 0, 64)
	if err != nil {
		return nil, errors.New("could not parse " + p.currentToken.Value + " as integer")
	}
	intVal.Value = value
	return intVal, nil
}

// -----------------------------------------------------------------------------
// Parsing Float Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseFloatValue() (ast.Expression, error) {
	floatVal := &ast.FloatValue{Token: p.currentToken}
	value, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		return nil, errors.New("could not parse " + p.currentToken.Value + " as float")
	}
	floatVal.Value = value
	return floatVal, nil
}

// -----------------------------------------------------------------------------
// Parsing Boolean Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseBooleanValue() (ast.Expression, error) {
	val := p.currentToken.Value
	if val == "true" {
		return &ast.BooleanValue{Token: p.currentToken, Value: true}, nil
	}
	return &ast.BooleanValue{Token: p.currentToken, Value: false}, nil
}

// -----------------------------------------------------------------------------
// Parsing String Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseStringValue() (ast.Expression, error) {
	return &ast.StringValue{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

// -----------------------------------------------------------------------------
// Parsing Char Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseCharValue() (ast.Expression, error) {
	return &ast.CharValue{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

// -----------------------------------------------------------------------------
// Parsing Array Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseArrayValues() (ast.Expression, error) {
	if !p.currTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		return nil, errors.New("expected an open square bracket (`[`) for an array, got: " + lexer.TokenKindString(p.currentToken.Kind))
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
				return nil, errors.New("expected a value after comma in an array, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			p.nextToken()
			break
		} else {
			return nil, errors.New("expected a comma (`,`) or a closing square bracket (`]`) after the value, got: " + lexer.TokenKindString(p.peekToken.Kind))
		}
	}
	array.Values = list
	return array, nil
}

// -----------------------------------------------------------------------------
// Parsing Hashmap Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseHashMapValues() (ast.Expression, error) {
	if !p.currTokenIsOk(lexer.OPEN_CURLY_BRACKET) {
		return nil, errors.New("expected an open curly bracket (`{`) for a hashmap, got: " + lexer.TokenKindString(p.currentToken.Kind))
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
			return nil, errors.New("expected a colon (`:`) after the key, got: " + lexer.TokenKindString(p.peekToken.Kind))
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
				return nil, errors.New("expected a value after comma, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
			p.nextToken()
			break
		} else {
			return nil, errors.New("expected a comma (`,`) or a closing curly bracket (`}`) after the value, got: " + lexer.TokenKindString(p.peekToken.Kind))
		}
	}
	hashMap.Pairs = pairs
	return hashMap, nil
}

// -----------------------------------------------------------------------------
// Parsing Prefix Expression
// -----------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------
// Parsing Infix Expression
// -----------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------
// Parsing Postfix Expression
// -----------------------------------------------------------------------------
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

	return expression, nil
}

// -----------------------------------------------------------------------------
// Parsing Assignment Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseAssignmentExpression(left ast.Expression) (ast.Expression, error) {
	identifier, ok := left.(*ast.Identifier)
	if !ok {
		return nil, errors.New("left side in an assignment operation must be an identifier, got: " + fmt.Sprintf("%T", left))
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

// -----------------------------------------------------------------------------
// Parsing Call Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseCallExpression(left ast.Expression) (ast.Expression, error) {
	exp := &ast.CallExpression{Token: p.currentToken, Name: left}
	expArgs, err := p.parseCallArgs(left)
	if err != nil {
		return nil, err
	}
	exp.Args = expArgs
	return exp, nil
}

func (p *Parser) parseCallArgs(left ast.Expression) ([]ast.Expression, error) {
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
		return nil, errors.New("expected a closing bracket (`)`) after the arguments in call expression " + left.(*ast.Identifier).Value + ", got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	return args, nil
}

// -----------------------------------------------------------------------------
// Parsing Index Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	exp := &ast.IndexExpression{Token: p.currentToken, Left: left}
	p.nextToken()
	idx, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	exp.Index = idx
	if !p.expectedPeekToken(lexer.CLOSE_SQUARE_BRACKET) {
		return nil, errors.New("expected a closing square bracket (`]`) after the index in index expression. got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	return exp, nil
}

// -----------------------------------------------------------------------------
// Parsing Grouped Expression
// -----------------------------------------------------------------------------
func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.nextToken()
	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("expected a closing bracket (`)`) after the grouped expression, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	return exp, nil
}

// -----------------------------------------------------------------------------
// Parsing Statements
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// Parsing Function Statement
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionStatement() (*ast.Function, error) {
	if inFunction {
		return nil, errors.New("can't declare a function inside a function")
	}
	inFunction = true
	stmt := &ast.Function{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after the `fun` keyword, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	if !p.expectedPeekToken(lexer.IDENTIFIER) {
		return nil, errors.New("expected an identifier(function name) after the colon (`:`), got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	if _, ok := FunctionMap[stmt.Name.Value]; ok && !p.inTesting {
		return nil, errors.New("can't declare a function twice, function with the same name `" + stmt.Name.Value + "` already exists")
	}
	if _, ok := builtinMap[stmt.Name.Value]; ok && !p.inTesting {
		return nil, errors.New("can't override a built-in function, function `" + stmt.Name.Value + "` already exists")
	}

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
		return nil, errors.New("expected an open curly bracket (`{`) after the function (`" + stmt.Name.Value + "`) signature, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	funBody, err := p.parseFunctionBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = funBody

	FunctionMap[stmt.Name.Value] = stmt
	inFunction = false

	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Function Parameters
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionParameters() ([]*ast.FunctionParameters, error) {
	var listToReturn []*ast.FunctionParameters

	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil, errors.New("expected an open bracket (`(`) after the function name, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	for !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		if !p.expectedPeekToken(lexer.IDENTIFIER) {
			return nil, errors.New("expected an identifier or a close bracket (`)`) after the open bracket (`(`) for function parameters, got: " + lexer.TokenKindString(p.peekToken.Kind))
		}
		param := &ast.FunctionParameters{ParameterName: &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}}

		if !p.expectedPeekToken(lexer.COLON) {
			return nil, errors.New("expected a colon (`:`) after the parameter " + param.ParameterName.Value + ", got: " + lexer.TokenKindString(p.peekToken.Kind))
		}

		if !p.expectedPeekToken(lexer.TYPE) {
			return nil, errors.New("expected a type for parameter " + param.ParameterName.Value + " after the colon (`:`), got: " + lexer.TokenKindString(p.peekToken.Kind))
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
					return nil, errors.New("expected a closing square bracket (`]`) after value type for hashmap, got: " + lexer.TokenKindString(p.peekToken.Kind))
				}
			} else {
				return nil, errors.New("Expected a closing square bracket (`]`) for array or a datatype for hashmap, got: " + lexer.TokenKindString(p.peekToken.Kind))
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
				return nil, errors.New("expected an identifier after comma (`,`) for function parameters, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
		} else {
			if !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				return nil, errors.New("expected a closing bracket (`)`) or a comma (`,`) after the parameter type, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
		}
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("expected a closing bracket (`)`) after function parameters, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	return listToReturn, nil
}

// -----------------------------------------------------------------------------
// Parsing Function Return Types
// -----------------------------------------------------------------------------
func (p *Parser) parseFunctionReturnTypes() ([]*ast.FunctionReturnType, error) {
	var listToReturn []*ast.FunctionReturnType

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after function parameters, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil, errors.New("expected an open bracket (`(`) after the colon (`:`), got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	for !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		if !p.expectedPeekToken(lexer.TYPE) {
			return nil, errors.New("expected a datatype after the open bracket (`(`), got: " + lexer.TokenKindString(p.peekToken.Kind))
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
					return nil, errors.New("expected a closing square bracket (`]`) after value datatype for hashmap, got: " + lexer.TokenKindString(p.peekToken.Kind))
				}
			} else {
				return nil, errors.New("expected a closing square bracket (`]`) for array or a datatype for hashmap, got: " + lexer.TokenKindString(p.peekToken.Kind))
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
				return nil, errors.New("expected a datatype after comma (`,`), got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
		} else {
			if !p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				return nil, errors.New("expected a closing bracket (`)`) or a comma (`,`) after the parameter datatype, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
		}
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("expected a closing bracket (`)`) after the return types, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	return listToReturn, nil
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
// Parsing VAR and CONST Statement
// -----------------------------------------------------------------------------
func (p *Parser) parseVarStmtSig() (*ast.VarStatement, error) {
	stmt := &ast.VarStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.IDENTIFIER) {
		if stmt.Token.Kind == lexer.CONST {
			return nil, errors.New("expected an identifier after `const` keyword, got: " + lexer.TokenKindString(p.peekToken.Kind))
		}
		return nil, errors.New("expected an identifier after `var` keyword, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after the identifier `" + stmt.Name.Value + "`, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	if !p.expectedPeekToken(lexer.TYPE) {
		return nil, errors.New("expected a datatype for identifier `" + stmt.Name.Value + "` after colon (`:`), got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	firstType := p.currentToken
	stmt.Type = &ast.Type{Token: firstType, Value: firstType.Value}

	if p.peekTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		p.nextToken()
		if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			stmt.Type.IsArray = true
			stmt.Type.IsHash = false
			stmt.Type.SubTypes = nil
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.TYPE) {
			p.nextToken()
			stmt.Type.IsArray = false
			stmt.Type.IsHash = true

			subTypes := []*ast.Type{}
			subTypes = append(subTypes, &ast.Type{Token: firstType, Value: firstType.Value, IsArray: false, IsHash: false, SubTypes: nil})
			subTypes = append(subTypes, &ast.Type{Token: p.currentToken, Value: p.currentToken.Value, IsArray: false, IsHash: false, SubTypes: nil})
			stmt.Type.SubTypes = subTypes

			if !p.expectedPeekToken(lexer.CLOSE_SQUARE_BRACKET) {
				return nil, errors.New("expected a closing square bracket (`]`) after value datatype for hashmap, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
		} else {
			return nil, errors.New("expected a closing square bracket (`]`) for array or a value datatype for hashmap, got: " + lexer.TokenKindString(p.peekToken.Kind))
		}
	} else {
		stmt.Type.IsArray = false
		stmt.Type.IsHash = false
		stmt.Type.SubTypes = nil
	}

	return stmt, nil
}

func (p *Parser) parseVarStatement() (ast.Statement, error) {
	if !inFunction && !p.inTesting {
		return nil, errors.New("can't declare a variable outside a function")
	}
	stmt, err := p.parseVarStmtSig()
	if err != nil {
		return nil, err
	}

	if p.peekTokenIsOk(lexer.COMMA) {
		list := []ast.Statement{}
		list = append(list, stmt)
		p.nextToken()
		p.nextToken()
		return p.parseMultipleAssignmentStatement(list)
	}

	if stmt.Type.IsArray && !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		return nil, errors.New("array `" + stmt.Name.Value + "` must always be initialized while declaring, for empty array use `[]`")
	} else if stmt.Type.IsHash && !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		return nil, errors.New("hashmap `" + stmt.Name.Value + "` must always be initialized while declaring, for empty hashmap use `{}`")
	} else if stmt.Token.Kind == lexer.CONST && !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		return nil, errors.New("const variable `" + stmt.Name.Value + "` must be initialized while declaring")
	}

	if p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		p.nextToken()
		p.nextToken()
		parsedExp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		stmt.Value = parsedExp

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
		return nil, errors.New("expected a semicolon (`;`) at the end of the statement after variable declaration, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Multiple Assignment Statement
// -----------------------------------------------------------------------------
func (p *Parser) parseMultipleAssignmentStatement(list []ast.Statement) (*ast.MultiValueAssignStmt, error) {
	if !inFunction && !p.inTesting {
		return nil, errors.New("can't declare a variable outside a function")
	}
	stmt := &ast.MultiValueAssignStmt{Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="}}

	if p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
		return nil, errors.New("expected an identifier or a `var`/`const` keyword after the comma (`,`), got: " + lexer.TokenKindString(p.currentToken.Kind))
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
			return nil, errors.New("expected an identifier or a `var`/`const` keyword, got: " + lexer.TokenKindString(p.currentToken.Kind))
		}

		if p.currTokenIsOk(lexer.COMMA) {
			if p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
				return nil, errors.New("expected an identifier or a `var`/`const` keyword after the comma (`,`), got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
			p.nextToken()
		} else if p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
			break
		} else {
			return nil, errors.New("expected a comma (`,`) or an equal sign (`=`) after identifier, got: " + lexer.TokenKindString(p.currentToken.Kind))
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
			if p.peekTokenIsOk(lexer.SEMI_COLON) {
				return nil, errors.New("expected an expression after comma (`,`) in multi-value assignment statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
			p.nextToken()
		} else {
			if !p.currTokenIsOk(lexer.SEMI_COLON) {
				return nil, errors.New("expected a comma (`,`) or a semicolon (`;`) after the expression, got: " + lexer.TokenKindString(p.currentToken.Kind))
			}
		}
	}

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
			return nil, errors.New("number of expressions on the right side of multi-assignment = 1, expected a function call, got: " + fmt.Sprintf("%T", valueList[0]))
		}
	} else {
		if len(valueList) != len(list) {
			return nil, errors.New("number of expression on the right side of multi-assignment do not match the number of declarations on the left, left=" + fmt.Sprint(len(list)) + ", right=" + fmt.Sprint(len(valueList)))
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
					return nil, errors.New("expected an identifier or a `var`/`const` keyword, got: " + fmt.Sprintf("%T", obj))
				}
			}
		}
	}

	stmt.Objects = list
	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Expression Statement
// -----------------------------------------------------------------------------
func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
	if !inFunction && !p.inTesting {
		return nil, errors.New("everything must be inside a function")
	}
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

	if !p.inTesting {
		if _, ok := parsedExp.(*ast.CallExpression); !ok {
			if _, ok = parsedExp.(*ast.PostfixExpression); !ok {
				if _, ok = parsedExp.(*ast.AssignmentExpression); !ok {
					return nil, errors.New("expected a function call, postfix expression or an assignment expression for expressions as statements, got: " + fmt.Sprintf("%T", parsedExp))
				}
			}
		}
	}

	stmt.Expression = parsedExp
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		if !p.inTesting {
			return nil, errors.New("expected a semicolon (`;`) at the end of the statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
		}
	}

	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Return Statement
// -----------------------------------------------------------------------------
func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	if !inFunction && !p.inTesting {
		return nil, errors.New("return statement can only be used inside a function")
	}

	stmt := &ast.ReturnStatement{Token: p.currentToken}

	if p.peekTokenIsOk(lexer.SEMI_COLON) {
		stmt.Value = nil
		p.nextToken()
		return stmt, nil
	}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after the `return` keyword, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	if p.peekTokenIsOk(lexer.OPEN_BRACKET) {
		p.nextToken()
		if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
			return nil, errors.New("expected values after open bracket (`(`) in `return` statement")
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
					return nil, errors.New("expected an expression after comma (`,`) in `return` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
				}
				p.nextToken()
			} else if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				p.nextToken()
				break
			} else {
				return nil, errors.New("expected a comma (`,`) or a closing bracket (`)`) after the expression, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
		}
	} else {
		p.nextToken()
		parsedExp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		stmt.Value = append(stmt.Value, parsedExp)
		if p.peekTokenIsOk(lexer.COMMA) {
			return nil, errors.New("expected a semicolon (`;`), got: " + lexer.TokenKindString(p.peekToken.Kind) + ". To return multiple values, use `return: (val1, val2, ...)`")
		}
	}

	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil, errors.New("expected a semicolon (`;`) at the end of `return` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Continue Statement
// -----------------------------------------------------------------------------
func (p *Parser) parseContinueStatement() (*ast.ContinueStatement, error) {
	if !inForLoop {
		return nil, errors.New("continue statement can only be used inside a `for loop`")
	}
	stmt := &ast.ContinueStatement{Token: p.currentToken}
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil, errors.New("expected a semicolon (`;`) at the end of `continue` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing Break Statement
// -----------------------------------------------------------------------------
func (p *Parser) parseBreakStatement() (*ast.BreakStatement, error) {
	if !inForLoop {
		return nil, errors.New("break statement can only be used inside a `for loop`")
	}
	stmt := &ast.BreakStatement{Token: p.currentToken}
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil, errors.New("expected a semicolon (`;`) at the end of `break` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	return stmt, nil
}

// -----------------------------------------------------------------------------
// Parsing If statement
// -----------------------------------------------------------------------------
func (p *Parser) parseIfStatement() (*ast.IfStatement, error) {
	if !inFunction && !p.inTesting {
		return nil, errors.New("if statement can only be used inside a function")
	}
	stmt := &ast.IfStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after `if` keyword, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil, errors.New("expected an open bracket (`(`) after the colon (`:`) in `if` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	groupedExpValue, err := p.parseGroupedExpression()
	if err != nil {
		return nil, err
	}
	stmt.Value = groupedExpValue

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after grouped expression for `if` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil, errors.New("expected an open curly bracket (`{`) after the colon (`:`) in `if` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
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
				return nil, errors.New("expected a colon (`:`) after the `else if` keyword, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
			if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
				return nil, errors.New("expected an open bracket (`(`) after the colon (`:`) in `else if` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}

			elseIfValue, err := p.parseGroupedExpression()
			if err != nil {
				return nil, err
			}
			elseIfStmt.Value = elseIfValue

			if !p.expectedPeekToken(lexer.COLON) {
				return nil, errors.New("expected a colon (`:`) after grouped expression for `else if` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
			}
			if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
				return nil, errors.New("expected an open curly bracket (`{`) after the colon (`:`) in `else if` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
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
			return nil, errors.New("expected a colon (`:`) after the `else` keyword, got: " + lexer.TokenKindString(p.peekToken.Kind))
		}
		if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
			return nil, errors.New("expected an open curly bracket (`{`) after the colon (`:`) in `else` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
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
	if !inFunction && !p.inTesting {
		return nil, errors.New("for loop can only be used inside a function")
	}
	inForLoop = true
	stmt := &ast.ForLoopStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after the `for` keyword, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil, errors.New("expected an open bracket (`(`) after the colon (`:`) in `for loop` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	p.nextToken()
	parsedStmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	if _, ok := parsedStmt.(*ast.VarStatement); !ok {
		if _, ok := parsedStmt.(*ast.ExpressionStatement); !ok {
			return nil, errors.New("expected a `var` statement or an assignment expression after an open bracket (`(`) in `for loop` statement, got: " + fmt.Sprintf("%T", parsedStmt))
		} else {
			assign, ok := parsedStmt.(*ast.ExpressionStatement).Expression.(*ast.AssignmentExpression)
			if !ok {
				return nil, errors.New("expected a `var` statement or an assignment expression after an open bracket (`(`) in `for loop` statement, got: " + fmt.Sprintf("%T", parsedStmt.(*ast.ExpressionStatement).Expression))
			} else {
				if assign.Operator != "=" {
					return nil, errors.New("expected assignment-equal operator (`=`) in assignment expression in `for loop` statement, got: " + assign.Operator)
				}
			}
		}
	}
	stmt.Left = parsedStmt

	p.nextToken()
	parsedExp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	midStmt, ok := parsedExp.(*ast.InfixExpression)
	if !ok {
		return nil, errors.New("expected an infix expression after the variable declaration in `for loop`, got: " + fmt.Sprintf("%T", parsedExp))
	}
	stmt.Middle = midStmt

	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil, errors.New("expected a semicolon (`;`) after infix expression in `for loop` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	p.nextToken()
	parsedExp, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if _, ok := parsedExp.(*ast.PostfixExpression); !ok {
		if _, ok := parsedExp.(*ast.AssignmentExpression); !ok {
			return nil, errors.New("expected a postfix expression or an assignment expression after the infix expression in `for loop`, got: " + fmt.Sprintf("%T", parsedExp))
		}
	}
	stmt.Right = parsedExp

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("expected a closing bracket (`)`) after the postfix expression in `for loop`, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after the closing bracket (`)`) in `for loop`, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil, errors.New("expected an open curly bracket (`{`) after the colon (`:`) in `for loop`, got: " + lexer.TokenKindString(p.peekToken.Kind))
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
// Parsing While Loop
// -----------------------------------------------------------------------------
func (p *Parser) parseWhileLoop() (*ast.WhileLoopStatement, error) {
	if !inFunction && !p.inTesting {
		return nil, errors.New("while loop can only be used inside a function")
	}
	inForLoop = true
	stmt := &ast.WhileLoopStatement{Token: p.currentToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after the `while` keyword, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil, errors.New("expected an open bracket (`(`) after the colon (`:`) in `while loop` statement, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}

	p.nextToken()
	parsedExp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Condition = parsedExp

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil, errors.New("expected a closing bracket (`)`) after break condition in `while loop`, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	if !p.expectedPeekToken(lexer.COLON) {
		return nil, errors.New("expected a colon (`:`) after the closing bracket (`)`) in `while loop`, got: " + lexer.TokenKindString(p.peekToken.Kind))
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil, errors.New("expected an open curly bracket (`{`) after the colon (`:`) in `while loop`, got: " + lexer.TokenKindString(p.peekToken.Kind))
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
// Parsing Types
// -----------------------------------------------------------------------------
func (p *Parser) parseType() (*ast.Type, error) {
	stmt := &ast.Type{Token: p.currentToken, Value: p.currentToken.Value}
	if !p.peekTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		p.nextToken()
		stmt.IsArray = false
		stmt.IsHash = false
		stmt.SubTypes = nil
		return stmt, nil
	} else {
		p.nextToken()
		// int[string[]]
		// int[string]
		// int[int[string]]
		// int[][]
		// int[string][]
		// int[string[]][]
		if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			p.nextToken()
			stmt.IsArray = true
			stmt.IsHash = false
			stmt.SubTypes = nil
			p.nextToken()
			return stmt, nil
		} else {
			p.nextToken()

			valueType, err := p.parseType()
			if err != nil {
				return nil, err
			}

			if !p.currTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
				return nil, errors.New("expected a closing square bracket (`]`) after value datatype for hashmap, got: " + lexer.TokenKindString(p.currentToken.Kind))
			}

			stmt.IsArray = false
			stmt.IsHash = true
			var subTypes []*ast.Type

			keyType := &ast.Type{Token: stmt.Token, Value: stmt.Value, IsArray: false, IsHash: false, SubTypes: nil}
			subTypes = append(subTypes, keyType)
			subTypes = append(subTypes, valueType)

			stmt.SubTypes = subTypes
			p.nextToken()
			return stmt, nil
		}
	}
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
