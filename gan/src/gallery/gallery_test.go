package gallery

import (
	"fmt"
	"gan/src/env"
	"gan/src/jwt"
	"gan/src/route"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/sdeoras/api"
)

// TestGen_Remote expects google cloud function to be up and running and it tests against that.
func TestGen_Remote(t *testing.T) {
	request := new(api.GanRequest)
	request.ModelName = "gan-mnist-generator"
	request.ModelVersion = "v1"
	request.Count = 2

	b, err := proto.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	req, err := jwt.Manager.Request(http.MethodPost, "https://"+
		filepath.Join(env.Domain, env.FuncName, route.Gallery),
		nil, b)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("%s:%s. Mesg:%s", "expected status 200 OK, got", resp.Status, string(b))
	}

	fmt.Println(string(b))
}
