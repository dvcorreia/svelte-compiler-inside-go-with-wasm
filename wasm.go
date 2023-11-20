package svelte

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
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

	wasi_snapshot_preview1.MustInstantiate(context.Background(), r)

	c := &WASM{
		Dev:          true,
		OutputFormat: outputFormat,

		r:            r,
		compiledWasm: compiledWasm,
	}

	return c, nil
}

type WASM struct {
	Dev          bool
	OutputFormat Generate

	r            wazero.Runtime
	compiledWasm wazero.CompiledModule
}

func (c *WASM) Compile(path string, code []byte) (*CompileResult, error) {
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

	bufIn, bufOut := new(bytes.Buffer), new(bytes.Buffer)

	eIn := json.NewEncoder(bufIn)
	eIn.SetEscapeHTML(false)

	if err := eIn.Encode(&opts); err != nil {
		return nil, err
	}

	eOut := json.NewDecoder(bufOut)

	configuration := wazero.NewModuleConfig().
		WithStdout(bufOut).
		WithStdin(bufIn).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime()

	_, err := c.r.InstantiateModule(
		context.Background(),
		c.compiledWasm,
		configuration,
	)
	if err != nil {
		return nil, err
	}

	var result CompileResult
	if err := eOut.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
