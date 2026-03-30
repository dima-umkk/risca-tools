package isa

import (
	"fmt"
)

type Parser struct {
	Instructions []Instruction
	Memory       []uint8
	Labels       map[string]uint32
	Constants    map[string]int32
	LineNumber   uint32
	CurAddress   uint32
}

func NewParser() *Parser {
	return &Parser{
		Instructions: make([]Instruction, 0, 1024),
		Memory:       make([]uint8, 0, 1024),
		Labels:       make(map[string]uint32),
		LineNumber:   1,
		CurAddress:   0,
	}
}

func (parser *Parser) ParseLine(line string) (Instruction, bool, error) {
	tokens, err := Tokenize(line)
	if err != nil { //Error parsing tokens
		return Instruction{}, false, err
	}
	if len(tokens) == 0 { //Skip line
		return Instruction{}, true, nil
	}

	var matched bool
	var instr = Instruction{}

	for len(tokens) > 0 {
		for _, rule := range syntaxRules {
			tokens, matched, err, instr = parser.checkRule(rule, tokens, instr)
			if err != nil {
				return Instruction{}, false, err
			}
			if matched {
				break
			}
		}
		if !matched {
			return Instruction{}, false, fmt.Errorf("Syntax error")
		}
	}
	return instr, false, err
}

func tokenIn(token Token, tokens []uint8) bool {
	for _, t := range tokens {
		if token.T == t {
			return true
		}
	}
	return false
}

func (parser *Parser) checkRule(rule Rule, tokens []Token, instr Instruction) ([]Token, bool, error, Instruction) {
	pos := 0
	err := error(nil)
	for pos+len(rule.Syntax)-1 < len(tokens) {
		match := true
		for i, sxtokens := range rule.Syntax {
			if !tokenIn(tokens[pos+i], sxtokens) {
				match = false
				break
			}
		}
		if match {
			tokens, instr, err = rule.ParseFunc(parser, rule, tokens, pos, instr)
			return tokens, true, err, instr
		} else {
			pos++
		}
	}
	return tokens, false, nil, instr
}

func (parser *Parser) ProcessLabels() error {
	for i, instr := range parser.Instructions {
		if instr.Label != "" {
			addr, ok := parser.Labels[instr.Label]
			if !ok {
				return fmt.Errorf("Label not defined: %s", instr.Label)
			}
			offset := int32(addr) - int32(instr.Address)
			parser.Instructions[i].Imm = int16(offset)
			switch instr.Opcode.Type {
			case OP_TYPE_7_IMM_REG:
				if offset < -64 || offset > 63 {
					return fmt.Errorf("Label too large for instruction! %s", instr)
				}
			case OP_TYPE_8_IMM_REG:
				if offset < -128 || offset > 127 {
					return fmt.Errorf("Label too large for instruction! %s", instr)
				}
			case OP_TYPE_13_IMM:
				if offset < -4096 || offset > 4095 {
					return fmt.Errorf("Label too large for instruction! %s", instr)
				}
			default:
				return fmt.Errorf("Label not applicable for instruction! %s", instr)
			}
		}
	}
	return nil
}

func parseRegisters(tokenrd, tokenrs string) (uint8, uint8, uint8, uint8, error) {
	bankd, banks := uint8(0), uint8(0)
	rd, found := GetRegisterNumber(tokenrd)
	if !found {
		return 0, 0, 0, 0, fmt.Errorf("Invalid register: %s", tokenrd)
	}
	rs, found := GetRegisterNumber(tokenrs)
	if !found {
		return 0, 0, 0, 0, fmt.Errorf("Invalid register: %s", tokenrs)
	}
	if rd > 7 {
		rd = rd - 7
		bankd = 1
	}
	if rs > 7 {
		rs = rs - 7
		banks = 1
	}
	return rd, bankd, rs, banks, nil
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

func ParseAluRegReg(parser *Parser, rule Rule, tokens []Token, tokenpos int, instr Instruction) ([]Token, Instruction, error) {
	var rd, bankd, rs, banks uint8
	var err error
	var func5 uint8

	aluT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	regST := tokens[tokenpos+3]

	if regDT.T != TK_REG_LR && regDT.T != TK_REG_SP {
		rd, bankd, err = parseRegister(regDT.Tk)
		if err != nil {
			return tokens, instr, err
		}
	}

	if regST.T != TK_REG_LR && regST.T != TK_REG_SP {
		rs, banks, err = parseRegister(regST.Tk)
		if err != nil {
			return tokens, instr, err
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
				return tokens, instr, fmt.Errorf("Source should be R0-R15, found %v", regST)
			}
			func5 = 13
		} else if regDT.T == TK_REG_SP { //LD SP, REG
			if regST.T != TK_REG {
				return tokens, instr, fmt.Errorf("Source should be R0-R15, found %v", regST)
			}
			func5 = 12
		} else if regST.T == TK_REG_LR { //LD REG, LR
			if regDT.T != TK_REG {
				return tokens, instr, fmt.Errorf("Destination should be R0-R15, found %v", regDT)
			}
			func5 = 11
		} else if regST.T == TK_REG_SP { //LD REG, SP
			if regDT.T != TK_REG {
				return tokens, instr, fmt.Errorf("Destination should be R0-R15, found %v", regDT)
			}
		}
	} else { //Alu instruction
		if err != nil {
			return tokens, instr, err
		}
		func5, err = getFunc5FromALU(tokens[0].Tk)
		if err != nil {
			return tokens, instr, err
		}
	}
	instr.Func5 = func5
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)

	parser.Instructions = append(parser.Instructions, instr)
	parser.CurAddress += 2
	return tokens, instr, nil
}

