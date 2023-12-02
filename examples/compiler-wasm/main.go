package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dvcorreia/go-svelte"
	"github.com/tetratelabs/wazero"
)

func main() {
	r := wazero.NewRuntime(context.Background())
	defer r.Close(context.Background())

	compiler, err := svelte.LoadWASM(r, svelte.GenerateSSR)
	if err != nil {
		log.Fatal(err)
	}

	code, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	ouput, err := compiler.Compile("test.svelte", code)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(os.Stdout, string(ouput.JS)) // ignore CSS
}
