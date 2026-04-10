package isa

import "fmt"

type Token struct {
	T      uint8
	Tk     string
	ValInt int32
	ValStr string
}

func (token Token) String() string {
	if token.ValInt != 0 {
		return fmt.Sprintf("%s('%s':%d)", GetTokenTypeString(token.T), token.Tk, token.ValInt)
	} else {
		return fmt.Sprintf("%s('%s')", GetTokenTypeString(token.T), token.Tk)
	}
}

const (
	TK_LD = iota
	TK_ST_BYTE
	TK_ST_WORD
	TK_LDI
	TK_LD_BYTE
	TK_LD_WORD
	TK_MOVI
	TK_MOVH
	TK_REG
	TK_COMMA
	TK_COLON
	TK_AT
	TK_ALU
	TK_NUMBER
	TK_PLUS
	TK_MINUS
	TK_L_SQBR
	TK_R_SQBR
	TK_PUSH
	TK_POP

	//control flow tokens
	TK_JMP
	TK_CALL
	TK_RET

	TK_INT
	TK_RETI

	TK_BEQZ
	TK_BNEZ
	TK_BGTZ
	TK_BLTZ

	TK_LABEL
	TK_EQU
	TK_BUCKS
	TK_ALIGN4

	//vars
	TK_DB
	TK_DD
	TK_STRING
)

var tokenTypeToString = map[uint8]string{
	TK_LD:      "TK_LD",
	TK_ST_WORD: "TK_ST_WORD",
	TK_ST_BYTE: "TK_ST_BYTE",
	TK_LDI:     "TK_LDI",
	TK_LD_BYTE: "TK_LD_BYTE",
	TK_LD_WORD: "TK_LD_WORD",
	TK_MOVI:    "TK_MOVI",
	TK_MOVH:    "TK_MOVH",
	TK_REG:     "TK_REG",
	TK_COMMA:   "TK_COMMA",
	TK_COLON:   "TK_COLON",
	TK_AT:      "TK_AT",
	TK_ALU:     "TK_ALU",
	TK_NUMBER:  "TK_NUMBER",
	TK_PLUS:    "TK_PLUS",
	TK_MINUS:   "TK_MINUS",
	TK_L_SQBR:  "TK_L_SQBR",
	TK_R_SQBR:  "TK_R_SQBR",
	TK_PUSH:    "TK_PUSH",
	TK_POP:     "TK_POP",
	TK_JMP:     "TK_JMP",
	TK_CALL:    "TK_CALL",
	TK_RET:     "TK_RET",
	TK_BEQZ:    "TK_BEQZ",
	TK_BNEZ:    "TK_BNEZ",
	TK_BGTZ:    "TK_BGTZ",
	TK_BLTZ:    "TK_BLTZ",
	TK_INT:     "TK_INT",
	TK_RETI:    "TK_RETI",
	TK_LABEL:   "TK_LABEL",
	TK_EQU:     "TK_EQU",
	TK_BUCKS:   "TK_BUCKS",
	TK_DB:      "TK_DB",
	TK_DD:      "TK_DD",
	TK_STRING:  "TK_STRING",
	TK_ALIGN4:  "TK_ALIGN4",
}

func GetTokenTypeString(tokenType uint8) string {
	if tokenString, exists := tokenTypeToString[tokenType]; exists {
		return tokenString
	}
	return "TK_UNKNOWN"
}

var mapRegisterToNumber = map[string]uint8{
	"R0":  0,
	"R1":  1,
	"R2":  2,
	"R3":  3,
	"R4":  4,
	"R5":  5,
	"R6":  6,
	"R7":  7,
	"R8":  8,
	"R9":  9,
	"R10": 10,
	"R11": 11,
	"R12": 12,
	"R13": 13,
	"R14": 14,
	"R15": 15,
}

func GetRegisterNumber(reg string) (uint8, bool) {
	if regNum, exists := mapRegisterToNumber[reg]; exists {
		return regNum, true
	}
	return 0, false
}

var mapWordToToken = map[string]uint8{
	"LD":      TK_LD,
	"STB":     TK_ST_BYTE,
	"STW":     TK_ST_WORD,
	"LDI":     TK_LDI,
	"LDB":     TK_LD_BYTE,
	"LDW":     TK_LD_WORD,
	"MOVI":    TK_MOVI,
	"MOVH":    TK_MOVH,
	"R0":      TK_REG,
	"R1":      TK_REG,
	"R2":      TK_REG,
	"R3":      TK_REG,
	"R4":      TK_REG,
	"R5":      TK_REG,
	"R6":      TK_REG,
	"R7":      TK_REG,
	"R8":      TK_REG,
	"R9":      TK_REG,
	"R10":     TK_REG,
	"R11":     TK_REG,
	"R12":     TK_REG,
	"R13":     TK_REG,
	"R14":     TK_REG,
	"R15":     TK_REG,
	"MOV":     TK_ALU,
	"ADD":     TK_ALU,
	"SUB":     TK_ALU,
	"SHL":     TK_ALU,
	"SHR":     TK_ALU,
	"AND":     TK_ALU,
	"OR":      TK_ALU,
	"XOR":     TK_ALU,
	"NOT":     TK_ALU,
	"MUL":     TK_ALU,
	"PUSH":    TK_PUSH,
	"POP":     TK_POP,
	"JMP":     TK_JMP,
	"CALL":    TK_CALL,
	"RET":     TK_RET,
	"BEQZ":    TK_BEQZ,
	"BNEZ":    TK_BNEZ,
	"BGTZ":    TK_BGTZ,
	"BLTZ":    TK_BLTZ,
	"INT":     TK_INT,
	"RETI":    TK_RETI,
	"EQU":     TK_EQU,
	"DB":      TK_DB,
	"DD":      TK_DD,
	".ALIGN4": TK_ALIGN4,
}

func GetTokenTypeByWord(word string) (uint8, bool) {
	if tokenType, exists := mapWordToToken[word]; exists {
		return tokenType, true
	}
	return 0, false
}
