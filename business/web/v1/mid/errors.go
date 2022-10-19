package mid

import (
	"context"
	"net/http"

	v1Web "github.com/pasdeta/go_service/business/web/v1"
	"github.com/pasdeta/go_service/foundation/web"
	"go.uber.org/zap"
)

func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, err := web.GetValues(ctx)
			if err != nil {

				return web.NewShutdownError("web value missing from context")
			}

			// Run the next handler and catch any propagated error.
			if err := handler(ctx, w, r); err != nil {

				// Log the error.
				log.Errorw("ERROR", "trace_id", v.TraceID, "message", err)

				// Build out the error response.
				var er v1Web.ErrorResponse
				var status int
				switch {
				case v1Web.IsRequestError(err):
					reqErr := v1Web.GetRequestError(err)
					er = v1Web.ErrorResponse{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				default:
					er = v1Web.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				// Respond with the error back to the client.
				if err := web.Respond(ctx, w, er, status); err != nil {

					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {

					return err
				}
			}

			return nil
		}

		return f
	}

	return m
}
