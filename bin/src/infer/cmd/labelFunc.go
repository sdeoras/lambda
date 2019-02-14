package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"sort"

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

	modelFile := viper.GetString("/label/modelFile")
	labelFile := viper.GetString("/label/labelFile")

	if len(args) == 0 {
		return fmt.Errorf("please provide an image to work with as argument")
	}

	fileName := args[0]

	// create operators to read from cloud and for working with images
	cloudOp, err := cloud.NewOperator(nil)
	if err != nil {
		return err
	}
	defer cloudOp.Close()

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

	// read image data
	imageData, err := cloudOp.Read(fileName)
	if err != nil {
		return err
	}
	// decode image
	im, err := imageOp.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Fatal(err)
	}

	imageRaw, err := imageOp.ResizeNormalize(299, 299, 0, 255, im)
	if err != nil {
		log.Fatal(err)
	}

	imT, err := tf.NewTensor(imageRaw)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal("session run:", err)
	}

	output, ok := out[0].Value().([][]float32)
	if !ok {
		log.Fatal("type inference error, expected [][]float32, got %T", out[0].Value())
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

		fmt.Println(fileName, labels[sOut[0]])
	}
	return nil
}
