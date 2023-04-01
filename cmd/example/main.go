package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/sfomuseum/go-http-wasm"
)

//go:embed index.html css/*.css javascript/*.js wasm/*.wasm
var FS embed.FS

func main() {

	host := flag.String("host", "localhost", "The host name to listen for requests on")
	port := flag.Int("port", 8080, "The host port to listen for requests on")

	flag.Parse()

	mux := http.NewServeMux()

	wasm_opts := wasm.DefaultWASMOptions()
	
	err := wasm.AppendAssetHandlers(mux)

	if err != nil {
		log.Fatalf("Failed to append wasm assets handler, %v", err)
	}

	http_fs := http.FS(FS)
	example_handler := http.FileServer(http_fs)

	example_handler = wasm.AppendResourcesHandler(example_handler, wasm_opts)

	mux.Handle("/", example_handler)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Listening for requests on %s\n", addr)

	err = http.ListenAndServe(addr, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
