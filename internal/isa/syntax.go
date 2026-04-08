package isa

import (
	"fmt"
	"strconv"
)

const (
	aluRegReg = iota
	regImm
	aluImm
	memLd
	memLdImm
	memSt
	memStImm
	ldi

	callRel
	callRelRd
	callRd
	ret
	retRd
	jmpRel
	jmpRd

	branch

	label
	equ
	evalConst
	evalLabelRef
	defineDBVar
	defineDDVar
)

type Rule struct {
	Type      uint8
	Syntax    [][]uint8
	Opcode    Opcode
	ParseFunc func(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error)
}

// Preprocessor
var labelSyntax = [][]uint8{{TK_LABEL}, {TK_COLON}}
var constEquSyntax = [][]uint8{{TK_LABEL}, {TK_EQU}, {TK_NUMBER}}
var evalConstSyntax = [][]uint8{{TK_BUCKS}, {TK_LABEL}}
var evalLabelRefSyntax = [][]uint8{{TK_AT}, {TK_LABEL}}
var defineVarSyntax = [][]uint8{{TK_LABEL}, {TK_DB}, {TK_NUMBER, TK_STRING}}
var defineVarDDSyntax = [][]uint8{{TK_LABEL}, {TK_DD}, {TK_NUMBER}}

// Instructions
var aluRegRegSyntax = [][]uint8{{TK_ALU}, {TK_REG}, {TK_COMMA}, {TK_REG}}
var regImmSyntax = [][]uint8{{TK_MOVI, TK_MOVH}, {TK_REG}, {TK_COMMA}, {TK_NUMBER}}
var aluImmSyntax = [][]uint8{{TK_ALU}, {TK_REG}, {TK_COMMA}, {TK_NUMBER, TK_LABEL}}

var memLdSyntax = [][]uint8{{TK_LD_BYTE, TK_LD_WORD}, {TK_REG}, {TK_COMMA}, {TK_L_SQBR}, {TK_REG}, {TK_R_SQBR}}
var memLdImmSyntax = [][]uint8{{TK_LD_BYTE, TK_LD_WORD}, {TK_REG}, {TK_COMMA}, {TK_L_SQBR}, {TK_REG}, {TK_NUMBER}, {TK_R_SQBR}}
var memStSyntax = [][]uint8{{TK_ST_BYTE, TK_ST_WORD}, {TK_L_SQBR}, {TK_REG}, {TK_R_SQBR}, {TK_COMMA}, {TK_REG}}
var memStImmSyntax = [][]uint8{{TK_ST_BYTE, TK_ST_WORD}, {TK_L_SQBR}, {TK_REG}, {TK_NUMBER}, {TK_R_SQBR}, {TK_COMMA}, {TK_REG}}

var ldiSyntax = [][]uint8{{TK_LDI}, {TK_REG}, {TK_COMMA}, {TK_LABEL, TK_NUMBER}}

var callRelSyntax = [][]uint8{{TK_CALL}, {TK_LABEL, TK_NUMBER}}
var callRelRdSyntax = [][]uint8{{TK_CALL}, {TK_LABEL, TK_NUMBER}, {TK_COMMA}, {TK_REG}}
var callRdSyntax = [][]uint8{{TK_CALL}, {TK_REG}}
var retSyntax = [][]uint8{{TK_RET}}
var retRdSyntax = [][]uint8{{TK_RET}, {TK_REG}}
var jmpRelSyntax = [][]uint8{{TK_JMP}, {TK_LABEL, TK_NUMBER}}
var jmpRdSyntax = [][]uint8{{TK_JMP}, {TK_REG}}

var branchSyntax = [][]uint8{{TK_BEQZ, TK_BNEZ, TK_BGTZ, TK_BLTZ}, {TK_REG}, {TK_COMMA}, {TK_LABEL, TK_NUMBER}}

