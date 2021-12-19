package main

// no need for a build tag, the filename limits this to windows
//go:generate cmd /C "set GOOS=js&& set GOARCH=wasm&& go build -o dist/goboy.wasm ."
