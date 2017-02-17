package validations

import (
	"errors"
	"math"
	"time"

	"encoding/json"
)

var (
	errOutsideTime = errors.New("timestamp difference was greater than allowed")
	errWrongApp    = errors.New("application IDs do not match")
)

type requestData struct {
	Request struct {
		Timestamp string `json:"timestamp"`
	} `json:"request"`

	Session struct {
		Application struct {
			ID string `json:"applicationId"`
		} `json:"application"`
	} `json:"session"`
}

// ValidateRequest ensures the request was made within TimeLimit and was for
// this AppID.
func ValidateRequest(data json.RawMessage) error {
	var r requestData

	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339, r.Request.Timestamp)
	if err != nil {
		return err
	}

	if math.Abs(time.Since(t).Seconds()) > TimeLimit {
		return errOutsideTime
	}

	if r.Session.Application.ID != AppID {
		return errWrongApp
	}

	return nil
}
