package lambda

import (
	"lambda/email"
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

// init defines the routes to route traffic to.
func init() {
	once.Do(func() {
		router = httprouter.NewRouter()
		router.Register("/"+email.Name, email.Send)
		router.Register("/health/"+email.Name, email.Health)
	})
}
