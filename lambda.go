package lambda

import (
	"lambda/src/email"
	"lambda/src/health"
	"lambda/src/infer"
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
	routeEmail  = "/email"
	routeInfer  = "/infer"
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
		router.Register(routeEmail, email.Send)
		router.Register(routeInfer, infer.InferImage)
	})
}
