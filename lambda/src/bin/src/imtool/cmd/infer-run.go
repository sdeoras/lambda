package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sdeoras/api"
	"github.com/sdeoras/comp/cloud"
	"github.com/sdeoras/comp/image"
	"github.com/sdeoras/dispatcher"
	"github.com/sdeoras/lsdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type score struct {
	Index int
	value float32
}

type scores []score

func (s scores) Len() int {
	return len(s)
}

func (s scores) Less(i, j int) bool {
	return s[i].value > s[j].value
}

func (s scores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func inferRunE(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag("/concurrency", rootCmd.Flags().Lookup("concurrency"))
	_ = viper.BindPFlag("/timeout", rootCmd.Flags().Lookup("timeout"))
	_ = viper.BindPFlag("/infer/modelFile", cmd.Flags().Lookup("model"))
	_ = viper.BindPFlag("/infer/labelFile", cmd.Flags().Lookup("label"))
	_ = viper.BindPFlag("/infer/file", cmd.Flags().Lookup("file"))

	n := viper.GetInt("/concurrency")
	t := viper.GetInt("/timeout")
	modelFile := viper.GetString("/infer/modelFile")
	labelFile := viper.GetString("/infer/labelFile")
	inputFiles := viper.GetStringSlice("/infer/file")
	inputFiles = append(inputFiles, args...)

	var files []string
	filesMap := make(map[string]struct{})

	// if not disk file put in map, else put in files slice
	for _, file := range inputFiles {
		if file == "-" || strings.Contains(file, "gs://") {
			filesMap[file] = struct{}{}
		} else {
			files = append(files, file)
		}
	}

	// do not show usage on error
	cmd.SilenceUsage = true

	if n <= 0 {
		return fmt.Errorf("concurrency value needs to be positive")
	}

	lister := lsdir.NewLister(true, "*")
	var err error
	files, err = lister.List(files...)
	if err != nil {
		return fmt.Errorf("error listing files:%v", err)
	}

	for _, file := range files {
		filesMap[file] = struct{}{}
	}

	files = make([]string, 0, len(filesMap))
	for file := range filesMap {
		files = append(files, file)
	}

	if len(files) == 0 {
		return fmt.Errorf("please provide at least an image to work with")
	}

	logrus.Infof("found %d files", len(files))

	// create operator to read from cloud
	cloudOp, err := cloud.NewOperator(nil)
	if err != nil {
		return err
	}
	defer cloudOp.Close()

	// create operator to work with images
	imageOp, err := image.NewOperator(nil)
	if err != nil {
		return err
	}
	defer imageOp.Close()

	// read labelFile file
	labelData, err := cloudOp.Read(labelFile)
	if err != nil {
		return err
	}

	labels := make(map[int]string)
	scanner := bufio.NewScanner(bytes.NewReader(labelData))
	i := 0
	for scanner.Scan() {
		labels[i] = scanner.Text()
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// read frozen graph from the modelFile
	graphDef, err := cloudOp.Read(modelFile)
	if err != nil {
		return err
	}

	graph := tf.NewGraph()
	if err := graph.Import(graphDef, ""); err != nil {
		return err
	}

	// start new session for compute
	sess, err := tf.NewSession(graph, nil)
	if err != nil {
		return fmt.Errorf("error starting new session:%v", err)
	}
	defer sess.Close()

	d := dispatcher.New(int32(n))
	c := make(chan string)

	// immediately spawn a go-routine that keeps reading from the channel and printing on stdout
	go func() {
		for {
			fmt.Println(<-c)
		}
	}()

	response := new(api.InferImageResponse)
	response.Outputs = make([]*api.InferOutput, 0, 0)
	mu := new(sync.Mutex)

	for _, fileName := range files {
		// do this if the var is being accessed from within a goroutine
		fileName := fileName

		d.Do(func() {
			l := logrus.WithField("file", fileName)

			// read image data
			var imageData []byte
			var err error

			if strings.Contains(fileName, "gs://") {
				imageData, err = cloudOp.Read(fileName)
				if err != nil {
					l.Errorf("error reading file %s:%v", fileName, err)
					return
				}
			} else {
				if fileName == "-" {
					imageData, err = ioutil.ReadAll(os.Stdin)
					if err != nil {
						l.Errorf("error reading input on stdin:%v", err)
						return
					}
				} else {
					imageData, err = ioutil.ReadFile(fileName)
					if err != nil {
						l.Errorf("error reading file %s:%v", fileName, err)
						return
					}
				}
			}

			// decode image
			im, err := imageOp.Decode(bytes.NewReader(imageData))
			if err != nil {
				l.Errorf("error decoding image:%v", err)
				return
			}

			imageRaw, err := imageOp.ResizeNormalize(299, 299, 0, 255, im)
			if err != nil {
				l.Errorf("error resizing image:%v", err)
				return
			}

			imT, err := tf.NewTensor(imageRaw)
			if err != nil {
				l.Errorf("error making image tensor:%v", err)
				return
			}

			feeds := make(map[tf.Output]*tf.Tensor)
			feeds[graph.Operation("Placeholder").Output(0)] = imT

			out, err := sess.Run(
				feeds,
				[]tf.Output{
					graph.Operation("final_result").Output(0),
				},
				nil,
			)
			if err != nil {
				l.Errorf("error running session:%v", err)
				return
			}

			output, ok := out[0].Value().([][]float32)
			if !ok {
				l.Errorf("type inference error, expected [][]float32, got %T", out[0].Value())
				return
			}

			for i := range output {
				s := make([]score, len(output[i]))
				for j := range output[i] {
					s[j].Index = j
					s[j].value = output[i][j]
				}

				sort.Sort(scores(s))

				out := new(api.InferOutput)
				out.Label = labels[s[0].Index]
				out.Name = fileName
				out.Probability = int64(s[0].value * 100)

				mu.Lock()
				response.Outputs = append(response.Outputs, out)
				mu.Unlock()

				break
			}
		})
	}

	// create a timeout
	timeout := time.After(time.Duration(t) * time.Second)

Loop:
	for {
		select {
		case <-timeout:
			logrus.Infof("timeout occurred. set to %d, use -t to change", t)
			return nil
		default:
			if !d.IsRunning() {
				break Loop
			} else {
				time.Sleep(time.Millisecond * 20)
			}
		}
	}

	jb, err := json.Marshal(response)
	if err != nil {
		return err
	}

	fmt.Println(string(jb))

	return nil
}
