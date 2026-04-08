package isa

import (
	"fmt"
)

type Instruction struct {
	Opcode  Opcode
	Rd      uint8
	Rs      uint8
	Func    uint8
	Imm     int16
	Address uint32
	Label   string
}

func (i Instruction) Pack() uint16 {
	var code uint16 = uint16(i.Opcode.Opc)

	//Pack instruction by opcode type
	switch i.Opcode.Opc {
	case OP_ALU_REG_REG:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Rs) << 7
		code |= uint16(i.Func) << 11
	case OP_ALU_IMM:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Func) << 7
		code |= uint16(i.Imm) << 9
	case OP_REG_IMM:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Func) << 7
		code |= uint16(i.Imm) << 8
	case OP_MEM:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Rs) << 7
		code |= uint16(i.Func) << 11
		code |= uint16(i.Imm) << 13
	case OP_BRANCH:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Func) << 7
		code |= uint16(i.Imm) << 9
	case OP_LDI:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Imm) << 7
	case OP_CALL_JUMP_RET:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Func) << 7
		code |= uint16(i.Imm) << 9
	case OP_INT:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Func) << 7
	case OP_DB:
		code = uint16(i.Imm)
	}

	return code
}

func Unpack(inst uint16) Instruction {
	i := Instruction{}
	i.Opcode = GetOpcode(uint8(inst) & 0b0000_0111)
	switch i.Opcode.Opc {
	case OP_ALU_REG_REG:
		i.Rd = getRd(inst)
		i.Rs = getRs(inst)
		i.Func = getFunc3(inst)
	case OP_ALU_IMM:
		i.Rd = getRd(inst)
		i.Func = getFunc2(inst)
		i.Imm = getImm7(inst)
	case OP_REG_IMM:
		i.Rd = getRd(inst)
		i.Func = getFunc1(inst)
		i.Imm = getImm8(inst)
	case OP_MEM:
		i.Rd = getRd(inst)
		i.Rs = getRs(inst)
		i.Func = getFunc22(inst)
		i.Imm = getImm3(inst)
	case OP_BRANCH:
		i.Rd = getRd(inst)
		i.Func = getFunc2(inst)
		i.Imm = getImm7(inst)
	case OP_LDI:
		i.Rd = getRd(inst)
		i.Imm = getImm9(inst)
	case OP_CALL_JUMP_RET:
		i.Rd = getRd(inst)
		i.Func = getFunc2(inst)
		i.Imm = getImm7(inst)
	case OP_INT:
		i.Rd = getRd(inst)
		i.Func = getFunc2(inst)
	}
	return i
}

func getRd(inst uint16) uint8 {
	return uint8(inst>>3) & 0b0000_1111
}

func getRs(inst uint16) uint8 {
	return uint8(inst>>7) & 0b0000_1111
}

func getFunc1(inst uint16) uint8 {
	return uint8(inst>>7) & 0b0000_0001
}

func getFunc2(inst uint16) uint8 {
	return uint8(inst>>7) & 0b0000_0011
}

func getFunc3(inst uint16) uint8 {
	return uint8(inst>>11) & 0b0000_0111
}

func getFunc22(inst uint16) uint8 {
	return uint8(inst>>11) & 0b0000_0011
}

func getImm3(inst uint16) int16 {
	return int16(int8(inst>>8) >> 5)
}

func getImm7(inst uint16) int16 {
	return int16(int8(inst>>8) >> 1)
}

func getImm8(inst uint16) int16 {
	return int16(int8(inst >> 8))
}

func getImm9(inst uint16) int16 {
	return int16(inst) >> 7
}

var mapFuncToAluRegReg = map[uint8]string{
	0: "MOV",
	1: "ADD",
	2: "SUB",
	3: "AND",
	4: "OR",
	5: "XOR",
	6: "NOT",
	7: "MUL",
}

var mapAluToFuncRegReg = map[string]uint8{
	"MOV": 0,
	"ADD": 1,
	"SUB": 2,
	"AND": 3,
	"OR":  4,
	"XOR": 5,
	"NOT": 6,
	"MUL": 7,
}

var mapFuncToAluImm = map[uint8]string{
	0: "SHL",
	1: "SHR",
	2: "ADD",
	3: "SUB",
}

var mapAluImmToFunc = map[string]uint8{
	"SHL": 0,
	"SHR": 1,
	"ADD": 2,
	"SUB": 3,
}

var mapFuncToRegImm = map[uint8]string{
	0: "MOVI",
	1: "MOVH",
}

var mapRegImmToFunc = map[string]uint8{
	"MOVI": 0,
	"MOVH": 1,
}

