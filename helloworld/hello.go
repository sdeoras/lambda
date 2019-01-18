// Package p contains an HTTP Cloud Function.
package p

// #cgo LDFLAGS: -Llib
//import "C"

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "%s\n", os.Getenv("LD_LIBRARY_PATH"))
	_, _ = fmt.Fprintf(w, "%s\n", os.Getenv("CODE_LOCATION"))

	names, err := explore(os.Getenv("CODE_LOCATION"))
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v\n", err)
		return
	}

	for _, file := range names {
		_, _ = fmt.Fprintf(w, "%s\n", file)
	}

	b, err := ioutil.ReadFile("/srv/worker.go")
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v\n", err)
		return
	}
	_, _ = fmt.Fprintf(w, "------------\n%s\n----------------\n", string(b))

	b, err = exec.Command("/srv/files/bin/a.out").Output()
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v\n", err)
		return
	}
	_, _ = fmt.Fprintf(w, "------------\n%s\n----------------\n", string(b))

	// require github.com/tensorflow/tensorflow v1.12.0
	//tf "github.com/tensorflow/tensorflow/tensorflow/go"
	//_, _ = fmt.Fprintf(w, "%s\n", tf.Version())

	var d struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		_, _ = fmt.Fprintf(w, "%s\n", "Hello World 1 - v13!")
		return
	}
	if d.Message == "" {
		_, _ = fmt.Fprintf(w, "%s\n", "Hello World 2 - v12!")
		return
	}
	_, _ = fmt.Fprint(w, html.EscapeString(d.Message))
}

func explore(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, file := range files {
		if file.IsDir() {
			out2, err := explore(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			for i := range out2 {
				out = append(out, filepath.Join(file.Name(), out2[i]))
			}
		}
		out = append(out, file.Name())
	}

	return out, nil
}
