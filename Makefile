bench: bench/v8 bench/wasm
bench/v8:
	go test -run=NONE -bench=^BenchmarkCompilers -count=6 -sveltec=v8 > v8.txt
bench/wasm:
	go test -run=NONE -bench=^BenchmarkCompilers -count=6 -sveltec=wasm > wasm.txt
cmp:
	benchstat v8.txt wasm.txt
