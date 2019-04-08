package kv

import (
	"context"
	"fmt"
	"os"
	"user/src/config"

	"github.com/sdeoras/kv"
)

var (
	nameSpace = "userProfile"
)

var (
	Db kv.KV
	//once sync.Once
)

func init() {
	var err error
	Db, _, err = kv.NewDataStoreKv(context.Background(),
		os.Getenv("GCP_PROJECT"), nameSpace)
	if err != nil {
		panic(fmt.Sprintf("v8:%v:%s:%s", err, config.Config.ProjectId,
			os.Getenv("GCP_PROJECT")))
	}
}
