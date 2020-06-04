package goboy

import (
	"fmt"
	"time"
)

func (gb *GameBoy) LoadROM(d []byte) {
	gb.romData = d
	gb.pc = 0
}

// ReadRom reads a byte from Rom at a given address, respecting Rom mapping
func (gb *GameBoy) ReadRom(address uint16) (value byte) {
	return gb.romData[address]
}

// WriteRom sets the value at a given address in Rom, respecting Rom mapping
func (gb *GameBoy) WriteRom(address uint16, value byte) {
	gb.romData[address] = value
}

// inspired by https://docs.libretro.com/development/cores/developing-cores/#retro_run
func (gb *GameBoy) RunFrame() {
	// this would probably continuously
	// run instructions until a VBlank interrupt
	fmt.Println("running first instruction")
	gb.RunInstruction()
	fmt.Println("running second instruction")
	gb.RunInstruction()
}

func (gb *GameBoy) RunInstruction() {
	start := time.Now()
	defer func() {
		fmt.Printf("instruction ET: %s\n", time.Since(start))
	}()
	// first byte of instruction might be a prefix
	fmt.Printf("PC: %.4X\n", gb.pc)

	prefix := gb.ReadRom(gb.pc)
	offset := uint16(1)

	var opbytes OpBytes
	var ok bool
	var displacement, opcode byte
	var immediate uint16

	// check for known prefixes
	if prefix == 0xcb {
		opcode = gb.ReadRom(gb.pc + offset)
		offset++

		opbytes, ok = cb[opcode]
	} else if prefix == 0xed {
		opcode = gb.ReadRom(gb.pc + offset)
		offset++

		opbytes, ok = ed[opcode]
	} else {
		// unprefixed, first byte must be opcode
		opcode = prefix
		prefix = 0

		opbytes, ok = unprefixed[opcode]
	}

	// no matching instruction
	if !ok {
		panic(fmt.Sprintf("unrecognized opcode %.2X at offset %.2X", opcode, gb.pc))
	}

	fmt.Printf("prefix: %.2X, opcode: %.2X, displacement: %t, immediate size: %d\n", prefix, opcode, opbytes.HasDisplacement, opbytes.ImmediateSize)

	if opbytes.HasDisplacement {
		// byte after opcode
		displacement = gb.ReadRom(gb.pc + offset)
		offset++
	}

	if opbytes.ImmediateSize == 1 {
		immediate = uint16(gb.ReadRom(gb.pc + offset))
	} else if opbytes.ImmediateSize == 2 {
		msb := gb.ReadRom(gb.pc + offset)
		lsb := gb.ReadRom(gb.pc + offset + 1)
		immediate = mergeBytes(msb, lsb)
	}
	offset += uint16(opbytes.ImmediateSize)

	fmt.Printf("displacement: %.2X, immediate: %.4X\n", displacement, immediate)

	opbytes.Operation(gb, prefix, opcode, displacement, immediate)

	gb.pc += offset

	fmt.Printf("Next PC: %.4X\n", gb.pc)
}
