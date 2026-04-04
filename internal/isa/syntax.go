package isa

import (
	"fmt"
	"strconv"
)

const (
	RuleALURegReg = iota
	LdRegImm
	AluRegImm
	Label
	Equ
	EvalConst
	EvalLabelRef
	DefineDBVar
	DefineDDVar
	JmpRel
	JmpRelCond
	CallRel
	CallJmpReg
	CallJmpCondReg
	RetCond
	Ret
)

type Rule struct {
	Type      uint8
	Syntax    [][]uint8
	Opcode    Opcode
	ParseFunc func(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error)
}

var aluRegRegSyntax = [][]uint8{{TK_ALU, TK_LD}, {TK_REG, TK_REG_SP, TK_REG_LR}, {TK_COMMA}, {TK_REG, TK_REG_SP, TK_REG_LR}}
var ldRegImmSyntax = [][]uint8{{TK_LD_0, TK_LD_1}, {TK_REG}, {TK_COMMA}, {TK_NUMBER}}
var labelSyntax = [][]uint8{{TK_LABEL}, {TK_COLON}}
var aluRegImmSyntax = [][]uint8{{TK_ALU, TK_LDI, TK_DJNZ}, {TK_REG}, {TK_COMMA}, {TK_NUMBER, TK_LABEL}}
var constEquSyntax = [][]uint8{{TK_LABEL}, {TK_EQU}, {TK_NUMBER}}
var evalConstSyntax = [][]uint8{{TK_BUCKS}, {TK_LABEL}}
var evalLabelRefSyntax = [][]uint8{{TK_AT}, {TK_LABEL}}
var defineVarSyntax = [][]uint8{{TK_LABEL}, {TK_DB}, {TK_NUMBER, TK_STRING}}
var defineVarDDSyntax = [][]uint8{{TK_LABEL}, {TK_DD}, {TK_NUMBER}}
var jmpRelSyntax = [][]uint8{{TK_JMP}, {TK_NUMBER, TK_LABEL}}
var jmpRelCondSyntax = [][]uint8{{TK_JMP}, {TK_REG}, {TK_CMP_EQ, TK_CMP_GT, TK_CMP_GTEQ, TK_CMP_LT, TK_CMP_LTEQ, TK_CMP_NEQ}, {TK_REG}, {TK_NUMBER, TK_LABEL}}
var callRelSyntax = [][]uint8{{TK_CALL}, {TK_LABEL, TK_NUMBER}}
var callJmpRegSyntax = [][]uint8{{TK_CALL, TK_JMP}, {TK_REG}}
var callJmpCondRegSyntax = [][]uint8{{TK_CALL, TK_JMP}, {TK_REG}, {TK_CMP_EQ, TK_CMP_NEQ, TK_CMP_GTEQ}, {TK_REG}, {TK_COMMA}, {TK_REG}}
var retCondSyntax = [][]uint8{{TK_RET}, {TK_REG}, {TK_CMP_EQ, TK_CMP_NEQ, TK_CMP_GTEQ}, {TK_REG}}
var retSyntax = [][]uint8{{TK_RET}}

