package parser

import (
	"errors"
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/environment"
	ktype "github.com/KhushPatibandha/Kolon/src/kType"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

// ------------------------------------------------------------------------------------------------------------------
// Statements
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currToken.Kind {
	case lexer.VAR, lexer.CONST:
		return p.parseVarConst()
	case lexer.RETURN:
		return p.parseReturn()
	case lexer.FUN:
		stmt, err := p.parseFunction()
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			return nil, nil
		}
		return stmt, nil
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
	body := &ast.Body{Token: p.currToken, Statements: []ast.Statement{}}
	p.nextToken()
	for !p.currTokenIsOk(lexer.CLOSE_CURLY_BRACKET) {
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
func (p *Parser) parseType() (*ktype.Type, error) {
	if !p.expectedPeekToken(lexer.TYPE) {
		return nil,
			errors.New(
				"expected a type, got: " + lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	stmt := ktype.NewBaseType(p.currToken.Value)

	for p.peekTokenIsOk(lexer.OPEN_SQUARE_BRACKET) {
		p.nextToken()

		if p.peekTokenIsOk(lexer.CLOSE_SQUARE_BRACKET) {
			stmt = ktype.NewArrayType(stmt)
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
		stmt = ktype.NewHashMapType(stmt, val)
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
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	if p.peekTokenIsOk(lexer.COMMA) {
		stmt.Expression = &ast.Assignment{
			Token:    lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
			Left:     &ast.Identifier{Token: p.currToken, Value: p.currToken.Value},
			Operator: "=",
		}
		list := []ast.Statement{}
		list = append(list, stmt)
		p.nextToken()
		p.nextToken()
		return p.parseMultiAssign(list)
	}

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	switch t := exp.(type) {
	case *ast.CallExpression:
		stmt.Expression = t
	case *ast.Postfix:
		stmt.Expression = t
	case *ast.Assignment:
		stmt.Expression = t
	default:
		return nil,
			errors.New(
				"expected a function call, postfix expression or an assignment " +
					"expression for expressions as statements, got: " +
					fmt.Sprintf("%T", exp),
			)
	}

	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		if !p.inTesting {
			return nil,
				errors.New(
					"expected a semicolon (`;`) at the end of the statement, got: " +
						lexer.TokenKindString(p.peekToken.Kind),
				)
		}
	}
	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// VarAndConst
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseVarConst() (ast.Statement, error) {
	if !p.inFunction && !p.inTesting {
		return nil, errors.New("can't declare a variable outside a function")
	}
	stmt, err := p.parseVarConstSig()
	if err != nil {
		return nil, err
	}

	if p.peekTokenIsOk(lexer.COMMA) {
		list := []ast.Statement{}
		list = append(list, stmt)
		p.nextToken()
		p.nextToken()
		return p.parseMultiAssign(list)
	}

	var value ast.Expression
	if p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
		p.nextToken()
		p.nextToken()
		value, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
	} else {
		value = p.assignDefaultValue(stmt.Type)
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

	err = typeCheckVarAndConst(stmt, p.stack.Top())
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseVarConstSig() (*ast.VarAndConst, error) {
	stmt := &ast.VarAndConst{Token: p.currToken}
	if !p.expectedPeekToken(lexer.IDENTIFIER) {
		return nil,
			errors.New(
				"expected an identifier after `" +
					p.currToken.Value + "` keyword, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Value}

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
	stmt.Name.Type = typ

	return stmt, nil
}

// ------------------------------------------------------------------------------------------------------------------
// Multi-Assignment
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseMultiAssign(list []ast.Statement) (*ast.MultiAssignment, error) {
	if !p.inFunction && !p.inTesting {
		return nil, errors.New("can't declare a variable outside a function")
	}
	stmt := &ast.MultiAssignment{Token: lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="}}

	if p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
		return nil,
			errors.New(
				"expected an identifier or a `var`/`const` keyword after the comma (`,`), got: " +
					lexer.TokenKindString(p.currToken.Kind),
			)
	}

	for !p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
		ele, err := p.parseMultiAssignLHS()
		if err != nil {
			return nil, err
		}
		list = append(list, ele)

		if p.currTokenIsOk(lexer.COMMA) {
			if p.peekTokenIsOk(lexer.EQUAL_ASSIGN) {
				return nil,
					errors.New(
						"expected an identifier or a `var`/`const` keyword after the comma (`,`), got: " +
							lexer.TokenKindString(p.peekToken.Kind),
					)
			}
			p.nextToken()
		} else if p.currTokenIsOk(lexer.EQUAL_ASSIGN) {
			break
		} else {
			return nil,
				errors.New(
					"expected a comma (`,`) or an equal sign (`=`), got: " +
						lexer.TokenKindString(p.currToken.Kind),
				)
		}
	}
	p.nextToken()

	values, err := p.parseMultiAssignRHS()
	if err != nil {
		return nil, err
	}
	if len(values) == 1 {
		stmt.SingleFunctionCall = true
	}

	if err := p.matchAssignments(list, values); err != nil {
		return nil, err
	}
	stmt.Objects = list

	err = typeCheckMultiAssign(stmt, p.stack.Top())
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseMultiAssignLHS() (ast.Statement, error) {
	switch {
	case p.currTokenIsOk(lexer.VAR), p.currTokenIsOk(lexer.CONST):
		ele, err := p.parseVarConstSig()
		if err != nil {
			return nil, err
		}
		p.nextToken()
		return ele, nil
	case p.currTokenIsOk(lexer.IDENTIFIER):
		ele := &ast.ExpressionStatement{
			Token: p.currToken,
			Expression: &ast.Assignment{
				Token:    lexer.Token{Kind: lexer.EQUAL_ASSIGN, Value: "="},
				Left:     &ast.Identifier{Token: p.currToken, Value: p.currToken.Value},
				Operator: "=",
			},
		}
		p.nextToken()
		return ele, nil
	default:
		return nil,
			errors.New(
				"expected an identifier or a `var`/`const` keyword, got: " +
					lexer.TokenKindString(p.currToken.Kind),
			)
	}
}

func (p *Parser) parseMultiAssignRHS() ([]ast.Expression, error) {
	values := []ast.Expression{}

	for {
		exp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		values = append(values, exp)
		p.nextToken()

		if p.currTokenIsOk(lexer.COMMA) {
			if p.peekTokenIsOk(lexer.SEMI_COLON) {
				return nil,
					errors.New(
						"expected an expression after comma (`,`) in " +
							"multi-value assignment statement, got: " +
							lexer.TokenKindString(p.peekToken.Kind))
			}
			p.nextToken()
		} else if p.currTokenIsOk(lexer.SEMI_COLON) {
			break
		} else {
			return nil,
				errors.New(
					"expected a comma (`,`) or a semicolon (`;`) after the expression, got: " +
						lexer.TokenKindString(p.currToken.Kind))
		}
	}
	return values, nil
}

func (p *Parser) matchAssignments(left []ast.Statement,
	right []ast.Expression,
) error {
	if len(right) == 1 {
		for _, ele := range left {
			switch t := ele.(type) {
			case *ast.VarAndConst:
				t.Value = right[0]
			case *ast.ExpressionStatement:
				if exp, ok := t.Expression.(*ast.Assignment); ok {
					exp.Right = right[0]
				}
			}
		}
		return nil
	}

	if len(left) != len(right) {
		return errors.New(
			"number of expression on the right side of multi-assignment do not " +
				"match the number of declarations on the left, left: " +
				fmt.Sprint(len(left)) + ", right: " + fmt.Sprint(len(right)),
		)
	}

	for i, ele := range left {
		switch t := ele.(type) {
		case *ast.VarAndConst:
			t.Value = right[i]
		case *ast.ExpressionStatement:
			if exp, ok := t.Expression.(*ast.Assignment); ok {
				exp.Right = right[i]
			}
		}
	}

	return nil
}

// ------------------------------------------------------------------------------------------------------------------
// Function
// ------------------------------------------------------------------------------------------------------------------
func (p *Parser) parseFunction() (*ast.Function, error) {
	if p.inFunction {
		return nil, errors.New("can't declare a function inside a function")
	}
	stmt := &ast.Function{Token: p.currToken, Parameters: nil, ReturnTypes: nil, Body: nil}

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
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Value}

	if existing, ok := p.env.GetFunc(stmt.Name.Value); ok &&
		(existing.Func.Builtin || existing.Func.Function.Body != nil) && !p.inTesting {
		return nil,
			errors.New(
				"can't declare a function twice, function with the same name `" +
					stmt.Name.Value + "` already exists",
			)
	}
	if existing, ok := p.env.GetFunc(stmt.Name.Value); ok &&
		existing.Func.Builtin && !p.inTesting {
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
	}

	if p.peekTokenIsOk(lexer.SEMI_COLON) {
		p.nextToken()
		if _, ok := p.env.GetFunc(stmt.Name.Value); ok {
			return nil,
				errors.New(
					"can't declare a function twice, function with the same name. `" +
						stmt.Name.Value + "` already exists",
				)
		}

		funcLocalEnv := p.BootstrapFuncEnv(stmt)
		p.env.Set(&environment.Symbol{
			IdentType: environment.FUNCTION,
			Ident:     stmt.Name,
			Func: &environment.FuncInfo{
				Function: stmt,
				Builtin:  false,
			},
			Env:  funcLocalEnv,
			Type: nil,
		})

		return stmt, nil
	}

	p.currFunction = stmt

	if !p.expectedPeekToken(lexer.OPEN_CURLY_BRACKET) {
		return nil,
			errors.New(
				"expected an open curly bracket (`{`) after the function (`" +
					stmt.Name.Value + "`) signature, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}

	existing, ok := p.env.GetFunc(stmt.Name.Value)
	if ok && !p.compareFunctionSig(existing.Func.Function, stmt) {
		return nil,
			errors.New(
				"function signature of " + stmt.Name.Value +
					" doesn't match with previously declared signature",
			)
	}

	if !ok {
		funcLocalEnv := p.BootstrapFuncEnv(stmt)
		p.env.Set(&environment.Symbol{
			IdentType: environment.FUNCTION,
			Ident:     stmt.Name,
			Func: &environment.FuncInfo{
				Function: stmt,
				Builtin:  false,
			},
			Type: nil,
			Env:  funcLocalEnv,
		})
	}
	f, _ := p.env.GetFunc(stmt.Name.Value)

	p.inFunction = true
	p.stack.Push(f.Env)
	funBody, err := p.parseBody()
	if err != nil {
		return nil, err
	}
	p.stack.Pop()
	p.inFunction = false
	p.currFunction = nil

	f.Func.Function.Body = funBody

	if !p.inTesting {
		err = typeCheckFunction(f.Func.Function)
		if err != nil {
			return nil, err
		}
	}

	if ok {
		return nil, nil
	}

	return f.Func.Function, nil
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
			ParameterName: &ast.Identifier{Token: p.currToken, Value: p.currToken.Value},
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
func (p *Parser) parseFunctionReturnTypes() ([]*ktype.Type, error) {
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

	var returnTypes []*ktype.Type

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
	stmt := &ast.If{Token: p.currToken}
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

	err = typeCheckBoolCon(stmt.Condition, "if")
	if err != nil {
		return nil, err
	}

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

	ifLocalEnv := environment.NewEnclosedEnvironment(p.stack.Top())
	p.stack.Push(ifLocalEnv)

	body, err := p.parseBody()
	if err != nil {
		return nil, err
	}

	p.stack.Pop()
	stmt.Body = body

	// ------------------------------------------------------------------------------------------------------------------
	// Else If
	// ------------------------------------------------------------------------------------------------------------------
	if p.peekTokenIsOk(lexer.ELSE_IF) {
		var elseIfList []*ast.ElseIf

		for p.peekTokenIsOk(lexer.ELSE_IF) {
			p.nextToken()
			elseIfStmt := &ast.ElseIf{Token: p.currToken}
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
			err = typeCheckBoolCon(elseIfStmt.Condition, "else if")
			if err != nil {
				return nil, err
			}

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

			ifElseLocalEnv := environment.NewEnclosedEnvironment(p.stack.Top())
			p.stack.Push(ifElseLocalEnv)

			elseIfBody, err := p.parseBody()
			if err != nil {
				return nil, err
			}

			p.stack.Pop()
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
		elseStmt := &ast.Else{Token: p.currToken}

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

		elseLocalEnv := environment.NewEnclosedEnvironment(p.stack.Top())
		p.stack.Push(elseLocalEnv)

		elseBody, err := p.parseBody()
		if err != nil {
			return nil, err
		}

		p.stack.Pop()

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
	stmt := &ast.Return{Token: p.currToken, Value: []ast.Expression{}}
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
						". to return multiple values, use `return: (val1, val2, ...)`",
				)
		}
	}

	err := typeCheckReturn(stmt, p.currFunction)
	if err != nil {
		return nil, err
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
	stmt := &ast.ForLoop{Token: p.currToken}

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

	forLoopLocalEnv := environment.NewEnclosedEnvironment(p.stack.Top())
	p.stack.Push(forLoopLocalEnv)

	p.nextToken()
	left, err := p.parseStatement()
	if err != nil {
		return nil, err
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

	err = typeCheckForLoop(stmt, forLoopLocalEnv)
	if err != nil {
		return nil, err
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

	p.stack.Pop()
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
	stmt := &ast.WhileLoop{Token: p.currToken}

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

	whileLoopLocalEnv := environment.NewEnclosedEnvironment(p.stack.Top())
	p.stack.Push(whileLoopLocalEnv)

	condition, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Condition = condition

	err = typeCheckBoolCon(stmt.Condition, "while")
	if err != nil {
		return nil, err
	}

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

	p.stack.Pop()
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
	stmt := &ast.Continue{Token: p.currToken}
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
	stmt := &ast.Break{Token: p.currToken}
	if !p.expectedPeekToken(lexer.SEMI_COLON) {
		return nil,
			errors.New(
				"expected a semicolon (`;`) at the end of `break` statement, got: " +
					lexer.TokenKindString(p.peekToken.Kind),
			)
	}
	return stmt, nil
}
