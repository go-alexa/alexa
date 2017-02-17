package validations

import (
	"bytes"
	"errors"
	"log"
	"strings"

	"net/http"
	"net/url"

	"crypto/x509"
	"encoding/pem"

	"github.com/boltdb/bolt"
)

var certBucket = []byte("certs")

var (
	errNoChain         = errors.New("unable to find certificate chain header")
	errUnacceptableURL = errors.New("url provided is not acceptable")
)

// hasCheckedDB prevents us having to always call initDB
var hasCheckedDB = false

// ValidateCertificate ensures that a request was from Amazon.
// If DB is not nil, it caches certificate chains to prevent unneeded downloads
// of the same certificate multiple times. It then returns the certificate so
// it can be used to verify the signature later.
func ValidateCertificate(r *http.Request) (*x509.Certificate, error) {
	// First, we need to extract the chain URL from the request
	chainURL, err := getChainURL(r)
	if err != nil {
		return nil, err
	}

	// Next, we need to make sure it's a valid URL
	err = verifyChainURL(chainURL)
	if err != nil {
		return nil, err
	}

	// If we haven't checked the DB yet, do that
	if !hasCheckedDB {
		if DB != nil {
			if err = initDB(); err != nil {
				return nil, err
			}
		} else {
			log.Println("No database is currently configured, it is highly suggested to do so")
		}

		hasCheckedDB = true
	}

	var certChain []byte

	// If we have a DB, see if we have the cert already cached
	if DB != nil {
		certChain = getCertFromCache(chainURL)
	}

	// If we don't or it's zero length, try and load it
	if certChain == nil || len(certChain) == 0 {
		certChain, err = loadCertChain(chainURL)
		if err != nil {
			return nil, err
		}
	}

	// If we have a DB, cache it so we don't have to load it again
	if DB != nil {
		if err = saveCertToCache(chainURL, certChain); err != nil {
			return nil, err
		}
	}

	// Return our signing certificate after verifying it
	return verifyCert(certChain)
}

// getChainURL attempts to get the certificate chain URL
// from the request headers.
func getChainURL(r *http.Request) (string, error) {
	chainURL := r.Header.Get("SignatureCertChainUrl")
	if chainURL == "" {
		return "", errNoChain
	}
	return chainURL, nil
}

// verifyChainURL verifies that the URL provided is
// acceptable for provided requirements.
func verifyChainURL(chainURL string) error {
	u, err := url.Parse(chainURL)
	if err != nil {
		return err
	}

	if u.Scheme != "https" || u.Host != "s3.amazonaws.com" || !strings.HasPrefix(u.Path, "/echo.api/") {
		return errUnacceptableURL
	}

	return nil
}

// loadCertChain attempts to load the certificate chain into a byte array.
func loadCertChain(chainURL string) ([]byte, error) {
	resp, err := http.Get(chainURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	_, err = buf.ReadFrom(resp.Body)

	return buf.Bytes(), err
}

// verifyCert verifies that the signing certificate is part of the chain,
// is valid for the name provided, and has not expired.
func verifyCert(certChain []byte) (*x509.Certificate, error) {
	// We need to use their root certificate chain
	roots := x509.NewCertPool()

	// First certificate in file is always signing cert
	signCert, remaining := pem.Decode(certChain)

	// Everything else is part of the root chain
	roots.AppendCertsFromPEM(remaining)

	cert, err := x509.ParseCertificate(signCert.Bytes)
	if err != nil {
		return nil, err
	}

	// We need to verify the certificate is valid for this name
	opts := x509.VerifyOptions{
		DNSName: "echo-api.amazon.com",
		Roots:   roots,
	}

	// Actually verify the chain
	_, err = cert.Verify(opts)

	// Return if there was any error and the certificate for later uses
	return cert, err
}

// initDB attempts to create the certificates bucket.
func initDB() error {
	return DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(certBucket)
		return err
	})
}

// getCertFromCache attempts to get the certificate from a bolt database.
func getCertFromCache(chainURL string) []byte {
	var data []byte

	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(certBucket)
		data = b.Get([]byte(chainURL))

		return nil
	})

	return data
}

// saveCertToCache puts a certificate into the cache.
func saveCertToCache(chainURL string, cert []byte) error {
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(certBucket)
		return b.Put([]byte(chainURL), cert)
	})
}
