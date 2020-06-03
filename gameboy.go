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
