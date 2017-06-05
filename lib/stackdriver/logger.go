package stackdriver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/yodasco/lambda-logger/lib/types"

	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
)

var (
	projectID string
)

func init() {
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsFile == "" {
		log.Fatal("Missing required ENV variable GOOGLE_APPLICATION_CREDENTIALS")
	}
	projectID = readProjectID(credentialsFile)
}

// LogEvents logs the events to stackdriver
func LogEvents(events types.CloudwatchLogEvents) error {
	client := mustCreateClient()
	defer closeClient(client)
	logger := client.Logger(events.LogGroup)
	for _, event := range events.LogEvents {
		labels := appendMap(event.Labels,
			"logStream", events.LogStream,
			"logGroup", events.LogGroup)
		logger.Log(logging.Entry{
			Payload:   event.Message,
			Severity:  severityFromLogLevel(event.Level),
			Timestamp: event.Timestamp,
			InsertID:  fmt.Sprintf("cloudwatch-%s", event.ID),
			Labels:    labels,
		})
	}
	log.Printf("Logged %d lines from %s\n", len(events.LogEvents), events.LogGroup)
	return nil
}

// appends to the map the list of key-values.
// Key-values must be an even numner of arguments representing string
// key-values
func appendMap(m map[string]string, keyValues ...string) map[string]string {
	if m == nil {
		m = make(map[string]string)
	}
	if len(keyValues)%2 != 0 {
		log.Fatal("There should be an even number of keyValues to append")
	}
	for i := 0; i < len(keyValues)/2; i++ {
		k := keyValues[i*2]
		v := keyValues[i*2+1]
		m[k] = v
	}
	return m
}

func severityFromLogLevel(level string) logging.Severity {
	switch level {
	case "debug":
		return logging.Debug
	case "info":
		return logging.Info
	case "warn", "warning":
		return logging.Warning
	case "error":
		return logging.Error
	case "fatal":
		return logging.Critical
	default:
		return logging.Default
	}
}
func mustCreateClient() *logging.Client {
	var err error
	ctx := context.Background()

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Log errors to stderr so that at least we know about it...
	client.OnError = func(e error) {
		log.Printf("Error sending to stackdriver: %v\n", e)
	}
	return client
}

func closeClient(client *logging.Client) {
	if err := client.Close(); err != nil {
		log.Printf("Error closing client: %v", err)
	}
}

// We only care about the project ID here so don't bother parsing all the file
type credentialsPartial struct {
	ProjectID string `json:"project_id"`
}

func readProjectID(credentialsFile string) string {
	file, e := ioutil.ReadFile(credentialsFile)
	if e != nil {
		log.Fatalf("File error: %v\n", e)
	}

	var cred credentialsPartial
	e = json.Unmarshal(file, &cred)
	if e != nil {
		log.Fatalf("Error unmarshaling file %s: %v\n", credentialsFile, e)
	}
	return cred.ProjectID
}
