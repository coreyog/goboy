package ram

import ()

var memory [65536]byte

// memory mapped addresses
const (
	JOYP = 65280
)

// ReadMemory reads a byte from memory at a given address
func ReadMemory(address int) (value byte) {
	return memory[address]
}

// WriteMemory sets the value at a given address in memory
func WriteMemory(address int, value byte) {
	memory[address] = value
}
