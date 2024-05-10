package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota // constant 풀에서 index 위치에 있는 constant를 가져오는 명령어
)

// Definition : For debugging purposes
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // 첫번째 인자가 2바이트의 크기를 가짐
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
