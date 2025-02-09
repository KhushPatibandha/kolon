package vm

import (
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/compiler/code"
	"github.com/KhushPatibandha/Kolon/src/compiler/compiler"
	"github.com/KhushPatibandha/Kolon/src/interpreter/object"
)

const StackSize = 2048

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	stackPointer int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		stackPointer: 0,
	}
}

func (vm *VM) push(obj object.Object) error {
	if vm.stackPointer >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.stackPointer] = obj
	vm.stackPointer++
	return nil
}

func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.stackPointer-1]
	vm.stackPointer--
	return obj
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}
	return vm.stack[vm.stackPointer-1]
}

func (vm *VM) LastPoppedStackEle() object.Object {
	return vm.stack[vm.stackPointer]
}

func (vm *VM) Run() error {
	for i := 0; i < len(vm.instructions); i++ {
		op := code.Opcode(vm.instructions[i])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[i+1:])
			i += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod, code.OpOr, code.OpAnd:
			err := vm.execBinaryOp(op)
			if err != nil {
				return err
			}
		case code.OpEqualEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanEqual, code.OpAndAnd, code.OpOrOr:
			err := vm.execComparisonOp(op)
			if err != nil {
				return err
			}
		case code.OpNot:
			err := vm.execNotOp()
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.execMinusOp()
			if err != nil {
				return err
			}
		case code.OpMinusMinus, code.OpPlusPlus:
			err := vm.execPostfixOp(op)
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(TRUE)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(FALSE)
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) execBinaryOp(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	lType := left.Type()
	rType := right.Type()

	if lType == object.INTEGER_OBJ && rType == object.INTEGER_OBJ {
		return vm.execBinaryIntegerOp(op, left, right)
	} else if lType == object.FLOAT_OBJ && rType == object.FLOAT_OBJ {
		return vm.execBinaryFloatOp(op, left, right)
	} else if (lType == object.FLOAT_OBJ && rType == object.INTEGER_OBJ) || (lType == object.INTEGER_OBJ && rType == object.FLOAT_OBJ) {
		leftVal := 0.0
		rightVal := 0.0
		if lType == object.INTEGER_OBJ {
			leftVal = float64(left.(*object.Integer).Value)
			rightVal = right.(*object.Float).Value
		} else {
			leftVal = left.(*object.Float).Value
			rightVal = float64(right.(*object.Integer).Value)
		}
		return vm.execBinaryFloatOp(op, &object.Float{Value: leftVal}, &object.Float{Value: rightVal})
	} else if lType == object.STRING_OBJ && rType == object.STRING_OBJ {
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		leftVal = leftVal[1 : len(leftVal)-1]
		rightVal = rightVal[1 : len(rightVal)-1]
		return vm.execBinaryStringAndCharOp(op, leftVal, rightVal)
	} else if lType == object.CHAR_OBJ && rType == object.CHAR_OBJ {
		leftVal := left.(*object.Char).Value
		rightVal := right.(*object.Char).Value
		leftVal = leftVal[1 : len(leftVal)-1]
		rightVal = rightVal[1 : len(rightVal)-1]
		return vm.execBinaryStringAndCharOp(op, leftVal, rightVal)
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", lType, rType)
}

func (vm *VM) execComparisonOp(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	lType := left.Type()
	rType := right.Type()

	if lType == object.INTEGER_OBJ && rType == object.INTEGER_OBJ {
		return vm.execComparisonIntegerOp(op, left, right)
	} else if lType == object.FLOAT_OBJ && rType == object.FLOAT_OBJ {
		return vm.execComparisonFloatOp(op, left, right)
	} else if lType == object.BOOLEAN_OBJ && rType == object.BOOLEAN_OBJ {
		return vm.execComparisonBooleanOp(op, left, right)
	} else if (lType == object.FLOAT_OBJ && rType == object.INTEGER_OBJ) || (lType == object.INTEGER_OBJ && rType == object.FLOAT_OBJ) {
		leftVal := 0.0
		rightVal := 0.0
		if lType == object.INTEGER_OBJ {
			leftVal = float64(left.(*object.Integer).Value)
			rightVal = right.(*object.Float).Value
		} else {
			leftVal = left.(*object.Float).Value
			rightVal = float64(right.(*object.Integer).Value)
		}
		return vm.execComparisonFloatOp(op, &object.Float{Value: leftVal}, &object.Float{Value: rightVal})
	} else if lType == object.STRING_OBJ && rType == object.STRING_OBJ {
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		leftVal = leftVal[1 : len(leftVal)-1]
		rightVal = rightVal[1 : len(rightVal)-1]
		return vm.execComparisonStringAndCharOp(op, leftVal, rightVal)
	} else if lType == object.CHAR_OBJ && rType == object.CHAR_OBJ {
		leftVal := left.(*object.Char).Value
		rightVal := right.(*object.Char).Value
		leftVal = leftVal[1 : len(leftVal)-1]
		rightVal = rightVal[1 : len(rightVal)-1]
		return vm.execComparisonStringAndCharOp(op, leftVal, rightVal)
	}
	return fmt.Errorf("unsupported types for comparison operation: %s %s", lType, rType)
}

func (vm *VM) execNotOp() error {
	operand := vm.pop()
	if operand == TRUE {
		return vm.push(FALSE)
	} else {
		return vm.push(TRUE)
	}
}

func (vm *VM) execMinusOp() error {
	right := vm.pop()
	if right.Type() == object.FLOAT_OBJ {
		return vm.push(&object.Float{Value: -right.(*object.Float).Value})
	}
	return vm.push(&object.Integer{Value: -right.(*object.Integer).Value})
}

func (vm *VM) execPostfixOp(op code.Opcode) error {
	left := vm.pop()
	if left.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		if op == code.OpPlusPlus {
			return vm.push(&object.Integer{Value: leftVal + 1})
		} else {
			return vm.push(&object.Integer{Value: leftVal - 1})
		}
	} else {
		leftVal := left.(*object.Float).Value
		if op == code.OpPlusPlus {
			return vm.push(&object.Float{Value: leftVal + 1})
		} else {
			return vm.push(&object.Float{Value: leftVal - 1})
		}
	}
}

