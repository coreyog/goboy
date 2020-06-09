package goboy

import ()

type OpCode uint8

const (
	NOP OpCode = iota
	STOP
	DJNZ
	JR
	LD
	LDD
	LDI
	ADD
	INC
	DEC
	RLCA
	RRCA
	RLA
	RRA
	DAA
	CPL
	SCF
	CCF
	HALT
	ALU
	RET
	POP
	JP
	EXX
	OUT
	IN
	DI
	EI
	CALL
	PUSH
	RST
	ROT
	BIT
	RES
	SET
	SBC
	ADC
	NEG
	RETN
	RETI
	IM
	RRD
	RLD
	BLI
)

type OpBytes struct {
	Code            OpCode
	HasDisplacement bool
	ImmediateSize   uint8 // 0, 1, 2
	Operation       func(gameboy *GameBoy, prefix uint8, opcode OpCode, displacement uint8, immediate uint16)
}

func (code OpCode) GetX() (x uint8) {
	c := uint8(code)
	return c >> 6
}

func (code OpCode) GetY() (y uint8) {
	c := uint8(code)
	return (c >> 3) & 0b0111
}

func (code OpCode) GetZ() (z uint8) {
	c := uint8(code)
	return c & 0b0111
}

func (code OpCode) GetPQ() (p uint8, q uint8) {
	c := uint8(code)
	p = (c & 0b00_110_000) >> 4
	q = (c & 0b00_001_000) >> 3

	return p, q
}

