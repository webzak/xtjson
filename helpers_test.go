package xtjson

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected: %v, Actual: %v", expected, actual)
	}
}

func assertNil(t *testing.T, value any) {
	t.Helper()
	if value != nil {
		t.Fatalf("Value: %v expected to be nil", value)
	}
}

func assertNotNil(t *testing.T, value any) {
	t.Helper()
	if value == nil {
		t.Fatalf("Value: %v expected to be not nil", value)
	}
}
