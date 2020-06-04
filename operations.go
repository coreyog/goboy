package goboy

import (
	"fmt"
	"math/bits"
)

func ld(gb *GameBoy, ext byte, opcode byte, displacement byte, immediate uint16) {
	x, y, z := OpCode(opcode).Split()
	p, q := splitY(y)

	if x == 0 && z == 1 && q == 0 {
		// opcode ~= 0b00_XX0_001
		switch p {
		case 0:
			gb.setBC(immediate)
		case 1:
			gb.setDE(immediate)
		case 2:
			gb.setHL(immediate)
		case 3:
			fmt.Println("LD SP, $EFFF")
			gb.sp = immediate
		}
	}

	// no flags to set
}

func xor(gb *GameBoy, prefix byte, opcode byte, displacement byte, immediate uint16) {
	// XOR A
	gb.a ^= gb.a
	fmt.Println("XOR A")

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
