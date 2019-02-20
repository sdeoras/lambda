package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sdeoras/dispatcher"
	"github.com/sirupsen/logrus"
	"log"
	"sort"
	"time"

	"github.com/sdeoras/lambda/api"

	"github.com/sdeoras/comp/cloud"
	"github.com/sdeoras/comp/image"
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

func label(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag("/label/modelFile", cmd.Flags().Lookup("model"))
	_ = viper.BindPFlag("/label/labelFile", cmd.Flags().Lookup("label"))
	_ = viper.BindPFlag("/concurrency", rootCmd.Flags().Lookup("concurrency"))
	_ = viper.BindPFlag("/timeout", rootCmd.Flags().Lookup("timeout"))

	modelFile := viper.GetString("/label/modelFile")
	labelFile := viper.GetString("/label/labelFile")
	n := viper.GetInt("/concurrency")
	t := viper.GetInt("/timeout")

	if n <= 0 {
		return fmt.Errorf("concurrency value needs to be positive")
	}

	if len(args) == 0 {
		return fmt.Errorf("please provide an image to work with as argument")
	}

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
		log.Fatal(err)
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

	for _, fileName := range args {
		// do this if the var is being accessed from within a goroutine
		fileName := fileName

		d.Do(func() {
			l := logrus.WithField("file", fileName)

			// read image data
			imageData, err := cloudOp.Read(fileName)
			if err != nil {
				l.Errorf("error reading file:%v", err)
				return
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

				sOut := make([]int, len(s))
				for j := range s {
					sOut[j] = s[j].Index
				}

				response := new(api.InferImageResponse)
				response.Label = labels[sOut[0]]
				jb, err := json.Marshal(response)
				if err != nil {
					l.Errorf("error marshaling json:%v", err)
					return
				}

				c <- fileName + ":" + string(jb)
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
			return fmt.Errorf("timeout occurred")
		default:
			if !d.IsRunning() {
				break Loop
			} else {
				time.Sleep(time.Millisecond*20)
			}
		}
	}

	return nil
}
