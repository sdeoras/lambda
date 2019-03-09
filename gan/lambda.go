package gan

import (
	"gan/src/gen"
	"gan/src/jwt"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/sdeoras/health"
	"github.com/sdeoras/httprouter"
)

var (
	once   sync.Once
	router httprouter.Router
)

const (
	routeGen = "gen"
)

// Lambda is the main entry point. It immediately calls router and exits.
func Lambda(w http.ResponseWriter, r *http.Request) {
	router.Route(w, r)
}

// init defines the routes to route traffic to.
func init() {
	once.Do(func() {
		f := func(input string) string {
			return filepath.Join("/", input)
		}

		h := health.NewProvider(health.OutputProto, jwt.Manager, nil)
		h.Register(f(routeGen), nil)

		router = httprouter.NewRouter()
		// register health check endpoint
		router.Register(health.StdRoute, h.Provide())

		// register services
		router.Register(f(routeGen), gen.GenerateImages)
	})
}
