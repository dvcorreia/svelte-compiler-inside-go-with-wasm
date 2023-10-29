// Package to run Svelte from Go.
//
// Adapted code from bud (ref: github.com/livebud/bud).
package svelte

//go:generate esbuild compiler.ts --format=iife --global-name=__svelte__ --bundle --platform=node --inject:shimssr.ts --external:url --outfile=compiler.js --log-level=warning

import (
	_ "embed"
)

// compiler.js is used to compile .svelte files into JS & CSS
//
//go:embed compiler.js
var compiler string

type SSR struct {
	JS  string
	CSS string
}

type Compiler interface {
	SSR(path string, code []byte) (*SSR, error)
	DOM(path string, code []byte) (*DOM, error)
}
