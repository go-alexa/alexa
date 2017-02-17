// Package events allows for a more structured approach to handle intents.
package events

import (
	"errors"

	"github.com/go-alexa/alexa/parser"
	"github.com/go-alexa/alexa/response"
)

const (
	// RequestLaunch is when there has been a request for the app launch.
	RequestLaunch = "LaunchRequest"
	// RequestIntent is when there has been a request for an intent.
	RequestIntent = "IntentRequest"
	// RequestEnded is when the session has ended.
	RequestEnded = "SessionEndedRequest"
)

var (
	// ErrNoHandler means there was an intent called without a handler defined.
	ErrNoHandler = errors.New("no handler was specificed for this intent")
)

// LaunchFunc is a func called when the app is launched.
type LaunchFunc func(*parser.Event) (*response.Response, error)

// IntentFunc is a func called when an intent has been called.
type IntentFunc func(*parser.Event) (*response.Response, error)

// EndedFunc is a func called when the session has ended.
type EndedFunc func(*parser.Event) (*response.Response, error)

// EventHandler is a handler for any Alexa events.
type EventHandler interface {
	Add(string, IntentFunc) EventHandler
	Event(*parser.Event) (*response.Response, error)
}

// Handler is a default implementation of the EventHandler.
type Handler struct {
	LaunchHandler LaunchFunc
	EndedHandler  EndedFunc

	IntentHandlers map[string]IntentFunc
}

// New creates a new Handler and initializes the map.
func New() *Handler {
	return &Handler{
		IntentHandlers: make(map[string]IntentFunc),
	}
}

// Add adds a new intent handler to the map for a specific intent name.
func (e *Handler) Add(intent string, handler IntentFunc) EventHandler {
	e.IntentHandlers[intent] = handler

	return e
}

// Event processes all event handlers for an event. It then returns the response
// or any errors that occurred while processing.
func (e *Handler) Event(ev *parser.Event) (*response.Response, error) {
	var resp *response.Response
	var err error

	switch ev.Request.Type {
	case RequestLaunch:
		if e.LaunchHandler != nil {
			resp, err = e.LaunchHandler(ev)
		} else {
			err = ErrNoHandler
		}

	case RequestEnded:
		if e.EndedHandler != nil {
			resp, err = e.EndedHandler(ev)
		} else {
			err = ErrNoHandler
		}

	case RequestIntent:
		if fn, ok := e.IntentHandlers[ev.Request.Intent.Name]; ok {
			resp, err = fn(ev)
		} else {
			err = ErrNoHandler
		}
	}

	return resp, err
}
