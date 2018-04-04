package main

import (
	"testing"
)

func TestDefaultConfiguration(t *testing.T) {
	test00 := DefaultConfiguration()
	if test00.Listen != "0.0.0.0:8080" {
		t.Errorf("test00: Listen is not the default value. Expected: %s Got: %s", "0.0.0.0:8080", test00.Listen)
	}
}
