GOMOD=vendor

wasm:
	GOOS=js GOARCH=wasm go build -mod $(GOMOD) -ldflags="-s -w" -o cmd/example/sfomuseum_pointinpolygon.wasm cmd/pointinpolygon/main.go

example:
	go run -mod $(GOMOD) cmd/example/main.go
