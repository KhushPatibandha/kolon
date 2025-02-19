package vm

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/compiler/code"
	"github.com/KhushPatibandha/Kolon/src/compiler/compiler"
	"github.com/KhushPatibandha/Kolon/src/object"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
)

var (
	TRUE      = &object.Boolean{Value: true}
	FALSE     = &object.Boolean{Value: false}
	lastPoped object.Object
)

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	stackPointer int
	globals      []object.Object
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		stackPointer: 0,
		globals:      make([]object.Object, GlobalsSize),
	}
}

func (vm *VM) push(obj object.Object) error {
	if vm.stackPointer >= StackSize {
		return errors.New("stack overflow")
	}
	vm.stack[vm.stackPointer] = obj
	vm.stackPointer++
	return nil
}

func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.stackPointer-1]
	vm.stackPointer--
	vm.stack[vm.stackPointer] = nil
	lastPoped = obj
	return obj
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}
	return vm.stack[vm.stackPointer-1]
}

func (vm *VM) LastPoppedStackEle() object.Object {
	return lastPoped
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
		case code.OpNot, code.OpMinus:
			err := vm.execPrefixOp(op)
			if err != nil {
				return err
			}
		case code.OpMinusMinus, code.OpPlusPlus:
			err := vm.execPostfixOp(op)
			if err != nil {
				return err
			}
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[i+1:]))
			i = pos - 1
		case code.OpJumpNotTrue:
			pos := int(code.ReadUint16(vm.instructions[i+1:]))
			i += 2
			condition := vm.pop()
			if condition == FALSE {
				i = pos - 1
			}
		case code.OpSetGlobal:
			globalIdx := code.ReadUint16(vm.instructions[i+1:])
			i += 2
			vm.globals[globalIdx] = vm.pop()
		case code.OpGetGlobal:
			globalIdx := code.ReadUint16(vm.instructions[i+1:])
			i += 2
			err := vm.push(vm.globals[globalIdx])
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
		case code.OpArray:
			totalEle := int(code.ReadUint16(vm.instructions[i+1:]))
			i += 2
			array := vm.buildArray(vm.stackPointer-totalEle, vm.stackPointer)
			vm.stackPointer -= totalEle
			err := vm.push(array)
			if err != nil {
				return err
			}
		case code.OpHash:
			totalEle := int(code.ReadUint16(vm.instructions[i+1:]))
			i += 2

			hash, err := vm.buildHash(vm.stackPointer-totalEle, vm.stackPointer)
			if err != nil {
				return err
			}
			vm.stackPointer -= totalEle
			err = vm.push(hash)
			if err != nil {
				return err
			}
		case code.OpIndex:
			idx := vm.pop()
			left := vm.pop()
			err := vm.execIndexOp(left, idx)
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
	return errors.New("invalid binary `infix` operation with variable types on left and right, got: `" + string(lType) + "` and `" + string(rType) + "`")
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
	return errors.New("invalid comparison `infix` operation with variable types on left and right, got: `" + string(lType) + "` and `" + string(rType) + "`")
}

func (vm *VM) execPrefixOp(op code.Opcode) error {
	if op == code.OpNot {
		return vm.execNotOp()
	} else {
		return vm.execMinusOp()
	}
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

func (vm *VM) execIndexOp(left object.Object, idx object.Object) error {
	if left.Type() == object.ARRAY_OBJ && idx.Type() == object.INTEGER_OBJ {
		return vm.execArrayIdxOp(left, idx)
	} else {
		return vm.execHashMapIdxOp(left, idx)
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
		return fmt.Errorf("can only use `+`, `-`, `*`, `/`, `%%`, `>`, `<`, `<=`, `>=`, `!=`, `==`, `|`, `&` infix operators with 2 `int`, got: %d", op)
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
		return fmt.Errorf("can only use `+`, `-`, `*`, `/`, `>`, `<`, `<=`, `>=`, `!=`, `==` infix operators with 2 `float`, got: %d", op)
	}
	return vm.push(&object.Float{Value: res})
}

func (vm *VM) execBinaryStringAndCharOp(op code.Opcode, leftVal string, rightVal string) error {
	switch op {
	case code.OpAdd:
		return vm.push(&object.String{Value: "\"" + leftVal + rightVal + "\""})
	default:
		if leftVal[0] == '\'' {
			return fmt.Errorf("can only use `+`, `==`, `!=` infix operators with 2 `char`, got: %d", op)
		} else {
			return fmt.Errorf("can only use `+`, `==`, `!=` infix operators with 2 `string`, got: %d", op)
		}
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
		return fmt.Errorf("can only use `+`, `-`, `*`, `/`, `%%`, `>`, `<`, `<=`, `>=`, `!=`, `==`, `|`, `&` infix operators with 2 `int`, got: %d", op)
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
		return fmt.Errorf("can only use `+`, `-`, `*`, `/`, `>`, `<`, `<=`, `>=`, `!=`, `==` infix operators with 2 `float`, got: %d", op)
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
	return fmt.Errorf("can only use `==`, `!=`, `&&`, `||` infix operators with 2 `bool`, got: %d", op)
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
		if leftVal[0] == '\'' {
			return fmt.Errorf("can only use `+`, `==`, `!=` infix operators with 2 `char`, got: %d", op)
		} else {
			return fmt.Errorf("can only use `+`, `==`, `!=` infix operators with 2 `string`, got: %d", op)
		}
	}
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

func (vm *VM) execArrayIdxOp(left object.Object, index object.Object) error {
	arrayObj := left.(*object.Array)
	idx := index.(*object.Integer).Value
	maxIdx := int64(len(arrayObj.Elements) - 1)
	if idx < 0 || idx > maxIdx {
		return errors.New("index out of range, index: " + strconv.FormatInt(idx, 10) + ", max index: " + strconv.FormatInt(maxIdx, 10) + ", min index: 0")
	}
	return vm.push(arrayObj.Elements[idx])
}

func (vm *VM) execHashMapIdxOp(left object.Object, index object.Object) error {
	hashObj := left.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return errors.New("unusable as hash key: " + string(index.Type()))
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return errors.New("key not found: " + index.Inspect())
	}
	return vm.push(pair.Value)
}

func (vm *VM) buildArray(sIdx int, eIdx int) object.Object {
	elements := make([]object.Object, eIdx-sIdx)
	for i := sIdx; i < eIdx; i++ {
		elements[i-sIdx] = vm.stack[i]
	}
	return &object.Array{Elements: elements}
}

func (vm *VM) buildHash(sIdx int, eIdx int) (object.Object, error) {
	pairs := make(map[object.HashKey]object.HashPair)

	for i := sIdx; i < eIdx; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, errors.New("unusable as hash key: " + string(key.Type()))
		}
		pairs[hashKey.HashKey()] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}, nil
}
