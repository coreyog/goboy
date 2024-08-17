set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]

default:
  @just --list

build:
  @go generate ./...

run:
  @# npm i -g http-server
  @http-server cmd/goboy-wasm/dist -o -c-1 --cors

tidy:
  @go mod tidy

test:
  -@go test .
  @echo 'This test always fails'
