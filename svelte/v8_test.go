package svelte_test

import (
	"strings"
	"testing"

	"github.com/dvcorreia/svelte-compiler-inside-go/svelte"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/matryer/is"
)

func TestV8SSR(t *testing.T) {
	is := is.New(t)
	vm, err := v8.Load()
	is.NoErr(err)
	compiler, err := svelte.LoadV8(vm)
	is.NoErr(err)
	ssr, err := compiler.SSR("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(ssr.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(ssr.JS, `<h1>hi world!</h1>`))
}

func TestV8SSRRecovery(t *testing.T) {
	is := is.New(t)
	vm, err := v8.Load()
	is.NoErr(err)
	compiler, err := svelte.LoadV8(vm)
	is.NoErr(err)
	ssr, err := compiler.SSR("test.svelte", []byte(`<h1>hi world!</h1></h1>`))
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), `</h1> attempted to close an element that was not open`))
	is.True(strings.Contains(err.Error(), `<h1>hi world!</h1></h1`))
	is.Equal(ssr, nil)
	ssr, err = compiler.SSR("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(ssr.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(ssr.JS, `<h1>hi world!</h1>`))
}

func TestV8DOM(t *testing.T) {
	is := is.New(t)
	vm, err := v8.Load()
	is.NoErr(err)
	compiler, err := svelte.LoadV8(vm)
	is.NoErr(err)
	dom, err := compiler.DOM("test.svelte", []byte(`<h1>hi world!</h1>`))
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
	compiler, err := svelte.LoadV8(vm)
	is.NoErr(err)
	dom, err := compiler.DOM("test.svelte", []byte(`<h1>hi world!</h1></h1>`))
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), `</h1> attempted to close an element that was not open`))
	is.True(strings.Contains(err.Error(), `<h1>hi world!</h1></h1`))
	is.Equal(dom, nil)
	dom, err = compiler.DOM("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(dom.JS, `from "svelte/internal"`))
	is.True(strings.Contains(dom.JS, `function create_fragment`))
	is.True(strings.Contains(dom.JS, `element("h1")`))
	is.True(strings.Contains(dom.JS, `text("hi world!")`))
}

// TODO: test compiler.Dev = false
