package svelte

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"

	"github.com/tetratelabs/wazero"
)

//go:generate esbuild wasm_compiler.ts --format=iife --bundle --inject:compiler_shim.ts --external:url --outfile=wasm_compiler.js --log-level=warning
//go:generate ./tools/javy-x86_64-macos-v1.1.2 compile wasm_compiler.js -o compiler.wasm

//go:embed compiler.wasm
var wasmCompiler []byte

func LoadWASM(r wazero.Runtime, outputFormat Generate) (Compiler, error) {
	compiledWasm, err := r.CompileModule(context.Background(), wasmCompiler)
	if err != nil {
		return nil, err
	}

	c := &WASM{
		Dev:          true,
		OutputFormat: outputFormat,

		runtime: r,
		wasm:    compiledWasm,
	}

	return c, nil
}

type WASM struct {
	Dev          bool
	OutputFormat Generate

	runtime wazero.Runtime
	wasm    wazero.CompiledModule
}

func (c WASM) Compile(path string, code []byte) (*CompileResult, error) {
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

	buf := new(bytes.Buffer)

	e := json.NewEncoder(buf)
	e.SetEscapeHTML(false)

	if err := e.Encode(&opts); err != nil {
		return nil, err
	}

	// TODO: run the code
	// instance, err := c.runtime.InstantiateModule(
	// 	context.Background(),
	// 	c.wasm,
	// 	wazero.NewModuleConfig().WithName(""),
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// g := instance.ExportedFunction("what name?")

	// _ = g

	return nil, nil
}
