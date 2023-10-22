package goboy

// memory mapped addresses
const (
	JOYP = 0xFF00 // 65280
)

// ReadMemory reads a byte from memory at a given address, respecting memory mapping
func (gb *GameBoy) ReadMemory(address uint16) (value byte) {
	return gb.memory[address]
}

// WriteMemory sets the value at a given address in memory, respecting memory mapping
func (gb *GameBoy) WriteMemory(address uint16, value byte) {
	gb.memory[address] = value
}

func (gb *GameBoy) PushStack(value uint16) {
	msb, lsb := splitBytes(value)
	gb.sp--
	gb.WriteMemory(gb.sp, msb)
	gb.sp--
	gb.WriteMemory(gb.sp, lsb)
}

func (gb *GameBoy) PopStack() (value uint16) {
	lsb := gb.ReadMemory(gb.sp)
	gb.sp++
	msb := gb.ReadMemory(gb.sp)
	gb.sp++

	return mergeBytes(msb, lsb)
}
