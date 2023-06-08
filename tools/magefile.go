//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() {
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		panic(err)
	}
	if err := sh.RunV("go", "run", "-mod=mod", "github.com/izumin5210/gex/cmd/gex", "--build"); err != nil {
		panic(err)
	}
}
