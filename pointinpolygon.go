package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"syscall/js"

	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
)

func PointInPolygonFunc(db database.SpatialDatabase) js.Func {

	log.Println("HELLO")
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		log.Println("WORLD")		
		// str_pt := args[0].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			log.Println("WOO")
			
			resolve := args[0]
			reject := args[1]

			go func() {

				log.Println("OMGWTF")
				
				ctx := context.Background()

				log.Println("STEP 1")
				q := url.Values{}

				log.Println("STEP 2")				
				f, err := filter.NewSPRFilterFromQuery(q)

				if err != nil {
					log.Println("SAD", 1, err)
					reject.Invoke(fmt.Sprintf("Failed to create filter, %v", err))
					return
				}

				log.Println("STEP 3")				
				coord, err := geo.NewCoordinate(-122.3830005850888, 37.61714514085811)

				if err != nil {
					log.Println("SAD", 2, err)
					reject.Invoke(fmt.Sprintf("Failed to create new coord, %v", err))
					return
				}

				// Blocks here...
				log.Println("STEP 4", coord)
				rsp, err := db.PointInPolygon(ctx, coord, f)

				log.Println("STEP 4a")
				
				if err != nil {
					log.Println("SAD", 3, err)
					reject.Invoke(fmt.Sprintf("Failed to do PIP, %v", err))
					return
				}

				log.Println("STEP 5")
				enc_rsp, err := json.Marshal(rsp)

				if err != nil {
					log.Println("SAD", 4, err)
					reject.Invoke(err.Error())
					return
				}

				log.Println("OK")
				resolve.Invoke(string(enc_rsp))
				return
			}()

			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
