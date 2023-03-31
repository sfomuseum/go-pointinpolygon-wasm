package wasm

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"time"

	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-feature/geometry"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
)

func NewSpatialDatabase(ctx context.Context, spatial_database_uri string, spatial_database_fs fs.FS) (database.SpatialDatabase, error) {

	log.Println("Create new database")

	db, err := database.NewSpatialDatabase(ctx, spatial_database_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create spatial db, %v", err)
	}

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		t1 := time.Now()

		defer func() {
			log.Printf("Time to index %s, %v\n", path, time.Since(t1))
		}()

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		geom_type, err := geometry.Type(body)

		if err != nil {
			return fmt.Errorf("Failed to determine geometry type for %s, %w", path, err)
		}

		switch geom_type {
		case "Polygon", "MultiPolygon":
			// okay
		default:
			return nil
		}

		err = db.IndexFeature(ctx, body)

		if err != nil {
			return fmt.Errorf("Failed to index %s, %w", path, err)
		}

		return nil
	}

	walk_func := func(path string, d fs.DirEntry, err error) error {

		if d.IsDir() {
			return nil
		}

		r, err := spatial_database_fs.Open(path)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", path, err)
		}

		rsc, err := ioutil.NewReadSeekCloser(r)

		if err != nil {
			return fmt.Errorf("Failed to create ReadSeekCloser for %s, %w", path, err)
		}

		defer rsc.Close()

		return iter_cb(ctx, path, rsc)
	}

	t1 := time.Now()

	defer func() {
		log.Printf("Time to index data, %v\n", time.Since(t1))
	}()

	err = fs.WalkDir(spatial_database_fs, "data", walk_func)

	if err != nil {
		return nil, fmt.Errorf("Failed to walk data, %w", err)
	}

	return db, nil
}
