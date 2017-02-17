# Golang library for Alexa Skills

This library is split into three parts. There are the server, validations,
and events packages. Each one tries to stay focused on what it does so it's easy
to implement just one part into your own project.

There likely are optimizations possible or ways to make things simpler. As this
project has not reached a major version yet, pull requests that make backwards
incompatible changes are still welcome. Also, this means you should not depend
on API stability yet.

## Example

```go
package main

import (
	"github.com/boltdb/bolt"

	"github.com/b00giZm/golexa"

	"github.com/go-alexa/alexa/events"
	"github.com/go-alexa/alexa/server"
	"github.com/go-alexa/alexa/validations"
)

func main() {
	golex := golexa.Default()

	d, err := bolt.Open("info.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	validations.DB = d

	events.Debug = true

	events.AddIntentHandler("HelloWorld",
		func(a *golexa.Alexa, intent *golexa.Intent, req *golexa.Request, session *golexa.Session) *golexa.Response {
			return a.
				Response().
				AddPlainTextSpeech("Hello, world!")
		})

	events.Register(golex)

	server.Run(golex)
}
```
