package emul

type AddressRange struct {
	Start uint32
	End   uint32
}

func (ar AddressRange) IsFor(addr uint32) bool {
	return addr <= ar.End && addr >= ar.Start
}

type AddressRangeSpecific struct {
	AddressList []uint32
}

func (ar AddressRangeSpecific) IsFor(addr uint32) bool {
	for _, a := range ar.AddressList {
		if a == addr {
			return true
		}
	}
	return false
}
