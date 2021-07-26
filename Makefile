.PHONY: serve wasm

web/wasm_exec.js:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" web

web/gowasm/static/main.wasm: web/gowasm/main.go
	GOOS=js GOARCH=wasm go build -o web/gowasm/static/main.wasm web/gowasm/main.go

wasm: web/wasm_exec.js web/gowasm/static/main.wasm
	go mod tidy && go mod vendor
	@echo "#### done ####"

serve: wasm
	go run main.go
