package svelte

import (
	"encoding/json"
	"fmt"

	"github.com/livebud/bud/package/js"
)

func LoadV8(vm js.VM) (*V8, error) {
	if err := vm.Script("svelte/compiler.js", compiler); err != nil {
		return nil, err
	}
	// TODO make dev configurable
	return &V8{vm, true}, nil
}

var _ Compiler = (*V8)(nil)

type V8 struct {
	VM  js.VM
	Dev bool
}

// Compile server-rendered code
func (c *V8) SSR(path string, code []byte) (*SSR, error) {
	expr := fmt.Sprintf(`;__svelte__.compile({ "path": %q, "code": %q, "target": "ssr", "dev": %t, "css": false })`, path, code, c.Dev)
	result, err := c.VM.Eval(path, expr)
	if err != nil {
		return nil, err
	}
	out := new(SSR)
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return nil, err
	}
	return out, nil
}

type DOM struct {
	JS  string
	CSS string
}

// Compile DOM code
func (c *V8) DOM(path string, code []byte) (*DOM, error) {
	expr := fmt.Sprintf(`;__svelte__.compile({ "path": %q, "code": %q, "target": "dom", "dev": %t, "css": true })`, path, code, c.Dev)
	result, err := c.VM.Eval(path, expr)
	if err != nil {
		return nil, err
	}
	out := new(DOM)
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return nil, err
	}
	return out, nil
}
