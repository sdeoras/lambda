// handlers_test.go
package payload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	handlerFuncA = iota
	handlerFuncB
)

func handlerToken(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "req body is nil", http.StatusBadRequest)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	// Set some claims
	claims["foo"] = "bar"
	//claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, _ = fmt.Fprintf(w, "%s", tokenString)
}

func handlerA(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "req body is nil", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	_, _ = fmt.Fprintf(w, "first called: %s", string(b))
}

func handlerB(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "req body is nil", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	_, _ = fmt.Fprintf(w, "second called: %s", string(b))
}

// TestFirstFunc tests an internally routed func.
func TestFirstFunc(t *testing.T) {
	registry[handlerFuncA] = handlerA
	registry[handlerFuncB] = handlerB

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	// Set some claims
	claims["foo"] = "bar"
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		t.Fatal(err)
	}

	data := new(Payload)
	data.FuncData = new(FuncData)
	data.FuncData.Id = handlerFuncA
	data.FuncData.Data = []byte("this is a test")
	data.TokenString = tokenString

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Route)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "first called: this is a test"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestSecondFunc tests an internally routed function.
func TestSecondFunc(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	// Set some claims
	claims["foo"] = "bar"
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		t.Fatal(err)
	}

	data := new(Payload)
	data.FuncData = new(FuncData)
	data.FuncData.Id = handlerFuncB
	data.FuncData.Data = []byte("this is a test")
	data.TokenString = tokenString

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Route)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "second called: this is a test"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestHandlerToken tests jwt token generation and parsing
func TestHandlerToken(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", bytes.NewReader([]byte("this is a test")))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerToken)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	tokenString := rr.Body.String()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if token == nil {
		t.Fatal("nil token")
	}

	if !token.Valid {
		t.Fatal("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("invalid claims")
	}

	if v, ok := claims["foo"]; !ok {
		t.Fatal("claim key not found")
	} else {
		if bar, ok := v.(string); !ok {
			t.Fatal("claim value is not a string")
		} else {
			if bar != "bar" {
				t.Fatal("claim value expected bar, got", bar)
			}
		}
	}
}

// TestRoute tests how to route after authenticating
func TestRoute(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	// Set some claims
	claims["user"] = "name"
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		t.Fatal(err)
	}

	data := new(Payload)
	data.FuncData = new(FuncData)
	data.FuncData.Id = HandlerHelloWorld
	data.FuncData.Data = []byte("i am authenticated with jwt!!")
	data.TokenString = tokenString

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "https://us-central1-"+
		os.Getenv("GCLOUD_PROJECT_NAME")+
		".cloudfunctions.net/router",
		bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	if resp.Status != "200 OK" {
		t.Fatal("expected status 200 OK, got:", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("response Body:", string(body))
}

// TestAuthOnly tests how to authenticate without routing
func TestAuthOnly(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	// Set some claims
	claims["user"] = "name"
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		t.Fatal(err)
	}

	data := new(Payload)
	data.FuncData = new(FuncData)
	data.FuncData.Id = HandlerAuthOnly
	data.TokenString = tokenString

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "https://us-central1-"+
		os.Getenv("GCLOUD_PROJECT_NAME")+
		".cloudfunctions.net/router",
		bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	if resp.Status != "200 OK" {
		t.Fatal("expected status 200 OK, got:", resp.Status)
	}
}
