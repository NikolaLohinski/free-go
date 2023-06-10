package spellbook

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/magefile/mage/sh"
)

func Step(message string) {
	fmt.Printf("├─ %s\n", message)
}

func RunSub(cmd string, args ...string) error {
	_, err := sh.Exec(
		nil,
		&prefixer{prefix: "│    ", writer: os.Stdout, trailingNewline: true},
		&prefixer{prefix: "│    ", writer: os.Stderr, trailingNewline: true},
		cmd,
		args...,
	)
	if err != nil {
		fmt.Println("└─ ❌")
	}
	return err
}

type prefixer struct {
	prefix          string
	writer          io.Writer
	trailingNewline bool
	buf             bytes.Buffer // reuse buffer to save allocations
}

func (p *prefixer) Write(payload []byte) (int, error) {
	p.buf.Reset() // clear the buffer

	for _, b := range payload {
		if p.trailingNewline {
			p.buf.WriteString(p.prefix)
			p.trailingNewline = false
		}

		p.buf.WriteByte(b)

		if b == '\n' {
			// do not print the prefix right after the newline character as this might
			// be the very last character of the stream and we want to avoid a trailing prefix.
			p.trailingNewline = true
		}
	}

	n, err := p.writer.Write(p.buf.Bytes())
	if err != nil {
		// never return more than original length to satisfy io.Writer interface
		if n > len(payload) {
			n = len(payload)
		}

		return n, err
	}

	// return original length to satisfy io.Writer interface
	return len(payload), nil
}
