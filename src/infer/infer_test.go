package infer

import (
	"fmt"
	"io/ioutil"
	"lambda/src/jwt"
	"net/http"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/sdeoras/lambda/api"
)

// TestInfer_Remote expects google cloud function to be up and running and it tests against that.
func TestInfer_Remote(t *testing.T) {
	b, err := ioutil.ReadFile("/Users/sdeoras/Downloads/sentinel.jpg")
	if err != nil {
		t.Fatal(err)
	}

	request := new(api.InferImageRequest)
	request.Images = make([]*api.Image, 1)
	request.Images[0] = new(api.Image)
	request.Images[0].Name = "xyz"
	request.Images[0].Data = b
	request.ModelName = "garage-door-checker"
	request.ModelVersion = "v1"

	b, err = proto.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	req, err := jwt.Requestor.Request(http.MethodPost, "https://"+os.Getenv("GOOGLE_GCF_DOMAIN")+
		"/"+ProjectName+"/"+Name, nil, b)
	req.Method = http.MethodPost

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

	response := new(api.InferImageResponse)
	if err := proto.Unmarshal(b, response); err != nil {
		t.Fatal(err)
	}

	fmt.Println("label:", response.Outputs[0].Label)
}
