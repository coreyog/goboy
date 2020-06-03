package goboy

func (gb *GameBoy) LoadROM(d []byte) {
	gb.romData = d
}

// inspired by https://docs.libretro.com/development/cores/developing-cores/#retro_run
func (gb *GameBoy) RunFrame() {

}
