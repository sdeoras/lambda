package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sdeoras/comp/cloud"
	"github.com/sdeoras/comp/image"
	"github.com/sdeoras/dispatcher"
	"github.com/sdeoras/lsdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func resizeRunE(cmd *cobra.Command, args []string) error {
	_ = viper.BindPFlag("/concurrency", rootCmd.Flags().Lookup("concurrency"))
	_ = viper.BindPFlag("/timeout", rootCmd.Flags().Lookup("timeout"))
	_ = viper.BindPFlag("/resize/file", cmd.Flags().Lookup("file"))
	_ = viper.BindPFlag("/resize/out", cmd.Flags().Lookup("out"))
	_ = viper.BindPFlag("/resize/height", cmd.Flags().Lookup("height"))
	_ = viper.BindPFlag("/resize/width", cmd.Flags().Lookup("width"))

	n := viper.GetInt("/concurrency")
	t := viper.GetInt("/timeout")
	diskFiles := viper.GetStringSlice("/resize/file")
	outPath := viper.GetString("/resize/out")
	height := viper.GetInt("/resize/height")
	width := viper.GetInt("/resize/width")

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

			h, w, _ := im.Size()

			if height <= 0 && width <= 0 {

			} else if height > 0 && width <= 0 {
				w = w / h * height
				h = height
			} else if height <= 0 && width > 0 {
				h = h / w * width
				w = width
			} else {
				h = height
				w = width
			}

			imagesRaw, err := imageOp.Resize(h, w, im)
			if err != nil {
				l.Errorf("error resizing image:%v", err)
				return
			}

			imageRaw := imagesRaw[0]
			x := make([][][]uint8, len(imageRaw))
			for i := range imageRaw {
				x[i] = make([][]uint8, len(imageRaw[i]))
				for j := range imageRaw[i] {
					x[i][j] = make([]uint8, len(imageRaw[i][j]))
					for k := range imageRaw[i][j] {
						x[i][j][k] = uint8(imageRaw[i][j][k])
					}
				}
			}

			b, err := imageOp.Encode(image.Image(x))
			if err != nil {
				l.Errorf("error encoding image:%v", err)
				return
			}

			var outFile string
			if strings.Contains(outPath, "://") {
				outFile = strings.TrimRight(outPath, "/") + "/" + filepath.Base(fileName)
			} else {
				outFile = filepath.Join(outPath, filepath.Base(fileName))
			}

			if err := cloudOp.Write(outFile, b); err != nil {
				l.Errorf("error encoding image:%v", err)
				return
			}

			l.Infof("resized to %d x %d in %v", h, w, time.Since(now))
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
