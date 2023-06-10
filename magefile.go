//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	// mage:import
	"github.com/nikolalohinski/free-go/spellbook"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Verify

func init() {
	name, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	p, err := filepath.Abs(path.Join(name, "tools", "bin"))
	if err != nil {
		panic(err)
	}

	if err = os.Setenv("PATH", fmt.Sprintf("%s:%s", p, os.Getenv("PATH"))); err != nil {
		panic(err)
	}

	if len(os.Args) == 1 {
		defaultTargeName := strings.ToLower(strings.TrimPrefix(runtime.FuncForPC(reflect.ValueOf(Default).Pointer()).Name(), "main."))
		fmt.Println("┌ mage", defaultTargeName)

		return
	}

	if os.Args[1] == "install" {
		return
	}

	fmt.Println("┌ mage", strings.Join(os.Args[1:], " "))
}

// Fetch and installs tooling for development
func Install() error {
	return sh.RunV("mage", "-d", "./tools")
}

// Validate code base
func Verify() {
	mg.SerialDeps(
		spellbook.Go.Tidy,
		spellbook.Go.Format,
		spellbook.Go.Lint,
		spellbook.Go.Test,
	)
}
