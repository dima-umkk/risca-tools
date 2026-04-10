package emul

const (
	UART_RX       = 0xFFFF0001
	UART_RX_READY = 0xFFFF0002
	UART_TX       = 0xFFFF0003
	UART_TX_BUSY  = 0xFFFF0004
)

type UartDevice struct {
	Range AddressRangeSpecific
	RX    []uint8
	TX    []uint8
}

func NewUartDevice() *UartDevice {
	return &UartDevice{Range: AddressRangeSpecific{AddressList: []uint32{UART_RX, UART_RX_READY, UART_TX, UART_TX_BUSY}}}
}

func (u *UartDevice) IsFor(addr uint32) bool {
	return u.Range.IsFor(addr)
}

func (u *UartDevice) Read8(addr uint32) uint8 {
	switch addr {
	case UART_RX:
		if len(u.RX) > 0 {
			d := u.RX[0]
			u.RX = u.RX[1:]
			return d
		}
	case UART_RX_READY:
		if len(u.RX) > 0 {
			return 1
		}
	case UART_TX:
		return 0
	case UART_TX_BUSY:
		return 0
	}
	return 0
}

func (u *UartDevice) Write8(addr uint32, value uint8) {
	if addr == UART_TX {
		u.TX = append(u.TX, value)
	}
}
