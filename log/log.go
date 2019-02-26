package log

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"cloud.google.com/go/logging"
)

const (
	logID = "gcf-lambda-logs"
)

var (
	Logger *logging.Logger
	once   sync.Once
)

func init() {
	once.Do(func() {
		// Create a Client
		ctx := context.Background()
		client, err := logging.NewClient(ctx,
			filepath.Join(
				"projects",
				os.Getenv("GCLOUD_PROJECT_NAME"),
			),
		)

		if err != nil {
			return
		}

		// Initialize a logger
		Logger = client.Logger(logID)
	})
}
