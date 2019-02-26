package infer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lambda/jwt"
	"lambda/log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"cloud.google.com/go/logging"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/sdeoras/lambda/api"
)

const (
	ProjectName = "lambda"
	Name        = "infer"

	// paths to run imtool binary
	imtoolPath = "/srv/files/bin/src/imtool"
	imtoolExec = imtoolPath + "/a.out"
	imtoolLib  = imtoolPath + "/lib"

	// why in /tmp?
	// pl. read: https://stackoverflow.com/questions/42719793/write-temporary-files-from-google-cloud-function
	imtoolModels = "/tmp" + "/" + modelDir

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

func copyModelIfNotExists(ctx context.Context, modelName, version string) error {
	localFolder := filepath.Join(imtoolModels, modelName, version)
	if _, err := os.Stat(localFolder); err == nil {
		if log.Logger != nil {
			log.Logger.Log(logging.Entry{Payload: "models found in tmp"})
		}
		return nil
	} else if os.IsNotExist(err) {
		if log.Logger != nil {
			log.Logger.Log(logging.Entry{Payload: "models not found in tmp"})
		}

		if err := os.MkdirAll(localFolder, 0755); err != nil {
			return err
		}

		client, err := storage.NewClient(ctx)
		if err != nil {
			return err
		}

		bucket := client.Bucket(os.Getenv("LAMBDA_BUCKET"))

		obj := bucket.Object(filepath.Join(modelDir, modelName, version, graphFile))
		r, err := obj.NewReader(ctx)
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath.Join(localFolder, graphFile), b, 0644); err != nil {
			return err
		}

		obj = bucket.Object(filepath.Join(modelDir, modelName, version, labelsFile))
		r, err = obj.NewReader(ctx)
		if err != nil {
			return err
		}

		b, err = ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath.Join(localFolder, labelsFile), b, 0644); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("unknown file existance status:%v", err)
	}

	return nil
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

	if log.Logger != nil {
		defer log.Logger.Flush()
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

	if err := copyModelIfNotExists(context.Background(), inferRequest.ModelName,
		inferRequest.ModelVersion); err != nil {
		http.Error(w,
			fmt.Sprintf("could not copy model or check existence:%v", err),
			http.StatusInternalServerError)
		return
	}

	localFolder := filepath.Join(imtoolModels, inferRequest.ModelName, inferRequest.ModelVersion)
	modelPath := filepath.Join(localFolder, graphFile)
	labelPath := filepath.Join(localFolder, labelsFile)

	response := new(api.InferImageResponse)

	for _, image := range inferRequest.Images {
		// buffer for writing data
		bb := new(bytes.Buffer)
		bw := bufio.NewWriter(bb)

		out := new(api.InferImageResponse)

		// define command and connect STDIN and STDOUT accordingly
		cmd := exec.Command(imtoolExec,
			[]string{
				"infer",
				"--model", modelPath,
				"--label", labelPath,
				"-f", "-", // receive data from STDIN
			}...)
		cmd.Stdin = bytes.NewReader(image.Data)
		cmd.Stdout = bw
		cmd.Stderr = ioutil.Discard

		// run command
		if err := cmd.Run(); err != nil {
			http.Error(w,
				fmt.Sprintf("could not successfully run imtool:%v", err),
				http.StatusInternalServerError)
			return
		}

		// flush writer
		if err := bw.Flush(); err != nil {
			http.Error(w,
				fmt.Sprintf("error flushing bufio writer:%v", err),
				http.StatusInternalServerError)
			return
		}

		// unmarshal output
		if err := json.Unmarshal(bb.Bytes(), out); err != nil {
			http.Error(w,
				fmt.Sprintf("could not successfull unmarshal into response:%v", err),
				http.StatusInternalServerError)
			return
		}

		// collect output
		response.Outputs = append(response.Outputs, out.Outputs...)
	}

	// serialize response as a protobuf
	b, err = proto.Marshal(response)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("could not successfull marshal response into proto:%v", err),
			http.StatusInternalServerError)
		return
	}

	// write that to http response writer
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
