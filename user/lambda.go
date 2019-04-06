package user

import (
	"net/http"
	"path/filepath"
	"sync"
	"user/src/info"
	"user/src/route"

	"github.com/sdeoras/health"
	"github.com/sdeoras/httprouter"
)

var (
	once   sync.Once
	router httprouter.Router
)

// Lambda is the main entry point. It immediately calls router and exits.
func Lambda(w http.ResponseWriter, r *http.Request) {
	switch router.IsRegistered(r.URL.Path) {
	case true:
		router.Route(w, r)
	default:
		http.FileServer(http.Dir(filepath.Join(
			"src", route.Register))).ServeHTTP(w, r)
	}
}

// init defines the routes to route traffic to.
func init() {
	once.Do(func() {
		f := func(input string) string {
			return filepath.Join("/", input)
		}

		h := health.NewProvider(health.OutputProto)
		h.Register(route.Register, nil)
		h.Register(route.Query, nil)

		router = httprouter.NewRouter()
		// register health check endpoint
		router.Register(health.StdRoute, h.NewHTTPHandler())

		// register services
		router.Register(f(route.Register), info.Register)
		router.Register(f(route.Query), info.Query)
	})
}
