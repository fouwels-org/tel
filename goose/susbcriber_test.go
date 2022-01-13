package goose

import (
	"testing"
)

func TestSubscriber(t *testing.T) {

	code := Start()
	if code != 0 {
		err := GetError()
		t.Fatalf("%v", err)
	}
}
