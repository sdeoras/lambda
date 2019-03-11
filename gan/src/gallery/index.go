package gallery

import (
	"fmt"
	"gan/src/jwt"
	"gan/src/route"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/sdeoras/api"
)

func Show(w http.ResponseWriter, r *http.Request) {
	// validate input request
	err := jwt.Manager.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check method
	if r.Method != http.MethodPost {
		http.Error(w,
			"error in gallery.Show: method not set to POST",
			http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("error reading http request body:%v", err),
			http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	request := new(api.GalleryRequest)
	if err := proto.Unmarshal(b, request); err != nil {
		http.Error(w,
			fmt.Sprintf("could not unmarshal image infer request:%v", err),
			http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("src", route.Gallery, "index.html"))
	if err != nil {
		http.Error(w, fmt.Sprintf("%s:%v", "error creating new template", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, request); err != nil {
		http.Error(w, fmt.Sprintf("%s:%v", "error executing template", err), http.StatusInternalServerError)
		return
	}
}
