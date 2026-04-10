package emul

type MemoryDevice struct {
	Range  AddressRange
	Data   []uint8
	Offset uint32
}

func NewMemoryDevice(start uint32, end uint32) *MemoryDevice {
	return &MemoryDevice{Range: AddressRange{Start: start, End: end}, Data: make([]uint8, end-start+1)}
}

func (m *MemoryDevice) IsFor(addr uint32) bool {
	return m.Range.IsFor(addr)
}

func (m *MemoryDevice) Read8(addr uint32) uint8 {
	memaddr := addr - m.Offset
	if memaddr < uint32(len(m.Data)) {
		return m.Data[memaddr]
	}
	return 0
}

func (m *MemoryDevice) Write8(addr uint32, value uint8) {
	memaddr := addr - m.Offset
	if memaddr < uint32(len(m.Data)) {
		m.Data[memaddr] = value
	}
}
