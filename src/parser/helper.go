package parser

import (
	"errors"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/environment"
	ktype "github.com/KhushPatibandha/Kolon/src/kType"
	"github.com/KhushPatibandha/Kolon/src/lexer"
)

var (
	defaultInt = &ast.Integer{
		Token: lexer.Token{Kind: lexer.INT, Value: "0"},
		Value: 0,
		Type:  ktype.NewBaseType("int"),
	}
	defaultFloat = &ast.Float{
		Token: lexer.Token{Kind: lexer.FLOAT, Value: "0.0"},
		Value: 0.0,
		Type:  ktype.NewBaseType("float"),
	}
	defaultBool = &ast.Bool{
		Token: lexer.Token{Kind: lexer.BOOL, Value: "false"},
		Value: false,
		Type:  ktype.NewBaseType("bool"),
	}
	defaultString = &ast.String{
		Token: lexer.Token{Kind: lexer.STRING, Value: "\"\""},
		Value: "\"\"",
		Type:  ktype.NewBaseType("string"),
	}
	defaultChar = &ast.Char{
		Token: lexer.Token{Kind: lexer.CHAR, Value: "''"},
		Value: "''",
		Type:  ktype.NewBaseType("char"),
	}
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

func (p *Parser) assignDefaultValue(t *ktype.Type) ast.Expression {
	switch t.Kind {
	case ktype.TypeBase:
		switch t.Name {
		case "int":
			return defaultInt
		case "float":
			return defaultFloat
		case "bool":
			return defaultBool
		case "string":
			return defaultString
		case "char":
			return defaultChar
		default:
			return nil
		}
	default:
		return nil
	}
}

func (p *Parser) loadBuiltins() {
	builtins := []string{
		"print", "println", "scan", "scanln", "len", "toString", "toFloat", "toInt",
		"push", "pop", "insert", "remove", "getIndex", "keys", "values", "containsKey",
		"typeOf", "slice",
	}
	for _, name := range builtins {
		p.env.FuncNameSpace[name] = &environment.Symbol{
			IdentType: environment.FUNCTION,
			Ident: &ast.Identifier{
				Token: lexer.Token{Kind: lexer.IDENTIFIER, Value: name},
				Value: name,
			},
			Func: &environment.FuncInfo{
				Builtin:  true,
				Function: nil,
			},
			Type: nil,
			Env:  nil,
		}
	}
}

func (p *Parser) BootstrapFuncEnv(stmt *ast.Function) *environment.Environment {
	funcLocalEnv := environment.NewEnclosedEnvironment(p.env)
	for _, param := range stmt.Parameters {
		funcLocalEnv.Set(&environment.Symbol{
			IdentType: environment.VAR,
			Ident:     param.ParameterName,
			Type:      param.ParameterType,
			Func:      nil,
			Env:       nil,
		})
	}
	return funcLocalEnv
}

func typeCheckBoolCon(exp ast.Expression, keyword string) error {
	con := exp.GetType()
	if con.TypeLen != 1 {
		return errors.New(
			"`" + keyword + "` condition must evaluate to a single boolean value, got: " +
				strconv.Itoa(con.TypeLen) +
				". in case of call expression, it must return a single value",
		)
	}
	b, _ := typeCheckBool()
	if !con.Types[0].Equals(b.Types[0]) {
		return errors.New(
			"condition for `" + keyword + "` statement must always result in a boolean value, got: " +
				con.Types[0].String(),
		)
	}
	return nil
}

func checkReturnAtTheEnd(stmt []ast.Statement) error {
	lastStmt := stmt[len(stmt)-1]
	switch n := lastStmt.(type) {
	case *ast.Return:
		return nil
	case *ast.If:
		if n.Alternate == nil {
			return errors.New("` must have a `return` statement at the end of all branches")
		}
		err := checkReturnAtTheEnd(n.Body.Statements)
		if err != nil {
			return err
		}
		if n.MultiConditionals != nil {
			for _, mc := range n.MultiConditionals {
				err = checkReturnAtTheEnd(mc.Body.Statements)
				if err != nil {
					return err
				}
			}
		}
		err = checkReturnAtTheEnd(n.Alternate.Body.Statements)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("` must have a `return` statement at the end of all branches")
	}
}

func (p *Parser) handleEOF() (ast.Expression, error) {
	return nil, errors.New("unexpected end of file")
}