var syntaxRules = []Rule{
	{Type: Equ, Syntax: constEquSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: parseConstEqu},
	{Type: EvalConst, Syntax: evalConstSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: parseEvalConst},
	{Type: EvalLabelRef, Syntax: evalLabelRefSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: parseEvalLabelRef},
	{Type: Label, Syntax: labelSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: parseLabel},
	{Type: DefineDBVar, Syntax: defineVarSyntax, Opcode: GetOpcode(OP_DB), ParseFunc: parseDefineDbVar},
	{Type: DefineDDVar, Syntax: defineVarDDSyntax, Opcode: GetOpcode(OP_DB), ParseFunc: parseDefineDdVar},
	{Type: RuleALURegReg, Syntax: aluRegRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG), ParseFunc: parseAluRegReg},
	{Type: LdRegImm, Syntax: ldRegImmSyntax, Opcode: GetOpcode(OP_LD_REG_IMM), ParseFunc: parseLdRegImm},
	{Type: AluRegImm, Syntax: aluRegImmSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: parseAluRegImm},
	{Type: JmpRel, Syntax: jmpRelSyntax, Opcode: GetOpcode(OP_JUMP_REL), ParseFunc: parseJmpRel},
	{Type: JmpRelCond, Syntax: jmpRelCondSyntax, Opcode: GetOpcode(OP_JUMP_REL), ParseFunc: parseJmpRelCond},
	{Type: CallJmpCondReg, Syntax: callJmpCondRegSyntax, Opcode: GetOpcode(OP_JUMP_CALL_RET_REG), ParseFunc: parseCallJmpCondReg},
	{Type: CallRel, Syntax: callRelSyntax, Opcode: GetOpcode(OP_CALL_REL), ParseFunc: parseCallRel},
	{Type: CallJmpReg, Syntax: callJmpRegSyntax, Opcode: GetOpcode(OP_JUMP_CALL_RET_REG), ParseFunc: parseCallJmpReg},
	{Type: RetCond, Syntax: retCondSyntax, Opcode: GetOpcode(OP_JUMP_CALL_RET_REG), ParseFunc: parseRetCond},
	{Type: Ret, Syntax: retSyntax, Opcode: GetOpcode(OP_JUMP_CALL_RET_REG), ParseFunc: parseRet},
}

func parseRegister(tokenrd string) (uint8, uint8, error) {
	bankd := uint8(0)
	rd, found := GetRegisterNumber(tokenrd)
	if !found {
		return 0, 0, fmt.Errorf("Invalid register: %s", tokenrd)
	}
	if rd > 7 {
		rd = rd - 7
		bankd = 1
	}
	return rd, bankd, nil
}

func parseRet(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	instr := Instruction{Opcode: rule.Opcode, Func2: 0, Ex: 2, Address: parser.CurAddress}
	tokens = append(tokens[:tokenpos], tokens[tokenpos+1:]...) //Remove parsed tokens from the list
	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

// Example: ADD r1, r10; LD r3, r13 ..
func parseAluRegReg(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd, bankd, rs, banks uint8
	var err error
	var func5 uint8

	aluT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	regST := tokens[tokenpos+3]
	instr := Instruction{}

	if regDT.T != TK_REG_LR && regDT.T != TK_REG_SP {
		rd, bankd, err = parseRegister(regDT.Tk)
		if err != nil {
			return tokens, err
		}
	}

	if regST.T != TK_REG_LR && regST.T != TK_REG_SP {
		rs, banks, err = parseRegister(regST.Tk)
		if err != nil {
			return tokens, err
		}
	}

	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Rs = rs
	instr = instr.makeEx(bankd, banks)
	instr.Address = parser.CurAddress
	func5 = 0 //LD for default LD REG, REG

	if aluT.T == TK_LD { // LD instruction
		if regDT.T == TK_REG_LR { // LD LR, REG
			if regST.T != TK_REG {
				return tokens, fmt.Errorf("Source should be R0-R15, found %v", regST)
			}
			func5 = 13
		} else if regDT.T == TK_REG_SP { //LD SP, REG
			if regST.T != TK_REG {
				return tokens, fmt.Errorf("Source should be R0-R15, found %v", regST)
			}
			func5 = 12
		} else if regST.T == TK_REG_LR { //LD REG, LR
			if regDT.T != TK_REG {
				return tokens, fmt.Errorf("Destination should be R0-R15, found %v", regDT)
			}
			func5 = 11
		} else if regST.T == TK_REG_SP { //LD REG, SP
			if regDT.T != TK_REG {
				return tokens, fmt.Errorf("Destination should be R0-R15, found %v", regDT)
			}
		}
	} else { //Alu instruction
		if err != nil {
			return tokens, err
		}
		func5, err = getFunc5FromALU(tokens[0].Tk)
		if err != nil {
			return tokens, err
		}
	}
	instr.Func5 = func5
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)

	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

// Example: LD R1,0xff
func parseLdRegImm(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd, bankd, func2 uint8
	var err error

	ldT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]
	instr := Instruction{}

	rd, bankd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}
	if ldT.T == TK_LD_1 {
		func2 = bankd<<1 | 1
	} else {
		func2 = bankd << 1
	}

	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Imm = int16(immT.ValInt)
	instr.Address = parser.CurAddress
	instr.Func2 = func2
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)

	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