var mapFuncToBranch = map[uint8]string{
	0: "BEQZ",
	1: "BNEZ",
	2: "BGTZ",
	3: "BLTZ",
}

var mapBranchToFunc = map[string]uint8{
	"BEQZ": 0,
	"BNEZ": 1,
	"BGTZ": 2,
	"BLTZ": 3,
}

func getFuncFromAlu(alu string) (uint8, error) {
	ifunc, exists := mapAluToFuncRegReg[alu]
	if !exists {
		return 0, fmt.Errorf("invalid ALU name")
	}
	return ifunc, nil
}

func getFuncFromRegImm(regimm string) (uint8, error) {
	ifunc, exists := mapRegImmToFunc[regimm]
	if !exists {
		return 0, fmt.Errorf("invalid name for REG IMM")
	}
	return ifunc, nil
}

func getFuncFromAluImm(aluimm string) (uint8, error) {
	ifunc, exists := mapAluImmToFunc[aluimm]
	if !exists {
		return 0, fmt.Errorf("Alu name for ALU IMM should be one of: SHL, SHR, ADD, SUB")
	}
	return ifunc, nil
}

func getFuncFromBranch(branch string) (uint8, error) {
	ifunc, exists := mapBranchToFunc[branch]
	if !exists {
		return 0, fmt.Errorf("invalid BRANCH func name")
	}
	return ifunc, nil
}

func (i Instruction) String() string {
	switch i.Opcode.Opc {
	case OP_ALU_REG_REG:
		if name, exists := mapFuncToAluRegReg[i.Func]; exists {
			return fmt.Sprintf("%s\tR%d, R%d", name, i.Rd, i.Rs)
		}
	case OP_ALU_IMM:
		if name, exists := mapFuncToAluImm[i.Func]; exists {
			return fmt.Sprintf("%s\tR%d, 0x%02X(%d)", name, i.Rd, uint8(i.Imm), i.Imm)
		}
	case OP_REG_IMM:
		if name, exists := mapFuncToRegImm[i.Func]; exists {
			return fmt.Sprintf("%s\tR%d, 0x%02X(%d)", name, i.Rd, uint8(i.Imm), i.Imm)
		}
	case OP_MEM:
		switch i.Func {
		case 0:
			return fmt.Sprintf("LDB\tR%d, [R%d+0x%02X(%d)]", i.Rd, i.Rs, uint8(i.Imm), i.Imm)
		case 1:
			return fmt.Sprintf("STB\t[R%d+0x%02X(%d)], R%d", i.Rs, uint8(i.Imm), i.Imm, i.Rd)
		case 2:
			return fmt.Sprintf("LDW\tR%d, [R%d+0x%02X(%d)]", i.Rd, i.Rs, uint8(i.Imm), i.Imm)
		case 3:
			return fmt.Sprintf("STW\t[R%d+0x%02X(%d)], R%d", i.Rs, uint8(i.Imm), i.Imm, i.Rd)
		}
	case OP_BRANCH:
		if name, exists := mapFuncToBranch[i.Func]; exists {
			return fmt.Sprintf("%s\tR%d, 0x%02X(%d)", name, i.Rd, uint8(i.Imm), i.Imm)
		}
	case OP_LDI:
		offset := i.Imm
		address := uint32(int32(i.Address) + int32(offset<<1))
		return fmt.Sprintf("LDI\tR%d, [0x%02X(%d)] -> %08X", i.Rd, uint8(i.Imm), i.Imm, address)
	case OP_CALL_JUMP_RET:
		offset := i.Imm
		instAddress := uint32(int32(i.Address) + int32(offset<<1))
		callAddress := uint32(int32(i.Address) + int32(offset<<2))
		switch i.Func {
		case 0:
			return fmt.Sprintf("CALL\t0x%02X(%d), R%d -> %08X", uint8(i.Imm), i.Imm, i.Rd, callAddress)
		case 1:
			return fmt.Sprintf("CALL\tR%d", i.Rd)
		case 2:
			if i.Rd == 14 { // Standard link register for RET - possible ret meaning
				return fmt.Sprintf("RET\tR%d", i.Rd)
			} else { // Possible JMP meaning
				return fmt.Sprintf("JMP\tR%d", i.Rd)
			}
		case 3:
			return fmt.Sprintf("JMP\t0x%02X(%d) -> %08X", uint8(i.Imm), i.Imm, instAddress)
		}
	case OP_INT:
		switch i.Func {
		case 0:
			return fmt.Sprintf("INT\tR%d", i.Rd)
		case 1:
			return "RETI"
		case 2:
			return fmt.Sprintf("MOV\tR%d, STS", i.Rd)
		case 3:
			return fmt.Sprintf("MOV\tSTS, R%d", i.Rd)
		}
	}
	return ""
}
