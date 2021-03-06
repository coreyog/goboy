package goboy

// Screen dimensions
const (
	ScreenWidth  = 160
	ScreenHeight = 144
)

// 4 colors
const (
	//          0xAARRGGBB
	Black     = 0xFF000000
	DarkGray  = 0xFF555555
	LightGray = 0xFFAAAAAA
	White     = 0xFFFFFFFF
)

// VRAM is an 8x8 set of tiles, memory mapped?
var VRAM [8][8]uint8

// memory mapped?
var scrollX uint8
var scrollY uint8