var syntaxRules = []Rule{
	// Preprocessor
	{Type: equ, Syntax: constEquSyntax, Opcode: GetOpcode(OP_DB), ParseFunc: parseConstEqu},
	{Type: evalConst, Syntax: evalConstSyntax, Opcode: GetOpcode(OP_DB), ParseFunc: parseEvalConst},
	{Type: evalLabelRef, Syntax: evalLabelRefSyntax, Opcode: GetOpcode(OP_DB), ParseFunc: parseEvalLabelRef},
	{Type: label, Syntax: labelSyntax, Opcode: GetOpcode(OP_DB), ParseFunc: parseLabel},
	{Type: defineDBVar, Syntax: defineVarSyntax, Opcode: GetOpcode(OP_DB), ParseFunc: parseDefineDbVar},
	{Type: defineDDVar, Syntax: defineVarDDSyntax, Opcode: GetOpcode(OP_DB), ParseFunc: parseDefineDdVar},

	// Instructions
	{Type: aluRegReg, Syntax: aluRegRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG), ParseFunc: parseAluRegReg},
	{Type: regImm, Syntax: regImmSyntax, Opcode: GetOpcode(OP_REG_IMM), ParseFunc: parseRegImm},
	{Type: aluImm, Syntax: aluImmSyntax, Opcode: GetOpcode(OP_ALU_IMM), ParseFunc: parseAluImm},
	//Mem
	{Type: memLd, Syntax: memLdSyntax, Opcode: GetOpcode(OP_MEM), ParseFunc: parseMemAll},
	{Type: memLdImm, Syntax: memLdImmSyntax, Opcode: GetOpcode(OP_MEM), ParseFunc: parseMemAll},
	{Type: memSt, Syntax: memStSyntax, Opcode: GetOpcode(OP_MEM), ParseFunc: parseMemAll},
	{Type: memStImm, Syntax: memStImmSyntax, Opcode: GetOpcode(OP_MEM), ParseFunc: parseMemAll},
	//LDI
	{Type: ldi, Syntax: ldiSyntax, Opcode: GetOpcode(OP_LDI), ParseFunc: parseLdi},
	//CALL/RET/JMP
	{Type: callRel, Syntax: callRelSyntax, Opcode: GetOpcode(OP_CALL_JUMP_RET), ParseFunc: parseCallJmpRetAll},
	{Type: callRelRd, Syntax: callRelRdSyntax, Opcode: GetOpcode(OP_CALL_JUMP_RET), ParseFunc: parseCallJmpRetAll},
	{Type: callRd, Syntax: callRdSyntax, Opcode: GetOpcode(OP_CALL_JUMP_RET), ParseFunc: parseCallJmpRetAll},
	{Type: ret, Syntax: retSyntax, Opcode: GetOpcode(OP_CALL_JUMP_RET), ParseFunc: parseCallJmpRetAll},
	{Type: retRd, Syntax: retRdSyntax, Opcode: GetOpcode(OP_CALL_JUMP_RET), ParseFunc: parseCallJmpRetAll},
	{Type: jmpRel, Syntax: jmpRelSyntax, Opcode: GetOpcode(OP_CALL_JUMP_RET), ParseFunc: parseCallJmpRetAll},
	{Type: jmpRd, Syntax: jmpRdSyntax, Opcode: GetOpcode(OP_CALL_JUMP_RET), ParseFunc: parseCallJmpRetAll},
	//BRANCH
	{Type: branch, Syntax: branchSyntax, Opcode: GetOpcode(OP_BRANCH), ParseFunc: parseBranch},
}

func parseRegister(tokenrd string) (uint8, error) {
	rd, found := GetRegisterNumber(tokenrd)
	if !found {
		return 0, fmt.Errorf("Invalid register: %s", tokenrd)
	}
	return rd, nil
}

// Example: ADD r1, r10; MOV r3, r13 ..
func parseAluRegReg(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd, rs uint8
	var err error

	aluT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	regST := tokens[tokenpos+3]
	instr := Instruction{}

	rd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}
	rs, err = parseRegister(regST.Tk)
	if err != nil {
		return tokens, err
	}

	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Rs = rs
	instr.Address = parser.CurAddress
	instr.Func, err = getFuncFromAlu(aluT.Tk)
	if err != nil {
		return tokens, err
	}
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)
	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

// Example: MOVI R1,0xff
func parseRegImm(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd uint8
	var err error

	ldT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]
	instr := Instruction{}

	rd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}
	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Imm = int16(immT.ValInt)
	instr.Address = parser.CurAddress
	instr.Func, err = getFuncFromRegImm(ldT.Tk)
	if err != nil {
		return tokens, err
	}
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)
	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

// Example: Add Rd, 0x10
func parseAluImm(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd uint8
	var err error

	aluT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]
	instr := Instruction{}

	rd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}
	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Imm = int16(immT.ValInt)
	instr.Address = parser.CurAddress
	instr.Func, err = getFuncFromAluImm(aluT.Tk)
	if err != nil {
		return tokens, err
	}
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)
	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

