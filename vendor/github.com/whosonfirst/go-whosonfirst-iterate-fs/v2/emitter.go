package fs

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	
	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/filters"
)

// FSEmitter implements the `Emitter` interface for crawling records in a fs.
type FSEmitter struct {
	emitter.Emitter
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
	fs      fs.FS
}

func NewFSEmitter(ctx context.Context, uri string, iterator_fs fs.FS) (emitter.Emitter, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create filters from query, %w", err)
	}

	idx := &FSEmitter{
		fs:      iterator_fs,
		filters: f,
	}

	return idx, nil
}

// WalkURI() walks (crawls) the fs named 'uri' and for each file (not excluded by any filters specified
// when `idx` was created) invokes 'index_cb'.
func (idx *FSEmitter) WalkURI(ctx context.Context, index_cb emitter.EmitterCallbackFunc, uri string) error {

	log.Printf("Walk %s\n", uri)
	
	var walk_func func(path string, d fs.DirEntry, err error) error

	walk_func = func(path string, d fs.DirEntry, err error) error {

		log.Printf("Path %s\n", path)
			
		if err != nil {
			return fmt.Errorf("Failed to walk %s, %w", path, err)
		}

		if d.IsDir() {
			return nil
		}

		r, err := idx.fs.Open(path)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", path, err)
		}

		rsc, err := ioutil.NewReadSeekCloser(r)

		if err != nil {
			return fmt.Errorf("Failed to create ReadSeekCloser for %s, %w", path, err)
		}

		defer rsc.Close()

		if idx.filters != nil {

			ok, err := idx.filters.Apply(ctx, rsc)

			if err != nil {
				return fmt.Errorf("Failed to apply filters for '%s', %w", path, err)
			}

			if !ok {
				return nil
			}

			_, err = rsc.Seek(0, 0)

			if err != nil {
				return fmt.Errorf("Failed to seek(0, 0) on reader for '%s', %w", path, err)
			}
		}

		err = index_cb(ctx, path, rsc)

		if err != nil {
			return fmt.Errorf("Failed to invoke callback for '%s', %w", path, err)
		}

		return nil
	}

	err := fs.WalkDir(idx.fs, uri, walk_func)

	if err != nil {
		return fmt.Errorf("Failed to walk filesystem, %w", err)
	}

	return nil
}
