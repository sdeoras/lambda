package gan

import (
	"gan/src/gen"
	"gan/src/health"
	"net/http"
	"sync"

	"github.com/sdeoras/httprouter"
)

var (
	once   sync.Once
	router httprouter.Router
)

const (
	routeHealth = "/health"
	routeGen    = "/gen"
)

// Lambda is the main entry point. It immediately calls router and exits.
func Lambda(w http.ResponseWriter, r *http.Request) {
	router.Route(w, r)
}

// init defines the routes to route traffic to.
func init() {
	once.Do(func() {
		router = httprouter.NewRouter()
		// register health check endpoint
		router.Register(routeHealth, health.Check)

		// register services
		router.Register(routeGen, gen.GenerateImages)
	})
}
