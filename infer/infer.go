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

	// paths to run imtool binary
	imtoolPath = "/srv/files/bin/src/imtool"
	imtoolExec = imtoolPath + "/a.out"
	imtoolLib  = imtoolPath + "/lib"

	// model location and convention
	modelDir   = "models"
	graphFile  = "output_graph.pb"
	labelsFile = "output_labels.txt"

	// image dir
	imageDir = "images"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		_ = os.Setenv("LD_LIBRARY_PATH",
			imtoolLib+":"+os.Getenv("LD_LIBRARY_PATH"))
	})
}

func writeToGS(ctx context.Context, bucketName, fileName string, buffer []byte) (int, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return 0, err
	}

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(fileName)
	w := obj.NewWriter(ctx)
	defer w.Close()
	return w.Write(buffer)
}

// InferImage provides image inferencing using imtool exec. It depends on a model and a
// labels file. Location of these files is determined using bkt name (env var) and model name
// in the http request buffer. The request buffer is a proto buffer based on an api defined
// in a proto file.
func InferImage(w http.ResponseWriter, r *http.Request) {
	// validate input request
	err := jwt.Validator.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check method
	if r.Method != http.MethodPost {
		http.Error(w,
			"method not set to POST",
			http.StatusBadRequest)
		return
	}

	// check env var.
	if len(os.Getenv("LAMBDA_BUCKET")) <= 0 {
		http.Error(w,
			"env var LAMBDA_BUCKET not set",
			http.StatusInternalServerError)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("error reading http request body:%v", err),
			http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	inferRequest := new(api.InferImageRequest)
	if err := proto.Unmarshal(b, inferRequest); err != nil {
		http.Error(w,
			fmt.Sprintf("could not unmarshal image infer request:%v", err),
			http.StatusBadRequest)
		return
	}

	// validate request
	if len(inferRequest.ModelName) == 0 ||
		len(inferRequest.ModelVersion) == 0 ||
		len(inferRequest.Images) == 0 {
		http.Error(w,
			fmt.Sprintf("invalid infer request, some fields are empty:%v", err),
			http.StatusBadRequest)
		return
	}

	modelFolder := fmt.Sprintf("gs://%s/%s/%s/%s",
		os.Getenv("LAMBDA_BUCKET"),
		modelDir,
		inferRequest.ModelName,
		inferRequest.ModelVersion)

	modelPath := fmt.Sprintf("%s/%s",
		modelFolder,
		graphFile)
	labelPath := fmt.Sprintf("%s/%s",
		modelFolder,
		labelsFile)

	id := uuid.New().String()
	files := make([]string, 0, 0)
	for _, image := range inferRequest.Images {
		fileName := image.Name + "_" + id + ".jpg"
		fileName = fmt.Sprintf("gs://%s/%s/%s/%s",
			os.Getenv("LAMBDA_BUCKET"),
			imageDir,
			inferRequest.ModelName,
			fileName)
		files = append(files, fileName)

		objName := fmt.Sprintf("%s/%s/%s",
			imageDir,
			inferRequest.ModelName,
			fileName)

		// write file to gcs
		if n, err := writeToGS(context.Background(),
			os.Getenv("LAMBDA_BUCKET"),
			objName, image.Data); err != nil {
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
	args := []string{"infer",
		"--model", modelPath,
		"--label", labelPath}
	b, err = exec.Command(imtoolExec,
		append(args, files...)...).Output()
	if err != nil {
		http.Error(w,
			fmt.Sprintf("could not successfully run infer:%v:imtool %v %v", err, args, files),
			http.StatusInternalServerError)
		return
	}

	response := new(api.InferImageResponse)
	if err := json.Unmarshal(b, response); err != nil {
		http.Error(w,
			fmt.Sprintf("could not successfull unmarshal into response:%v", err),
			http.StatusInternalServerError)
		return
	}

	b, err = proto.Marshal(response)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("could not successfull marshal response into proto:%v", err),
			http.StatusInternalServerError)
		return
	}

	n, err := w.Write(b)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("could not successfull write to response writer:%v", err),
			http.StatusInternalServerError)
		return
	}

	if n != len(b) {
		http.Error(w,
			fmt.Sprintf("could not successfull write all data to response writer:%v of %v bytes",
				n, len(b)),
			http.StatusInternalServerError)
		return
	}
}