func ParseLdRegImm(parser *Parser, rule Rule, tokens []Token, tokenpos int, instr Instruction) ([]Token, Instruction, error) {
	var rd, bankd, func2 uint8
	var err error

	ldT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]

	rd, bankd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, instr, err
	}
	if ldT.T == TK_LD_1 {
		func2 = bankd<<1 | 1
	} else {
		func2 = bankd<<1 | 0
	}

	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Imm = int16(immT.ValInt)
	instr.Address = parser.CurAddress
	instr.Func2 = func2
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)

	parser.Instructions = append(parser.Instructions, instr)
	parser.CurAddress += 2
	return tokens, instr, nil
}

// ALU REG, IMM (7 bit Immediate operations)
//
//	func(2 bit) = register bank (0 or 1)
//	func(0-1 bits):
//
// 0. ADD/SUB: Rd = Rd + (signed(IMM))
// 1. SHL/SHR: Rd = << or >> signed(IMM & 31)
// 2. LDI Rd = [PC + signed(IMM)]; 32 bit constant loading, IMM in 32 bit dword (-512 ... +512 bytes)
// 3. DJNZ Rd, PC + signed(IMM); Rd-- if not zero, jump taken
func ParseAluRegImm(parser *Parser, rule Rule, tokens []Token, tokenpos int, instr Instruction) ([]Token, Instruction, error) {
	var rd, bankd, func3 uint8
	var err error

	aluT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]

	rd, bankd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, instr, err
	}
	func3 = bankd << 2

	if immT.T == TK_LABEL {
		instr.Label = immT.Tk
		// labelAddr, ok := parser.Labels[immT.Tk]
		// if !ok {
		// 	return tokens, instr, fmt.Errorf("label '%s' not found", immT.Tk)
		// }
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
		return tokens, instr, fmt.Errorf("Instruction should be one of: ADD, SUB, SHL, SHR, LDI, DJNZ! '%s' not found", aluT.Tk)
	}

	instr.Opcode = rule.Opcode
	instr.Rd = rd

	instr.Address = parser.CurAddress
	instr.Func3 = func3
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)

	parser.Instructions = append(parser.Instructions, instr)
	parser.CurAddress += 2
	return tokens, instr, nil
}

func ParseLabel(parser *Parser, rule Rule, tokens []Token, tokenpos int, instr Instruction) ([]Token, Instruction, error) {
	labelT := tokens[tokenpos]
	parser.Labels[labelT.Tk] = parser.CurAddress

	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+2:]...)
	return tokens, instr, nil
}

const (
	RuleALURegReg = iota
	LdRegImm
	AluRegImm
	Label
)

type Rule struct {
	Type      uint8
	Syntax    [][]uint8
	Opcode    Opcode
	ParseFunc func(parser *Parser, rule Rule, tokens []Token, tokenpos int, instr Instruction) ([]Token, Instruction, error)
}

var aluRegRegSyntax = [][]uint8{{TK_ALU, TK_LD}, {TK_REG, TK_REG_SP, TK_REG_LR}, {TK_COMMA}, {TK_REG, TK_REG_SP, TK_REG_LR}}
var ldRegImmSyntax = [][]uint8{{TK_LD_0, TK_LD_1}, {TK_REG}, {TK_COMMA}, {TK_NUMBER}}
var labelSyntax = [][]uint8{{TK_LABEL}, {TK_COLON}}
var aluRegImmSyntax = [][]uint8{{TK_ALU, TK_LDI, TK_DJNZ}, {TK_REG}, {TK_COMMA}, {TK_NUMBER, TK_LABEL}}

var syntaxRules = []Rule{
	{Type: RuleALURegReg, Syntax: aluRegRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG), ParseFunc: ParseAluRegReg},
	{Type: LdRegImm, Syntax: ldRegImmSyntax, Opcode: GetOpcode(OP_LD_REG_IMM), ParseFunc: ParseLdRegImm},
	{Type: AluRegImm, Syntax: aluRegImmSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: ParseAluRegImm},
	{Type: Label, Syntax: labelSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: ParseLabel},
}
