package spellbook

import (
	"context"
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
func (Go) Test(ctx context.Context) error {
	return Run(Invoke(ctx, "Running unit tests"), "ginkgo", "--skip-package", "integration", "./...")
}

// Runs ginkgo for integration test
func (Go) Integration(ctx context.Context) error {
	return Run(Invoke(ctx, "Running integration tests"), "ginkgo", "-tags=integration", "./integration/...")
}

// Cleans dependencies and imports
func (Go) Tidy(ctx context.Context) error {
	return Run(Invoke(ctx, "Cleaning dependencies and imports in code files"), "go", "mod", "tidy", "-v")
}

// Runs a linter on the source code
func (Go) Lint(ctx context.Context) error {
	return Run(Invoke(ctx, "Running golang-ci linter on code base"), "golangci-lint", "run", "--verbose", "--fix")
}

// Runs formatting tools on the code base
func (Go) Format(ctx context.Context) error {
	ctx = Invoke(ctx, "Running gofumpt formatter on code base")
	args := append([]string{"-w"}, getGoFiles()...)
	return Run(ctx, "gofumpt", args...)
}

// Builds and opens a coverage report
func (Go) Cover(ctx context.Context) error {
	ctx = Invoke(ctx, "Generating coverage report")
	if err := Run(ctx, "go", "test", "-v", "-coverprofile", "cover.out", "./..."); err != nil {
		panic(err)
	}
	if err := sh.Run("go", "tool", "cover", "-html", "cover.out", "-o", "cover.html"); err != nil {
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
