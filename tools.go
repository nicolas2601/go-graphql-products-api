//go:build tools
// +build tools

// Package tools pins the code-generation tooling as an explicit module dependency.
package tools

import (
	_ "github.com/99designs/gqlgen"
)
