package gallery

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gan/src/config"
	"gan/src/jwt"
	"gan/src/log"
	"gan/src/login"
	"gan/src/route"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/sdeoras/api/pb"
	"github.com/sdeoras/oauth"
)

const (
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

func writeToGS(ctx context.Context, bucketName, fileName string, buffer []byte, public bool) (int, error) {
	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return 0, err
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	obj := bucket.Object(fileName)
	w := obj.NewWriter(ctx)
	if public {
		w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	}
	defer w.Close()
	return w.Write(buffer)
}

func copyModelIfNotExists(ctx context.Context, modelName, version string) error {
	localFolder := filepath.Join(toolModels, modelName, version)
	if _, err := os.Stat(localFolder); err == nil {
		log.Stdout().Println("models found in tmp")
		return nil
	} else if os.IsNotExist(err) {
		log.Stdout().Println("models not found in tmp")

		if err := os.MkdirAll(localFolder, 0755); err != nil {
			return err
		}

		client, err := storage.NewClient(ctx)
		if err != nil {
			return err
		}

		bucket := client.Bucket(config.Config().BucketName)

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

func GenerateDriver(w http.ResponseWriter, r *http.Request) {
	// check if this is a callback from auth provider, else redirect to login page
	content, err := login.Provider[login.GoogleProvider].GetUserInfo(r)
	if err != nil {
		mesg := err.Error()
		url := "https://" + filepath.Join(
			config.Config().Domain,
			config.Config().FuncName,
			route.Root,
		)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			http.Error(w,
				fmt.Sprintf("%v:error in gan.GenerateDriver:%s:%s",
					http.StatusInternalServerError, mesg, err.Error()),
				http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, url, http.StatusPermanentRedirect)
		return
	}

	// unmarshal contents into a struct
	ac := new(oauth.GoogleAuthContent)
	if err := ac.Unmarshal(content); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:%v",
				http.StatusInternalServerError, err.Error()),
			http.StatusInternalServerError)
		return
	}

	request := new(pb.GanRequest)
	request.Count = 10
	request.ModelName = "gan-mnist-generator"
	request.ModelVersion = "v1"

	b, err := proto.Marshal(request)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not marshal gan request",
				http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	// Pl. see the link below to understand why jwt passed in header
	// is not preserved during http redirect call
	// https://stackoverflow.com/questions/36345696/golang-http-redirect-with-headers
	// Hence, we pass it in URL
	url := "https://" + filepath.Join(
		config.Config().Domain,
		config.Config().FuncName,
		route.Gallery,
	)

	req, err := jwt.Manager().NewHTTPRequest(http.MethodPost, url, nil, b)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not successfull create http request:%v",
				http.StatusInternalServerError, err),
			http.StatusInternalServerError)
		return
	}

	GenerateImages(w, req)
}

// GenerateImages is a GAN based image generator
func GenerateImages(w http.ResponseWriter, r *http.Request) {
	// validate input request
	if err := jwt.Manager().Validate(r); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:%s", http.StatusBadRequest, err.Error()),
			http.StatusBadRequest)
		return
	}

	// check method
	if r.Method != http.MethodPost {
		http.Error(w,
			fmt.Sprintf("%v:%s", http.StatusBadRequest,
				"error in gen.GenerateImages: method not set to POST"),
			http.StatusBadRequest)
		return
	}

	// check env var.
	if len(config.Config().BucketName) <= 0 {
		http.Error(w,
			fmt.Sprintf("%v:%s", http.StatusInternalServerError,
				"env var for GCS bucket name is not set"),
			http.StatusInternalServerError)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v:error reading http request body:%v",
				http.StatusBadRequest, err),
			http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	request := new(pb.GanRequest)
	if err := proto.Unmarshal(b, request); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not unmarshal image infer request:%v",
				http.StatusBadRequest, err),
			http.StatusBadRequest)
		return
	}

	// validate request
	if len(request.ModelName) == 0 ||
		len(request.ModelVersion) == 0 ||
		request.Count <= 0 {
		http.Error(w,
			fmt.Sprintf("%v:invalid infer request, some fields are empty",
				http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	if err := copyModelIfNotExists(context.Background(), request.ModelName,
		request.ModelVersion); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not copy model or check existence:%v",
				http.StatusInternalServerError, err),
			http.StatusInternalServerError)
		return
	}

	localFolder := filepath.Join(toolModels, request.ModelName, request.ModelVersion)
	modelPath := filepath.Join(localFolder, checkpointFile)

	response := new(pb.GanResponse)

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
			fmt.Sprintf("%v:could not successfully run imtool:%v",
				http.StatusInternalServerError, err),
			http.StatusInternalServerError)
		return
	}

	// flush writer
	if err := bw.Flush(); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:error flushing bufio writer:%v",
				http.StatusInternalServerError, err),
			http.StatusInternalServerError)
		return
	}

	// unmarshal output
	if err := json.Unmarshal(bb.Bytes(), response); err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not successfull unmarshal into response:%v",
				http.StatusInternalServerError, err),
			http.StatusInternalServerError)
		return
	}

	id := filepath.Join("images", config.Config().FuncName, uuid.New().String())

	for i := range response.Images {
		if _, err := writeToGS(
			context.Background(),
			config.Config().BucketName,
			filepath.Join(id, fmt.Sprintf("image-%d.jpg", i)),
			response.Images[i].Data,
			true); err != nil {
			http.Error(w,
				fmt.Sprintf("%v:could not successfull write to gcs:%v",
					http.StatusInternalServerError, err),
				http.StatusInternalServerError)
			return
		}
	}

	galleryRequest := new(pb.GalleryRequest)
	galleryRequest.GalleryItems = make([]*pb.GalleryItem, len(response.Images))
	for i := range response.Images {
		galleryRequest.GalleryItems[i] = &pb.GalleryItem{
			Id:         int64(i),
			FileName:   filepath.Join(id, fmt.Sprintf("image-%d.jpg", i)),
			Title:      "MNIST GAN image",
			Caption:    "a randomly generated MNIST image using a neural net based generative adversarial network (GAN)",
			BucketName: config.Config().BucketName,
		}
	}

	// serialize response as a protobuf
	b, err = proto.Marshal(galleryRequest)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not successfull marshal response into proto:%v",
				http.StatusInternalServerError, err),
			http.StatusInternalServerError)
		return
	}

	// Pl. see the link below to understand why jwt passed in header
	// is not preserved during http redirect call
	// https://stackoverflow.com/questions/36345696/golang-http-redirect-with-headers
	// Hence, we pass it in URL
	url := "https://" +
		filepath.Join(config.Config().Domain,
			config.Config().FuncName,
			route.Gallery,
		)

	req, err := jwt.Manager().NewHTTPRequest(http.MethodPost, url, nil, b)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("%v:could not successfull create http request:%v",
				http.StatusInternalServerError, err),
			http.StatusInternalServerError)
		return
	}

	Show(w, req)
}
