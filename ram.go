package goboy

// memory mapped addresses
const (
	JOYP = 0xFF00 // 65280
)

// ReadMemory reads a byte from memory at a given address, respecting memory mapping
func (gb *GameBoy) ReadMemory(address int) (value byte) {
	return gb.memory[address]
}

// WriteMemory sets the value at a given address in memory, respecting memory mapping
func (gb *GameBoy) WriteMemory(address int, value byte) {
	gb.memory[address] = value
}
