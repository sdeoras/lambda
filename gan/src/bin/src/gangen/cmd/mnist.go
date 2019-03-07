// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"gan/src/gan/mnist"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mnistCmd represents the mnist command
var mnistCmd = &cobra.Command{
	Use:   "mnist",
	Short: "Generate MNIST images",
	Long:  `Generate MNIST images using pre-trained GAN model`,
	RunE:  RunMnist,
}

const (
	mnistModel = "model"
	mnistCount = "count"
	mnistOut   = "out"
)

func RunMnist(cmd *cobra.Command, args []string) error {
	// define a func literal to derive viper keys
	f := func(input string) string {
		return filepath.Join("/mnist", input)
	}

	// bind cmd flags to viper
	_ = viper.BindPFlag(f(mnistModel), cmd.Flags().Lookup(mnistModel))
	_ = viper.BindPFlag(f(mnistCount), cmd.Flags().Lookup(mnistCount))
	_ = viper.BindPFlag(f(mnistOut), cmd.Flags().Lookup(mnistOut))

	// get flag values
	model := viper.GetString(f(mnistModel))
	count := viper.GetInt(f(mnistCount))
	out := viper.GetString(f(mnistOut))

	g, err := mnist.NewGenerator(model)
	if err != nil {
		return err
	}

	output, err := g.Generate(count)
	if err != nil {
		return err
	}

	switch out {
	case "-":
		if jb, err := json.Marshal(output); err != nil {
			return err
		} else {
			fmt.Println(string(jb))
		}
	default:
	}

	return nil
}

func init() {
	rootCmd.AddCommand(mnistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mnistCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mnistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	mnistCmd.Flags().StringP(mnistModel, "m", "", "path to model pb file")
	mnistCmd.Flags().IntP(mnistCount, "n", 1, "number of images to generate")
	mnistCmd.Flags().StringP(mnistOut, "o", "-", "output path or - for STDOUT")
}
