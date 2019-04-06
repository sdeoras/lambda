package info

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"user/src/jwt"
	"user/src/kv"

	"github.com/golang/protobuf/proto"
	"github.com/sdeoras/api/pb"
)

func Register(w http.ResponseWriter, r *http.Request) {
	// validate input request
	if err := jwt.Manager.Validate(r); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:%s", http.StatusBadRequest, err.Error()),
			http.StatusBadRequest)
		return
	}

	// check method
	if r.Method != http.MethodPost {
		http.Error(w,
			fmt.Sprintf("%v:%s", http.StatusBadRequest,
				"error in gen.GenerateImages: method not set to POST"),
			http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v:error reading http request body:%v",
				http.StatusBadRequest, err),
			http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	request := new(pb.SetUserInfoRequest)
	if err := proto.Unmarshal(b, request); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not unmarshal image infer request:%v",
				http.StatusBadRequest, err),
			http.StatusBadRequest)
		return
	}

	if err := kv.Db.Set(request.UserMeta.UserName, b); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:error storing user info in kvdb:%v",
				http.StatusInternalServerError, err),
			http.StatusInternalServerError)
		return
	}
}

func Query(w http.ResponseWriter, r *http.Request) {
	// validate input request
	if err := jwt.Manager.Validate(r); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:%s", http.StatusBadRequest, err.Error()),
			http.StatusBadRequest)
		return
	}

	// check method
	if r.Method != http.MethodPost {
		http.Error(w,
			fmt.Sprintf("%v:%s", http.StatusBadRequest,
				"error in gen.GenerateImages: method not set to POST"),
			http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v:error reading http request body:%v",
				http.StatusBadRequest, err),
			http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	request := new(pb.GetUserInfoRequest)
	if err := proto.Unmarshal(b, request); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not unmarshal image infer request:%v",
				http.StatusBadRequest, err),
			http.StatusBadRequest)
		return
	}

	val, err := kv.Db.Get(request.UserMeta.UserName)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not get user info:%v",
				http.StatusInternalServerError, err),
			http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(val)
}
