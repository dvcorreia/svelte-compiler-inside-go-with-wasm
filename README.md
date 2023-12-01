## Results

```console
goos: darwin
goarch: amd64
pkg: github.com/dvcorreia/go-svelte
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
               │    v8.txt    │                 wasm.txt                 │
               │    sec/op    │     sec/op      vs base                  │
CompilersSSR-8   271.5µ ± 10%   58743.4µ ±  6%  +21537.52% (p=0.002 n=6)
CompilersDOM-8   2.223m ± 14%   117.804m ± 10%   +5200.00% (p=0.002 n=6)
geomean          776.8µ           83.19m        +10608.82%
```

You can install benchstat:

```console
go install golang.org/x/perf/cmd/...@latest
```

## Profiling

`go test -run=NONE -bench=^BenchmarkWasmSSR -count=5 -benchmem -cpuprofile cpuprofile.out`

`go tool pprof cpuprofile.out`

Running external code ~=64% of the time, so no optimize will fix the extremelly low performance we are seeing.

```console
Showing nodes accounting for 6800ms, 87.07% of 7810ms total
Dropped 141 nodes (cum <= 39.05ms)
Showing top 10 nodes out of 78
      flat  flat%   sum%        cum   cum%
    4970ms 63.64% 63.64%     4970ms 63.64%  runtime._ExternalCode
     370ms  4.74% 68.37%      370ms  4.74%  runtime.memmove
     340ms  4.35% 72.73%      340ms  4.35%  runtime.madvise
     280ms  3.59% 76.31%      420ms  5.38%  github.com/tetratelabs/wazero/internal/asm/amd64.(*AssemblerImpl).encodeMemoryToRegister
     230ms  2.94% 79.26%      240ms  3.07%  github.com/tetratelabs/wazero/internal/asm/amd64.appendUint32 (inline)
     160ms  2.05% 81.31%      370ms  4.74%  github.com/tetratelabs/wazero/internal/asm/amd64.(*AssemblerImpl).encodeRegisterToMemory
     150ms  1.92% 83.23%      240ms  3.07%  github.com/tetratelabs/wazero/internal/asm/amd64.appendConst (inline)
     140ms  1.79% 85.02%      140ms  1.79%  runtime.memclrNoHeapPointers
      80ms  1.02% 86.04%      230ms  2.94%  github.com/tetratelabs/wazero/internal/asm/amd64.(*AssemblerImpl).encodeConstToRegister
      80ms  1.02% 87.07%      200ms  2.56%  github.com/tetratelabs/wazero/internal/asm/amd64.(*AssemblerImpl).encodeRelativeJump
```

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