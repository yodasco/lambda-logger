package main

import (
	"encoding/json"

	"github.com/apex/apex"
	"github.com/yodasco/lambda-logger/lib/parser"
	"github.com/yodasco/lambda-logger/lib/stackdriver"
)

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		// Parse incoming event
		logs, err := parser.Parse(event)
		if err != nil {
			return nil, err
		}

		// Send to stackdriver
		stackdriver.LogEvents(logs)
		return "ok", nil
	})
}
