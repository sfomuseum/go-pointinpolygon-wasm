# go-whosonfirst-iterate-fs

Go package implementing the `go-whosonfirst-iterate/v2` interfaces for iterating through `io/fs.FS` instances.

## Documentation

Documentation is incomplete at this time. For the time being have a look at the [emitter_test.go](emitter_test.go) and [iterator_test](iterator_test.go) tests.

## Important

This package highlights many of the shortcomings with the way that the `whosonfirst/go-whosonfirst-iterate/v2` package is structured. Hopefully these problems will be addressed in an as-yet unreleased "v3" package. In the meantime it's kind of "smelly" but it works...

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://pkg.go.dev/io/fs