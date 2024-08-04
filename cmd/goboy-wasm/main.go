//go:build js && wasm
// +build js,wasm

package main

import (
	"container/ring"
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

var ( // constant-like variables
	// JS Constants
	window = js.Global()
	JSNULL = js.Null()

	// JS Types
	Uint8ClampedArray = window.Get("Uint8ClampedArray")
	ImageData         = window.Get("ImageData")
	AudioContext      = window.Get("AudioContext")

	// JS Functions
	requestAnimationFrame = window.Get("requestAnimationFrame")

	// JS Global Objects
	document = window.Get("document")
	console  = window.Get("console")
)

var (
	img            *image.RGBA // THE frame buffer
	ctx            js.Value    // CanvasRenderingContext2D
	jsOnFrame      js.Func
	audioCtx       js.Value // AudioContext
	audioCtxDest   js.Value // AudioDestinationNode
	oscillator     js.Value // OscillatorNode
	gain           js.Value // GainNode
	curGain        float32  = 0.05
	pixelData      js.Value // Uint8ClampedArray
	fps            js.Value // HTMLSpanElement
	progress       float64
	killSwitch     chan struct{} = make(chan struct{}, 1)
	closing        bool
	prevTS         float64
	calcFPS        bool = true
	fpsSum         float64
	fpsHistory     *ring.Ring      // history of the last [fpsHistorySize] frame times in seconds
	fpsHistorySize int        = 10 // size of the history ring
	frameCount     uint64

	gb *goboy.GameBoy = &goboy.GameBoy{}
)

func init() {
	fpsHistory = ring.New(fpsHistorySize)
	for range fpsHistorySize {
		fpsHistory.Value = float64(0.0)
		fpsHistory = fpsHistory.Next()
	}

	jsOnFrame = js.FuncOf(onFrame)
	window.Set("stopWASM", js.FuncOf(stopWASM))
	window.Set("loadROM", js.FuncOf(loadROM))
	window.Set("_toggleFPS", js.FuncOf(toggleFPS))
}

func main() {
	// prep state
	canvas := document.Call("getElementById", "target")

	fn := canvas.Get("getContext")
	if !fn.Truthy() {
		console.Call("log", "getContext not found, closing")
		return
	}

	fps = document.Call("getElementById", "fps")
	ctx = canvas.Call("getContext", "2d")

	// setup audio stuff
	audioCtx = AudioContext.New()
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
	pixelData = Uint8ClampedArray.New(len(img.Pix))

	// 1px solid black border
	draw.Draw(img, img.Bounds(), image.NewUniform(colornames.Black), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(1, 1, width-1, height-1), image.NewUniform(colornames.White), image.Point{}, draw.Src)

	// draw frame, kick off RAF loop
	onFrame(JSNULL, []js.Value{js.ValueOf(0)})

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
	ts := args[0].Float()      // time since start in milliseconds
	dt := (ts - prevTS) / 1000 // time since last frame in seconds
	frameCount++
	prevTS = ts

	// update FPS in DOM
	updateFPS(dt)

	if calcFPS {
		if frameCount%30 == 0 {
			calc := fpsSum / float64(fpsHistorySize)
			text := fmt.Sprintf("fps: %0.0f", 1/calc)
			fps.Set("innerHTML", text)
		}
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

	return JSNULL
}

func drawImage(ctx js.Value, img *image.RGBA) {
	// copy to JS
	js.CopyBytesToJS(pixelData, img.Pix)

	// make it Image Data
	imgData := ImageData.New(pixelData, width)

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
		return JSNULL
	}

	array := args[0]
	conName := array.Get("constructor").Get("name").String()

	if !strings.EqualFold(conName, "uint8array") {
		fmt.Printf("invalid argument, expected: Uint8Array, actual: %s\n", conName)
		return JSNULL
	}

	size := int(array.Get("byteLength").Float())
	data := make([]byte, size)

	js.CopyBytesToGo(data, array)

	gb.LoadROM(data)

	// TODO: run multiple frames
	gb.RunFrame()

	return JSNULL
}

func stopWASM(this js.Value, args []js.Value) interface{} {
	closing = true
	return JSNULL
}

func updateFPS(sinceLastFrame float64) {
	fpsSum += sinceLastFrame - fpsHistory.Value.(float64)

	fpsHistory.Value = sinceLastFrame
	fpsHistory = fpsHistory.Next()
}

func toggleFPS(_ js.Value, _ []js.Value) interface{} {
	calcFPS = !calcFPS

	if calcFPS {
		fps.Set("innerHTML", "fps: -")
	} else {
		fps.Set("innerHTML", "---")
	}

	return JSNULL
}
