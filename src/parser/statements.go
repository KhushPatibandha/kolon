package parser

import (
	"errors"
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Statements
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currToken.Kind {
	case lexer.VAR:
		return p.parseVarConst()
	case lexer.CONST:
		return p.parseVarConst()
	case lexer.RETURN:
		return p.parseReturn()
	case lexer.FUN:
		return p.parseFunction()
	case lexer.IF:
		return p.parseIf()
	case lexer.FOR:
		return p.parseForLoop()
	case lexer.WHILE:
		return p.parseWhileLoop()
	case lexer.CONTINUE:
		return p.parseContinue()
	case lexer.BREAK:
		return p.parseBreak()
	default:
		return p.parseExpressionStatement()
	}
}

// ------------------------------------------------------------------------------------------------------------------
// Body
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseBody() (*ast.Body, error) {
	body := &ast.Body{Token: &p.currToken, Statements: []ast.Statement{}}
	p.nextToken()
	for !p.currTokenIsOk(lexer.CLOSE_CURLY_BRACKET) && !p.currTokenIsOk(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		body.Statements = append(body.Statements, stmt)
		p.nextToken()
	}
	return body, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Types
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseType() (*ast.Type, error) {
	if !p.expectedPeekToken(lexer.TYPE) {
		return nil,
			errors.New(
				"expected a type, got: " + lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	stmt := &ast.Type{
		Kind:        ast.TypeBase,
		Token:       &p.currToken,
		Name:        p.currToken.Value,
		ElementType: nil,
		KeyType:     nil,
		ValueType:   nil,
	}

	for p.peekTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		p.nextToken()

		if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			stmt = &ast.Type{
				Kind:        ast.TypeArray,
				Token:       stmt.Token,
				ElementType: stmt,
				Name:        "",
				KeyType:     nil,
				ValueType:   nil,
			}
			p.nextToken()
			continue
		} else if !p.peekTokenIsOk(lexer.TYPE) {
			return nil,
				errors.New(
					"expected a closing square bracket (`]`) " +
						"for array or a datatype for hashmap, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
		val, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if !p.expectedPeekToken(lexer.CLOSE_SQUARE_BRACKET) {
			return nil,
				errors.New(
					"expected a closing square bracket (`]`) after value type for hashmap, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
		stmt = &ast.Type{
			Kind:        ast.TypeHashMap,
			Token:       stmt.Token,
			KeyType:     stmt,
			ValueType:   val,
			Name:        "",
			ElementType: nil,
		}
	}
	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Expression Statements
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
	if !p.inFunction && !p.inTesting {
		return nil, errors.New("everything must be inside a function")
	}
	stmt := &ast.ExpressionStatement{Token: &p.currToken}

	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// VarAndConst
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseVarConst() (*ast.VarAndConst, error) {
	if !p.inFunction && !p.inTesting {
		return nil, errors.New("can't declare a variable outside a function")
	}
	stmt, err := p.parseVarConstSig()
	if err != nil {
		return nil, err
	}

	switch stmt.Type.Kind {
	case ast.TypeArray:
		if !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
			return nil,
				errors.New(
					"array `" + stmt.Name.Value + "` must always be " +
						"initialized while declaring, for empty array use `[]`")
		}
	case ast.TypeHashMap:
		if !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
			return nil,
				errors.New(
					"hashmap `" + stmt.Name.Value +
						"` must always be initialized while declaring, for empty hashmap use `{}`")
		}
	default:
		if stmt.Token.Kind == lexer.CONST && !p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
			return nil,
				errors.New(
					"const variable `" + stmt.Name.Value +
						"` must be initialized while declaring")
		}
	}

	var value ast.Expression
	if p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		p.nextToken()
		p.nextToken()
		value, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		p.assignTypeToValue(value, stmt.Type)
	} else {
		value, err = p.assignDefaultValue(stmt.Type)
		if err != nil {
			return nil, err
		}
	}
	stmt.Value = value

	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil,
			errors.New(
				"expected a semicolon (`;`) at the end of the statement after " +
					"variable declaration, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	return stmt, nil
}

func (p *Parser) parseVarConstSig() (*ast.VarAndConst, error) {
	stmt := &ast.VarAndConst{Token: &p.currToken}
	if !p.expectedPeekToken(lexer.IDENTIFIER) {
		return nil,
			errors.New(
				"expected an identifier after `" +
					p.currToken.Value + "` keyword, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	stmt.Name = &ast.Identifier{Token: &p.currToken, Value: p.currToken.Value}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after the identifier `" +
					stmt.Name.Value + "`, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}

	typ, err := p.parseType()
	if err != nil {
		return nil, err
	}
	stmt.Type = typ

	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Function
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseFunction() (*ast.Function, error) {
	if p.inFunction {
		return nil, errors.New("can't declare a function inside a function")
	}
	stmt := &ast.Function{Token: &p.currToken, Parameters: nil, ReturnTypes: nil, Body: nil}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after the `fun` keyword, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.IDENTIFIER) {
		return nil,
			errors.New(
				"expected an identifier(function name) after the colon (`:`), got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	stmt.Name = &ast.Identifier{Token: &p.currToken, Value: p.currToken.Value}

	if existing, ok := p.functionMap[stmt.Name.Value]; ok && existing.Body != nil && !p.inTesting {
		return nil,
			errors.New(
				"can't declare a function twice, function with the same name `" +
					stmt.Name.Value + "` already exists",
			)
	}
	if _, ok := p.builtinMap[stmt.Name.Value]; ok && !p.inTesting {
		return nil,
			errors.New(
				"can't override a built-in function, function `" +
					stmt.Name.Value + "` already exists",
			)
	}

	params, err := p.parseFunctionParams()
	if err != nil {
		return nil, err
	}
	stmt.Parameters = params

	if p.peekTokenIsOk(lexer.COLON) {
		ret, err := p.parseFunctionReturnTypes()
		if err != nil {
			return nil, err
		}
		stmt.ReturnTypes = ret
	} else if p.peekTokenIsOk(lexer.SEMI_COLON) {
		p.nextToken()
		if _, ok := p.functionMap[stmt.Name.Value]; ok {
			return nil,
				errors.New(
					"can't declare a function twice, function with the same name `" +
						stmt.Name.Value + "` already exists",
				)
		}
		p.functionMap[stmt.Name.Value] = stmt
		return stmt, nil
	}

	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil,
			errors.New(
				"expected an open curly bracket (`{`) after the function (`" +
					stmt.Name.Value + "`) signature, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}

	existing, ok := p.functionMap[stmt.Name.Value]
	if ok && !p.compareFunctionSig(existing, stmt) {
		return nil,
			errors.New(
				"function signature of " + stmt.Name.Value +
					" doesn't match with previously declared signature",
			)
	}

	p.inFunction = true
	funBody, err := p.parseBody()
	if err != nil {
		return nil, err
	}
	p.inFunction = false

	if ok {
		existing.Body = funBody
		return existing, nil
	}

	stmt.Body = funBody
	p.functionMap[stmt.Name.Value] = stmt
	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Function Params
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseFunctionParams() ([]*ast.FunctionParameter, error) {
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil,
			errors.New(
				"expected an open bracket (`(`) after the function name, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		p.nextToken()
		return nil, nil
	}
	params := []*ast.FunctionParameter{}
	for {
		if !p.expectedPeekToken(lexer.IDENTIFIER) {
			return nil,
				errors.New(
					"expected an identifier or a close bracket (`)`) after " +
						"the open bracket (`(`) for function parameters, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
		param := &ast.FunctionParameter{
			ParameterName: &ast.Identifier{Token: &p.currToken, Value: p.currToken.Value},
		}
		if !p.expectedPeekToken(lexer.COLON) {
			return nil,
				errors.New(
					"expected a colon (`:`) after the parameter " +
						param.ParameterName.Value +
						", got: " + lexer.TokenKindString(p.peekToken.Kind),
				)
		}
		paramType, err := p.parseType()
		if err != nil {
			return nil, err
		}
		param.ParameterType = paramType
		params = append(params, param)

		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				return nil,
					errors.New(
						"expected an identifier after comma (`,`) for function parameters, got: " +
							lexer.TokenKindString(p.peekToken.Kind),
					)
			}
		} else if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
			break
		} else {
			return nil,
				errors.New(
					"expected a closing bracket (`)`) or a comma (`,`) after the parameter type, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
	}
	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil,
			errors.New(
				"expected a closing bracket (`)`) after function parameters, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	return params, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Function Return Types
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseFunctionReturnTypes() ([]*ast.Type, error) {
	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after function parameters, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil,
			errors.New(
				"expected an open bracket (`(`) after the colon (`:`), got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
		return nil,
			errors.New(
				"expected at least one return type after open bracket (`(`), got: CLOSE_BRACKET",
			)
	}

	var returnTypes []*ast.Type

	for {
		retType, err := p.parseType()
		if err != nil {
			return nil, err
		}
		returnTypes = append(returnTypes, retType)

		if p.peekTokenIsOk(lexer.COMMA) {
			p.nextToken()
			if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				return nil,
					errors.New(
						"expected a datatype after comma (`,`), got: " +
							lexer.TokenKindString(p.peekToken.Kind),
					)
			}
		} else if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
			break
		} else {
			return nil,
				errors.New(
					"expected a closing bracket (`)`) or a comma (`,`) " +
						"after the parameter datatype, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
	}
	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil,
			errors.New(
				"expected a closing bracket (`)`) after the return types, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	return returnTypes, nil
}

// ------------------------------------------------------------------------------------------------------------------
// If
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseIf() (*ast.If, error) {
	if !p.inFunction && !p.inTesting {
		return nil, errors.New("if statement can only be used inside a function")
	}
	stmt := &ast.If{Token: &p.currToken}
	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after `if` keyword, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil,
			errors.New(
				"expected an open bracket (`(`) after the colon (`:`) in `if` statement, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	condition, err := p.parseGroupedExp()
	if err != nil {
		return nil, err
	}
	stmt.Condition = condition

	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after grouped expression for `if` statement, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil,
			errors.New("expected an open curly bracket (`{`) after the " +
				"colon (`:`) in `if` statement, got: " +
				lexer.TokenKindString(p.peekToken.Kind),
			)
	}

	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	// ------------------------------------------------------------------------------------------------------------------
	// Else If
	// ------------------------------------------------------------------------------------------------------------------
	if p.peekTokenIsOk(lexer.ELSE_IF) {
		var elseIfList []*ast.ElseIf

		for p.peekTokenIsOk(lexer.ELSE_IF) {
			p.nextToken()
			elseIfStmt := &ast.ElseIf{Token: &p.currToken}
			if !p.expectedPeekToken(lexer.COLON) {
				return nil,
					errors.New(
						"expected a colon (`:`) after the `else if` keyword, got: " +
							lexer.TokenKindString(p.peekToken.Kind),
					)
			}
			if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
				return nil,
					errors.New(
						"expected an open bracket (`(`) after the colon (`:`) in `else if` statement, got: " +
							lexer.TokenKindString(p.peekToken.Kind),
					)
			}
			elseIfCondition, err := p.parseGroupedExp()
			if err != nil {
				return nil, err
			}
			elseIfStmt.Condition = elseIfCondition

			if !p.expectedPeekToken(lexer.COLON) {
				return nil,
					errors.New(
						"expected a colon (`:`) after grouped " +
							"expression for `else if` statement, got: " +
							lexer.TokenKindString(p.peekToken.Kind),
					)
			}
			if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
				return nil,
					errors.New(
						"expected an open curly bracket (`{`) after " +
							"the colon (`:`) in `else if` statement, got: " +
							lexer.TokenKindString(p.peekToken.Kind),
					)
			}

			elseIfBody, err := p.parseBody()
			if err != nil {
				return nil, err
			}
			elseIfStmt.Body = elseIfBody
			elseIfList = append(elseIfList, elseIfStmt)
		}
		stmt.MultiConditionals = elseIfList
	} else {
		stmt.MultiConditionals = nil
	}

	// ------------------------------------------------------------------------------------------------------------------
	// Else
	// ------------------------------------------------------------------------------------------------------------------
	if p.peekTokenIsOk(lexer.ELSE) {
		p.nextToken()
		elseStmt := &ast.Else{Token: &p.currToken}

		if !p.expectedPeekToken(lexer.COLON) {
			return nil,
				errors.New(
					"expected a colon (`:`) after the `else` keyword, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
		if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
			return nil,
				errors.New(
					"expected an open curly bracket (`{`) after the " +
						"colon (`:`) in `else` statement, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}

		elseBody, err := p.parseBody()
		if err != nil {
			return nil, err
		}
		elseStmt.Body = elseBody
		stmt.Alternate = elseStmt
	} else {
		stmt.Alternate = nil
	}

	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Return
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseReturn() (*ast.Return, error) {
	if !p.inFunction && !p.inTesting {
		return nil, errors.New("return statement can only be used inside a function")
	}
	stmt := &ast.Return{Token: &p.currToken, Value: []ast.Expression{}}
	if p.peekTokenIsOk(lexer.SEMI_COLON) {
		stmt.Value = nil
		p.nextToken()
		return stmt, nil
	}
	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after the `return` keyword, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}

	if p.peekTokenIsOk(lexer.OPEN_BRACKET) {
		p.nextToken()
		if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
			return nil, errors.New("expected values after open bracket (`(`) in `return` statement")
		}
		p.nextToken()
		for {
			exp, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			stmt.Value = append(stmt.Value, exp)
			if p.peekTokenIsOk(lexer.COMMA) {
				p.nextToken()
				if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
					return nil,
						errors.New(
							"expected an expression after comma (`,`) in `return` statement, got: " +
								lexer.TokenKindString(p.peekToken.Kind),
						)
				}
				p.nextToken()
			} else if p.peekTokenIsOk(lexer.CLOSE_BRACKET) {
				p.nextToken()
				break
			} else {
				return nil,
					errors.New("expected a comma (`,`) or a closing bracket (`)`) after the expression, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
					)
			}
		}
	} else {
		p.nextToken()
		exp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		stmt.Value = append(stmt.Value, exp)
		if p.peekTokenIsOk(lexer.COMMA) {
			return nil,
				errors.New(
					"expected a semicolon (`;`), got: " +
						lexer.TokenKindString(p.peekToken.Kind) +
						". To return multiple values, use `return: (val1, val2, ...)`",
				)
		}
	}

	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil,
			errors.New(
				"expected a semicolon (`;`) at the end of `return` statement, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// ForLoop
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseForLoop() (*ast.ForLoop, error) {
	if !p.inFunction && !p.inTesting {
		return nil, errors.New("for loop can only be used inside a function")
	}
	p.inLoop = true
	stmt := &ast.ForLoop{Token: &p.currToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after the `for` keyword, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil,
			errors.New(
				"expected an open bracket (`(`) after the colon (`:`) in `for loop` statement, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}

	p.nextToken()
	left, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	if _, ok := left.(*ast.VarAndConst); !ok {
		if _, ok := left.(*ast.ExpressionStatement); !ok {
			return nil,
				errors.New(
					"expected a `var` statement or an assignment expression " +
						"after an open bracket (`(`) in `for loop` statement, got: " +
						fmt.Sprintf("%T", left),
				)
		} else {
			assign, ok := left.(*ast.ExpressionStatement).Expression.(*ast.Assignment)
			if !ok {
				return nil,
					errors.New(
						"expected a `var` statement or an assignment expression " +
							"after an open bracket (`(`) in `for loop` statement, got: " +
							fmt.Sprintf("%T", left.(*ast.ExpressionStatement).Expression),
					)
			} else {
				if assign.Operator != "=" {
					return nil,
						errors.New(
							"expected assignment-equal operator (`=`) in assignment " +
								"expression in `for loop` statement, got: " +
								assign.Operator,
						)
				}
			}
		}
	}
	stmt.Left = left

	p.nextToken()
	middle, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if _, ok := middle.(*ast.Infix); !ok {
		return nil,
			errors.New(
				"expected an infix expression after the variable declaration in `for loop`, got: " +
					fmt.Sprintf("%T", middle),
			)
	}
	stmt.Middle = middle.(*ast.Infix)
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil,
			errors.New(
				"expected a semicolon (`;`) after infix expression in `for loop` statement, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}

	p.nextToken()
	right, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if _, ok := right.(*ast.Postfix); !ok {
		if _, ok := right.(*ast.Assignment); !ok {
			return nil,
				errors.New(
					"expected a postfix expression or an assignment expression " +
						"after the infix expression in `for loop`, got: " +
						fmt.Sprintf("%T", right),
				)
		}
		stmt.Right = right.(*ast.Assignment)
	} else {
		stmt.Right = right.(*ast.Postfix)
	}
	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil,
			errors.New(
				"expected a closing bracket (`)`) after the postfix expression in `for loop`, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after the closing bracket (`)`) in `for loop`, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil,
			errors.New(
				"expected an open curly bracket (`{`) after the colon (`:`) in `for loop`, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	p.inLoop = false
	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// WhileLoop
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseWhileLoop() (*ast.WhileLoop, error) {
	if !p.inFunction && !p.inTesting {
		return nil, errors.New("while loop can only be used inside a function")
	}
	p.inLoop = true
	stmt := &ast.WhileLoop{Token: &p.currToken}

	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after the `while` keyword, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.OPEN_BRACKET) {
		return nil,
			errors.New(
				"expected an open bracket (`(`) after the colon (`:`) in `while loop` statement, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	p.nextToken()
	condition, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Condition = condition

	if !p.expectedPeekToken(lexer.CLOSE_BRACKET) {
		return nil,
			errors.New(
				"expected a closing bracket (`)`) after break condition in `while loop`, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.COLON) {
		return nil,
			errors.New(
				"expected a colon (`:`) after the closing bracket (`)`) in `while loop`, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil,
			errors.New(
				"expected an open curly bracket (`{`) after the colon (`:`) in `while loop`, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	p.inLoop = false
	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Continue
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseContinue() (*ast.Continue, error) {
	if !p.inLoop {
		return nil, errors.New("continue statement can only be used inside a `for loop`")
	}
	stmt := &ast.Continue{Token: &p.currToken}
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil,
			errors.New(
				"expected a semicolon (`;`) at the end of `continue` statement, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Break
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseBreak() (*ast.Break, error) {
	if !p.inLoop {
		return nil, errors.New("break statement can only be used inside a `for loop`")
	}
	stmt := &ast.Break{Token: &p.currToken}
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil,
			errors.New(
				"expected a semicolon (`;`) at the end of `break` statement, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	return stmt, nil
}
