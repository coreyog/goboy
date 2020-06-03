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
	// zero flag mask - set if the result of the operation is zero
	MaskZeroFlag uint8 = 128

	// operation flag mask - set if the operation was subtraction
	MaskOpFlag uint8 = 64

	// high carry flag mask - set if there was an overflow in the lower half of the result
	MaskHighCarryFlag uint8 = 32

	// carry flag mask - set if there was an overflow in the result
	MaskCarryFlag uint8 = 16
)
