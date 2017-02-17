// Package validations allows for verifying a request came from Amazon.
// It is required to properly validate requests in order to submit a Skill.
package validations

import (
	"github.com/boltdb/bolt"
)

// DB must be set to a valid database in order to properly cache certificates
// instead of downloading them for every request.
var DB *bolt.DB

// TimeLimit is the maximum variance allowed in the timestamp from current time.
// It must be under 150 seconds to submit to Amazon.
var TimeLimit float64 = 60

// AppID is your Skill's ID. It must be set to verify App ID in requests.
var AppID string
