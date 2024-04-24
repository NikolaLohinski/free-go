//go:build mage

package main

import (
	"os"
	"path"

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
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := os.Rename(path.Join(wd, "bin", "v2"), path.Join(wd, "bin", "ginkgo")); err != nil {
		panic(err)
	}
}
