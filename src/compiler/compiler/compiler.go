package compiler

import (
	"errors"
	"sort"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/compiler/code"
	"github.com/KhushPatibandha/Kolon/src/object"
)

type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
	InTesting           bool
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	OpCode   code.Opcode
	Position int
}

func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymTable(),
		InTesting:           false,
	}
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.FunctionBody:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.IntegerValue:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConst(integer))
	case *ast.FloatValue:
		flaot := &object.Float{Value: node.Value}
		c.emit(code.OpConstant, c.addConst(flaot))
	case *ast.BooleanValue:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.StringValue:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConst(str))
	case *ast.CharValue:
		char := &object.Char{Value: node.Value}
		c.emit(code.OpConstant, c.addConst(char))
	case *ast.ArrayValue:
		for _, el := range node.Values {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Values))
	case *ast.HashMap:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		if c.InTesting {
			sort.Slice(keys, func(i int, j int) bool {
				return keys[i].String() < keys[j].String()
			})
		}
		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}
		c.emit(code.OpHash, len(node.Pairs)*2)
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return errors.New("undefined variable: " + node.Value)
		}
		c.emit(code.OpGetGlobal, symbol.Index)
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		if _, ok := node.Expression.(*ast.AssignmentExpression); !ok {
			c.emit(code.OpPop)
		}
	case *ast.InfixExpression:
		if node.Operator == "<" || node.Operator == "<=" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			if node.Operator == "<" {
				c.emit(code.OpGreaterThan)
			} else {
				c.emit(code.OpGreaterThanEqual)
			}
			return nil
		}
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "%":
			c.emit(code.OpMod)
		case "|":
			c.emit(code.OpOr)
		case "&":
			c.emit(code.OpAnd)
		case ">":
			c.emit(code.OpGreaterThan)
		case ">=":
			c.emit(code.OpGreaterThanEqual)
		case "==":
			c.emit(code.OpEqualEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case "&&":
			c.emit(code.OpAndAnd)
		case "||":
			c.emit(code.OpOrOr)
		default:
			return errors.New("unknown operator for infix operation. can only use `+`, `-`, `*`, `/`, `%`, `>`, `<`, `<=`, `>=`, `!=`, `==`, `|`, `&`, `&&`, `||` infix operators, got: " + node.Operator)
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpNot)
		case "-":
			c.emit(code.OpMinus)
		default:
			return errors.New("unknown operator for prefix operation. can only use `!`, `-` prefix operators, got: " + node.Operator)
		}
	case *ast.PostfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "++":
			c.emit(code.OpPlusPlus)
		case "--":
			c.emit(code.OpMinusMinus)
		default:
			return errors.New("unknown operator for postfix operation. can only use `++`, `--` postfix operators, got: " + node.Operator)
		}
	case *ast.AssignmentExpression:
		symbol, ok := c.symbolTable.Resolve(node.Left.Value)
		if !ok {
			return errors.New("undefined variable: " + node.Left.Value)
		}

		if node.Operator == "=" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
		} else {
			c.emit(code.OpGetGlobal, symbol.Index)

			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			switch node.Operator {
			case "+=":
				c.emit(code.OpAdd)
			case "-=":
				c.emit(code.OpSub)
			case "/=":
				c.emit(code.OpDiv)
			case "*=":
				c.emit(code.OpMul)
			case "%=":
				c.emit(code.OpMod)
			default:
				return errors.New("unknown operator for assignment operation. can only use `=`, `+=`, `-=`, `*=`, `/=`, `%=`, got: " + node.Operator)
			}
		}
		c.emit(code.OpSetGlobal, symbol.Index)
	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)
	case *ast.VarStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)
	case *ast.IfStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		jumpNotTruePos := c.emit(code.OpJumpNotTrue, 9999)
		err = c.Compile(node.Body)
		if err != nil {
			return err
		}

		var allJumpPos []int

		if node.Consequence == nil && node.MultiConseq == nil {
			afterIfBodyPos := len(c.instructions)
			c.changeOperand(jumpNotTruePos, afterIfBodyPos)
		} else if node.MultiConseq != nil {
			jumpPos := c.emit(code.OpJump, 9999)
			allJumpPos = append(allJumpPos, jumpPos)
			afterIfBodyPos := len(c.instructions)
			c.changeOperand(jumpNotTruePos, afterIfBodyPos)

			for _, elseIf := range node.MultiConseq {
				err := c.Compile(elseIf.Value)
				if err != nil {
					return err
				}
				jumpPosNotTrueForElseIf := c.emit(code.OpJumpNotTrue, 9999)
				err = c.Compile(elseIf.Body)
				if err != nil {
					return err
				}
				jumpPos := c.emit(code.OpJump, 9999)
				allJumpPos = append(allJumpPos, jumpPos)
				afterElseIfBodyPos := len(c.instructions)
				c.changeOperand(jumpPosNotTrueForElseIf, afterElseIfBodyPos)
			}

		}
		if node.Consequence != nil {
			if node.MultiConseq == nil {
				jumpPos := c.emit(code.OpJump, 9999)
				allJumpPos = append(allJumpPos, jumpPos)
				afterIfBodyPos := len(c.instructions)
				c.changeOperand(jumpNotTruePos, afterIfBodyPos)
			}

			err := c.Compile(node.Consequence.Body)
			if err != nil {
				return err
			}
		}
		endPos := len(c.instructions)
		for _, jumpPos := range allJumpPos {
			c.changeOperand(jumpPos, endPos)
		}
	}
	return nil
}

func (c *Compiler) addConst(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(opcode code.Opcode, operands ...int) int {
	ins := code.Make(opcode, operands...)
	pos := c.addIns(ins)
	c.setLastIns(opcode, pos)
	return pos
}

func (c *Compiler) addIns(ins []byte) int {
	newInsPos := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return newInsPos
}

func (c *Compiler) setLastIns(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{OpCode: op, Position: pos}
	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) replaceIns(pos int, newIns []byte) {
	for i := 0; i < len(newIns); i++ {
		c.instructions[pos+i] = newIns[i]
	}
}

func (c *Compiler) changeOperand(pos int, operand int) {
	op := code.Opcode(c.instructions[pos])
	newIns := code.Make(op, operand)
	c.replaceIns(pos, newIns)
}
