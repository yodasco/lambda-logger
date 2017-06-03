package parser

import (
	"encoding/json"
	"time"

	"github.com/yodasco/lambda-logger/lib/types"
)

// Parse log lines from AWS Cloudwatch logs log groups.
func Parse(bytes []byte) (events types.CloudwatchLogEvents, err error) {
	err = json.Unmarshal(bytes, &events)
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
