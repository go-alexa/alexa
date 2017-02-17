// Package server is a server for an Alexa Skill that performs validations.
package server

import (
	"net/http"
	"os"

	"encoding/json"

	"github.com/gorilla/handlers"

	"github.com/go-alexa/alexa/events"
	"github.com/go-alexa/alexa/parser"
	"github.com/go-alexa/alexa/validations"
)

// Host is the host for the HTTP server to listen on. By default, it uses the
// environment variable of HTTP_HOST.
var Host = os.Getenv("HTTP_HOST")

// Events is the event handler.
var Events events.EventHandler

// writeBadRequest writes a http.StatusBadRequest error and message
func writeBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(http.StatusText(http.StatusBadRequest)))
}

// writeServerError writes a http.StatusInternalServerError error and message
func writeServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// Handler is the function that handles the Alexa HTTP request.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Verify certificate is good
	cert, err := validations.ValidateCertificate(r)
	if err != nil {
		writeBadRequest(w)
		return
	}

	// Verify signature is good
	body, err := validations.ValidateSignature(r, cert)
	if err != nil {
		writeBadRequest(w)
		return
	}

	var data json.RawMessage

	err = json.Unmarshal(body, &data)
	if err != nil {
		writeBadRequest(w)
		return
	}

	ev, err := parser.Parse(data)
	if err != nil {
		writeBadRequest(w)
		return
	}

	// Make sure the request is good
	if err = validations.ValidateRequest(ev); err != nil {
		writeBadRequest(w)
		return
	}

	// Try and process the request
	resp, err := Events.Event(ev)
	if err != nil {
		writeServerError(w)
		return
	}

	// Convert the data into bytes to send back
	b, err := json.Marshal(resp)
	if err != nil {
		writeServerError(w)
		return
	}

	// Send back the data
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// Run starts the server. It mounts the handler on /alexa.
func Run(ev events.EventHandler) error {
	Events = ev

	http.HandleFunc("/alexa", Handler)

	return http.ListenAndServe(Host, handlers.LoggingHandler(os.Stdout,
		handlers.ProxyHeaders(http.DefaultServeMux)))
}
