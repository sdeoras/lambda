package health

import (
	"fmt"
	"gan/src/jwt"
	"io/ioutil"
	"net/http"

	"github.com/gogo/protobuf/proto"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Check responds with health check status
func Check(w http.ResponseWriter, r *http.Request) {
	// validate input request
	err := jwt.Validator.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%s:%v", "could not read from http request", err),
			http.StatusBadRequest)
		return
	}

	request := new(healthpb.HealthCheckRequest)
	if err := proto.Unmarshal(b, request); err != nil {
		http.Error(w,
			fmt.Sprintf("%s:%v", "could not unmarshal request", err),
			http.StatusBadRequest)
		return
	}

	response := new(healthpb.HealthCheckResponse)
	switch request.Service {
	case "infer":
		response.Status = healthpb.HealthCheckResponse_SERVING
	case "email":
		response.Status = healthpb.HealthCheckResponse_SERVING
	default:
		response.Status = healthpb.HealthCheckResponse_SERVICE_UNKNOWN
	}

	b, err = proto.Marshal(response)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%s:%v", "could not marshal response", err),
			http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(b)
}
