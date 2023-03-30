# go-pointinpolygon-wasm

Experimental work to embed a whosonfirst/go-whosonfirst-spatial/database.SpatialDatabase instance and all its data in a WebAssembly binary.

Does not work yet. Specifically it is blocking on calls to

```
rsp, err := db.PointInPolygon(ctx, coord, f)
```

My current thinking is that this might have to do with channels and wait groups but I have not tested that yet.

## See also

* https://github.com/sfomuseum/go-pmtiles-wasm