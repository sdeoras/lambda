package jwt

import (
	"lambda/src/env"
	"sync"

	"github.com/sdeoras/jwt"
)

var (
	once    sync.Once
	Manager jwt.Manager
)

// init initializes secret key based on environment variable and creates a new
// jwt token Manager. It also registers some functions to route the traffic to.
func init() {
	once.Do(func() {
		Manager = jwt.NewManager(env.JwtSecret,
			jwt.EnforceExpiration())
	})
}
