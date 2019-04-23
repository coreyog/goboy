package cpu

import ()

// Program Counter - The memory address of the next instruction to fetch
var pc uint16

// Stack Pointer - the memory address of the top of the stack
var sp uint8

// Flags
var f uint8

// zero flag big in f - set if the result of the operation is zero
var zfMask uint8 = 128

// operation flag mask - set if the operation was subtraction
var ofMask uint8 = 64

// high carry flag mask - set if there was an overflow in the lower half of the result
var hcfMask uint8 = 32

// carry flag mask - set if there was an overflow in the result
var cfMask uint8 = 16

// General purpose registers
var a uint8
var b uint8
var c uint8
var d uint8
var e uint8
var h uint8
var l uint8

// Number of elapsed ticks since the start of execution
var tickCount uint64

// interrupts
const (
	VBlank    = 0x0040
	LCDStatus = 0x0048
	Timer     = 0x0050
	Serial    = 0x0058
	Joypad    = 0x0060
)

// Fetch retrieves the next instruction to execute
func Fetch() {}

// Decode determines what is needed to execute the fetched instruction
func Decode() {}

// Execute executes the instruction that has been fetched and decoded
func Execute() {}
