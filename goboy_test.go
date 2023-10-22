package goboy_test

import (
	"os"
	"testing"

	"github.com/coreyog/goboy"
	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	// utility test for debugging goboy internals without WASM
	rom, err := os.ReadFile("cmd/goboy-wasm/dist/DMG_ROM.bin")
	assert.NoError(t, err)

	gb := &goboy.GameBoy{}
	gb.Debug = true

	gb.LoadROM(rom)

	gb.RunFrame()

	t.Skip()
}
