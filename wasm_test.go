package svelte_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/dvcorreia/go-svelte"
	"github.com/matryer/is"
	"github.com/tetratelabs/wazero"
)

func BenchmarkWasmSSR(b *testing.B) {
	benchmarkWasmSSR(b)
}

func benchmarkWasmSSR(b *testing.B) {
	r := wazero.NewRuntime(context.Background())
	defer r.Close(context.Background())

	compiler, err := svelte.LoadWASM(r, svelte.GenerateSSR)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := compiler.Compile("test.svelte", benchCode)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWasmDOM(b *testing.B) {
	benchmarkWasmDOM(b)
}

func benchmarkWasmDOM(b *testing.B) {
	r := wazero.NewRuntime(context.Background())
	defer r.Close(context.Background())

	compiler, err := svelte.LoadWASM(r, svelte.GenerateDOM)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := compiler.Compile("test.svelte", benchCode)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestWasmSSR(t *testing.T) {
	is := is.New(t)
	r := wazero.NewRuntime(context.Background())
	defer r.Close(context.Background())
	compiler, err := svelte.LoadWASM(r, svelte.GenerateSSR)
	is.NoErr(err)
	res, err := compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(res.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(res.JS, `<h1>hi world!</h1>`))
}

func TestWasmSSRRecovery(t *testing.T) {
	t.Skip("currently it can not recover from errors")
	is := is.New(t)
	r := wazero.NewRuntime(context.Background())
	defer r.Close(context.Background())
	compiler, err := svelte.LoadWASM(r, svelte.GenerateSSR)
	is.NoErr(err)
	ssr, err := compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1></h1>`))
	is.True(err != nil)
	fmt.Println(err.Error())
	is.True(strings.Contains(err.Error(), `</h1> attempted to close an element that was not open`))
	is.True(strings.Contains(err.Error(), `<h1>hi world!</h1></h1`))
	is.Equal(ssr, nil)
	ssr, err = compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(ssr.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(ssr.JS, `<h1>hi world!</h1>`))
}

func TestWasmDOM(t *testing.T) {
	is := is.New(t)
	r := wazero.NewRuntime(context.Background())
	defer r.Close(context.Background())
	compiler, err := svelte.LoadWASM(r, svelte.GenerateDOM)
	is.NoErr(err)
	dom, err := compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(dom.JS, `from "svelte/internal"`))
	is.True(strings.Contains(dom.JS, `function create_fragment`))
	is.True(strings.Contains(dom.JS, `element("h1")`))
	is.True(strings.Contains(dom.JS, `text("hi world!")`))
}

func TestWasmDOMRecovery(t *testing.T) {
	t.Skip("currently it can not recover from errors")
	is := is.New(t)
	r := wazero.NewRuntime(context.Background())
	defer r.Close(context.Background())
	compiler, err := svelte.LoadWASM(r, svelte.GenerateDOM)
	is.NoErr(err)
	dom, err := compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1></h1>`))
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), `</h1> attempted to close an element that was not open`))
	is.True(strings.Contains(err.Error(), `<h1>hi world!</h1></h1`))
	is.Equal(dom, nil)
	dom, err = compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(dom.JS, `from "svelte/internal"`))
	is.True(strings.Contains(dom.JS, `function create_fragment`))
	is.True(strings.Contains(dom.JS, `element("h1")`))
	is.True(strings.Contains(dom.JS, `text("hi world!")`))
}
