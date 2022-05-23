wasm:
	GOOS=js GOARCH=wasm go build -o demo/public/main.wasm ./sdk-wasm/wasm.go