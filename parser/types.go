package parser

import (
	"time"
)

// Event is the base type for any request from Amazon.
type Event struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Request Request `json:"request"`
}

// Session is information about the user, any set session data, or the app.
type Session struct {
	ID          string            `json:"sessionId"`
	IsNew       bool              `json:"new"`
	Attributes  SessionAttributes `json:"attributes"`
	Application Application       `json:"application"`
	User        User              `json:"user"`
}

// SessionAttributes are arbitrary data set in a previous request. They only
// last for the duration of a single session.
type SessionAttributes map[string]interface{}

// Application is information about what application is being called.
type Application struct {
	ID string `json:"applicationId"`
}

// User is information about the user, including access token if one has been
// set through linking an account.
type User struct {
	ID          string `json:"userId"`
	AccessToken string `json:"accessToken"`
}

// Request is information about the request, including the intent and data.
type Request struct {
	ID        string `json:"requestId"`
	Type      string `json:"type"`
	Timestamp Time   `json:"timestamp"`
	Locale    string `json:"locale,omitempty"`
	Intent    Intent `json:"intent,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

// Time is a timestamp of the request.
type Time time.Time

// MarshalJSON allows for encoding the timestamp in the correct format.
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(time.RFC3339Nano)), nil
}

// UnmarshalJSON allows for decoding the time in the correct format.
func (t Time) UnmarshalJSON(b []byte) error {
	time, err := time.Parse("\""+time.RFC3339Nano+"\"", string(b))
	if err != nil {
		return err
	}

	t = Time(time)

	return nil
}

// ToTime allows for converting a Time struct into a standard time.Time.
func (t Time) ToTime() time.Time {
	return time.Time(t)
}

// Intent is information about the intent, including its name and slots.
type Intent struct {
	Name  string          `json:"name"`
	Slots map[string]Slot `json:"slots"`
}

// Slot is the data for an intent.
type Slot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
