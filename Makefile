.PHONY: build-wasm serve clean

build-wasm:
	GOOS=js GOARCH=wasm go build -o web-react/public/g8emu.wasm ./cmd/wasm/main.go
	# @cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js web/
	@cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js web-react/public/

serve:
	@echo "Starting server at http://localhost:8080"
	@python -m http.server 8080 -d web/

clean:
	@rm -f web/g8emu.wasm web/wasm_exec.js