func (vm *VM) execBinaryIntegerOp(op code.Opcode, left object.Object, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	var res int64

	switch op {
	case code.OpAdd:
		res = leftVal + rightVal
	case code.OpSub:
		res = leftVal - rightVal
	case code.OpMul:
		res = leftVal * rightVal
	case code.OpDiv:
		res = leftVal / rightVal
	case code.OpMod:
		res = leftVal % rightVal
	case code.OpAnd:
		res = leftVal & rightVal
	case code.OpOr:
		res = leftVal | rightVal
	default:
		return fmt.Errorf("unknown integer operator %d", op)
	}
	return vm.push(&object.Integer{Value: res})
}

func (vm *VM) execBinaryFloatOp(op code.Opcode, left object.Object, right object.Object) error {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	var res float64

	switch op {
	case code.OpAdd:
		res = leftVal + rightVal
	case code.OpSub:
		res = leftVal - rightVal
	case code.OpMul:
		res = leftVal * rightVal
	case code.OpDiv:
		res = leftVal / rightVal
	default:
		return fmt.Errorf("unknown float operator %d", op)
	}
	return vm.push(&object.Float{Value: res})
}

func (vm *VM) execBinaryStringAndCharOp(op code.Opcode, leftVal string, rightVal string) error {
	switch op {
	case code.OpAdd:
		return vm.push(&object.String{Value: "\"" + leftVal + rightVal + "\""})
	default:
		return fmt.Errorf("unknown string operator %d", op)
	}
}

func (vm *VM) execComparisonIntegerOp(op code.Opcode, left object.Object, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case code.OpEqualEqual:
		if leftVal == rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpNotEqual:
		if leftVal != rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpGreaterThan:
		if leftVal > rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpGreaterThanEqual:
		if leftVal >= rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) execComparisonFloatOp(op code.Opcode, left object.Object, right object.Object) error {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch op {
	case code.OpEqualEqual:
		if leftVal == rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpNotEqual:
		if leftVal != rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpGreaterThan:
		if leftVal > rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpGreaterThanEqual:
		if leftVal >= rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) execComparisonBooleanOp(op code.Opcode, left object.Object, right object.Object) error {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch op {
	case code.OpEqualEqual:
		if leftVal == rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpNotEqual:
		if leftVal != rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpAndAnd:
		if leftVal && rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpOrOr:
		if leftVal || rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	}
	return fmt.Errorf("unknown operator: %d", op)
}

func (vm *VM) execComparisonStringAndCharOp(op code.Opcode, leftVal string, rightVal string) error {
	switch op {
	case code.OpEqualEqual:
		if leftVal == rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	case code.OpNotEqual:
		if leftVal != rightVal {
			return vm.push(TRUE)
		}
		return vm.push(FALSE)
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}
