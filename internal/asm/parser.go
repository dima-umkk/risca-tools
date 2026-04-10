package asm

import (
	"fmt"

	"github.com/dima-kgd/risca-tools/internal/isa"
)

type Parser struct {
	Instructions []isa.Instruction
	Memory       map[uint32]isa.Instruction
	Labels       map[string]uint32
	Constants    map[string]int32
	CurAddress   uint32
}

func NewParser() *Parser {
	return &Parser{
		Instructions: make([]isa.Instruction, 0, 1024),
		Memory:       make(map[uint32]isa.Instruction),
		Labels:       make(map[string]uint32),
		Constants:    make(map[string]int32),
		CurAddress:   0,
	}
}

func (parser *Parser) ParseLine(line string) (bool, error) {
	tokens, err := Tokenize(line)
	if err != nil { //Error parsing tokens
		return false, fmt.Errorf("Format error: %v", err)
	}
	if len(tokens) == 0 { //Skip line
		return true, nil
	}

	var matched bool

	for len(tokens) > 0 {
		for _, rule := range syntaxRules {
			tokens, matched, err = parser.applyRule(rule, tokens)
			if err != nil {
				return false, fmt.Errorf("Syntax error: %v", err)
			}
			if matched {
				break
			}
		}
		if !matched {
			return false, fmt.Errorf("Syntax error: Unknown syntax: %s", line)
		}
	}
	return false, err
}

func tokenIn(token Token, tokens []uint8) bool {
	for _, t := range tokens {
		if token.T == t {
			return true
		}
	}
	return false
}

func (parser *Parser) addInstruction(instruction isa.Instruction) error {
	if _, ok := parser.Memory[instruction.Address]; ok {
		return fmt.Errorf("Duplicate instruction at address: %d. Instructinon: %s", instruction.Address, instruction)
	}
	parser.Instructions = append(parser.Instructions, instruction)
	parser.Memory[instruction.Address] = instruction
	return nil
}

func (parser *Parser) applyRule(rule Rule, tokens []Token) ([]Token, bool, error) {
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
			tokens, err = rule.ParseFunc(parser, rule, tokens, pos)
			return tokens, true, err
		} else {
			pos++
		}
	}
	return tokens, false, nil
}

func (parser *Parser) ProcessLabels() error {
	for i, instr := range parser.Instructions {
		if instr.Label != "" {
			addr, ok := parser.Labels[instr.Label]
			if !ok {
				return fmt.Errorf("Label not defined: %s", instr.Label)
			}
			offset := (int32(addr) - int32(instr.Address)) >> 1     //offset in instructions (2 bytes per instruction)
			offsetCall := (int32(addr) - int32(instr.Address)) >> 2 //offset in words (4 bytes)
			parser.Instructions[i].Imm = int16(offset)
			switch instr.Opcode.Opc {
			case isa.OP_BRANCH:
				if offset < -64 || offset > 63 {
					return fmt.Errorf("Label too large for instruction! %s, %s", instr.Label, instr)
				}
			case isa.OP_LDI:
				if offset < -512 || offset > 511 {
					return fmt.Errorf("Label too large for instruction! %s, %s", instr.Label, instr)
				}
			case isa.OP_CALL_JUMP_RET:
				if instr.Func == 0 { // CALL Imm(7), Rd
					if addr&0b0000_0011 != 0 {
						return fmt.Errorf("CALL address must be aligned to 4 bytes! %s(0x%08X), %s", instr.Label, addr, instr)
					}
					parser.Instructions[i].Imm = int16(offsetCall)
				}
				if offset < -64 || offset > 63 {
					return fmt.Errorf("Label too large for instruction! %s, %s", instr.Label, instr)
				}
			default:
				return fmt.Errorf("Label not applicable for instruction! %s, %s", instr.Label, instr)
			}
			parser.Memory[instr.Address] = parser.Instructions[i]
		}
	}
	return nil
}
