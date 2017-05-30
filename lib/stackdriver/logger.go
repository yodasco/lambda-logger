package stackdriver

import (
	"log"

	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
)

func Log() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "yodas-test"

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name of the log to write to.
	logName := "my-log"

	// Selects the log to write to.
	logger := client.Logger(logName)

	// Sets the data to log.
	text := "Hello, world!"

	// Adds an entry to the log buffer.
	logger.Log(logging.Entry{Payload: text})

	// Closes the client and flushes the buffer to the Stackdriver Logging
	// service.
	if err := client.Close(); err != nil {
		log.Fatalf("Failed to close client: %v", err)
	}

	log.Printf("Logged: %v\n", text)
}
