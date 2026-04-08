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
			name:     "define string reference",
			input:    "MystrR db @Mystr",
			expected: []Token{{T: TK_LABEL, Tk: "MYSTRR"}, {T: TK_DB, Tk: "DB"}, {T: TK_AT, Tk: "@"}, {T: TK_LABEL, Tk: "MYSTR"}},
		},
		{
			name:     "define string",
			input:    "Mystr db 'Hello',0",
			expected: []Token{{T: TK_LABEL, Tk: "MYSTR"}, {T: TK_DB, Tk: "DB"}, {T: TK_STRING, Tk: "Hello", ValStr: "Hello"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "0", ValInt: 0}},
		},
		{
			name:     "string",
			input:    "'Hello'",
			expected: []Token{{T: TK_STRING, Tk: "Hello", ValStr: "Hello"}},
		},
		{
			name:     "const read",
			input:    "add r1, $const1",
			expected: []Token{{T: TK_ALU, Tk: "ADD"}, {T: TK_REG, Tk: "R1"}, {T: TK_COMMA, Tk: ","}, {T: TK_BUCKS, Tk: "$"}, {T: TK_LABEL, Tk: "CONST1"}},
		},
		{
			name:     "equ const define",
			input:    "const1 equ 1234567890",
			expected: []Token{{T: TK_LABEL, Tk: "CONST1"}, {T: TK_EQU, Tk: "EQU"}, {T: TK_NUMBER, Tk: "1234567890", ValInt: 1234567890}},
		},
		{
			name:     "equ const define",
			input:    "equ const1 1234567890",
			expected: []Token{{T: TK_EQU, Tk: "EQU"}, {T: TK_LABEL, Tk: "CONST1"}, {T: TK_NUMBER, Tk: "1234567890", ValInt: 1234567890}},
		},
		{
			name:     "just label with colon",
			input:    "label1:",
			expected: []Token{{T: TK_LABEL, Tk: "LABEL1"}, {T: TK_COLON, Tk: ":"}},
		},
		{
			name:     "label with colon and instruction",
			input:    "label1:\tLD R9, R0",
			expected: []Token{{T: TK_LABEL, Tk: "LABEL1"}, {T: TK_COLON, Tk: ":"}, {T: TK_LD, Tk: "LD"}, {T: TK_REG, Tk: "R9"}, {T: TK_COMMA, Tk: ","}, {T: TK_REG, Tk: "R0"}},
		},
		{
			name:     "ld reg reg comment",
			input:    "\tLD R9, R0; some comment",
			expected: []Token{{T: TK_LD, Tk: "LD"}, {T: TK_REG, Tk: "R9"}, {T: TK_COMMA, Tk: ","}, {T: TK_REG, Tk: "R0"}},
		},
		{
			name:     "full line comment",
			input:    "; this is a comment",
			expected: nil,
		},
		{
			name:     "ldi reg label",
			input:    "ldi r0 variable1",
			expected: []Token{{T: TK_LDI, Tk: "LDI"}, {T: TK_REG, Tk: "R0"}, {T: TK_LABEL, Tk: "VARIABLE1"}},
		},
		{
			name:     "ret",
			input:    "\tret\t",
			expected: []Token{{T: TK_RET, Tk: "RET"}},
		},
		{
			name:     "call reg",
			input:    "call r7",
			expected: []Token{{T: TK_CALL, Tk: "CALL"}, {T: TK_REG, Tk: "R7"}},
		},
		{
			name:     "jmp reg",
			input:    "jmp r4",
			expected: []Token{{T: TK_JMP, Tk: "JMP"}, {T: TK_REG, Tk: "R4"}},
		},
		{
			name:     "ldb reg, [reg+imm]",
			input:    "ldb r10, [r12+5]",
			expected: []Token{{T: TK_LD_BYTE, Tk: "LDB"}, {T: TK_REG, Tk: "R10"}, {T: TK_COMMA, Tk: ","}, {T: TK_L_SQBR, Tk: "["}, {T: TK_REG, Tk: "R12"}, {T: TK_PLUS, Tk: "+"}, {T: TK_NUMBER, Tk: "5", ValInt: 5}, {T: TK_R_SQBR, Tk: "]"}},
		},
		{
			name:     "ldw reg, [reg++]",
			input:    "ldw r10, [r12++]",
			expected: []Token{{T: TK_LD_WORD, Tk: "LDW"}, {T: TK_REG, Tk: "R10"}, {T: TK_COMMA, Tk: ","}, {T: TK_L_SQBR, Tk: "["}, {T: TK_REG, Tk: "R12"}, {T: TK_PLUS, Tk: "+"}, {T: TK_PLUS, Tk: "+"}, {T: TK_R_SQBR, Tk: "]"}},
		},
		{
			name:     "movi reg, imm",
			input:    "movi r6, 10",
			expected: []Token{{T: TK_MOVI, Tk: "MOVI"}, {T: TK_REG, Tk: "R6"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "10", ValInt: 10}},
		},
		{
			name:     "add reg, imm",
			input:    "add r6, 1",
			expected: []Token{{T: TK_ALU, Tk: "ADD"}, {T: TK_REG, Tk: "R6"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "1", ValInt: 1}},
		},
		{
			name:     "ldb reg, [reg]",
			input:    "ldb r11, [r12]",
			expected: []Token{{T: TK_LD_BYTE, Tk: "LDB"}, {T: TK_REG, Tk: "R11"}, {T: TK_COMMA, Tk: ","}, {T: TK_L_SQBR, Tk: "["}, {T: TK_REG, Tk: "R12"}, {T: TK_R_SQBR, Tk: "]"}},
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
