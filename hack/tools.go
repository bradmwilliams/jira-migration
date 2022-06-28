//go:build tools
// +build tools

// This package imports things required by build scripts, to force `go mod` to see them as dependencies
package tools

import (
	_ "github.com/openshift/build-machinery-go"
	_ "k8s.io/code-generator"
)
