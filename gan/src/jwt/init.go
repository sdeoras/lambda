package jwt

import (
	"gan/src/config"
	"sync"
	"time"

	"github.com/sdeoras/jwt"
)

var (
	once    sync.Once
	manager jwt.Manager
)

// initialize initializes manager instance once per lifetime
func initialize() {
	once.Do(func() {
		manager = jwt.NewManager(config.Config().JwtSecret,
			jwt.EnforceExpiration(),      // on the server side ensure jwt token has expiry
			jwt.SetLifeSpan(time.Minute), // on the client side put expiry in jwt token
		)
	})
}

// Manager provides access to the instance
func Manager() jwt.Manager {
	initialize()
	return manager
}
