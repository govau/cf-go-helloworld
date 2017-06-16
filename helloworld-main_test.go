package main

import (
	"net/http"
	"testing"
)

func testGuess(t *testing.T, expected string, ips []string) {
	if guessIPofRequester(&http.Request{
		Header: http.Header(map[string][]string{
			"X-Forwarded-For": ips,
		}),
	}) != expected {
		t.Fail()
	}
}

func TestIPAddressGuesser(t *testing.T) {
	testGuess(t, "bar", []string{"foo", "bar"})
	testGuess(t, "foo", []string{"foo"})
	testGuess(t, "unknown", []string{})
	testGuess(t, "unknown", []string{"10.0.0.0"})
	testGuess(t, "baz", []string{"baz", "10.0.0.0"})
	testGuess(t, "101.1.2.3", []string{"101.1.2.3, 10.0.0.138"})
}