// ALU REG, IMM (7 bit Immediate operations)
//
//	func(2 bit) = register bank (0 or 1)
//	func(0-1 bits):
//
// 0) ADD/SUB: Rd = Rd + (signed(IMM))
// 1) SHL/SHR: Rd = << or >> signed(IMM & 31)
// 2) 2) LDI Rd = PC - IMM; get 32 bit address relative to PC, IMM in instructions (-128 ... 0 instructions)
// 3) DJNZ Rd, PC + signed(IMM); Rd-- if not zero, jump taken. IMM in instructions (-64 ... +63 instructions)
func parseAluRegImm(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd, bankd, func3 uint8
	var err error

	aluT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]
	instr := Instruction{}

	rd, bankd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}

	if immT.T == TK_LABEL {
		instr.Label = immT.Tk
	} else {
		instr.Imm = int16(immT.ValInt)
	}

	switch aluT.Tk {
	case "ADD":
		func3 = 0
	case "SUB":
		instr.Imm = ^instr.Imm + 1
		func3 = 0
	case "SHL":
		func3 = 1
	case "SHR":
		instr.Imm = ^instr.Imm + 1
		func3 = 1
	case "LDI":
		func3 = 2
	case "DJNZ":
		func3 = 3
	default:
		return tokens, fmt.Errorf("Instruction should be one of: ADD, SUB, SHL, SHR, LDI, DJNZ! '%s' not found", aluT.Tk)
	}

	func3 = func3 | bankd<<2
	instr.Opcode = rule.Opcode
	instr.Rd = rd

	instr.Address = parser.CurAddress
	instr.Func3 = func3
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)

	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

// Example: call func1
func parseCallRel(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	immT := tokens[tokenpos+1]
	instr := Instruction{}
	instr.Rd = 0
	if immT.T == TK_LABEL {
		instr.Label = immT.Tk
	} else {
		instr.Imm = int16(immT.ValInt)
	}
	instr.Opcode = rule.Opcode
	instr.Address = parser.CurAddress
	parser.addInstruction(instr)
	parser.CurAddress += 2

	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+2:]...)
	return tokens, nil
}

// Example: jmp loop1
func parseJmpRel(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	immT := tokens[tokenpos+1]
	instr := Instruction{}
	instr.Rd = 0
	if immT.T == TK_LABEL {
		instr.Label = immT.Tk
	} else {
		instr.Rs = uint8(immT.ValInt >> 7)
		instr.Ex = uint8(immT.ValInt>>5) & 0b0000_0011
		instr.Func2 = uint8(immT.ValInt>>3) & 0b0000_0011
		instr.Rx = uint8(immT.ValInt) & 0b0000_0111
	}
	instr.Opcode = rule.Opcode
	instr.Address = parser.CurAddress
	parser.addInstruction(instr)
	parser.CurAddress += 2

	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+2:]...)
	return tokens, nil
}

// Example: jmp r1 != r2 loop1
func parseJmpRelCond(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	regsT := tokens[tokenpos+1]
	condT := tokens[tokenpos+2]
	regxT := tokens[tokenpos+3]
	immT := tokens[tokenpos+4]

	rs, banks, err := parseRegister(regsT.Tk)
	if err != nil {
		return tokens, err
	}
	rx, bankx, err := parseRegister(regxT.Tk)
	if err != nil {
		return tokens, err
	}
	if banks == 1 || bankx == 1 {
		return tokens, fmt.Errorf("Only registers from 1 bank available (R0-R7) for JMP REG COND REG LABEL. Got %v and %v", regsT, regxT)
	}

	instr := Instruction{}
	instr.Rs = rs
	instr.Rx = rx
	switch condT.T {
	case TK_CMP_EQ:
		instr.Rd = 1
	case TK_CMP_NEQ:
		instr.Rd = 2
	case TK_CMP_GT:
		instr.Rd = 3
	case TK_CMP_GTEQ:
		instr.Rd = 4
	case TK_CMP_LT:
		instr.Rd = 5
	case TK_CMP_LTEQ:
		instr.Rd = 6
	}

	if immT.T == TK_LABEL {
		instr.Label = immT.Tk
	} else {
		instr.Rs = uint8(immT.ValInt >> 7)
		instr.Ex = uint8(immT.ValInt>>5) & 0b0000_0011
		instr.Func2 = uint8(immT.ValInt>>3) & 0b0000_0011
		instr.Rx = uint8(immT.ValInt) & 0b0000_0111
	}
	instr.Opcode = rule.Opcode
	instr.Address = parser.CurAddress
	parser.addInstruction(instr)
	parser.CurAddress += 2

	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+5:]...)
	return tokens, nil
}

