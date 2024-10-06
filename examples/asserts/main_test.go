package main

import (
	"testing"

	"github.com/turtak/go-kit/testing/asserts"
)

func TestManuallyTrue(t *testing.T) {
	asserts.NotPanics(t, func() {
		caller1()
	})
}
