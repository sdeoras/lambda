package gen

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gan/src/jwt"
	"gan/src/log"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/sdeoras/api"
)

const (
	ProjectName = "gan"
	Name        = "gen"

	// paths to run imtool binary
	toolPath = "/srv/files/src/bin/src/gangen"
	toolBin  = toolPath + "/a.out"
	toolLib  = toolPath + "/lib"

	// why in /tmp?
	// pl. read: https://stackoverflow.com/questions/42719793/write-temporary-files-from-google-cloud-function
	toolModels = "/tmp" + "/" + modelDir

	// model location and convention
	modelDir       = "models"
	checkpointFile = "cp.pb"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		_ = os.Setenv("LD_LIBRARY_PATH",
			toolLib+":"+os.Getenv("LD_LIBRARY_PATH"))
	})
}

func copyModelIfNotExists(ctx context.Context, modelName, version string) error {
	localFolder := filepath.Join(toolModels, modelName, version)
	if _, err := os.Stat(localFolder); err == nil {
		log.Out.Println("models found in tmp")
		return nil
	} else if os.IsNotExist(err) {
		log.Out.Println("models not found in tmp")

		if err := os.MkdirAll(localFolder, 0755); err != nil {
			return err
		}

		client, err := storage.NewClient(ctx)
		if err != nil {
			return err
		}

		bucket := client.Bucket(os.Getenv("LAMBDA_BUCKET"))

		obj := bucket.Object(filepath.Join(modelDir, modelName, version, checkpointFile))
		r, err := obj.NewReader(ctx)
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath.Join(localFolder, checkpointFile), b, 0644); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("unknown file existance status:%v", err)
	}

	return nil
}

// GenerateImages is a GAN based image generator
func GenerateImages(w http.ResponseWriter, r *http.Request) {
	// validate input request
	err := jwt.Validator.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	request := new(api.GanRequest)
	if err := proto.Unmarshal(b, request); err != nil {
		http.Error(w,
			fmt.Sprintf("could not unmarshal image infer request:%v", err),
			http.StatusBadRequest)
		return
	}

	// validate request
	if len(request.ModelName) == 0 ||
		len(request.ModelVersion) == 0 ||
		request.Count <= 0 {
		http.Error(w,
			fmt.Sprintf("invalid infer request, some fields are empty"),
			http.StatusBadRequest)
		return
	}

	if err := copyModelIfNotExists(context.Background(), request.ModelName,
		request.ModelVersion); err != nil {
		http.Error(w,
			fmt.Sprintf("could not copy model or check existence:%v", err),
			http.StatusInternalServerError)
		return
	}

	localFolder := filepath.Join(toolModels, request.ModelName, request.ModelVersion)
	modelPath := filepath.Join(localFolder, checkpointFile)

	response := new(api.GanResponse)

	// buffer for writing data
	bb := new(bytes.Buffer)
	bw := bufio.NewWriter(bb)

	// define command and connect STDIN and STDOUT accordingly
	cmd := exec.Command(toolBin,
		[]string{
			"mnist",
			"--model", modelPath,
			"--count", fmt.Sprintf("%d", request.Count),
			"--out", "-", // send data to STDOUT
		}...)

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
	if err := json.Unmarshal(bb.Bytes(), response); err != nil {
		http.Error(w,
			fmt.Sprintf("could not successfull unmarshal into response:%v", err),
			http.StatusInternalServerError)
		return
	}

	// serialize response as a protobuf
	b, err = json.MarshalIndent(response, "", "  ")
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
