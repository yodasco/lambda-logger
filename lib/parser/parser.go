package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"github.com/yodasco/lambda-logger/lib/types"
)

// AwsLambdaLogGroupEvent is used to decompose the incoming logs events.
// This is the structure of an incoming log event
// Incoming event looks like:
// {"awslogs":{"data":"===base64 zipped json array ==="}}
type AwsLambdaLogGroupEvent struct {
	AwsLogs struct {
		Data string `json:"data"`
	} `json:"awslogs"`
}

// Parse log lines from AWS Cloudwatch logs log groups.
func Parse(event json.RawMessage) (events types.CloudwatchLogEvents, err error) {
	var lambdaEvent AwsLambdaLogGroupEvent
	if err = json.Unmarshal(event, &lambdaEvent); err != nil {
		return
	}

	decoded, err := decode(lambdaEvent.AwsLogs.Data)
	if err != nil {
		return
	}

	err = json.Unmarshal(decoded, &events)
	if err != nil {
		return
	}
	// Work each entry to set proper timestamp, message, severity, labels
	for i := range events.LogEvents {
		e := events.LogEvents[i]
		// Fix timestamp
		e.Timestamp = fromUnixMilli(e.TimestampInt)

		// Try to parse the message as json.
		var asJSON map[string]string
		err := json.Unmarshal([]byte(e.Message), &asJSON)
		if err == nil {
			// Success, it's a json message, so treat its fields separately
			e.Message = asJSON["msg"]
			delete(asJSON, "msg")
			e.Level = asJSON["level"]
			delete(asJSON, "level")

			// And the rest of the json fields would go into labels
			e.Labels = asJSON
		} else {
			e.Labels = make(map[string]string)
		}
		events.LogEvents[i] = e
	}
	return
}

// Go does't convery unix milli, so we have to that that...
func fromUnixMilli(ms int64) time.Time {
	return time.Unix(ms/1000, 0).Add(
		time.Duration(ms%1000) * time.Millisecond)
}

// Decodes a base64 gzipped string
// That's how AWS log evnets show up
func decode(b64z string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(b64z)
	if err != nil {
		return nil, err
	}

	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	return buf.Bytes(), err
}
