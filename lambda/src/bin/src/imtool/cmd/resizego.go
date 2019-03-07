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
	"github.com/spf13/cobra"
)

// resizegoCmd represents the resizego command
var resizegoCmd = &cobra.Command{
	Use:   "resizego",
	Short: "resize using go lib",
	Long:  ``,
	RunE:  resizegoRunE,
}

func init() {
	rootCmd.AddCommand(resizegoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resizegoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resizegoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	resizegoCmd.Flags().Int("height", -1, "height of image")
	resizegoCmd.Flags().Int("width", -1, "width of image")
	resizegoCmd.Flags().StringSliceP("file", "f", []string{}, "file or folder on disk")
	resizegoCmd.Flags().StringP("out", "o", "", "output folder or bucket path")
}
