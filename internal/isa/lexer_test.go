package isa

import (
	"reflect"
	"testing"
)

func TestLexer_Tokenize(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:     "add reg, imm",
			input:    "add r6, 1",
			expected: []Token{{T: TK_ALU, Tk: "ADD"}, {T: TK_REG, Tk: "R6"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "1", ValInt: 1}},
		},
		{
			name:     "ld.b reg, [reg]",
			input:    "ld.b r11, [r12]",
			expected: []Token{{T: TK_LD_BYTE, Tk: "LD.B"}, {T: TK_REG, Tk: "R11"}, {T: TK_COMMA, Tk: ","}, {T: TK_L_SQBR, Tk: "["}, {T: TK_REG, Tk: "R12"}, {T: TK_R_SQBR, Tk: "]"}},
		},
		{
			name:     "pop lr",
			input:    "pop lr",
			expected: []Token{{T: TK_POP, Tk: "POP"}, {T: TK_REG_LR, Tk: "LR"}},
		},
		{
			name:     "push reg",
			input:    "push r11",
			expected: []Token{{T: TK_PUSH, Tk: "PUSH"}, {T: TK_REG, Tk: "R11"}},
		},
		{
			name:     "add reg reg",
			input:    "\tADD R1, R2",
			expected: []Token{{T: TK_ALU, Tk: "ADD"}, {T: TK_REG, Tk: "R1"}, {T: TK_COMMA, Tk: ","}, {T: TK_REG, Tk: "R2"}},
		},
		{
			name:     "ld reg reg",
			input:    "\tLD R9, R0",
			expected: []Token{{T: TK_LD, Tk: "LD"}, {T: TK_REG, Tk: "R9"}, {T: TK_COMMA, Tk: ","}, {T: TK_REG, Tk: "R0"}},
		},
		{
			name:     "ld reg imm dec",
			input:    "\tLD R3, 255",
			expected: []Token{{T: TK_LD, Tk: "LD"}, {T: TK_REG, Tk: "R3"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "255", ValInt: 255}},
		},
		{
			name:     "ld reg imm hex",
			input:    "\tLD R3, 0xFF",
			expected: []Token{{T: TK_LD, Tk: "LD"}, {T: TK_REG, Tk: "R3"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "0XFF", ValInt: 255}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := Tokenize(tc.input)
			if err != nil {
				t.Errorf("Tokenize(%q) error: %v", tc.input, err)
			}
			if !reflect.DeepEqual(tokens, tc.expected) {
				t.Errorf("Tokenize(%q) = %v; want %v", tc.input, tokens, tc.expected)
			}
		})
	}
}
