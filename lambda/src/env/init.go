package env

import (
	"os"
	"sync"
)

var (
	once         sync.Once
	FuncName     string
	Domain       string
	ClientId     string
	ClientSecret string
	JwtSecret    string
	CodeLocation string
	Bucket       string
)

func init() {
	once.Do(func() {
		FuncName = os.Getenv("FUNCTION_NAME")
		Domain = os.Getenv("GOOGLE_GCF_DOMAIN")
		ClientId = os.Getenv("GOOGLE_CLIENT_ID")
		ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
		JwtSecret = os.Getenv("JWT_SECRET_KEY")
		CodeLocation = os.Getenv("CODE_LOCATION")
		Bucket = os.Getenv("CLOUD_FUNCTIONS_BUCKET")
	})
}
