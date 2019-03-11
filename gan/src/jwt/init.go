package jwt

import (
	"sync"

	"github.com/google/uuid"
	"github.com/sdeoras/jwt"
)

var (
	once    sync.Once
	Manager jwt.Manager
)

func init() {
	once.Do(func() {
		Manager = jwt.NewManager(uuid.New().String())
	})
}