var unprefixed = map[byte]OpBytes{
	//XX_YYY_ZZZ
	//   PPQ
	0b00_000_000: {NOP, false, 0, nil},
	0b00_001_000: {LD, false, 2, nil},
	0b00_010_000: {STOP, false, 0, nil},
	0b00_011_000: {JR, true, 0, jr},
	0b00_100_000: {JR, true, 0, jr},
	0b00_101_000: {JR, true, 0, jr},
	0b00_110_000: {JR, true, 0, jr},
	0b00_111_000: {JR, true, 0, jr},

	0b00_000_001: {LD, false, 2, nil},
	0b00_010_001: {LD, false, 2, nil},
	0b00_100_001: {LD, false, 2, ld},
	0b00_110_001: {LD, false, 2, ld},

	0b00_001_001: {ADD, false, 0, nil},
	0b00_011_001: {ADD, false, 0, nil},
	0b00_101_001: {ADD, false, 0, nil},
	0b00_111_001: {ADD, false, 0, nil},

	0b00_000_010: {LD, false, 0, ldid},
	0b00_010_010: {LD, false, 0, ldid},
	0b00_100_010: {LDI, false, 0, ldid},
	0b00_110_010: {LDD, false, 0, ldid},

	0b00_001_010: {LD, false, 0, ldid},
	0b00_011_010: {LD, false, 0, ldid},
	0b00_101_010: {LDI, false, 0, ldid},
	0b00_111_010: {LDD, false, 0, ldid},

	0b00_000_011: {INC, false, 0, nil},
	0b00_010_011: {INC, false, 0, nil},
	0b00_100_011: {INC, false, 0, nil},
	0b00_110_011: {INC, false, 0, nil},

	0b00_001_011: {DEC, false, 0, nil},
	0b00_011_011: {DEC, false, 0, nil},
	0b00_101_011: {DEC, false, 0, nil},
	0b00_111_011: {DEC, false, 0, nil},

	0b00_000_100: {INC, false, 0, nil},
	0b00_001_100: {INC, false, 0, nil},
	0b00_010_100: {INC, false, 0, nil},
	0b00_011_100: {INC, false, 0, nil},
	0b00_100_100: {INC, false, 0, nil},
	0b00_101_100: {INC, false, 0, nil},
	0b00_110_100: {INC, false, 0, nil},
	0b00_111_100: {INC, false, 0, nil},
	//XX_YYY_ZZZ
	//   PPQ
	0b00_000_101: {DEC, false, 0, nil},
	0b00_001_101: {DEC, false, 0, nil},
	0b00_010_101: {DEC, false, 0, nil},
	0b00_011_101: {DEC, false, 0, nil},
	0b00_100_101: {DEC, false, 0, nil},
	0b00_101_101: {DEC, false, 0, nil},
	0b00_110_101: {DEC, false, 0, nil},
	0b00_111_101: {DEC, false, 0, nil},

	0b00_000_110: {LD, false, 1, nil},
	0b00_001_110: {LD, false, 1, nil},
	0b00_010_110: {LD, false, 1, nil},
	0b00_011_110: {LD, false, 1, nil},
	0b00_100_110: {LD, false, 1, nil},
	0b00_101_110: {LD, false, 1, nil},
	0b00_110_110: {LD, false, 1, nil},
	0b00_111_110: {LD, false, 1, nil},

	0b00_000_111: {RLCA, false, 0, nil},
	0b00_001_111: {RRCA, false, 0, nil},
	0b00_010_111: {RLA, false, 0, nil},
	0b00_011_111: {RRA, false, 0, nil},
	0b00_100_111: {DAA, false, 0, nil},
	0b00_101_111: {CPL, false, 0, nil},
	0b00_110_111: {SCF, false, 0, nil},
	0b00_111_111: {CCF, false, 0, nil},

	0b01_000_000: {LD, false, 0, nil},
	0b01_000_001: {LD, false, 0, nil},
	0b01_000_010: {LD, false, 0, nil},
	0b01_000_011: {LD, false, 0, nil},
	0b01_000_100: {LD, false, 0, nil},
	0b01_000_101: {LD, false, 0, nil},
	0b01_000_110: {LD, false, 0, nil},
	0b01_000_111: {LD, false, 0, nil},
	0b01_001_000: {LD, false, 0, nil},
	0b01_001_001: {LD, false, 0, nil},
	0b01_001_010: {LD, false, 0, nil},
	0b01_001_011: {LD, false, 0, nil},
	0b01_001_100: {LD, false, 0, nil},
	0b01_001_101: {LD, false, 0, nil},
	0b01_001_110: {LD, false, 0, nil},
	0b01_001_111: {LD, false, 0, nil},
	0b01_010_000: {LD, false, 0, nil},
	0b01_010_001: {LD, false, 0, nil},
	0b01_010_010: {LD, false, 0, nil},
	0b01_010_011: {LD, false, 0, nil},
	0b01_010_100: {LD, false, 0, nil},
	0b01_010_101: {LD, false, 0, nil},
	0b01_010_110: {LD, false, 0, nil},
	0b01_010_111: {LD, false, 0, nil},
	0b01_011_000: {LD, false, 0, nil},
	0b01_011_001: {LD, false, 0, nil},
	0b01_011_010: {LD, false, 0, nil},
	0b01_011_011: {LD, false, 0, nil},
	0b01_011_100: {LD, false, 0, nil},
	0b01_011_101: {LD, false, 0, nil},
	0b01_011_110: {LD, false, 0, nil},
	0b01_011_111: {LD, false, 0, nil},
	0b01_100_000: {LD, false, 0, nil},
	0b01_100_001: {LD, false, 0, nil},
	0b01_100_010: {LD, false, 0, nil},
	0b01_100_011: {LD, false, 0, nil},
	0b01_100_100: {LD, false, 0, nil},
	0b01_100_101: {LD, false, 0, nil},
	0b01_100_110: {LD, false, 0, nil},
	0b01_100_111: {LD, false, 0, nil},
	0b01_101_000: {LD, false, 0, nil},
	0b01_101_001: {LD, false, 0, nil},
	0b01_101_010: {LD, false, 0, nil},
	0b01_101_011: {LD, false, 0, nil},
	0b01_101_100: {LD, false, 0, nil},
	0b01_101_101: {LD, false, 0, nil},
	0b01_101_110: {LD, false, 0, nil},
	0b01_101_111: {LD, false, 0, nil},
	0b01_110_000: {LD, false, 0, nil},
	0b01_110_001: {LD, false, 0, nil},
	0b01_110_010: {LD, false, 0, nil},
	0b01_110_011: {LD, false, 0, nil},
	0b01_110_100: {LD, false, 0, nil},
	0b01_110_101: {LD, false, 0, nil},
	0b01_110_110: {HALT, false, 0, nil},
	0b01_110_111: {LD, false, 0, nil},
	0b01_111_000: {LD, false, 0, nil},
	0b01_111_001: {LD, false, 0, nil},
	0b01_111_010: {LD, false, 0, nil},
	0b01_111_011: {LD, false, 0, nil},
	0b01_111_100: {LD, false, 0, nil},
	0b01_111_101: {LD, false, 0, nil},
	0b01_111_110: {LD, false, 0, nil},
	0b01_111_111: {LD, false, 0, nil},
	//XX_YYY_ZZZ
	//   PPQ
	0b10_000_000: {ALU, false, 0, nil},
	0b10_000_001: {ALU, false, 0, nil},
	0b10_000_010: {ALU, false, 0, nil},
	0b10_000_011: {ALU, false, 0, nil},
	0b10_000_100: {ALU, false, 0, nil},
	0b10_000_101: {ALU, false, 0, nil},
	0b10_000_110: {ALU, false, 0, nil},
	0b10_000_111: {ALU, false, 0, nil},
	0b10_001_000: {ALU, false, 0, nil},
	0b10_001_001: {ALU, false, 0, nil},
	0b10_001_010: {ALU, false, 0, nil},
	0b10_001_011: {ALU, false, 0, nil},
	0b10_001_100: {ALU, false, 0, nil},
	0b10_001_101: {ALU, false, 0, nil},
	0b10_001_110: {ALU, false, 0, nil},
	0b10_001_111: {ALU, false, 0, nil},
	0b10_010_000: {ALU, false, 0, nil},
	0b10_010_001: {ALU, false, 0, nil},
	0b10_010_010: {ALU, false, 0, nil},
	0b10_010_011: {ALU, false, 0, nil},
	0b10_010_100: {ALU, false, 0, nil},
	0b10_010_101: {ALU, false, 0, nil},
	0b10_010_110: {ALU, false, 0, nil},
	0b10_010_111: {ALU, false, 0, nil},
	0b10_011_000: {ALU, false, 0, nil},
	0b10_011_001: {ALU, false, 0, nil},
	0b10_011_010: {ALU, false, 0, nil},
	0b10_011_011: {ALU, false, 0, nil},
	0b10_011_100: {ALU, false, 0, nil},
	0b10_011_101: {ALU, false, 0, nil},
	0b10_011_110: {ALU, false, 0, nil},
	0b10_011_111: {ALU, false, 0, nil},
	0b10_100_000: {ALU, false, 0, nil},
	0b10_100_001: {ALU, false, 0, nil},
	0b10_100_010: {ALU, false, 0, nil},
	0b10_100_011: {ALU, false, 0, nil},
	0b10_100_100: {ALU, false, 0, nil},
	0b10_100_101: {ALU, false, 0, nil},
	0b10_100_110: {ALU, false, 0, nil},
	0b10_100_111: {ALU, false, 0, nil},
	0b10_101_000: {ALU, false, 0, nil},
	0b10_101_001: {ALU, false, 0, nil},
	0b10_101_010: {ALU, false, 0, nil},
	0b10_101_011: {ALU, false, 0, nil},
	0b10_101_100: {ALU, false, 0, nil},
	0b10_101_101: {ALU, false, 0, nil},
	0b10_101_110: {ALU, false, 0, nil},
	0b10_101_111: {ALU, false, 0, xor},
	0b10_110_000: {ALU, false, 0, nil},
	0b10_110_001: {ALU, false, 0, nil},
	0b10_110_010: {ALU, false, 0, nil},
	0b10_110_011: {ALU, false, 0, nil},
	0b10_110_100: {ALU, false, 0, nil},
	0b10_110_101: {ALU, false, 0, nil},
	0b10_110_110: {ALU, false, 0, nil},
	0b10_110_111: {ALU, false, 0, nil},
	0b10_111_000: {ALU, false, 0, nil},
	0b10_111_001: {ALU, false, 0, nil},
	0b10_111_010: {ALU, false, 0, nil},
	0b10_111_011: {ALU, false, 0, nil},
	0b10_111_100: {ALU, false, 0, nil},
	0b10_111_101: {ALU, false, 0, nil},
	0b10_111_110: {ALU, false, 0, nil},
	0b10_111_111: {ALU, false, 0, nil},
	//XX_YYY_ZZZ
	//   PPQ
	0b11_000_000: {RET, false, 0, nil},
	0b11_001_000: {RET, false, 0, nil},
	0b11_010_000: {RET, false, 0, nil},
	0b11_011_000: {RET, false, 0, nil},
	0b11_100_000: {LD, false, 2, nil},
	0b11_101_000: {ADD, true, 0, nil},
	0b11_110_000: {LD, false, 1, nil},
	0b11_111_000: {LD, true, 0, nil},

	0b11_000_001: {POP, false, 0, nil},
	0b11_010_001: {POP, false, 0, nil},
	0b11_100_001: {POP, false, 0, nil},
	0b11_110_001: {POP, false, 0, nil},

	0b11_001_001: {RET, false, 0, nil},
	0b11_011_001: {RETI, false, 0, nil},
	0b11_101_001: {JP, false, 0, nil},
	0b11_111_001: {LD, false, 0, nil},

	0b11_000_010: {JP, false, 2, nil},
	0b11_001_010: {JP, false, 2, nil},
	0b11_010_010: {JP, false, 2, nil},
	0b11_011_010: {JP, false, 2, nil},
	0b11_100_010: {LD, false, 0, nil},
	0b11_101_010: {LD, false, 2, nil},
	0b11_110_010: {LD, false, 0, nil},
	0b11_111_010: {LD, false, 2, nil},

	0b11_000_011: {JP, false, 2, nil},
	// gap for CB prefix and removed instructions
	0b11_110_011: {DI, false, 0, nil},
	0b11_111_011: {EI, false, 0, nil},

	0b11_000_100: {CALL, false, 2, nil},
	0b11_001_100: {CALL, false, 2, nil},
	0b11_010_100: {CALL, false, 2, nil},
	0b11_011_100: {CALL, false, 2, nil},
	// gap for removed instructions

	0b11_000_101: {PUSH, false, 0, nil},
	0b11_010_101: {PUSH, false, 0, nil},
	0b11_100_101: {PUSH, false, 0, nil},
	0b11_110_101: {PUSH, false, 0, nil},
	//XX_YYY_ZZZ
	//   PPQ
	0b11_001_101: {CALL, false, 2, nil},
	// gap for removed instructions

	0b11_000_110: {ALU, false, 1, nil},
	0b11_001_110: {ALU, false, 1, nil},
	0b11_010_110: {ALU, false, 1, nil},
	0b11_011_110: {ALU, false, 1, nil},
	0b11_100_110: {ALU, false, 1, nil},
	0b11_101_110: {ALU, false, 1, nil},
	0b11_110_110: {ALU, false, 1, nil},
	0b11_111_110: {ALU, false, 1, nil},

	0b11_000_111: {RST, false, 0, nil},
	0b11_001_111: {RST, false, 0, nil},
	0b11_010_111: {RST, false, 0, nil},
	0b11_011_111: {RST, false, 0, nil},
	0b11_100_111: {RST, false, 0, nil},
	0b11_101_111: {RST, false, 0, nil},
	0b11_110_111: {RST, false, 0, nil},
	0b11_111_111: {RST, false, 0, nil},
}

