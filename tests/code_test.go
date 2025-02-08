package tests

import (
	"testing"

	"github.com/KhushPatibandha/Kolon/src/compiler/code"
)

func Test48(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		operands []int
		expected []byte
	}{
		{code.OpConstant, []int{65534}, []byte{byte(code.OpConstant), 255, 254}},
		{code.OpAdd, []int{}, []byte{byte(code.OpAdd)}},
	}

	for _, tt := range tests {
		instructions := code.Make(tt.op, tt.operands...)
		if len(instructions) != len(tt.expected) {
			t.Errorf("instruction has worng length. want=%d, got=%d", len(tt.expected), len(instructions))
		}
		for i, b := range tt.expected {
			if instructions[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instructions[i])
			}
		}
	}
}

func Test50(t *testing.T) {
	instructions := []code.Instructions{
		code.Make(code.OpAdd),
		code.Make(code.OpConstant, 2),
		code.Make(code.OpConstant, 65535),
	}
	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

	concatted := code.Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}
	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q", expected, concatted.String())
	}
}

func Test51(t *testing.T) {
	tests := []struct {
		op        code.Opcode
		operands  []int
		bytesRead int
	}{
		{code.OpConstant, []int{65535}, 2},
	}
	for _, tt := range tests {
		instruction := code.Make(tt.op, tt.operands...)
		def, err := code.Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("defination not found: %q\n", err)
		}
		operandsRead, n := code.ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%d, go=%d", tt.bytesRead, n)
		}
		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
			}
		}
	}
}