// Example: LD Rd, [Rs Imm]; ST [Rs Imm], Rd
func parseMemAll(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd, rs uint8
	var err error

	ldStT := tokens[tokenpos]
	var regDT Token
	var regST Token
	var imm int16
	var isLd uint8
	var isByte uint8

	switch ldStT.T {
	case TK_LD_BYTE:
		isLd = 1
		isByte = 1
	case TK_LD_WORD:
		isLd = 1
		isByte = 0
	case TK_ST_BYTE:
		isLd = 0
		isByte = 1
	case TK_ST_WORD:
		isLd = 0
		isByte = 0
	}

	if tokens[tokenpos+1].T == TK_REG { // Rd, [Rs Imm]
		regDT = tokens[tokenpos+1]
		regST = tokens[tokenpos+4]
		if tokens[tokenpos+5].T == TK_NUMBER { // Rd, [Rs Imm]
			imm = int16(tokens[tokenpos+5].ValInt)
		} else { // Rd, [Rs]
		}
	} else { // [Rs Imm], Rd
		regST = tokens[tokenpos+2]
		if tokens[tokenpos+3].T == TK_NUMBER { // [Rs Imm], Rd
			imm = int16(tokens[tokenpos+3].ValInt)
			regDT = tokens[tokenpos+6]
		} else { // [Rs], Rd
			regDT = tokens[tokenpos+5]
		}
	}

	rd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}
	rs, err = parseRegister(regST.Tk)
	if err != nil {
		return tokens, err
	}

	instr := Instruction{}
	instr.Rd = rd
	instr.Rs = rs
	instr.Imm = imm
	instr.Func = isByte<<1 | isLd
	instr.Opcode = rule.Opcode
	instr.Address = parser.CurAddress
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+int(len(rule.Syntax)):]...)
	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

// Example: LDI Rd, imm; LDI Rd, LABEL
func parseLdi(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd uint8
	var err error

	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]
	instr := Instruction{}

	if immT.T == TK_LABEL {
		instr.Label = immT.Tk
	} else {
		instr.Imm = int16(immT.ValInt)
	}

	rd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}
	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Address = parser.CurAddress
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)
	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

func parseCallJmpRetAll(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd uint8 = 14 // Default Link Register if not specified
	var err error
	var regDT, immT Token

	instr := Instruction{}

	switch rule.Type {
	case callRel: // CALL Imm(7)
		immT = tokens[tokenpos+1]
		instr.Func = 0
	case callRelRd: // CALL Imm(7), Rd
		immT = tokens[tokenpos+1]
		regDT = tokens[tokenpos+3]
		rd, err = parseRegister(regDT.Tk)
		if err != nil {
			return tokens, err
		}
		instr.Func = 0
	case callRd: // CALL Rd
		regDT = tokens[tokenpos+1]
		rd, err = parseRegister(regDT.Tk)
		if err != nil {
			return tokens, err
		}
		instr.Func = 1
	case ret: // RET
		instr.Func = 2
	case retRd: // RET  Rd
		regDT = tokens[tokenpos+1]
		rd, err = parseRegister(regDT.Tk)
		if err != nil {
			return tokens, err
		}
		instr.Func = 2
	case jmpRel: // JR signed(Imm(7))
		immT = tokens[tokenpos+1]
		instr.Func = 3
	case jmpRd: // JMP Rd; same as RET Rd
		regDT = tokens[tokenpos+1]
		rd, err = parseRegister(regDT.Tk)
		if err != nil {
			return tokens, err
		}
		instr.Func = 2
	}

	if instr.Func == 0 || instr.Func == 3 { //CALL Imm or JR Imm
		if immT.T == TK_LABEL {
			instr.Label = immT.Tk
		} else {
			instr.Imm = int16(immT.ValInt)
		}
	}

	instr.Rd = rd
	instr.Opcode = rule.Opcode
	instr.Address = parser.CurAddress
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+int(len(rule.Syntax)):]...)
	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

// Example: BEQZ R1, LABEL
func parseBranch(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd uint8
	var err error

	brT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]
	instr := Instruction{}

	rd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}
	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Imm = int16(immT.ValInt)
	instr.Address = parser.CurAddress
	instr.Func, err = getFuncFromBranch(brT.Tk)
	if err != nil {
		return tokens, err
	}
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)
	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

func addLabel(parser *Parser, token Token) error {
	if _, ok := parser.Labels[token.Tk]; ok {
		return fmt.Errorf("Label '%s' already defined! On address: 0x%08x", token.Tk, parser.Labels[token.Tk])
	}
	parser.Labels[token.Tk] = parser.CurAddress
	return nil
}

// Example: label1:
func parseLabel(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	if err := addLabel(parser, tokens[tokenpos]); err != nil {
		return tokens, err
	}
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+2:]...)
	return tokens, nil
}

// Example: myconst1 equ 0x1234
func parseConstEqu(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	constT := tokens[tokenpos]
	numberT := tokens[tokenpos+2]
	if _, ok := parser.Constants[constT.Tk]; ok {
		return tokens, fmt.Errorf("Constant '%s' already defined! Value: 0x%08x", constT.Tk, parser.Constants[constT.Tk])
	}
	parser.Constants[constT.Tk] = numberT.ValInt
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+3:]...)
	return tokens, nil
}

