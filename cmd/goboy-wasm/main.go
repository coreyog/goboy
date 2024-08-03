//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"image"
	"image/draw"
	"strings"
	"syscall/js"

	"github.com/coreyog/goboy"

	"golang.org/x/image/colornames"
)

const (
	width  = 160
	height = 144
)

var ctx js.Value
var requestAnimationFrame js.Value
var jsOnFrame js.Func
var console js.Value
var audioCtx js.Value
var audioCtxDest js.Value
var oscillator js.Value
var gain js.Value
var curGain float32 = 0.05
var fps js.Value
var img *image.RGBA
var progress float64
var prevTS float64
var killSwitch chan struct{}
var closing bool
var updateFPS bool = false

var gb *goboy.GameBoy

func main() {
	gb = &goboy.GameBoy{}

	// prep state
	killSwitch = make(chan struct{})
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "target")

	fn := canvas.Get("getContext")
	if !fn.Truthy() {
		return
	}

	// setup video stuff
	fps = document.Call("getElementById", "fps")
	console = js.Global().Get("console")
	ctx = canvas.Call("getContext", "2d")
	requestAnimationFrame = js.Global().Get("requestAnimationFrame")
	jsOnFrame = js.FuncOf(onFrame)

	// setup audio stuff
	audioCtxFunc := js.Global().Get("AudioContext")
	audioCtx = audioCtxFunc.New()
	audioCtxDest = audioCtx.Get("destination")
	oscillator = audioCtx.Call("createOscillator")
	oscillator.Set("type", "square")
	gain = audioCtx.Call("createGain")
	oscillator.Call("connect", gain)
	gain.Call("connect", audioCtxDest)
	gain.Get("gain").Set("value", curGain)
	// oscillator.Call("start")

	// create image
	img = image.NewRGBA(image.Rect(0, 0, width, height))

	// 1px solid black border
	draw.Draw(img, img.Bounds(), image.NewUniform(colornames.Black), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(1, 1, width-1, height-1), image.NewUniform(colornames.White), image.Point{}, draw.Src)

	// draw frame
	onFrame(js.Null(), []js.Value{js.ValueOf(0)})

	// register kill switch
	js.Global().Set("stopWASM", js.FuncOf(stopWASM))
	js.Global().Set("loadROM", js.FuncOf(loadROM))

	// wait for call to stopWASM
	<-killSwitch

	// oscillator.Call("stop")

	// fill the image with white and clear the canvas
	draw.Draw(img, image.Rect(0, 0, width, height), image.NewUniform(colornames.White), image.Point{}, draw.Src)
	drawImage(ctx, img)

	fps.Set("innerHTML", "---")
}

func onFrame(this js.Value, args []js.Value) interface{} {
	// determine timestamp and delta time
	ts := args[0].Float()      // in milliseconds since start
	dt := (ts - prevTS) / 1000 // in seconds since last frame
	prevTS = ts

	// update FPS in DOM
	if updateFPS {
		text := fmt.Sprintf("fps: %0.0f\n", 1/dt)
		fps.Set("innerHTML", text)
	}

	// inset colored rectangle
	draw.Draw(img, image.Rect(10, 10, width-10, height-10), image.NewUniform(Keypoints.GetInterpolatedColorFor(progress)), image.Point{}, draw.Src)
	drawImage(ctx, img)

	// increment progress through the gradient
	progress += dt / 3
	if progress > 1 {
		progress -= 1
	}

	// playAudio(ts, dt)

	if !closing {
		requestAnimationFrame.Invoke(jsOnFrame)
	} else {
		killSwitch <- struct{}{}
	}

	return js.Null()
}

func drawImage(ctx js.Value, img *image.RGBA) {
	// copy to JS
	Uint8Array := js.Global().Get("Uint8Array")
	jsData := Uint8Array.New(len(img.Pix))
	js.CopyBytesToJS(jsData, img.Pix)

	// clamp the data
	Uint8ClampedArray := js.Global().Get("Uint8ClampedArray")
	jsClampedData := Uint8ClampedArray.New(jsData) // view, don't use "Uint8ClampedArray.from(...)"

	// make it Image Data
	ImageData := js.Global().Get("ImageData")
	imgData := ImageData.New(jsClampedData, width)

	// put it on the canvas
	ctx.Call("putImageData", imgData, 0, 0)
}

// func playAudio(ts float64, dt float64) {
// 	if int(ts) != int(ts-dt) {
// 		if curGain != 0 {
// 			curGain = 0
// 		} else {
// 			curGain = 0.05
// 		}
// 		gain.Get("gain").Set("value", curGain)
// 	}
// }

func loadROM(this js.Value, args []js.Value) interface{} {
	fmt.Printf("WASM - loading ROM (%d)\n", len(args))
	if len(args) != 1 {
		fmt.Printf("invalid number of args, expected 1, got %d", len(args))
		return js.Null()
	}

	array := args[0]
	conName := array.Get("constructor").Get("name").String()

	if !strings.EqualFold(conName, "uint8array") {
		fmt.Printf("invalid argument, expected: Uint8Array, actual: %s\n", conName)
		return js.Null()
	}

	size := int(array.Get("byteLength").Float())
	data := make([]byte, size)

	js.CopyBytesToGo(data, array)

	gb.LoadROM(data)

	// TODO: run multiple frames
	gb.RunFrame()

	return js.Null()
}

func stopWASM(this js.Value, args []js.Value) interface{} {
	closing = true
	return js.Null()
}
