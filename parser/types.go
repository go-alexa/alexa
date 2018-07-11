package parser

import (
	"time"
)

// Event is the base type for any request from Amazon.
type Event struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Context Context `json:"context,omitempty"`
	Request Request `json:"request"`
}

// Session is information about the user, any set session data, or the app.
type Session struct {
	ID          string            `json:"sessionId"`
	IsNew       bool              `json:"new"`
	Attributes  SessionAttributes `json:"attributes,omitempty"`
	Application Application       `json:"application,omitempty"`
	User        User              `json:"user,omitempty"`
}

// SessionAttributes are arbitrary data set in a previous request. They only
// last for the duration of a single session.
type SessionAttributes map[string]interface{}

// Application is information about what application is being called.
type Application struct {
	ID string `json:"applicationId"`
}

// Permissions hold the consent token which can retrieve additional information regarding
// the current user
type Permissions struct {
	ConsentToken string `json:"consentToken,omitempty"`
}

// User is information about the user, including access token if one has been
// set through linking an account.
type User struct {
	ID          string      `json:"userId"`
	AccessToken string      `json:"accessToken"`
	Permissions Permissions `json:"permissions,omitempty"`
}

// Request is information about the request, including the intent and data.
type Request struct {
	ID        string `json:"requestId"`
	Type      string `json:"type"`
	Locale    string `json:"locale,omitempty"`
	Timestamp *Time  `json:"timestamp"`
	Intent    Intent `json:"intent,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

// Time is a timestamp of the request.
type Time time.Time

// MarshalJSON allows for encoding the timestamp in the correct format.
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(time.RFC3339)), nil
}

// UnmarshalJSON allows for decoding the time in the correct format.
func (t *Time) UnmarshalJSON(b []byte) error {

	parsedTime, err := time.Parse("\""+time.RFC3339+"\"", string(b))
	if err != nil {
		return err
	}

	*t = Time(parsedTime)

	return nil
}

// ToTime allows for converting a Time struct into a standard time.Time.
func (t Time) ToTime() time.Time {
	return time.Time(t)
}

// Intent is information about the intent, including its name and slots.
type Intent struct {
	Name               string          `json:"name"`
	Slots              map[string]Slot `json:"slots,omitempty"`
	ConfirmationStatus string          `json:"confirmationStatus,omitempty"`
}

// Slot is the data for an intent.
type Slot struct {
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Resolutions Resolutions `json:"resolutions,omitempty"`
}

// Resolutions is the data for a resolution contained within a slot that contains an array of resolution
type Resolutions struct {
	ResolutionsPerAuthority []Resolution `json:"resolutionsPerAuthority"`
}

// Resolution is the data used within the resolutions data array. It information about the resolved slot value
type Resolution struct {
	Authority string            `json:"authority"`
	Status    ResolutionStatus  `json:"status"`
	Values    []ResolutionValue `json:"values,omitempty"`
}

// ResolutionStatus is the resolution status used within resolution
type ResolutionStatus struct {
	Code string `json:"code"`
}

/*
ResolutionValue is the resolution value used within the resolution
 The values struct is really stupid. It is the following:
"values": [
	{
		"value": {
			"name": "Powerball",
			"id": "b29a2c0b8e34ca0d91642d1eae81f5cf"
		}
	}
]
*/
type ResolutionValue struct {
	Value ResolutionValueValue `json:"value"`
}

// ResolutionValueValue is the resolution value used within the resolution
type ResolutionValueValue struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// AudioPlayer holds info about the audioplayer usage of the users device
type AudioPlayer struct {
	PlayerActivity string `json:"playerActivity,omitempty"`
}

// SupportedInterfaces Holds information regarding supported interfaces (future oriented - unneccessary atm)
type SupportedInterfaces struct {
	AudioPlayer AudioPlayer `json:"AudioPlayer,omitempty"`
}

// Device holds information regarding the users device
type Device struct {
	ID         string              `json:"deviceId,omitempty"`
	Interfaces SupportedInterfaces `json:"supportedInterfaces,omitempty"`
}

// System holds information regarding the users system (dot, alexa ... future stuff)
type System struct {
	Application Application `json:"application,omitempty"`
	User        User        `json:"user,omitempty"`
	Device      Device      `json:"device,omitempty"`
	APIEndpoint string      `json:"apiEndpoint,omitempty"`
}

// Context  holds more context to the users setup
type Context struct {
	AudioPlayer AudioPlayer `json:"AudioPlayer,omitempty"`
	System      System      `json:"System,omitempty"`
}