// Example: add r1, $const1; replace $const1 by number
func parseEvalConst(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	constT := tokens[tokenpos+1]
	constVal, ok := parser.Constants[constT.Tk]
	if !ok {
		return tokens, fmt.Errorf("Constant '%s' is not defined!", constT.Tk)
	}
	//Remove second token and replace first by number
	tokens[tokenpos] = Token{T: TK_NUMBER, Tk: strconv.Itoa(int(constVal)), ValInt: constVal}
	tokens = append(tokens[:tokenpos+1], tokens[tokenpos+2:]...)
	return tokens, nil
}

// Example: mystrref db @mystr
func parseEvalLabelRef(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	labelT := tokens[tokenpos+1]
	labelVal, ok := parser.Labels[labelT.Tk]
	if !ok {
		return tokens, fmt.Errorf("Label '%s' is not defined!", labelT.Tk)
	}
	//Remove second token and replace first by number
	tokens[tokenpos] = Token{T: TK_NUMBER, Tk: strconv.Itoa(int(labelVal)), ValInt: int32(labelVal)}
	tokens = append(tokens[:tokenpos+1], tokens[tokenpos+2:]...)
	return tokens, nil
}

// Example: Mystr db 'Hello',0
func parseDefineDbVar(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	labelT := tokens[tokenpos]
	instr := Instruction{}
	instrByte := 1

	if err := addLabel(parser, labelT); err != nil {
		return tokens, err
	}

	for i := tokenpos + 2; i < len(tokens); i++ {
		numberT := tokens[i]
		switch numberT.T {
		case TK_COMMA:
			continue
		case TK_NUMBER:
			// if numberT.ValInt < 0 || numberT.ValInt > 255 {
			// 	return tokens, fmt.Errorf("Number out of range(0..255) for db: 0x%08X", uint32(numberT.ValInt))
			// }
			for number := uint32(numberT.ValInt); number > 0; number >>= 8 {
				if instrByte == 1 {
					instr.Opcode = rule.Opcode
					instr.Address = parser.CurAddress
					instr.Imm = instr.Imm | int16(number&0x000000FF)<<8
					instrByte = 0
				} else {
					instr.Imm = instr.Imm | int16(number&0x000000FF)
					instrByte = 1
					parser.addInstruction(instr)
					parser.CurAddress += 2
					instr = Instruction{}
				}
			}
		case TK_STRING:
			for _, char := range []byte(numberT.ValStr) {
				if instrByte == 1 {
					instr.Opcode = rule.Opcode
					instr.Address = parser.CurAddress
					instr.Imm = instr.Imm | int16(char)<<8
					instrByte = 0
				} else {
					instr.Imm = instr.Imm | int16(char)
					instrByte = 1
					parser.addInstruction(instr)
					parser.CurAddress += 2
					instr = Instruction{}
				}
			}
		default:
			return tokens, fmt.Errorf("Unexpected token: %s", numberT)
		}
	}
	if instrByte == 0 {
		parser.addInstruction(instr)
		parser.CurAddress += 2
	}
	tokens = []Token{}
	return tokens, nil
}

// Example: Mystrref dd @Mystr, 0xFFBBCCDD
func parseDefineDdVar(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	labelT := tokens[tokenpos]
	instr := Instruction{}
	instrByte := 1

	if err := addLabel(parser, labelT); err != nil {
		return tokens, err
	}

	for i := tokenpos + 2; i < len(tokens); i++ {
		numberT := tokens[i]
		switch numberT.T {
		case TK_COMMA:
			continue
		case TK_NUMBER:
			// if numberT.ValInt < 0 || numberT.ValInt > 255 {
			// 	return tokens, fmt.Errorf("Number out of range(0..255) for db: 0x%08X", uint32(numberT.ValInt))
			// }
			number := uint32(numberT.ValInt)
			for range 3 {
				if instrByte == 1 {
					instr.Opcode = rule.Opcode
					instr.Address = parser.CurAddress
					instr.Imm = instr.Imm | int16(number&0x000000FF)<<8
					instrByte = 0
				} else {
					instr.Imm = instr.Imm | int16(number&0x000000FF)
					instrByte = 1
					parser.addInstruction(instr)
					parser.CurAddress += 2
					instr = Instruction{}
				}
				number >>= 8
			}
		default:
			return tokens, fmt.Errorf("Unexpected token: %s", numberT)
		}
	}
	if instrByte == 0 {
		parser.addInstruction(instr)
		parser.CurAddress += 2
	}
	tokens = []Token{}
	return tokens, nil
}
