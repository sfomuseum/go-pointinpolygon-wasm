package fs

import (
	"context"
	"fmt"
	"io/fs"
	"net/url"

	"github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
)

func NewIteratorWithFS(ctx context.Context, emitter_uri string, emitter_cb emitter.EmitterCallbackFunc, iterator_fs fs.FS) (*iterator.Iterator, error) {

	fs_emitter, err := NewFSEmitter(ctx, emitter_uri, iterator_fs)

	if err != nil {
		return nil, fmt.Errorf("Failed to create FS emitter, %w", err)
	}

	null_u, _ := url.Parse(emitter_uri)
	null_u.Scheme = "null"

	iter, err := iterator.NewIterator(ctx, null_u.String(), emitter_cb)

	if err != nil {
		return nil, fmt.Errorf("Failed to create stub iterator, %w", err)
	}

	iter.Emitter = fs_emitter
	return iter, nil
}
