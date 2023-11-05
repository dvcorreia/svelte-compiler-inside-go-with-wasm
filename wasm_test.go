package svelte

import (
	"context"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/tetratelabs/wazero"
)

func TestWasmDOM(t *testing.T) {
	is := is.New(t)
	r := wazero.NewRuntime(context.Background())
	compiler, err := LoadWASM(r, GenerateDOM)
	is.NoErr(err)
	dom, err := compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(dom.JS, `from "svelte/internal"`))
	is.True(strings.Contains(dom.JS, `function create_fragment`))
	is.True(strings.Contains(dom.JS, `element("h1")`))
	is.True(strings.Contains(dom.JS, `text("hi world!")`))
}
