package goboy

func ldsp(gb *GameBoy, displacement byte, immediate uint16) {
	// LD SP, $FFFF
	gb.sp = immediate
	// no flags to set
}

func xora(gb *GameBoy, displacement byte, immediate uint16) {
	gb.a ^= gb.a
	// some flags to set, unclear
}
