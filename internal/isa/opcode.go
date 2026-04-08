package isa

type Opcode struct {
	Opc  uint8
	Name string
}

func (o Opcode) String() string {
	return o.Name
}

const (
	OP_ALU_REG_REG   = 0x00
	OP_ALU_IMM       = 0x01
	OP_REG_IMM       = 0x02
	OP_MEM           = 0x03
	OP_BRANCH        = 0x04
	OP_LDI           = 0x05
	OP_CALL_JUMP_RET = 0x06
	OP_INT           = 0x07
	OP_DB            = 0x08 //virtual opcode to define variables in memory
)

var (
	OpAluRegReg   = Opcode{Opc: OP_ALU_REG_REG, Name: "ALU_REG_REG"}
	OpAluImm      = Opcode{Opc: OP_ALU_IMM, Name: "OP_ALU_IMM"}
	OpRegImm      = Opcode{Opc: OP_REG_IMM, Name: "OP_REG_IMM"}
	OpMem         = Opcode{Opc: OP_MEM, Name: "OP_MEM"}
	OpBrach       = Opcode{Opc: OP_BRANCH, Name: "OP_BRANCH"}
	OpLDI         = Opcode{Opc: OP_LDI, Name: "OP_LDI"}
	OpCallJumpRet = Opcode{Opc: OP_CALL_JUMP_RET, Name: "OP_CALL_JUMP_RET"}
	OpINT         = Opcode{Opc: OP_INT, Name: "OP_INT"}
	OpDB          = Opcode{Opc: OP_DB, Name: "DB"}
)

var opcodeMap = map[uint8]Opcode{
	OP_ALU_REG_REG:   OpAluRegReg,
	OP_ALU_IMM:       OpAluImm,
	OP_REG_IMM:       OpRegImm,
	OP_MEM:           OpMem,
	OP_BRANCH:        OpBrach,
	OP_LDI:           OpLDI,
	OP_CALL_JUMP_RET: OpCallJumpRet,
	OP_INT:           OpINT,
	OP_DB:            OpDB,
}

func GetOpcode(opc uint8) Opcode {
	opcode, _ := opcodeMap[opc]
	return opcode
}
