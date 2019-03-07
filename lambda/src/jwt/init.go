package jwt

import (
	"os"
	"sync"

	"github.com/sdeoras/jwt"
)

var (
	once      sync.Once
	Validator jwt.Validator
	Requestor jwt.Requestor
)

// init initializes secret key based on environment variable and creates a new
// jwt token Validator. It also registers some functions to route the traffic to.
func init() {
	once.Do(func() {
		Validator = jwt.NewValidator(os.Getenv("JWT_SECRET_KEY"))
		Requestor = jwt.NewRequestor(os.Getenv("JWT_SECRET_KEY"))
	})
}
