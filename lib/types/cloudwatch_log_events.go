package types

import "time"

// CloudwatchLogEvents defiles the structure of a cloudwatch log event that
// contains multiple LogEvents inside
type CloudwatchLogEvents struct {
	MessageType         string     `json:"messageType"`
	Owner               string     `json:"owner"`
	LogGroup            string     `json:"logGroup"`
	LogStream           string     `json:"logStream"`
	SubscriptionFilters []string   `json:"subscriptionFilters"`
	LogEvents           []LogEvent `json:"logEvents"`
}

// LogEvent is the single atomic log line representation
type LogEvent struct {
	ID           string `json:"id"`
	TimestampInt int64  `json:"timestamp"`
	Timestamp    time.Time
	Message      string `json:"message"`
	Level string
	Labels map[string]string
}
