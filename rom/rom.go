package rom

import ()

var data []byte

func Load(d []byte) {
	data = d
	preprocess()
}

func preprocess() {
	pos := 0x100
	for pos < len(data) {
		d := data[pos]
		pos++

		var bytes OpBytes
		switch d {
		case 0xCB:
		case 0xED:
		case 0xDD:
		case 0xFD:
		default:
			// no prefix opcode
			bytes = unprefixed[d]
		}

		if bytes.HasDisplacement {
			d = data[pos]
			pos++
		}

		// imm := make([]int, bytes.ImmediateSize)
		// copy(imm, data[pos:pos+int(bytes.ImmediateSize)])
	}
}

func processOpCode(code byte) (x byte, y byte, z byte) {
	x = code >> 6
	y = (code >> 3) & 0b0111
	z = code & 0b0111

	return x, y, z
}

func splitY(y byte) (p byte, q byte) {
	p = (y >> 1) & 0b0011
	q = y & 0b0001

	return p, q
}
