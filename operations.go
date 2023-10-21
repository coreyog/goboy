package goboy

import (
	"fmt"
	"math/bits"
)

func ld(gb *GameBoy, ext uint8, opcode OpCode, displacement uint8, immediate uint16) {
	gb.debugLnF("operation LD")

	x, z := opcode.GetX(), opcode.GetZ()
	p, q := opcode.GetPQ()

	if x == 0 && z == 1 && q == 0 {
		// opcode ~= 0b00_XX0_001
		switch p {
		case 0:
			gb.setBC(immediate)
		case 1:
			gb.setDE(immediate)
		case 2:
			// LD HL, $9FFF
			gb.setHL(immediate)
		case 3:
			// LD SP, $EFFF
			gb.sp = immediate
		}
	} else if x == 0 && z == 6 {
		// LD r[y], n
		y := opcode.GetY()
		switch y {
		case 0:
			gb.b = byte(immediate)
		case 1:
			gb.c = byte(immediate)
		case 2:
			gb.d = byte(immediate)
		case 3:
			gb.e = byte(immediate)
		case 4:
			gb.h = byte(immediate)
		case 5:
			gb.l = byte(immediate)
		case 6:
			gb.WriteMemory(gb.readHL(), byte(immediate))
		case 7:
			gb.a = byte(immediate)
		}
	} else if x == 3 && z == 2 {
		y := opcode.GetY()
		switch y {
		case 4:
			gb.WriteMemory(0xFF00+uint16(gb.c), gb.a)
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
	case 2:
		memLoc = gb.readHL()
	case 3:
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
		switch y {
		case 0:
			gb.b++
		case 1:
			gb.c++
		case 2:
			gb.d++
		case 3:
			gb.e++
		case 4:
			gb.h++
		case 5:
			gb.l++
		case 6:
			hl := gb.readHL()
			gb.WriteMemory(hl, gb.ReadMemory(hl)+1)
		case 7:
			gb.a++
		}
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

	gb.f &= ^MaskSubtractionFlag // reset
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

	panic(fmt.Sprintf("unexpected Table R lookup value: %d", z))
}
