package compiler

import (
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/compiler/code"
	"github.com/KhushPatibandha/Kolon/src/interpreter/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
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
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
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
			return fmt.Errorf("unknown operator %s", node.Operator)
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
	return pos
}

func (c *Compiler) addIns(ins []byte) int {
	newInsPos := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return newInsPos
}
