package lambda

import (
	"lambda/src/email"
	"lambda/src/infer"
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
	routeEmail = "email"
	routeInfer = "infer"
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

		h := health.NewProvider(health.OutputProto)
		h.Register(routeEmail, nil)
		h.Register(routeInfer, nil)

		router = httprouter.NewRouter()
		// register health check endpoint
		router.Register(health.StdRoute, h.NewHTTPHandler())

		// register services
		router.Register(f(routeEmail), email.Send)
		router.Register(f(routeInfer), infer.InferImage)
	})
}
