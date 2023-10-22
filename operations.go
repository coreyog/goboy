package goboy

import (
	"fmt"
	"math/bits"

	"github.com/pkg/errors"
)

func ld(gb *GameBoy, ext uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation LD")

	x, z := opcode.GetX(), opcode.GetZ()
	p, q := opcode.GetPQ()

	if x == 0 {
		if z == 1 && q == 0 {
			tableRPWrite(gb, p, immediate)
		} else if z == 6 {
			// LD r[y], n
			y := opcode.GetY()
			tableRWrite(gb, y, byte(immediate))
		}
	} else if x == 1 {
		value := tableRRead(gb, z)

		y := opcode.GetY()
		tableRWrite(gb, y, value)
	} else if x == 3 {
		y := opcode.GetY()
		if z == 0 {
			switch y {
			case 4:
				gb.WriteMemory(0xFF00+immediate, gb.a)
			}
		} else if z == 2 {
			switch y {
			case 4:
				gb.WriteMemory(0xFF00+uint16(gb.c), gb.a)
			}
		}
	}

	// no flags to change
}

// ldid covers both LDD and LDI
func ldid(gb *GameBoy, ext uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation LDD/LDI")

	p, q := opcode.GetPQ()

	var memLoc uint16
	switch p {
	case 0:
		memLoc = gb.readBC()
	case 1:
		memLoc = gb.readDE()
	case 2, 3:
		memLoc = gb.readHL()
	}

	if q == 1 {
		gb.a = gb.ReadMemory(memLoc)
	} else {
		gb.WriteMemory(memLoc, gb.a)
	}

	if p == 2 {
		gb.setHL(gb.readHL() + 1)
	} else if p == 3 {
		gb.setHL(gb.readHL() - 1)
	}
}

func inc(gb *GameBoy, prefix uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation INC")

	y, z := opcode.GetY(), opcode.GetZ()

	if z == 4 {
		tableRWrite(gb, y, tableRRead(gb, y)+1)
	}
}

func xor(gb *GameBoy, prefix uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation XOR")

	z := opcode.GetZ()

	value := tableRRead(gb, z)

	gb.a ^= value

	// set sign flag
	if gb.a&0b1000_0000 != 0 {
		gb.f |= MaskSignFlag
	} else {
		gb.f &= ^MaskSignFlag
	}

	// set zero flag
	if gb.a == 0 {
		gb.f |= MaskZeroFlag
	} else {
		gb.f &= ^MaskZeroFlag
	}

	// set parity flag
	if bits.OnesCount8(gb.a)%2 == 0 {
		gb.f |= MaskParityOverflowFlag
	} else {
		gb.f &= ^MaskParityOverflowFlag
	}

	// clear carry, high carry, and subtraction flags
	gb.f &= ^(MaskHighCarryFlag | MaskSubtractionFlag | MaskCarryFlag)
}

func bit(gb *GameBoy, ext uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation BIT")

	y, z := opcode.GetY(), opcode.GetZ()
	r := tableRRead(gb, z)

	b := (r >> y) & 1

	if b == 0 {
		gb.f |= MaskZeroFlag
	} else {
		gb.f &= ^MaskZeroFlag
	}

	gb.f &= ^MaskSubtractionFlag // clear
	gb.f |= MaskHighCarryFlag    // set
}

func jr(gb *GameBoy, ext uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation JR")

	y := opcode.GetY()

	var jump bool

	if y == 3 {
		jump = true
	} else {
		// y = 4..7
		switch y {
		case 4: // NZ
			jump = gb.f&MaskZeroFlag == 0
		case 5: // Z
			jump = gb.f&MaskZeroFlag != 0
		case 6: // NC
			jump = gb.f&MaskCarryFlag == 0
		case 7: // C
			jump = gb.f&MaskCarryFlag != 0
		}
	}

	if jump {
		signedEnlargedDisplacement := int16(int8(displacement))
		gb.pc = uint16(int16(gb.pc) + signedEnlargedDisplacement)
	}
}

func push(gb *GameBoy, ext uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation PUSH")

	p, _ := opcode.GetPQ()

	rp2 := tableRP2Read(gb, p)

	gb.PushStack(rp2)
}

func call(gb *GameBoy, ext uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation CALL")

	// the call instruction is 3 bytes long and the PC still points to it as the
	// current instruction
	nextPC := gb.pc + 3

	gb.PushStack(nextPC)
	gb.pc = immediate - 3 // -3 because the PC is incremented by 3 after the instruction is executed
}

func rl(gb *GameBoy, ext uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation RL")

	z := opcode.GetZ()

	// var bit bool
	r := tableRRead(gb, z)
	bit := r&1<<7 != 0
	tableRWrite(gb, z, r<<1)

	if bit {
		// set carry
		gb.f |= MaskCarryFlag
	} else {
		// clear carry
		gb.f &= ^MaskCarryFlag
	}
}

func tableRRead(gb *GameBoy, z uint8) (value uint8) {
	// table r from https://gb-archive.github.io/salvage/decoding_gbz80_opcodes/Decoding%20Gamboy%20Z80%20Opcodes.html
	switch z {
	case 0:
		return gb.b
	case 1:
		return gb.c
	case 2:
		return gb.d
	case 3:
		return gb.e
	case 4:
		return gb.h
	case 5:
		return gb.l
	case 6:
		return gb.ReadMemory(gb.readHL())
	case 7:
		return gb.a
	}

	panic(errors.New(fmt.Sprintf("unexpected Table R lookup value: %d", z)))
}

func tableRWrite(gb *GameBoy, z uint8, value uint8) {
	// table r from https://gb-archive.github.io/salvage/decoding_gbz80_opcodes/Decoding%20Gamboy%20Z80%20Opcodes.html
	switch z {
	case 0:
		gb.b = value
	case 1:
		gb.c = value
	case 2:
		gb.d = value
	case 3:
		gb.e = value
	case 4:
		gb.h = value
	case 5:
		gb.l = value
	case 6:
		gb.WriteMemory(gb.readHL(), value)
	case 7:
		gb.a = value
	default:
		panic(errors.New(fmt.Sprintf("unexpected Table R lookup value: %d", z)))
	}
}

func tableRPWrite(gb *GameBoy, p uint8, value uint16) {
	switch p {
	case 0:
		gb.setBC(value)
	case 1:
		gb.setDE(value)
	case 2:
		// LD HL, $9FFF
		gb.setHL(value)
	case 3:
		// LD SP, $EFFF
		gb.sp = value
	}
}

func tableRP2Read(gb *GameBoy, p uint8) (value uint16) {
	switch p {
	case 0:
		return gb.readBC()
	case 1:
		return gb.readDE()
	case 2:
		return gb.readHL()
	case 3:
		return gb.readAF()
	}

	panic(errors.New(fmt.Sprintf("unexpected Table RP2 lookup value: %d", p)))
}

func tableRP2Write(gb *GameBoy, p uint8, value uint16) {
	switch p {
	case 0:
		gb.setBC(value)
	case 1:
		gb.setDE(value)
	case 2:
		gb.setHL(value)
	case 3:
		gb.setAF(value)
	default:
		panic(errors.New(fmt.Sprintf("unexpected Table RP2 lookup value: %d", p)))
	}
}
