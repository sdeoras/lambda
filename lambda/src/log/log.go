package log

import (
	"log"
	"os"
	"sync"
)

var (
	once sync.Once
	Out  *log.Logger
	Err  *log.Logger
)

func init() {
	once.Do(func() {
		Out = log.New(os.Stdout, "", 0)
		Err = log.New(os.Stderr, "", 0)
	})
}
