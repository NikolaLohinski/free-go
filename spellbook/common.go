package spellbook

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type MagicalKey string

const MagicalWord MagicalKey = "mage"

var (
	MagicalContext context.Context
	magicalLock    sync.Mutex
)

func init() {
	MagicalContext = context.Background()
}

func Invoke(ctx context.Context, message string) context.Context {
	pc, _, _, _ := runtime.Caller(1)
	spell := strings.Split(runtime.FuncForPC(pc).Name(), "spellbook.")[1]
	depth := 0
	incantation, ok := ctx.Value(MagicalWord).(map[string]int)
	magicalLock.Lock()
	if ok {
		depth, ok = incantation[spell]
		if !ok {
			depth = 0
		}
	} else {
		incantation = map[string]int{spell: depth}
	}
	ctx = context.WithValue(ctx, MagicalWord, incantation)
	magicalLock.Unlock()
	prefix := "├─"
	for i := 0; i < depth; i++ {
		prefix = "│ " + prefix
	}

	fmt.Println(prefix + "🪄 [\033[1;94m" + strings.ToLower(strings.ReplaceAll(spell, ".", ":")) + "\033[0m] " + message)

	return ctx
}

func Run(ctx context.Context, cmd string, args ...string) error {
	pc, _, _, _ := runtime.Caller(1)
	spell := strings.Split(runtime.FuncForPC(pc).Name(), "spellbook.")[1]
	depth := 0
	incantation, ok := ctx.Value(MagicalWord).(map[string]int)
	if ok {
		magicalLock.Lock()
		depth, ok = incantation[spell]
		if !ok {
			depth = 0
		}
		magicalLock.Unlock()
	}
	prefix := "│ "
	for i := 0; i < depth; i++ {
		prefix = "│ " + prefix
	}

	_, err := sh.Exec(
		nil,
		&prefixer{prefix: prefix + "│ ", writer: os.Stdout, trailingNewline: true},
		&prefixer{prefix: prefix + "│ ", writer: os.Stderr, trailingNewline: true},
		cmd,
		args...,
	)
	magicalLock.Lock()
	delete(incantation, spell)
	magicalLock.Unlock()
	prefix = prefix + "└ "
	if err != nil {
		fmt.Println(prefix + "💥 \033[1;91mSpell failed\033[0m")
		fmt.Println("└ 💀 \033[1;91mSnape got me\033[0m")
	} else {
		fmt.Println(prefix + "✨ \033[1;92mSpell succeeded\033[0m")
		magicalLock.Lock()
		if len(incantation) == 0 {
			fmt.Println("└ 💫 \033[1;92mMischief managed\033[0m")
		}
		magicalLock.Unlock()
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

func Combine(ctx context.Context, spells ...interface{}) {
	depth := 0
	incantation, ok := ctx.Value(MagicalWord).(map[string]int)
	if !ok || incantation == nil {
		incantation = make(map[string]int)
	} else {
		for _, d := range incantation {
			depth = d + 1
			break
		}
	}
	for _, spell := range spells {
		name := strings.Split(runtime.FuncForPC(reflect.ValueOf(spell).Pointer()).Name(), "spellbook.")[1]
		magicalLock.Lock()
		incantation[name] = depth
		magicalLock.Unlock()
	}
	ctx = context.WithValue(ctx, MagicalWord, incantation)
	mg.SerialCtxDeps(ctx, spells...)
}
