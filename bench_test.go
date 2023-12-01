package svelte_test

import (
	_ "embed"
	"flag"
	"testing"
)

//go:embed benchmark-component.svelte
var benchCode []byte

var svelteC *string = flag.String("sveltec", "wasm", "the svelte compiler to benchmark (e.g. v8, wasm)")

func BenchmarkCompilersSSR(b *testing.B) {
	switch *svelteC {
	case "v8":
		benchmarkV8SSR(b)
	case "wasm":
		benchmarkWasmSSR(b)
	}
}

func BenchmarkCompilersDOM(b *testing.B) {
	switch *svelteC {
	case "v8":
		benchmarkV8DOM(b)
	case "wasm":
		benchmarkWasmDOM(b)
	}
}
