// Package parser allows for parsing an incoming request.
package parser

import (
	"encoding/json"
)

// Parse takes a raw JSON message and returns it as a formatted Event.
func Parse(msg json.RawMessage) (*Event, error) {
	var event Event

	if err := json.Unmarshal(msg, &event); err != nil {
		return nil, err
	}

	return &event, nil
}
