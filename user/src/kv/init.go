package kv

import (
	"context"
	"fmt"
	"os"
	"sync"
	"user/src/config"

	"github.com/sdeoras/kv"
)

var (
	nameSpace = "userProfile"
)

var (
	kvdb kv.KV
	once sync.Once
)

// initialize initializes kvdb instance once per lifetime
func initialize() {
	once.Do(func() {
		var err error
		projectID, ok := os.LookupEnv("GCP_PROJECT")
		if !ok {
			panic("GCP_PROJECT env var not set")
		}

		kvdb, _, err = kv.NewDataStoreKv(context.Background(),
			projectID, nameSpace)
		if err != nil {
			mesg := fmt.Sprintf("%v:%s:%s", err, config.Config().ProjectId,
				os.Getenv("GCP_PROJECT"))
			panic(mesg)
		}
	})
}

// KV should be called to access instance since it is being initialized.
// Using init() func to initialize does not seem to work when env vars
// are required.
func KV() kv.KV {
	initialize()
	return kvdb
}
