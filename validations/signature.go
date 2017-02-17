package validations

import (
	"bytes"
	"errors"

	"io/ioutil"

	"net/http"

	"crypto/x509"
	"encoding/base64"
)

var (
	errNoSignature = errors.New("unable to find signature header")
)

// ValidateSignature ensures that the request body was made with the
// certificate provided in the header.
func ValidateSignature(r *http.Request, cert *x509.Certificate) error {
	// First, get the signature from the headers
	sig, err := getSignature(r)
	if err != nil {
		return err
	}

	// Then get the request body
	body, err := getBody(r)
	if err != nil {
		return err
	}

	// Return if the signature properly verified
	return verifySignature(sig, body, cert)
}

// getSignature attempts to get the signature from the headers.
func getSignature(r *http.Request) (string, error) {
	sig := r.Header.Get("Signature")
	if sig == "" {
		return sig, errNoSignature
	}
	return sig, nil
}

// getBody gets the contents of the body as a byte array and returns it.
func getBody(r *http.Request) ([]byte, error) {
	// Read the entire body in
	body, err := ioutil.ReadAll(r.Body)

	// Return the body contents so we can work with it later
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	return body, err
}

// verifySignature actually checks that the body matches the signature.
func verifySignature(sig string, body []byte, cert *x509.Certificate) error {
	decodedSig, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	return cert.CheckSignature(x509.SHA1WithRSA, body, decodedSig)
}
