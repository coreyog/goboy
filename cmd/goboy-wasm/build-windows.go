//go:build windows
// +build windows

package main

//go:generate cmd /C "set GOOS=js&& set GOARCH=wasm&& go build -o dist/goboy.wasm ."
