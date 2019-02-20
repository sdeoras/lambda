package infer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lambda/jwt"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/sdeoras/lambda/api"
)

const (
	ProjectName = "lambda"
	Name        = "infer"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		_ = os.Setenv("LD_LIBRARY_PATH", "/srv/files/bin/src/infer/lib:"+
			os.Getenv("LD_LIBRARY_PATH"))
	})
}

func writeToGS(ctx context.Context, bucketName, fileName string, buffer []byte) (int, error) {
	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return 0, err
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	obj := bucket.Object(fileName)
	w := obj.NewWriter(ctx)
	defer w.Close()
	return w.Write(buffer)
}

func InferImage(w http.ResponseWriter, r *http.Request) {
	// validate input request
	err := jwt.Validator.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error readhing http request body in SendEMail:%v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	inferRequest := new(api.InferImageRequest)
	if err := proto.Unmarshal(b, inferRequest); err != nil {
		http.Error(w, fmt.Sprintf("could not unmarshal image infer request:%v", err), http.StatusBadRequest)
		return
	}

	inferRequest.ModelPath = "gs://" + os.Getenv("LAMBDA_BUCKET") + "/" + inferRequest.ModelPath
	inferRequest.LabelPath = "gs://" + os.Getenv("LAMBDA_BUCKET") + "/" + inferRequest.LabelPath

	id := uuid.New().String()
	files := make([]string, 0, 0)
	for _, image := range inferRequest.Images {
		fileName := image.Name + "_" + id + ".jpg"
		files = append(files, "gs://"+os.Getenv("LAMBDA_BUCKET")+"/"+fileName)

		// write file to gcs
		if n, err := writeToGS(context.Background(),
			os.Getenv("LAMBDA_BUCKET"),
			fileName, image.Data); err != nil {
			http.Error(w, fmt.Sprintf("could not successfull write to gcs bucket:%v", err), http.StatusInternalServerError)
			return
		} else {
			if n != len(image.Data) {
				http.Error(w, fmt.Sprintf("could not successfull write all data to gcs bucket:%v of %v", n, len(image.Data)), http.StatusInternalServerError)
				return
			}
		}
	}

	// executing the shell binary a.out produces json output that can be unmarshal'ed into
	// infer response object
	args := []string{"label",
		"--model", inferRequest.ModelPath,
		"--label", inferRequest.LabelPath}
	b, err = exec.Command("/srv/files/bin/src/infer/a.out",
		append(args, files...)...).Output()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not successfully run infer:%v", err), http.StatusInternalServerError)
		return
	}

	response := new(api.InferImageResponse)
	if err := json.Unmarshal(b, response); err != nil {
		http.Error(w, fmt.Sprintf("could not successfull unmarshal into response:%v", err), http.StatusInternalServerError)
		return
	}

	b, err = proto.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not successfull marshal response into proto:%v", err), http.StatusInternalServerError)
		return
	}

	n, err := w.Write(b)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not successfull write to response writer:%v", err), http.StatusInternalServerError)
		return
	}

	if n != len(b) {
		http.Error(w, fmt.Sprintf("could not successfull write all data to response writer:%v of %v bytes", n, len(b)), http.StatusInternalServerError)
		return
	}
}
