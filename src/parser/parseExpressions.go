package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Expressions
// ------------------------------------------------------------------------------------------------------------------
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

type (
	prefixParseFn  func() (ast.Expression, error)
	infixParseFn   func(ast.Expression) (ast.Expression, error)
	postfixParseFn func(ast.Expression) (ast.Expression, error)
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

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefixParseFns[p.currToken.Kind]
	if prefix == nil {
		return nil, errors.New("no prefix parse function for: " +
			lexer.TokenKindString(p.currToken.Kind))
	}
	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for !p.peekTokenIsOk(lexer.SEMI_COLON) && precedence < p.peekPrecedence() {
		if p.postfixParseFns[p.peekToken.Kind] != nil {
			postfix := p.postfixParseFns[p.peekToken.Kind]
			if postfix == nil {
				return nil, errors.New("no postfix parse function for: " +
					lexer.TokenKindString(p.peekToken.Kind))
			}
			p.nextToken()
			leftExp, err = postfix(leftExp)
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

// ------------------------------------------------------------------------------------------------------------------
// Identifier
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseIdentifier() (ast.Expression, error) {
	exp := &ast.Identifier{Token: p.currToken, Value: p.currToken.Value}
	if p.peekTokenIsOk(lexer.OPEN_BRACKET) {
		return exp, nil
	}
	t, err := typeCheckIdent(exp, p.stack.Top())
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Integer
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseInteger() (ast.Expression, error) {
	exp := &ast.Integer{Token: p.currToken}
	val, err := strconv.ParseInt(p.currToken.Value, 0, 64)
	if err != nil {
		return nil, errors.New("could not parse " + p.currToken.Value + " as integer")
	}
	exp.Value = val
	t, err := typeCheckInteger()
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Float
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseFloat() (ast.Expression, error) {
	exp := &ast.Float{Token: p.currToken}
	val, err := strconv.ParseFloat(p.currToken.Value, 64)
	if err != nil {
		return nil, errors.New("could not parse " + p.currToken.Value + " as float")
	}
	exp.Value = val
	t, err := typeCheckFloat()
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Boolean
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseBoolean() (ast.Expression, error) {
	exp := &ast.Bool{Token: p.currToken}
	if p.currToken.Value == "true" {
		exp.Value = true
	} else {
		exp.Value = false
	}
	t, err := typeCheckBool()
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// String
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseString() (ast.Expression, error) {
	exp := &ast.String{Token: p.currToken, Value: p.currToken.Value}
	t, err := typeCheckString()
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Char
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseChar() (ast.Expression, error) {
	exp := &ast.Char{Token: p.currToken, Value: p.currToken.Value}
	t, err := typeCheckChar()
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// HashMap
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseHashMap() (ast.Expression, error) {
	if !p.currTokenIsOk(lexer.OPEN_CURLY_BRACKET) {
		return nil,
			errors.New(
				"expected an open curly bracket (`{`) for a hashmap, got: " +
					lexer.TokenKindString(p.currToken.Kind),
			)
	}
	exp := &ast.HashMap{
		Token:     p.currToken,
		KeyType:   nil,
		ValueType: nil,
		Pairs:     map[ast.BaseType]ast.Expression{},
	}
	p.nextToken()
	if p.currTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
		return exp, nil
	}

	for {
		kExp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		k, ok := kExp.(ast.BaseType)
		if !ok {
			return nil,
				errors.New(
					"key in a hashmap can only be of `BaseType`, got: " +
						fmt.Sprintf("%T", kExp),
				)
		}
		if !p.expectedPeekToken(lexer.COLON) {
			return nil,
				errors.New(
					"expected a colon (`:`) after the key, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
		p.nextToken()

		vExp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		exp.Pairs[k] = vExp

		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
				return nil,
					errors.New(
						"expected a value after comma, got: " +
							lexer.TokenKindString(p.peekToken.Kind),
					)
			}
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
			p.nextToken()
			break
		} else {
			return nil,
				errors.New(
					"expected a comma (`,`) or a closing curly bracket (`}`) after the value, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
	}
	t, err := typeCheckHashMap(exp, p.stack.Top())
	if err != nil {
		return nil, err
	}
	exp.KeyType = t.Types[0].KeyType
	exp.ValueType = t.Types[0].ValueType
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Array
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseArray() (ast.Expression, error) {
	if !p.currTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		return nil,
			errors.New(
				"expected an open square bracket (`[`) for an array, got: " +
					lexer.TokenKindString(p.currToken.Kind),
			)
	}
	exp := &ast.Array{Token: p.currToken, Values: []ast.Expression{}, Type: nil}
	p.nextToken()
	if p.currTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
		return exp, nil
	}
	for {
		var val ast.Expression
		var err error

		if p.currTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
			val, err = p.parseArray()
		} else {
			val, err = p.parseExpression(LOWEST)
		}
		if err != nil {
			return nil, err
		}

		exp.Values = append(exp.Values, val)

		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
				return nil,
					errors.New(
						"expected a value after comma in an array, got: " +
							lexer.TokenKindString(p.peekToken.Kind),
					)
			}
			p.nextToken()
		} else if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			p.nextToken()
			break
		} else {
			return nil,
				errors.New(
					"expected a comma (`,`) or a closing square bracket (`]`) after the value, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
	}
	t, err := typeCheckArray(exp, p.stack.Top())
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0].ElementType
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Prefix
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parsePrefix() (ast.Expression, error) {
	exp := &ast.Prefix{Token: p.currToken, Operator: p.currToken.Value}
	p.nextToken()
	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}
	exp.Right = right
	t, err := typeCheckPrefix(exp, p.stack.Top())
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Infix
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseInfix(left ast.Expression) (ast.Expression, error) {
	exp := &ast.Infix{Token: p.currToken, Operator: p.currToken.Value, Left: left}
	precedence := p.currentPrecedence()
	p.nextToken()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	exp.Right = right
	t, err := typeCheckInfix(exp, p.stack.Top())
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Postfix
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parsePostfix(left ast.Expression) (ast.Expression, error) {
	exp := &ast.Postfix{Token: p.currToken, Operator: p.currToken.Value, Left: left}
	t, err := typeCheckPostfix(exp, p.stack.Top())
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Assignment
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseAssignment(left ast.Expression) (ast.Expression, error) {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return nil, errors.New(
			"left side in an assignment operation must be an identifier, got: " +
				fmt.Sprintf("%T", left),
		)
	}
	exp := &ast.Assignment{Token: p.currToken, Left: ident, Operator: p.currToken.Value}
	p.nextToken()
	right, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	exp.Right = right
	t, err := typeCheckAssignment(exp, p.stack.Top())
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// CallExpression
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseCall(left ast.Expression) (ast.Expression, error) {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return nil, errors.New(
			"function name in a call expression must be an identifier, got: " +
				fmt.Sprintf("%T", left),
		)
	}
	exp := &ast.CallExpression{Token: p.currToken, Name: ident, Args: nil}
	args, err := p.parseCallArgs(ident)
	if err != nil {
		return nil, err
	}
	exp.Args = args
	t, err := typeCheckCallExp(exp, p.stack.Top())
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types
	return exp, nil
}

func (p *Parser) parseCallArgs(left *ast.Identifier) ([]ast.Expression, error) {
	if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		p.nextToken()
		return nil, nil
	}
	var args []ast.Expression

	p.nextToken()
	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	args = append(args, exp)

	for p.peekTokenIsOk(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		exp, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, exp)
	}

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil,
			errors.New(
				"expected a closing bracket (`)`) after the arguments in call expression " +
					left.Value + ", got: " + lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	return args, nil
}

// ------------------------------------------------------------------------------------------------------------------
// IndexExpression
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseIndex(left ast.Expression) (ast.Expression, error) {
	exp := &ast.IndexExpression{Token: p.currToken, Left: left}
	p.nextToken()
	idx, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	exp.Index = idx
	if !p.expectedPeekToken(lexer.CLOSE_SQUARE_BRACKET) {
		return nil,
			errors.New(
				"expected a closing square bracket (`]`) after the index in index expression. got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	t, err := typeCheckIndexExp(exp, p.stack.Top())
	if err != nil {
		return nil, err
	}
	exp.Type = t.Types[0]
	return exp, nil
}

// ------------------------------------------------------------------------------------------------------------------
// GroupedExp
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseGroupedExp() (ast.Expression, error) {
	p.nextToken()
	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil,
			errors.New(
				"expected a closing bracket (`)`) after the grouped expression, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	return exp, nil
}
