package emul

type BusDevice struct {
	Devices []Device
}

func (b *BusDevice) Add(d Device) {
	b.Devices = append(b.Devices, d)
}

func (b *BusDevice) find(addr uint32) (bool, Device) {
	for _, d := range b.Devices {
		if d.IsFor(addr) {
			return true, d
		}
	}
	return false, nil
}

func (b *BusDevice) Read8(addr uint32) uint8 {
	found, dev := b.find(addr)
	if !found {
		return 0
	}
	return dev.Read8(addr)
}

func (b *BusDevice) Write8(addr uint32, value uint8) {
	found, dev := b.find(addr)
	if found {
		dev.Write8(addr, value)
	}
}

func (b *BusDevice) IsFor(addr uint32) bool {
	return true
}

func (b *BusDevice) Read16(addr uint32) uint16 {
	l := b.Read8(addr)
	h := b.Read8(addr + 1)
	return uint16(h)<<8 | uint16(l)
}

func (b *BusDevice) Read32(addr uint32) uint32 {
	b1 := b.Read8(addr)
	b2 := b.Read8(addr + 1)
	b3 := b.Read8(addr + 2)
	b4 := b.Read8(addr + 3)
	return uint32(b4)<<24 | uint32(b3)<<16 | uint32(b2)<<8 | uint32(b1)
}

func (b *BusDevice) Write32(addr uint32, value uint32) {
	b.Write8(addr, uint8(value&0xFF))
	b.Write8(addr+1, uint8((value>>8)&0xFF))
	b.Write8(addr+2, uint8((value>>16)&0xFF))
	b.Write8(addr+3, uint8((value>>24)&0xFF))
}
