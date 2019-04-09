package info

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"user/src/config"
	"user/src/jwt"
	"user/src/route"

	"github.com/golang/protobuf/proto"
	"github.com/sdeoras/api/pb"
)

func TestRegister(t *testing.T) {
	if len(config.Config().Domain) == 0 ||
		len(config.Config().FuncName) == 0 {
		t.Fatal("not all env vars set properly")
	}

	userName := "spiderman"
	setRequest := &pb.SetUserInfoRequest{
		UserMeta: &pb.UserMeta{
			UserName:  userName,
			UserEmail: "spiderman@city.com",
			UserId:    0,
		},
		HealthCheckEndPoints: &pb.HealthCheckEndPoints{
			Url: []string{
				"http://url1",
				"http://url2",
			},
		},
	}

	b, err := proto.Marshal(setRequest)
	if err != nil {
		t.Fatal(err)
	}

	req, err := jwt.Manager().NewHTTPRequest(http.MethodPost, "https://"+
		filepath.Join(config.Config().Domain, config.Config().FuncName, route.Register),
		nil, b)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("%s:%s. Mesg:%s", "expected status 200 OK, got", resp.Status, string(b))
	}

	_ = resp.Body.Close()

	getRequest := &pb.GetUserInfoRequest{
		UserMeta: &pb.UserMeta{
			UserName: userName,
		},
	}

	b, err = proto.Marshal(getRequest)
	if err != nil {
		t.Fatal(err)
	}

	req, err = jwt.Manager().NewHTTPRequest(http.MethodPost, "https://"+
		filepath.Join(config.Config().Domain, config.Config().FuncName, route.Query),
		nil, b)

	client = &http.Client{}
	resp, err = client.Do(req)
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

	getResponse := new(pb.GetUserInfoResponse)

	if err := proto.Unmarshal(b, getResponse); err != nil {
		//fmt.Println(hex.EncodeToString(b))
		hex.Dumper(os.Stdout).Write(b)
		t.Fatal(err)
	}

	if jb, err := json.MarshalIndent(getResponse, "", "  "); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(string(jb))
	}
}

func TestQuery(t *testing.T) {
	if len(config.Config().Domain) == 0 ||
		len(config.Config().FuncName) == 0 {
		t.Fatal("not all env vars set properly")
	}

	userName := "spiderman"
	getRequest := &pb.GetUserInfoRequest{
		UserMeta: &pb.UserMeta{
			UserName: userName,
		},
	}

	b, err := proto.Marshal(getRequest)
	if err != nil {
		t.Fatal(err)
	}

	req, err := jwt.Manager().NewHTTPRequest(http.MethodPost, "https://"+
		filepath.Join(config.Config().Domain, config.Config().FuncName, route.Query),
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

	getResponse := new(pb.GetUserInfoResponse)

	if err := proto.Unmarshal(b, getResponse); err != nil {
		t.Fatal(err)
	}

	if jb, err := json.MarshalIndent(getResponse, "", "  "); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(string(jb))
	}
}
