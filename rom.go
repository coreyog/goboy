package goboy

import (
	"fmt"
	"time"
)

func (gb *GameBoy) LoadROM(d []byte) {
	gb.romData = d
	gb.pc = 0
}

// ReadRom8 reads a byte from Rom at a given address, respecting Rom mapping
func (gb *GameBoy) ReadRom8(address uint16) (value byte) {
	return gb.romData[address]
}

func (gb *GameBoy) ReadRom16(address uint16) (value uint16) {
	lsb := gb.romData[address] // little endian
	msb := gb.romData[address+1]
	return mergeBytes(msb, lsb)
}

// WriteRom sets the value at a given address in Rom, respecting Rom mapping
func (gb *GameBoy) WriteRom(address uint16, value byte) {
	gb.romData[address] = value
}

// inspired by https://docs.libretro.com/development/cores/developing-cores/#retro_run
func (gb *GameBoy) RunFrame() {
	// this would probably continuously
	// run instructions until a VBlank interrupt

	// i=24579 - exit loop zeroing vram
	for i := 0; i < 24585; i++ {
		gb.debugLnF("instruction #%d", i)

		if gb.pc == 0x0016 {
			gb.debugLnF("breakpoint")
		}

		gb.RunInstruction()
	}
}

func (gb *GameBoy) RunInstruction() {
	start := time.Now()
	defer func() {
		gb.debugLnF("instruction ET: %s\n", time.Since(start))
	}()

	// first byte of instruction might be a prefix
	gb.debugLnF("PC: %.4X", gb.pc)

	prefix := gb.ReadRom8(gb.pc)
	offset := uint16(1)

	var opbytes OpBytes
	var ok bool
	var displacement, opcode byte
	var immediate uint16

	// check for known prefixes
	if prefix == 0xCB {
		opcode = gb.ReadRom8(gb.pc + offset)
		offset++

		opbytes, ok = cb[opcode]
	} else {
		// unprefixed, first byte must be opcode
		opcode = prefix
		prefix = 0

		opbytes, ok = unprefixed[opcode]
	}

	// no matching instruction
	if !ok {
		// the actual gameboy does nothing, no crash, no change of flags, nothing
		// TODO: remove this panic, eventually
		panic(fmt.Sprintf("unrecognized opcode %.2X at offset %.2X", opcode, gb.pc))
	}

	// gb.debugPrintlnf("prefix: %.2X, opcode: %.2X, has displacement: %t, immediate size: %d", prefix, opcode, opbytes.HasDisplacement, opbytes.ImmediateSize)

	if opbytes.HasDisplacement {
		// byte after opcode
		displacement = gb.ReadRom8(gb.pc + offset)
		offset++
	}

	if opbytes.ImmediateSize == 1 {
		immediate = uint16(gb.ReadRom8(gb.pc + offset))
	} else if opbytes.ImmediateSize == 2 {
		immediate = gb.ReadRom16(gb.pc + offset)
	}

	offset += uint16(opbytes.ImmediateSize)

	// gb.debugPrintlnf("displacement: %.2X, immediate: %.4X", displacement, immediate)

	if opbytes.Operation != nil {
		opbytes.Operation(gb, prefix, OpCode(opcode), displacement, immediate)
	} else if opcode != 0 {
		gb.debugLnF("unknown opcode %.2X at PC %.4X", opcode, gb.pc)
	}

	gb.pc += offset

	gb.debugLnF("HL: %.4X", gb.readHL())
	gb.debugLnF("next PC: %.4X", gb.pc)
}

func (gb *GameBoy) debugLnF(format string, a ...interface{}) {
	if gb.Debug {
		fmt.Printf(format, a...)
		fmt.Println()
	}
}
