package payload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

// init initializes secret key based on environment variable and
// creates a new instance of handler function registry
func init() {
	once.Do(func() {
		secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
		registry = make(map[int]func(w http.ResponseWriter, r *http.Request))
		registry[HandlerHelloWorld] = HelloWorld
	})
}

// HelloWorld is called after authentication via Route func.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
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

	_, _ = fmt.Fprintf(w, "hello world called with: %s", string(b))
}

// Route authenticates request based on JWT token and routes request to registered
// function.
func Route(w http.ResponseWriter, r *http.Request) {
	if len(secretKey) == 0 {
		http.Error(w, "jwt secret is invalid", http.StatusInternalServerError)
		return
	}

	p := new(Payload)

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

	if err := json.Unmarshal(b, p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(p.TokenString) == 0 {
		http.Error(w, "token is not valid", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(p.TokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("%s:%v", "invalid token", err), http.StatusBadRequest)
		return
	}

	if token == nil {
		http.Error(w, fmt.Sprintf("%s:%v", "invalid token", "nil"), http.StatusBadRequest)
		return
	}

	if !token.Valid {
		http.Error(w, fmt.Sprintf("%s:%s", "invalid token", "invalid"), http.StatusBadRequest)
		return
	}

	if p.FuncData == nil {
		http.Error(w, "func data is nil", http.StatusBadRequest)
		return
	}

	if f, ok := registry[p.FuncData.Id]; ok {
		r.Body = ioutil.NopCloser(bytes.NewReader(p.FuncData.Data))
		f(w, r)
	}
}
