package main

import (
	"context"
	"log"
	"syscall/js"

	"github.com/sfomuseum/go-pointinpolygon-wasm"
	"github.com/sfomuseum/go-pointinpolygon-wasm/static"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-rtree"
)

func main() {

	ctx := context.Background()

	db_uri := "rtree://?strict=false"

	db, err := wasm.NewSpatialDatabase(ctx, db_uri, static.FS)

	if err != nil {
		log.Fatalf("Failed to create spatial database, %v", err)
	}

	pip_func := wasm.PointInPolygonFunc(db)
	defer pip_func.Release()

	js.Global().Set("sfomuseum_pointinpolygon", pip_func)

	c := make(chan struct{}, 0)

	log.Printf("Point in polygon binary initialized")
	<-c

}
