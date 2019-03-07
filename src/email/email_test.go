package email

import (
	"io/ioutil"
	"lambda/src/jwt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/sdeoras/api"
)

// TestSend_Local tests using locally spawned http test server.
func TestSend_Local(t *testing.T) {
	sendRequest := &api.EmailRequest{
		FromName:  os.Getenv("EMAIL_FROM_NAME"),
		FromEmail: os.Getenv("EMAIL_FROM_EMAIL"),
		ToName:    os.Getenv("EMAIL_TO_NAME"),
		ToEmail:   os.Getenv("EMAIL_TO_EMAIL"),
		Subject:   "email test via http endpoint",
		Body:      []byte("<strong>to check if all is going through</strong>"),
	}

	b, err := proto.Marshal(sendRequest)
	if err != nil {
		t.Fatal(err)
	}

	req, err := jwt.Requestor.Request(http.MethodPost, "/", nil, b)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Send)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	b, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	sendResponse := new(api.EmailResponse)
	if err := proto.Unmarshal(b, sendResponse); err != nil {
		t.Fatal(err)
	}

	if sendResponse.StatusCode != 202 {
		t.Fatal("sending email failed with status code:", sendResponse.StatusCode)
	}
}

// TestSend_Remote expects google cloud function to be up and running and it tests against that.
func TestSend_Remote(t *testing.T) {
	sendRequest := &api.EmailRequest{
		FromName:  os.Getenv("EMAIL_FROM_NAME"),
		FromEmail: os.Getenv("EMAIL_FROM_EMAIL"),
		ToName:    os.Getenv("EMAIL_TO_NAME"),
		ToEmail:   os.Getenv("EMAIL_TO_EMAIL"),
		Subject:   "email test via http endpoint",
		Body:      []byte("<strong>to check if all is going through</strong>"),
	}

	b, err := proto.Marshal(sendRequest)
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

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("%s:%s", "expected status 200 OK, got", resp.Status)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	sendResponse := new(api.EmailResponse)
	if err := proto.Unmarshal(b, sendResponse); err != nil {
		t.Fatal(err)
	}

	if sendResponse.StatusCode != 202 {
		t.Fatal("sending email failed with status code:", sendResponse.StatusCode)
	}
}
