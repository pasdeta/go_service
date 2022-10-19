package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/pasdeta/go_service/business/web/metrics"
	"github.com/pasdeta/go_service/foundation/web"
)

func Panics() web.Middleware {

	m := func(handler web.Handler) web.Handler {
		f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {

					// Stack trace will be provided.
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

					metrics.AddPanics(ctx)
				}
			}()

			// Call the next handler and set its return value in the err variable.
			return handler(ctx, w, r)
		}

		return f
	}

	return m
}
