package code

import "testing"

// Make : opcode와 operands를 받아서 명령어를 생성하는 함수
// Big Endian 방식으로 생성
func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}}, // 65534 -> 0xfffe
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at position %d. want=%d, got=%d",
					i, b, instruction[i])
			}
		}
	}
}
