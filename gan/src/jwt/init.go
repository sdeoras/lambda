package jwt

import (
	"gan/src/config"
	"sync"
	"time"

	"github.com/sdeoras/jwt"
)

var (
	once    sync.Once
	Manager jwt.Manager
)

func init() {
	once.Do(func() {
		Manager = jwt.NewManager(config.Config.JwtSecret,
			jwt.EnforceExpiration(),      // on the server side ensure jwt token has expiry
			jwt.SetLifeSpan(time.Second), // on the client side put expiry in jwt token
		)
	})
}
