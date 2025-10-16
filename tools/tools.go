//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/onsi/ginkgo/v2/ginkgo"
	_ "github.com/x-izumin/gex/cmd/gex"
	_ "mvdan.cc/gofumpt"
)
