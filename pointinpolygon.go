package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"syscall/js"

	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
)

func PointInPolygonFunc(db database.SpatialDatabase) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		// str_pt := args[0].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {

				ctx := context.Background()
				q := url.Values{}
				
				f, err := filter.NewSPRFilterFromQuery(q)

				if err != nil {
					reject.Invoke(fmt.Sprintf("Failed to create filter, %v", err))
					return
				}

				coord, err := geo.NewCoordinate(-122.3830005850888, 37.61714514085811)

				if err != nil {
					reject.Invoke(fmt.Sprintf("Failed to create new coord, %v", err))
					return
				}

				// Blocks here...

				rsp, err := db.PointInPolygon(ctx, coord, f)

				if err != nil {
					reject.Invoke(fmt.Sprintf("Failed to do PIP, %v", err))
					return
				}

				enc_rsp, err := json.Marshal(rsp)

				if err != nil {
					reject.Invoke(err.Error())
					return
				}

				resolve.Invoke(string(enc_rsp))
				return
			}()

			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
