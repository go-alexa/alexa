// Package server is a server for an Alexa Skill that performs validations.
package server

import (
	"net/http"
	"os"

	"encoding/json"

	"github.com/gorilla/handlers"

	"github.com/b00giZm/golexa"

	"github.com/ixchi/alexa-colorful/validations"
)

// Host is the host for the HTTP server to listen on. By default, it uses the
// environment variable of HTTP_HOST.
var Host = os.Getenv("HTTP_HOST")

// g is the Alexa instance.
var g *golexa.Alexa

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
	if err = validations.ValidateSignature(r, cert); err != nil {
		writeBadRequest(w)
		return
	}

	// Now that we know those things are good, try to decode the request payload
	decoder := json.NewDecoder(r.Body)

	var data json.RawMessage

	err = decoder.Decode(&data)
	if err != nil {
		writeBadRequest(w)
		return
	}

	// Make sure the request is good
	if err = validations.ValidateRequest(data); err != nil {
		writeBadRequest(w)
		return
	}

	// Try and process the request
	resp, err := g.Process(data)
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
func Run(golex *golexa.Alexa) error {
	g = golex

	http.HandleFunc("/alexa", Handler)

	return http.ListenAndServe(Host, handlers.LoggingHandler(os.Stdout,
		handlers.ProxyHeaders(http.DefaultServeMux)))
}
