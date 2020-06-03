package goboy

import (
	"math/bits"
)

func ldsp(gb *GameBoy, displacement byte, immediate uint16) {
	// LD SP, $FFFF
	gb.sp = immediate
	// no flags to set
}

func xora(gb *GameBoy, displacement byte, immediate uint16) {
	// XOR A
	gb.a ^= gb.a

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
