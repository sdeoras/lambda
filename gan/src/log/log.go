package log

import (
	"log"
	"os"
	"sync"
)

var (
	once sync.Once
	out  *log.Logger
	err  *log.Logger
)

// initialize initializes loggers once per lifetime
func initialize() {
	once.Do(func() {
		out = log.New(os.Stdout, "", 0)
		err = log.New(os.Stderr, "", 0)
	})
}

func Stdout() *log.Logger {
	initialize()
	return out
}

func Stderr() *log.Logger {
	initialize()
	return err
}
