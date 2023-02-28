package main

import "testing"

func Hello() {
	result := "Hello"
	expected := "Hello"

	if result != expected {
		t.Errorf("Suspicious")

	}
}
