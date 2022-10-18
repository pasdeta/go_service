// Package web contains a small web framework extension.
package web

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {

	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (app *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(app.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx = context.WithValue(ctx, key, &v)

		if err := handler(ctx, w, r); err != nil {
			// Dont know what to do here

			return
		}
		// Code After
	}
	app.ContextMux.Handle(method, path, h)
}
