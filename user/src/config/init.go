package config

import (
	"sync"

	"github.com/sdeoras/api/lambda/config"
	"github.com/sdeoras/api/pb"
)

var (
	once   sync.Once
	Config *pb.ConfigResponse
)

func init() {
	once.Do(func() {
		Config = config.NewConfigFromEnv()
	})
}
