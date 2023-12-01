package svelte_test

import (
	"strings"
	"testing"

	"github.com/dvcorreia/go-svelte"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/matryer/is"
)

func BenchmarkV8SSR(b *testing.B) {
	benchmarkV8SSR(b)
}

func benchmarkV8SSR(b *testing.B) {
	vm, err := v8.Load()
	if err != nil {
		b.Fatal(err)
	}

	compiler, err := svelte.LoadV8(vm, svelte.GenerateSSR)
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

func BenchmarkV8DOM(b *testing.B) {
	benchmarkV8DOM(b)
}

func benchmarkV8DOM(b *testing.B) {
	vm, err := v8.Load()
	if err != nil {
		b.Fatal(err)
	}

	compiler, err := svelte.LoadV8(vm, svelte.GenerateDOM)
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

func TestV8SSR(t *testing.T) {
	is := is.New(t)
	vm, err := v8.Load()
	is.NoErr(err)
	compiler, err := svelte.LoadV8(vm, svelte.GenerateSSR)
	is.NoErr(err)
	r, err := compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(r.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(r.JS, `<h1>hi world!</h1>`))
}

func TestV8SSRRecovery(t *testing.T) {
	is := is.New(t)
	vm, err := v8.Load()
	is.NoErr(err)
	compiler, err := svelte.LoadV8(vm, svelte.GenerateSSR)
	is.NoErr(err)
	ssr, err := compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1></h1>`))
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), `</h1> attempted to close an element that was not open`))
	is.True(strings.Contains(err.Error(), `<h1>hi world!</h1></h1`))
	is.Equal(ssr, nil)
	ssr, err = compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(ssr.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(ssr.JS, `<h1>hi world!</h1>`))
}

func TestV8DOM(t *testing.T) {
	is := is.New(t)
	vm, err := v8.Load()
	is.NoErr(err)
	compiler, err := svelte.LoadV8(vm, svelte.GenerateDOM)
	is.NoErr(err)
	dom, err := compiler.Compile("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(dom.JS, `from "svelte/internal"`))
	is.True(strings.Contains(dom.JS, `function create_fragment`))
	is.True(strings.Contains(dom.JS, `element("h1")`))
	is.True(strings.Contains(dom.JS, `text("hi world!")`))
}

func TestV8DOMRecovery(t *testing.T) {
	is := is.New(t)
	vm, err := v8.Load()
	is.NoErr(err)
	compiler, err := svelte.LoadV8(vm, svelte.GenerateDOM)
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
