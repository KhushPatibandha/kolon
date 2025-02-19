package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type (
	Instructions []byte
	Opcode       byte
)

const (
	OpConstant Opcode = iota
	OpPop
	OpTrue
	OpFalse
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpAnd
	OpOr
	OpEqualEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterThanEqual
	OpAndAnd
	OpOrOr
	OpMinus
	OpNot
	OpPlusPlus
	OpMinusMinus
	OpJumpNotTrue
	OpJump
	OpGetGlobal
	OpSetGlobal
	OpArray
	OpHash
	OpIndex
)

type Defination struct {
	Name          string
	OperandWidths []int
}

var definations = map[Opcode]*Defination{
	OpConstant:         {"OpConstant", []int{2}},
	OpPop:              {"OpPop", []int{}},
	OpTrue:             {"OpTrue", []int{}},
	OpFalse:            {"OpFalse", []int{}},
	OpAdd:              {"OpAdd", []int{}},
	OpSub:              {"OpSub", []int{}},
	OpMul:              {"OpMul", []int{}},
	OpDiv:              {"OpDiv", []int{}},
	OpMod:              {"OpMod", []int{}},
	OpAnd:              {"OpAnd", []int{}},
	OpOr:               {"OpOr", []int{}},
	OpEqualEqual:       {"OpEqualEqual", []int{}},
	OpNotEqual:         {"OpNotEqual", []int{}},
	OpGreaterThan:      {"OpGreaterThan", []int{}},
	OpGreaterThanEqual: {"OpGreaterThanEqual", []int{}},
	OpAndAnd:           {"OpAndAnd", []int{}},
	OpOrOr:             {"OpOrOr", []int{}},
	OpMinus:            {"OpMinus", []int{}},
	OpNot:              {"OpNot", []int{}},
	OpPlusPlus:         {"OpPlusPlus", []int{}},
	OpMinusMinus:       {"OpMinusMinus", []int{}},
	OpJumpNotTrue:      {"OpJumpNotTrue", []int{2}},
	OpJump:             {"OpJump", []int{2}},
	OpGetGlobal:        {"OpGetGlobal", []int{2}},
	OpSetGlobal:        {"OpSetGlobal", []int{2}},
	OpArray:            {"OpArray", []int{2}},
	OpHash:             {"OpHash", []int{2}},
	OpIndex:            {"OpIndex", []int{}},
}

func Lookup(op byte) (*Defination, error) {
	def, ok := definations[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definations[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}

func ReadOperands(def *Defination, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0
	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "Error: %s\n", err)
			continue
		}
		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

func (ins Instructions) fmtInstruction(def *Defination, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}
	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("Error: unhandled operandCount for %s\n", def.Name)
}
