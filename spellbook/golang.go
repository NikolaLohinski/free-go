package spellbook

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Go mg.Namespace

// Runs ginkgo for unit tests
func (Go) Test() error {
	Step("Running unit tests")
	return RunSub("ginkgo", "./...")
}

// Runs ginkgo for integration test
func (Go) Integration() error {
	Step("Running integration tests")
	return RunSub("ginkgo", "-tags=integration", "./integration/...")
}

// Clean dependencies and imports
func (Go) Tidy() error {
	Step("Cleaning dependencies and imports in code files")
	return sh.RunV("go", "mod", "tidy", "-v")
}

// Run linter on code
func (Go) Lint() error {
	Step("Running golang-ci linter on code base")
	return RunSub("golangci-lint", "run", "--verbose", "--fix")
}

// Run formatting tools on code base
func (Go) Format() error {
	Step("Running gofumpt formatter on code base")
	args := append([]string{"-w"}, getGoFiles()...)
	return sh.RunV("gofumpt", args...)
}

// Build and open coverage report
func (Go) Cover() error {
	Step("Generating coverage report")
	if err := RunSub("go", "test", "-v", "-coverprofile", "cover.out", "./..."); err != nil {
		panic(err)
	}
	if err := sh.RunV("go", "tool", "cover", "-html", "cover.out", "-o", "cover.html"); err != nil {
		panic(err)
	}
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	default:
		cmd = "xdg-open"
	}
	args = append(args, "cover.html")

	return sh.Run(cmd, args...)
}

func getGoFiles() []string {
	var goFiles []string

	err := filepath.Walk(os.Getenv("PWD"), func(path string, info os.FileInfo, err error) error {
		if regexp.MustCompile("(vendor|tools|spellbook)/").MatchString(path) {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		absPath := strings.Replace(path, os.Getenv("PWD"), ".", 1)
		goFiles = append(goFiles, absPath)

		return nil
	})
	if err != nil {
		panic(err)
	}

	return goFiles
}
