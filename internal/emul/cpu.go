package emul

import (
	"github.com/dima-kgd/risca-tools/internal/isa"
)

type Device interface {
	IsFor(addr uint32) bool
	Read8(addr uint32) uint8
	Write8(addr uint32, value uint8)
}

type CPU struct {
	Registers [16]uint32
	PC        uint32
	Bus       BusDevice
}

func (c *CPU) Peek(offset int32) isa.Instruction {
	instrBin := c.Bus.Read16(c.PC + uint32(offset<<1))
	i := isa.Unpack(instrBin)
	i.Address = c.PC + uint32(offset<<1)
	return i
}

func (c *CPU) Step() {
	i := c.Peek(0)
	rd := i.Rd
	rs := i.Rs

	switch i.Opcode.Opc {
	case isa.OP_ALU_REG_REG:
		switch i.Func {
		case 0: // MOV
			c.Registers[rd] = c.Registers[rs]
		case 1: // ADD
			c.Registers[rd] = c.Registers[rd] + c.Registers[rs]
		case 2: // SUB
			c.Registers[rd] = c.Registers[rd] - c.Registers[rs]
		case 3: // AND
			c.Registers[rd] = c.Registers[rd] & c.Registers[rs]
		case 4: // OR
			c.Registers[rd] = c.Registers[rd] | c.Registers[rs]
		case 5: // XOR
			c.Registers[rd] = c.Registers[rd] ^ c.Registers[rs]
		case 6: // NOT
			c.Registers[rd] = ^c.Registers[rs]
		case 7: // MUL
			c.Registers[rd] = c.Registers[rd] * c.Registers[rs]
		}
		c.PC += 2
	case isa.OP_ALU_IMM:
		switch i.Func {
		case 0: // SHL
			c.Registers[rd] = c.Registers[rd] << uint32(i.Imm)
		case 1: // SHR
			c.Registers[rd] = c.Registers[rd] >> uint32(i.Imm)
		case 2: // ADD
			c.Registers[rd] = c.Registers[rd] + uint32(i.Imm)
		case 3: // SUB
			c.Registers[rd] = c.Registers[rd] - uint32(i.Imm)
		}
		c.PC += 2
	case isa.OP_REG_IMM:
		switch i.Func {
		case 0: // MOVI
			c.Registers[rd] = uint32(i.Imm)
		case 1: // MOVH
			c.Registers[rd] = uint32(i.Imm) << 8
		}
		c.PC += 2
	case isa.OP_MEM:
		switch i.Func {
		case 0: // LDB
			c.Registers[rd] = uint32(c.Bus.Read8(c.Registers[rs] + uint32(i.Imm)))
		case 1: // STB
			c.Bus.Write8(c.Registers[rs]+uint32(i.Imm), uint8(c.Registers[rd]))
		case 2: // LDW
			c.Registers[rd] = c.Bus.Read32(c.Registers[rs] + uint32(i.Imm))
		case 3: // STW
			c.Bus.Write32(c.Registers[rs]+uint32(i.Imm), c.Registers[rd])
		}
		c.PC += 2
	case isa.OP_BRANCH:
		switch i.Func {
		case 0: // BEQZ
			if c.Registers[rd] == 0 {
				c.PC = c.PC + uint32(i.Imm)<<1
			} else {
				c.PC += 2
			}
		case 1: // BNEZ
			if c.Registers[rd] != 0 {
				c.PC = c.PC + uint32(i.Imm)<<1
			} else {
				c.PC += 2
			}
		case 2: // BGTZ
			if int32(c.Registers[rd]) > 0 {
				c.PC = c.PC + uint32(i.Imm)<<1
			} else {
				c.PC += 2
			}
		case 3: // BLTZ
			if int32(c.Registers[rd]) < 0 {
				c.PC = c.PC + uint32(i.Imm)<<1
			} else {
				c.PC += 2
			}
		}
	case isa.OP_LDI:
		c.Registers[rd] = c.Bus.Read32(c.PC + uint32(i.Imm)<<1)
		c.PC += 2
	case isa.OP_CALL_JUMP_RET:
		switch i.Func {
		case 0: // CALL Imm, Rd
			c.Registers[rd] = c.PC + 2
			c.PC = c.PC + uint32(i.Imm)<<2
		case 1: // CALL Rd
			c.Registers[14] = c.PC + 2
			c.PC = c.Registers[rd]
		case 2: // RET Rd
			c.PC = c.Registers[rd]
		case 3: // JR Imm
			c.PC = c.PC + uint32(i.Imm)<<1
		}
	}
}
