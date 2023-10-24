package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/coreyog/goboy"

	"golang.design/x/clipboard"
)

func main() {
	// handle ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		os.Exit(0)
	}()

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	// easy input
	reader := bufio.NewReader(os.Stdin)

	// reference materials
	unprefixed, cb, opnames := goboy.GetOpLookups()

	for {
		// prompt
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// ctrl+c or ctrl+d (EOF)
				return
			}

			panic(err)
		}

		// clean off newline from reader
		input := strings.TrimSpace(line)

		// looking for number...
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("invalid input")
			continue
		}

		// between [0-255]
		if num < 0 || num > 255 {
			// better explain rejection
			fmt.Println("expected number between 0 and 255")
			continue
		}

		b := byte(num)

		// anatomy of an opcode
		x := b >> 6
		y := (b >> 3) & 7
		z := b & 7
		p := y >> 1
		q := y & 1

		// output
		underscores := fmt.Sprintf("%02b_%03b_%03b", x, y, z)
		fmt.Printf("0x%02X\n", b)
		fmt.Printf("0b%s\n", underscores)
		fmt.Printf("   %d   %d   %d\n", x, y, z)
		fmt.Printf("     %d %d\n", p, q)
		fmt.Printf("un: %s (%t)\n", opnames[unprefixed[b].Code], unprefixed[b].Operation != nil)
		fmt.Printf("cb: %s (%t)\n\n", opnames[cb[b].Code], cb[b].Operation != nil)

		// update clipboard, useful for my "flow"
		clipboard.Write(clipboard.FmtText, []byte(underscores))
	}
}