var cb = map[byte]OpBytes{
	//XX_YYY_ZZZ
	//   PPQ
	0b00_000_000: {ROT, false, 0, nil},
	0b00_000_001: {ROT, false, 0, nil},
	0b00_000_010: {ROT, false, 0, nil},
	0b00_000_011: {ROT, false, 0, nil},
	0b00_000_100: {ROT, false, 0, nil},
	0b00_000_101: {ROT, false, 0, nil},
	0b00_000_110: {ROT, false, 0, nil},
	0b00_000_111: {ROT, false, 0, nil},
	0b00_001_000: {ROT, false, 0, nil},
	0b00_001_001: {ROT, false, 0, nil},
	0b00_001_010: {ROT, false, 0, nil},
	0b00_001_011: {ROT, false, 0, nil},
	0b00_001_100: {ROT, false, 0, nil},
	0b00_001_101: {ROT, false, 0, nil},
	0b00_001_110: {ROT, false, 0, nil},
	0b00_001_111: {ROT, false, 0, nil},
	0b00_010_000: {ROT, false, 0, nil},
	0b00_010_001: {ROT, false, 0, nil},
	0b00_010_010: {ROT, false, 0, nil},
	0b00_010_011: {ROT, false, 0, nil},
	0b00_010_100: {ROT, false, 0, nil},
	0b00_010_101: {ROT, false, 0, nil},
	0b00_010_110: {ROT, false, 0, nil},
	0b00_010_111: {ROT, false, 0, nil},
	0b00_011_000: {ROT, false, 0, nil},
	0b00_011_001: {ROT, false, 0, nil},
	0b00_011_010: {ROT, false, 0, nil},
	0b00_011_011: {ROT, false, 0, nil},
	0b00_011_100: {ROT, false, 0, nil},
	0b00_011_101: {ROT, false, 0, nil},
	0b00_011_110: {ROT, false, 0, nil},
	0b00_011_111: {ROT, false, 0, nil},
	0b00_100_000: {ROT, false, 0, nil},
	0b00_100_001: {ROT, false, 0, nil},
	0b00_100_010: {ROT, false, 0, nil},
	0b00_100_011: {ROT, false, 0, nil},
	0b00_100_100: {ROT, false, 0, nil},
	0b00_100_101: {ROT, false, 0, nil},
	0b00_100_110: {ROT, false, 0, nil},
	0b00_100_111: {ROT, false, 0, nil},
	0b00_101_000: {ROT, false, 0, nil},
	0b00_101_001: {ROT, false, 0, nil},
	0b00_101_010: {ROT, false, 0, nil},
	0b00_101_011: {ROT, false, 0, nil},
	0b00_101_100: {ROT, false, 0, nil},
	0b00_101_101: {ROT, false, 0, nil},
	0b00_101_110: {ROT, false, 0, nil},
	0b00_101_111: {ROT, false, 0, nil},
	0b00_110_000: {ROT, false, 0, nil},
	0b00_110_001: {ROT, false, 0, nil},
	0b00_110_010: {ROT, false, 0, nil},
	0b00_110_011: {ROT, false, 0, nil},
	0b00_110_100: {ROT, false, 0, nil},
	0b00_110_101: {ROT, false, 0, nil},
	0b00_110_110: {ROT, false, 0, nil},
	0b00_110_111: {ROT, false, 0, nil},
	0b00_111_000: {ROT, false, 0, nil},
	0b00_111_001: {ROT, false, 0, nil},
	0b00_111_010: {ROT, false, 0, nil},
	0b00_111_011: {ROT, false, 0, nil},
	0b00_111_100: {ROT, false, 0, nil},
	0b00_111_101: {ROT, false, 0, nil},
	0b00_111_110: {ROT, false, 0, nil},
	0b00_111_111: {ROT, false, 0, nil},
	//XX_YYY_ZZZ
	//   PPQ
	0b01_000_000: {BIT, false, 0, bit},
	0b01_000_001: {BIT, false, 0, bit},
	0b01_000_010: {BIT, false, 0, bit},
	0b01_000_011: {BIT, false, 0, bit},
	0b01_000_100: {BIT, false, 0, bit},
	0b01_000_101: {BIT, false, 0, bit},
	0b01_000_110: {BIT, false, 0, bit},
	0b01_000_111: {BIT, false, 0, bit},
	0b01_001_000: {BIT, false, 0, bit},
	0b01_001_001: {BIT, false, 0, bit},
	0b01_001_010: {BIT, false, 0, bit},
	0b01_001_011: {BIT, false, 0, bit},
	0b01_001_100: {BIT, false, 0, bit},
	0b01_001_101: {BIT, false, 0, bit},
	0b01_001_110: {BIT, false, 0, bit},
	0b01_001_111: {BIT, false, 0, bit},
	0b01_010_000: {BIT, false, 0, bit},
	0b01_010_001: {BIT, false, 0, bit},
	0b01_010_010: {BIT, false, 0, bit},
	0b01_010_011: {BIT, false, 0, bit},
	0b01_010_100: {BIT, false, 0, bit},
	0b01_010_101: {BIT, false, 0, bit},
	0b01_010_110: {BIT, false, 0, bit},
	0b01_010_111: {BIT, false, 0, bit},
	0b01_011_000: {BIT, false, 0, bit},
	0b01_011_001: {BIT, false, 0, bit},
	0b01_011_010: {BIT, false, 0, bit},
	0b01_011_011: {BIT, false, 0, bit},
	0b01_011_100: {BIT, false, 0, bit},
	0b01_011_101: {BIT, false, 0, bit},
	0b01_011_110: {BIT, false, 0, bit},
	0b01_011_111: {BIT, false, 0, bit},
	0b01_100_000: {BIT, false, 0, bit},
	0b01_100_001: {BIT, false, 0, bit},
	0b01_100_010: {BIT, false, 0, bit},
	0b01_100_011: {BIT, false, 0, bit},
	0b01_100_100: {BIT, false, 0, bit},
	0b01_100_101: {BIT, false, 0, bit},
	0b01_100_110: {BIT, false, 0, bit},
	0b01_100_111: {BIT, false, 0, bit},
	0b01_101_000: {BIT, false, 0, bit},
	0b01_101_001: {BIT, false, 0, bit},
	0b01_101_010: {BIT, false, 0, bit},
	0b01_101_011: {BIT, false, 0, bit},
	0b01_101_100: {BIT, false, 0, bit},
	0b01_101_101: {BIT, false, 0, bit},
	0b01_101_110: {BIT, false, 0, bit},
	0b01_101_111: {BIT, false, 0, bit},
	0b01_110_000: {BIT, false, 0, bit},
	0b01_110_001: {BIT, false, 0, bit},
	0b01_110_010: {BIT, false, 0, bit},
	0b01_110_011: {BIT, false, 0, bit},
	0b01_110_100: {BIT, false, 0, bit},
	0b01_110_101: {BIT, false, 0, bit},
	0b01_110_110: {BIT, false, 0, bit},
	0b01_110_111: {BIT, false, 0, bit},
	0b01_111_000: {BIT, false, 0, bit},
	0b01_111_001: {BIT, false, 0, bit},
	0b01_111_010: {BIT, false, 0, bit},
	0b01_111_011: {BIT, false, 0, bit},
	0b01_111_100: {BIT, false, 0, bit},
	0b01_111_101: {BIT, false, 0, bit},
	0b01_111_110: {BIT, false, 0, bit},
	0b01_111_111: {BIT, false, 0, bit},
	//XX_YYY_ZZZ
	//   PPQ
	0b10_000_000: {RES, false, 0, nil},
	0b10_000_001: {RES, false, 0, nil},
	0b10_000_010: {RES, false, 0, nil},
	0b10_000_011: {RES, false, 0, nil},
	0b10_000_100: {RES, false, 0, nil},
	0b10_000_101: {RES, false, 0, nil},
	0b10_000_110: {RES, false, 0, nil},
	0b10_000_111: {RES, false, 0, nil},
	0b10_001_000: {RES, false, 0, nil},
	0b10_001_001: {RES, false, 0, nil},
	0b10_001_010: {RES, false, 0, nil},
	0b10_001_011: {RES, false, 0, nil},
	0b10_001_100: {RES, false, 0, nil},
	0b10_001_101: {RES, false, 0, nil},
	0b10_001_110: {RES, false, 0, nil},
	0b10_001_111: {RES, false, 0, nil},
	0b10_010_000: {RES, false, 0, nil},
	0b10_010_001: {RES, false, 0, nil},
	0b10_010_010: {RES, false, 0, nil},
	0b10_010_011: {RES, false, 0, nil},
	0b10_010_100: {RES, false, 0, nil},
	0b10_010_101: {RES, false, 0, nil},
	0b10_010_110: {RES, false, 0, nil},
	0b10_010_111: {RES, false, 0, nil},
	0b10_011_000: {RES, false, 0, nil},
	0b10_011_001: {RES, false, 0, nil},
	0b10_011_010: {RES, false, 0, nil},
	0b10_011_011: {RES, false, 0, nil},
	0b10_011_100: {RES, false, 0, nil},
	0b10_011_101: {RES, false, 0, nil},
	0b10_011_110: {RES, false, 0, nil},
	0b10_011_111: {RES, false, 0, nil},
	0b10_100_000: {RES, false, 0, nil},
	0b10_100_001: {RES, false, 0, nil},
	0b10_100_010: {RES, false, 0, nil},
	0b10_100_011: {RES, false, 0, nil},
	0b10_100_100: {RES, false, 0, nil},
	0b10_100_101: {RES, false, 0, nil},
	0b10_100_110: {RES, false, 0, nil},
	0b10_100_111: {RES, false, 0, nil},
	0b10_101_000: {RES, false, 0, nil},
	0b10_101_001: {RES, false, 0, nil},
	0b10_101_010: {RES, false, 0, nil},
	0b10_101_011: {RES, false, 0, nil},
	0b10_101_100: {RES, false, 0, nil},
	0b10_101_101: {RES, false, 0, nil},
	0b10_101_110: {RES, false, 0, nil},
	0b10_101_111: {RES, false, 0, nil},
	0b10_110_000: {RES, false, 0, nil},
	0b10_110_001: {RES, false, 0, nil},
	0b10_110_010: {RES, false, 0, nil},
	0b10_110_011: {RES, false, 0, nil},
	0b10_110_100: {RES, false, 0, nil},
	0b10_110_101: {RES, false, 0, nil},
	0b10_110_110: {RES, false, 0, nil},
	0b10_110_111: {RES, false, 0, nil},
	0b10_111_000: {RES, false, 0, nil},
	0b10_111_001: {RES, false, 0, nil},
	0b10_111_010: {RES, false, 0, nil},
	0b10_111_011: {RES, false, 0, nil},
	0b10_111_100: {RES, false, 0, nil},
	0b10_111_101: {RES, false, 0, nil},
	0b10_111_110: {RES, false, 0, nil},
	0b10_111_111: {RES, false, 0, nil},
	//XX_YYY_ZZZ
	//   PPQ
	0b11_000_000: {SET, false, 0, nil},
	0b11_000_001: {SET, false, 0, nil},
	0b11_000_010: {SET, false, 0, nil},
	0b11_000_011: {SET, false, 0, nil},
	0b11_000_100: {SET, false, 0, nil},
	0b11_000_101: {SET, false, 0, nil},
	0b11_000_110: {SET, false, 0, nil},
	0b11_000_111: {SET, false, 0, nil},
	0b11_001_000: {SET, false, 0, nil},
	0b11_001_001: {SET, false, 0, nil},
	0b11_001_010: {SET, false, 0, nil},
	0b11_001_011: {SET, false, 0, nil},
	0b11_001_100: {SET, false, 0, nil},
	0b11_001_101: {SET, false, 0, nil},
	0b11_001_110: {SET, false, 0, nil},
	0b11_001_111: {SET, false, 0, nil},
	0b11_010_000: {SET, false, 0, nil},
	0b11_010_001: {SET, false, 0, nil},
	0b11_010_010: {SET, false, 0, nil},
	0b11_010_011: {SET, false, 0, nil},
	0b11_010_100: {SET, false, 0, nil},
	0b11_010_101: {SET, false, 0, nil},
	0b11_010_110: {SET, false, 0, nil},
	0b11_010_111: {SET, false, 0, nil},
	0b11_011_000: {SET, false, 0, nil},
	0b11_011_001: {SET, false, 0, nil},
	0b11_011_010: {SET, false, 0, nil},
	0b11_011_011: {SET, false, 0, nil},
	0b11_011_100: {SET, false, 0, nil},
	0b11_011_101: {SET, false, 0, nil},
	0b11_011_110: {SET, false, 0, nil},
	0b11_011_111: {SET, false, 0, nil},
	0b11_100_000: {SET, false, 0, nil},
	0b11_100_001: {SET, false, 0, nil},
	0b11_100_010: {SET, false, 0, nil},
	0b11_100_011: {SET, false, 0, nil},
	0b11_100_100: {SET, false, 0, nil},
	0b11_100_101: {SET, false, 0, nil},
	0b11_100_110: {SET, false, 0, nil},
	0b11_100_111: {SET, false, 0, nil},
	0b11_101_000: {SET, false, 0, nil},
	0b11_101_001: {SET, false, 0, nil},
	0b11_101_010: {SET, false, 0, nil},
	0b11_101_011: {SET, false, 0, nil},
	0b11_101_100: {SET, false, 0, nil},
	0b11_101_101: {SET, false, 0, nil},
	0b11_101_110: {SET, false, 0, nil},
	0b11_101_111: {SET, false, 0, nil},
	0b11_110_000: {SET, false, 0, nil},
	0b11_110_001: {SET, false, 0, nil},
	0b11_110_010: {SET, false, 0, nil},
	0b11_110_011: {SET, false, 0, nil},
	0b11_110_100: {SET, false, 0, nil},
	0b11_110_101: {SET, false, 0, nil},
	0b11_110_110: {SET, false, 0, nil},
	0b11_110_111: {SET, false, 0, nil},
	0b11_111_000: {SET, false, 0, nil},
	0b11_111_001: {SET, false, 0, nil},
	0b11_111_010: {SET, false, 0, nil},
	0b11_111_011: {SET, false, 0, nil},
	0b11_111_100: {SET, false, 0, nil},
	0b11_111_101: {SET, false, 0, nil},
	0b11_111_110: {SET, false, 0, nil},
	0b11_111_111: {SET, false, 0, nil},
}
