// Package to run Svelte from Go.
package svelte

import (
	_ "embed"
)

// Compiler abstraction for the Svelte compiler.
type Compiler interface {
	Compile(path string, code []byte) (*CompileResult, error)
}

// CompileResult is the compilation result of the svelte compiler.
// Others fields available for the compilation results omited for simplicity.
type CompileResult struct {
	JS  string `json:"js"`
	CSS string `json:"css"`
}

type CompileOptions struct {
	FileName string   `json:"filename,omitempty"`
	Generate Generate `json:"generate,omitempty"`

	// If true, causes extra code to be added to components that will
	// perform runtime checks and provide debugging information during development.
	Dev bool `json:"dev,omitempty"`

	CSS CSS `json:"css,omitempty"`
}

type Generate string

const (
	// Svelte emits a JavaScript class for mounting to the DOM.
	GenerateDOM Generate = "dom"

	// Svelte emits an object with a render method suitable for server-side rendering.
	GenerateSSR Generate = "ssr"
)

type CSS string

const (
	// Styles will be included in the JavaScript class
	// and injected at runtime for the components actually rendered.
	CSSInjected CSS = "injected"

	//  the CSS will be returned in the css field of the compilation result.
	// Most Svelte bundler plugins will set this to 'external' and use the CSS that is
	// statically generated for better performance, as it will result in smaller JavaScript
	// bundles and the output can be served as cacheable .css files.
	CSSExternal CSS = "external"

	// styles are completely avoided and no CSS output is generated.
	CSSNone CSS = "none"
)
