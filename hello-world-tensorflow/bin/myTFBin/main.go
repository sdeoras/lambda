// Package p contains an HTTP Cloud Function.
package main

import (
	"fmt"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", tf.Version())
}
