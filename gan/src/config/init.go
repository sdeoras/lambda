package config

import (
	"sync"

	"github.com/sdeoras/api/lambda/config"
	"github.com/sdeoras/api/pb"
)

var (
	once sync.Once
	conf *pb.ConfigResponse
)

// initialize initializes config once per lifetime
func initialize() {
	once.Do(func() {
		conf = config.NewConfigFromEnv()
	})
}

func Config() *pb.ConfigResponse {
	initialize()
	return conf
}
