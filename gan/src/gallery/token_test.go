package gallery

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gogo/protobuf/proto"
	"github.com/sdeoras/api"
)

func TestToken(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(tokenString)

	request := new(api.GanRequest)

	request.Count = 10
	request.ModelName = "gan-mnist-generator"
	request.ModelVersion = "v1"

	b, err := proto.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("/tmp/gangen.pb", b, 0644); err != nil {
		t.Fatal(err)
	}
}
