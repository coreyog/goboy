package goboy

type GameBoy struct {
	memory [65536]byte // 64K RAM
	pc     uint16      // Program Counter - The memory address of the next instruction to fetch
	sp     uint16      // Stack Pointer - the memory address of the top of the stack
	f      uint8       // Flags

	// General purpose registers, f is reserved for the flags register
	a uint8
	b uint8
	c uint8
	d uint8
	e uint8
	h uint8
	l uint8

	tickCount uint64 // Number of elapsed ticks since the start of execution

	romData []byte
}

const (
	MaskSignFlag           uint8 = 0b1000_0000 // sign flag mask - set to the state of the most-significant bit of the Accumulator (bit 7)
	MaskZeroFlag           uint8 = 0b0100_0000 // zero flag mask - set if the result of the operation is zero
	MaskHighCarryFlag      uint8 = 0b0001_0000 // high carry flag mask - set if there was an overflow in the lower half of the result
	MaskParityOverflowFlag uint8 = 0b0000_0100 // parity / overflow mask - set if a math operation overflows for most operations, for some operations this acts as a parity bit (1 if result has an even number of 1s, 0 if result has an odd number of 1s)
	MaskSubtractionFlag    uint8 = 0b0000_0010 // operation flag mask - set if the operation was subtraction
	MaskCarryFlag          uint8 = 0b0000_0001 // carry flag mask - set if there was an overflow in the result
)

func mergeBytes(msb uint8, lsb uint8) uint16 {
	return (uint16(msb) << 8) | uint16(lsb)
}

func splitBytes(value uint16) (msb uint8, lsb uint8) {
	msb = uint8((value & 0xFF00) >> 8)
	lsb = uint8(value & 0x00FF)
	return msb, lsb
}

func (gb *GameBoy) setAF(x uint16) {
	gb.a, gb.f = splitBytes(x)
}

func (gb *GameBoy) readAF() (af uint16) {
	return mergeBytes(gb.a, gb.f)
}

func (gb *GameBoy) setBC(x uint16) {
	gb.b, gb.c = splitBytes(x)
}

func (gb *GameBoy) readBC() (bc uint16) {
	return mergeBytes(gb.b, gb.c)
}

func (gb *GameBoy) setDE(x uint16) {
	gb.d, gb.e = splitBytes(x)
}

func (gb *GameBoy) readDE() (bc uint16) {
	return mergeBytes(gb.d, gb.e)
}

func (gb *GameBoy) setHL(x uint16) {
	gb.h, gb.l = splitBytes(x)
}

func (gb *GameBoy) readHL() (bc uint16) {
	return mergeBytes(gb.h, gb.l)
}
