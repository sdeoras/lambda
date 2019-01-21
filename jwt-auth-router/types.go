package payload

import (
	"net/http"
	"sync"
)

var once sync.Once
var secretKey []byte
var registry map[int]func(w http.ResponseWriter, r *http.Request)

const (
	HandlerAuthOnly = iota // do not register any function against this key
	HandlerHelloWorld
)

// FuncData is a payload for a particular function.
// Id is registered function id and data is populated in the http request body
// prior to calling that particular function handler.
type FuncData struct {
	Id   int
	Data []byte
}

// Payload contains payload for a function and a jwt token for authenticating.
type Payload struct {
	FuncData    *FuncData
	TokenString string
}
