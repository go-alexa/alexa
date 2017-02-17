package validations

import (
	"errors"
	"math"
	"time"

	"github.com/go-alexa/alexa/parser"
)

var (
	errOutsideTime = errors.New("timestamp difference was greater than allowed")
	errWrongApp    = errors.New("application IDs do not match")
)

// ValidateRequest ensures the request was made within TimeLimit and was for
// this AppID.
func ValidateRequest(ev *parser.Event) error {
	if math.Abs(time.Since(ev.Request.Timestamp.ToTime()).Seconds()) > TimeLimit {
		return errOutsideTime
	}

	if ev.Session.Application.ID != AppID {
		return errWrongApp
	}

	return nil
}
