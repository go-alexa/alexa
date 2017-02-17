# Golang library for Alexa Skills

This library is split into a few packages. There are the server, validations,
events, response, and parser packages. Each one tries to stay focused on what
it does so it's easy to implement a minimal amount into your own project.

There likely are optimizations possible or ways to make things simpler. As this
project has not reached a major version yet, pull requests that make backwards
incompatible changes are still welcome. Also, this means you should not depend
on API stability yet.

## Example

Look at [server/server_test.go](server/server_test.go).
