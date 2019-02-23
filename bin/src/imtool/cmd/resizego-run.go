package cmd

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"
	"github.com/sdeoras/dispatcher"
	"github.com/sdeoras/lsdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func resizegoRunE(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag("/concurrency", rootCmd.Flags().Lookup("concurrency"))
	_ = viper.BindPFlag("/timeout", rootCmd.Flags().Lookup("timeout"))
	_ = viper.BindPFlag("/resizego/file", cmd.Flags().Lookup("file"))
	_ = viper.BindPFlag("/resizego/out", cmd.Flags().Lookup("out"))
	_ = viper.BindPFlag("/resizego/height", cmd.Flags().Lookup("height"))
	_ = viper.BindPFlag("/resizego/width", cmd.Flags().Lookup("width"))

	n := viper.GetInt("/concurrency")
	t := viper.GetInt("/timeout")
	diskFiles := viper.GetStringSlice("/resizego/file")
	outPath := viper.GetString("/resizego/out")
	height := viper.GetInt("/resizego/height")
	width := viper.GetInt("/resizego/width")

	// do not show usage on error
	cmd.SilenceUsage = true

	if n <= 0 {
		return fmt.Errorf("concurrency value needs to be positive")
	}

	lister := lsdir.NewLister(true, "*")
	files, err := lister.List(diskFiles...)
	if err != nil {
		return fmt.Errorf("error listing files:%v", err)
	}

	files = append(files, args...)

	if len(files) == 0 {
		return fmt.Errorf("please provide at least an image to work with")
	}

	logrus.Infof("found %d files", len(files))

	d := dispatcher.New(int32(n))
	c := make(chan string)

	// immediately spawn a go-routine that keeps reading from the channel and printing on stdout
	go func() {
		for {
			fmt.Println(<-c)
		}
	}()

	for _, fileName := range files {
		// do this if the var is being accessed from within a goroutine
		fileName := fileName

		d.Do(func() {
			now := time.Now()
			l := logrus.WithField("file", fileName)

			// open "test.jpg"
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}

			// decode jpeg into image.Image
			img, err := jpeg.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()

			// resize to width 1000 using Lanczos resampling
			// and preserve aspect ratio
			m := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

			out, err := os.Create(filepath.Join(outPath, filepath.Base(fileName)))
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			// write new image to file
			jpeg.Encode(out, m, nil)

			l.Infof("resized to %d x %d in %v", height, width, time.Since(now))
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

	return nil
}
