package lambda

import (
	"fmt"
	"lambda/email"
	"lambda/infer"
	"net/http"
	"sync"

	"github.com/sdeoras/httprouter"
)

var once sync.Once
var router httprouter.Router

// Lambda is the main entry point. It immediately calls router and exits.
func Lambda(w http.ResponseWriter, r *http.Request) {
	router.Route(w, r)
}

// Health returns ok
func Health(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "ok")
}

// init defines the routes to route traffic to.
func init() {
	once.Do(func() {
		router = httprouter.NewRouter()
		// register health check endpoint
		router.Register("/health", Health)

		// register services
		router.Register("/"+email.Name, email.Send)
		router.Register("/"+infer.Name, infer.InferImage)
	})
}
