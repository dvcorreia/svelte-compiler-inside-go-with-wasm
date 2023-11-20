## Notes
 
### Running on Nix currently does not work (on MacOS at least) due to a strange error!

Running:

```console
CGO_ENABLED=1 go test -v -run ^TestWasmDOM$ github.com/dvcorreia/go-svelte
```

Gives:

```console
# github.com/dvcorreia/go-svelte.test
/nix/store/8zv2194v9gl0hlx2ac8p91ihnjwdwm6g-go-1.20.10/share/go/pkg/tool/darwin_amd64/link: running clang++ failed: exit status 1
Undefined symbols for architecture x86_64:
  "_futimens", referenced from:
      _libc_futimens_trampoline in go.o
ld: symbol(s) not found for architecture x86_64
clang-11: error: linker command failed with exit code 1 (use -v to see invocation)

FAIL    github.com/dvcorreia/go-svelte [build failed]
FAIL
```

It works with my Go local install and I don't understand how to fix this with Nix packages.