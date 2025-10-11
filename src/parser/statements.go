package parser

import (
	"errors"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Statements
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currToken.Kind {
	case lexer.VAR:
		return nil, nil
	case lexer.CONST:
		return nil, nil
	case lexer.RETURN:
		return nil, nil
	case lexer.FUN:
		return p.parseFunction()
	case lexer.IF:
		return p.parseIf()
	case lexer.FOR:
		return nil, nil
	case lexer.WHILE:
		return nil, nil
	case lexer.CONTINUE:
		return p.parseContinue()
	case lexer.BREAK:
		return p.parseBreak()
	default:
		return nil, nil
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
		Kind:  ast.TypeBase,
		Token: &p.currToken,
		Name:  p.currToken.Value,
	}

	for p.peekTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		p.nextToken()

		if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			stmt = &ast.Type{
				Kind:        ast.TypeArray,
				Token:       stmt.Token,
				ElementType: stmt,
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
			Kind:      ast.TypeHashMap,
			Token:     stmt.Token,
			KeyType:   stmt,
			ValueType: val,
		}
	}
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
