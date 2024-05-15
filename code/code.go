package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

func (ins Instructions) fmtInstruction(def *Definition, operands []int) any {
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

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

type Opcode byte

const (
	OpConstant Opcode = iota // constant 풀에서 index 위치에 있는 constant를 가져오는 명령어
	OpAdd                    // 스택에서 상위 2개의 값을 꺼내서 더한 결과를 다시 스택에 넣는 명령어
	OpPop                    // 스택에서 값을 꺼내는 명령어
	OpSub                    // 스택에서 상위 2개의 값을 꺼내서 뺀 결과를 다시 스택에 넣는 명령어
	OpMul                    // 스택에서 상위 2개의 값을 꺼내서 곱한 결과를 다시 스택에 넣는 명령어
	OpDiv                    // 스택에서 상위 2개의 값을 꺼내서 나눈 결과를 다시 스택에 넣는 명령어
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
)

// Definition : For debugging purposes
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}}, // 첫번째 인자가 2바이트의 크기를 가짐
	OpAdd:           {"OpAdd", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil

}

// Make : opcode와 operands를 받아서 명령어를 생성하는 함수
// Big Endian 방식으로 생성
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// 각 operand의 크기를 더해서 최종 instruction의 크기를 구함
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// instruction을 저장할 byte slice 생성
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	// operand를 byte slice에 저장
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			// Big Endian 방식으로 저장
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}

// ReadOperands : Make가 생성한 instruction을 다시 decode
// 복호화된 피연산자와 이를 읽은 바이트 수를 반환
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
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
