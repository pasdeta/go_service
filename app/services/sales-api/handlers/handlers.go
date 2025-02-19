// Package handlers manages the different versions of the API.
package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pasdeta/go_service/app/services/sales-api/handlers/probegrp"
	"github.com/pasdeta/go_service/business/web/auth"
	"github.com/pasdeta/go_service/business/web/v1/mid"
	"github.com/pasdeta/go_service/foundation/web"
	"go.uber.org/zap"
)

// Options represent optional parameters.
type Options struct {
	corsOrigin string
}

// WithCORS provides configuration options for CORS.
func WithCORS(origin string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origin
	}
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Build    string
	Auth     *auth.Auth
	DB       *sqlx.DB
	// Tracer   trace.Tracer
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	authen := mid.Authenticate(cfg.Auth)
	admin := mid.Authorize(auth.RoleAdmin)

	probegrp := probegrp.Handlers{
		Log:   cfg.Log,
		Build: cfg.Build,
		DB:    cfg.DB,
	}
	app.Handle(http.MethodGet, "/liveness", probegrp.Liveness)
	app.Handle(http.MethodGet, "/readiness", probegrp.Readiness)

	app.Handle(http.MethodGet, "/test400", probegrp.TestError400)
	app.Handle(http.MethodGet, "/test500", probegrp.TestError500)
	app.Handle(http.MethodGet, "/testpanic", probegrp.TestPanic)

	app.Handle(http.MethodGet, "/testauth", probegrp.TestAuth, authen, admin)

	return app
}

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}
