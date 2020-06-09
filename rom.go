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

	// i=24579 - exit loop zeroing vram
	for i := 0; i < 24580; i++ {
		gb.debugPrintlnf("instruction #%d", i)

		if gb.pc == 0x000C {
			gb.debugPrintlnf("breakpoint")
		}

		gb.RunInstruction()
	}
}

func (gb *GameBoy) RunInstruction() {
	if gb.Debug {
		start := time.Now()
		defer func() {
			gb.debugPrintlnf("instruction ET: %s\n", time.Since(start))
		}()
	}

	// first byte of instruction might be a prefix
	gb.debugPrintlnf("PC: %.4X", gb.pc)

	prefix := gb.ReadRom(gb.pc)
	offset := uint16(1)

	var opbytes OpBytes
	var ok bool
	var displacement, opcode byte
	var immediate uint16

	// check for known prefixes
	if prefix == 0xCB {
		opcode = gb.ReadRom(gb.pc + offset)
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
		displacement = gb.ReadRom(gb.pc + offset)
		offset++
	}

	if opbytes.ImmediateSize == 1 {
		immediate = uint16(gb.ReadRom(gb.pc + offset))
	} else if opbytes.ImmediateSize == 2 {
		lsb := gb.ReadRom(gb.pc + offset) // little endian
		msb := gb.ReadRom(gb.pc + offset + 1)
		immediate = mergeBytes(msb, lsb)
	}

	offset += uint16(opbytes.ImmediateSize)

	// gb.debugPrintlnf("displacement: %.2X, immediate: %.4X", displacement, immediate)

	if opbytes.Operation != nil {
		opbytes.Operation(gb, prefix, OpCode(opcode), displacement, immediate)
	} else if opcode != 0 {
		gb.debugPrintlnf("unknown opcode %.2X at PC %.4X", opcode, gb.pc)
	}

	gb.pc += offset

	gb.debugPrintlnf("HL: %.4X", gb.readHL())
	gb.debugPrintlnf("next PC: %.4X", gb.pc)
}

func (gb *GameBoy) debugPrintlnf(format string, a ...interface{}) {
	if gb.Debug {
		fmt.Printf(format, a...)
		fmt.Println()
	}
}
