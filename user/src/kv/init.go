package kv

import (
	"sync"

	"github.com/sdeoras/kv"
)

var (
	fileName  = "/tmp/bolt.db"
	nameSpace = "userProfile"
)

var (
	Db   kv.KV
	once sync.Once
)

func init() {
	once.Do(func() {
		var err error
		Db, _, err = kv.NewBoltKv(fileName, nameSpace)
		if err != nil {
			panic(err)
		}
	})
}
