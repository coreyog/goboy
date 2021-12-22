//go:build !windows
// +build !windows

package main

//go:generate $SHELL -c "GOOS=js GOARCH=wasm go build -o dist/goboy.wasm ."
