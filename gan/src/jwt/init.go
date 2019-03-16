package jwt

import (
	"gan/src/env"
	"sync"

	"github.com/sdeoras/jwt"
)

var (
	once    sync.Once
	Manager jwt.Manager
)

func init() {
	once.Do(func() {
		Manager = jwt.NewManager(env.JwtSecret,
			jwt.EnforceExpiration(),
		)
	})
}
