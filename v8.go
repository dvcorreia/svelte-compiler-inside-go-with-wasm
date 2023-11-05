package svelte

import (
	"bytes"
	"encoding/json"

	_ "embed"

	"github.com/livebud/bud/package/js"
)

// Adapted code from bud (ref: github.com/livebud/bud).

//go:generate esbuild compiler.ts --format=iife --global-name=__svelte__ --bundle --platform=node --inject:compiler_shim.ts --external:url --outfile=compiler.js --log-level=warning

//go:embed compiler.js
var compiler string

func LoadV8(vm js.VM, outputFormat Generate) (Compiler, error) {
	if err := vm.Script("svelte/compiler.js", compiler); err != nil {
		return nil, err
	}

	compiler := &V8{
		VM:           vm,
		Dev:          true,
		OutputFormat: outputFormat,
	}

	// TODO make dev configurable
	return compiler, nil
}

type V8 struct {
	VM           js.VM
	Dev          bool
	OutputFormat Generate
}

func (c *V8) Compile(path string, code []byte) (*CompileResult, error) {
	type input struct {
		Code   string   `json:"code"`
		Path   string   `json:"path"`
		Target Generate `json:"target"`
		Dev    bool     `json:"dev"`
		Css    string   `json:"css"`
	}

	var css CSS
	switch c.OutputFormat {
	case GenerateDOM:
		css = CSSInjected
	case GenerateSSR:
		css = CSSExternal
	default:
		css = CSSNone
	}

	opts := input{
		Code:   string(code),
		Path:   path,
		Target: c.OutputFormat,
		Dev:    c.Dev,
		Css:    string(css),
	}

	buf := bytes.NewBufferString(";__svelte__.compile('")

	e := json.NewEncoder(buf)
	e.SetEscapeHTML(false)

	if err := e.Encode(&opts); err != nil {
		return nil, err
	}

	// unread the \n set by the json encoder
	buf.Truncate(buf.Len() - 1)

	if _, err := buf.WriteString("')"); err != nil {
		return nil, err
	}

	result, err := c.VM.Eval(path, buf.String())
	if err != nil {
		return nil, err
	}

	var out CompileResult
	if err := json.Unmarshal([]byte(result), &out); err != nil {
		return nil, err
	}

	return &out, nil
}
