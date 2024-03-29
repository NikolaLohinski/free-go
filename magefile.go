//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	// mage:import
	"github.com/nikolalohinski/free-go/spellbook"

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

	if len(os.Args) == 2 && os.Args[1] == "install" {
		return
	}

	fmt.Println("🧙‍ \033[1;94mI solemnly swear that I am up to no good\033[0m")
}

// Fetch and installs tooling for development
func Install() error {
	return sh.RunV("mage", "-d", "./tools")
}

// Validate code base
func Verify() {
	spellbook.Combine(
		spellbook.MagicalContext,
		spellbook.Go.Tidy,
		spellbook.Go.Format,
		spellbook.Go.Lint,
		spellbook.Go.Test,
	)
}