// Example: call r6; jmp r3
func parseCallJmpReg(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	func2, rs, rx := uint8(0), uint8(0), uint8(0)
	rd, bank, err := parseRegister(tokens[tokenpos+1].Tk)
	if err != nil {
		return tokens, err
	}
	if bank == 1 {
		return tokens, fmt.Errorf("Only registers from 1 bank available (R0-R7) for JMP REG. Got %v", tokens[tokenpos+1])
	}
	ex := uint8(0) //JMP
	if tokens[tokenpos].T == TK_CALL {
		ex = 1
	}
	instr := Instruction{Opcode: rule.Opcode, Func2: func2, Rs: rs, Rx: rx, Rd: rd, Ex: ex, Address: parser.CurAddress}
	parser.addInstruction(instr)
	parser.CurAddress += 2

	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+2:]...)
	return tokens, nil
}

// Example: RET r0 == r1
func parseRetCond(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	regsT := tokens[tokenpos+1]
	condT := tokens[tokenpos+2]
	regxT := tokens[tokenpos+3]
	ex := uint8(2) //RET
	rs, banks, err := parseRegister(regsT.Tk)
	if err != nil {
		return tokens, err
	}
	rx, bankx, err := parseRegister(regxT.Tk)
	if err != nil {
		return tokens, err
	}
	if banks == 1 || bankx == 1 {
		return tokens, fmt.Errorf("Only registers from 1 bank available (R0-R7) for JMP/CALL REG COND REG , REG. Got %v, %v", regsT, regxT)
	}
	func2 := uint8(0)
	switch condT.T {
	case TK_CMP_EQ:
		func2 = 1
	case TK_CMP_NEQ:
		func2 = 2
	case TK_CMP_GTEQ:
		func2 = 3
	}
	instr := Instruction{Opcode: rule.Opcode, Func2: func2, Rs: rs, Rx: rx, Rd: 0, Ex: ex, Address: parser.CurAddress}
	parser.addInstruction(instr)
	parser.CurAddress += 2
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)
	return tokens, nil
}

// Example: call r1 == r2, r3; jmp r1 != r2, r3
func parseCallJmpCondReg(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	regsT := tokens[tokenpos+1]
	condT := tokens[tokenpos+2]
	regxT := tokens[tokenpos+3]
	regdT := tokens[tokenpos+5]
	ex := uint8(0) //JMP
	if tokens[tokenpos].T == TK_CALL {
		ex = 1
	}
	rs, banks, err := parseRegister(regsT.Tk)
	if err != nil {
		return tokens, err
	}
	rx, bankx, err := parseRegister(regxT.Tk)
	if err != nil {
		return tokens, err
	}
	rd, bankd, err := parseRegister(regdT.Tk)
	if err != nil {
		return tokens, err
	}
	if banks == 1 || bankx == 1 || bankd == 1 {
		return tokens, fmt.Errorf("Only registers from 1 bank available (R0-R7) for JMP/CALL REG COND REG , REG. Got %v, %v, %v", regsT, regxT, regdT)
	}
	func2 := uint8(0)
	switch condT.T {
	case TK_CMP_EQ:
		func2 = 1
	case TK_CMP_NEQ:
		func2 = 2
	case TK_CMP_GTEQ:
		func2 = 3
	}
	instr := Instruction{Opcode: rule.Opcode, Func2: func2, Rs: rs, Rx: rx, Rd: rd, Ex: ex, Address: parser.CurAddress}
	parser.addInstruction(instr)
	parser.CurAddress += 2
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+6:]...)
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
