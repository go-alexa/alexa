// Package events allows for a more structured approach to handle intents.
package events

import (
	"log"

	"github.com/b00giZm/golexa"
)

// LaunchFunc is a function that is called when the Skill is opened.
type LaunchFunc func(*golexa.Alexa, *golexa.Request, *golexa.Session) *golexa.Response

// EndedFunc is a function that is called when the session is ended.
type EndedFunc func(*golexa.Alexa, *golexa.Request, *golexa.Session) *golexa.Response

// IntentFunc is a function that is called when a specific intent is requested.
type IntentFunc func(*golexa.Alexa, *golexa.Intent, *golexa.Request, *golexa.Session) *golexa.Response

// LaunchHandler is called whenever the Skill is opened.
var LaunchHandler LaunchFunc

// EndedHandler is called when the session is ended.
var EndedHandler EndedFunc

// UnhandledHandler is called whenever an intent does not have a handler.
var UnhandledHandler IntentFunc

// IntentHandler is called whenever an intent with the specific name is
// requested. Note that this should remain nil if you want to use the normal
// intent handler system.
var IntentHandler IntentFunc

// intentHandlers is a map of the handlers registered.
var intentHandlers map[string]IntentFunc

var Debug = false

// init allocates the map.
func init() {
	intentHandlers = make(map[string]IntentFunc)
}

// AddIntentHandler adds a handler for a specific intent.
func AddIntentHandler(intent string, fn IntentFunc) {
	if Debug {
		log.Printf("Added handler for intent: %s\n", intent)
	}

	intentHandlers[intent] = fn
}

// handleIntent is the default handler for passing intents to the right place.
func handleIntent(a *golexa.Alexa, i *golexa.Intent, r *golexa.Request, s *golexa.Session) *golexa.Response {
	if fn, ok := intentHandlers[i.Name]; ok {
		if Debug {
			log.Printf("Running handler for intent: %s\n", i.Name)
		}

		resp := fn(a, i, r, s)

		if Debug && resp == nil {
			log.Printf("Handler for %s returned nil!\n", i.Name)
		}

		return resp
	} else if UnhandledHandler != nil {
		if Debug {
			log.Println("Running handler for unhandled intent")
		}

		return UnhandledHandler(a, i, r, s)
	}

	return nil
}

// Register attaches the event system to a standard golexa.Alexa instance.
func Register(golex *golexa.Alexa) {
	if LaunchHandler != nil {
		if Debug {
			log.Println("Running handler for launch")
		}

		golex.OnLaunch(LaunchHandler)
	}

	if EndedHandler != nil {
		if Debug {
			log.Println("Running handler for session end")
		}

		golex.OnSessionEnded(EndedHandler)
	}

	if IntentHandler != nil {
		if Debug {
			log.Println("Running handler for intent")
		}

		golex.OnIntent(IntentHandler)
	} else {
		golex.OnIntent(handleIntent)
	}
}
