package cpu

import ()

var instructions map[byte]func()

func init() {
	instructions = make(map[byte]func())
}
