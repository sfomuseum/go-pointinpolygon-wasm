package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/whosonfirst/go-whosonfirst-spatial-pip"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
)

func PointInPolygonFunc(db database.SpatialDatabase) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		str_request := args[0].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {

				ctx := context.Background()

				var pip_request *pip.PointInPolygonRequest

				err := json.Unmarshal([]byte(str_request), &pip_request)

				if err != nil {
					reject.Invoke(fmt.Sprintf("Failed to parse request, %v", err))
					return
				}

				f, err := pip.NewSPRFilterFromPointInPolygonRequest(pip_request)

				if err != nil {
					reject.Invoke(fmt.Sprintf("Failed to create filter, %v", err))
					return
				}

				coord, err := geo.NewCoordinate(pip_request.Longitude, pip_request.Latitude)

				if err != nil {
					reject.Invoke(fmt.Sprintf("Failed to create new coord, %v", err))
					return
				}

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
