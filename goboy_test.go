package goboy_test

import (
	"io/ioutil"
	"testing"

	"github.com/coreyog/goboy"
)

func TestDebug(t *testing.T) {
	// utility test for debugging goboy internals without WASM
	rom, err := ioutil.ReadFile("cmd/goboy-wasm/dist/dmg0_rom.bin")
	if err != nil {
		panic(err)
	}

	gb := &goboy.GameBoy{}
	gb.Debug = true

	gb.LoadROM(rom)

	gb.RunFrame()

	t.Skip()
}
