package xtjson

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, Actual: %v", expected, actual)
	}
}

func assertNil(t *testing.T, value any) {
	t.Helper()
	if value == nil {
		return
	}
	iv := reflect.ValueOf(value)
	if !iv.IsValid() {
		return
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		if iv.IsNil() {
			return
		}
	}
	t.Fatalf("value: %v expected to be nil", value)
}

func assertParsed(t *testing.T, node any, err any) {
	t.Helper()
	if node == nil {
		t.Fatalf("node %v expected to be not nil", node)
	}
	if err != nil {
		t.Fatalf("error %v expected to be nil", err)
	}
}

func assertString(t *testing.T, expected string, node *Node) {
	t.Helper()
	value, err := node.String()
	if err != nil {
		t.Fatalf("error %v expected to be nil", err)
	}
	if expected != value {
		t.Fatalf("value expected: %v, actual: %v", expected, value)
	}
}

func assertBool(t *testing.T, expected bool, node *Node) {
	t.Helper()
	value, err := node.Bool()
	if err != nil {
		t.Fatalf("error %v expected to be nil", err)
	}
	if expected != value {
		t.Fatalf("value expected: %v, actual: %v", expected, value)
	}
}

func assertNumber(t *testing.T, expected float64, node *Node) {
	t.Helper()
	value, err := node.Number()
	if err != nil {
		t.Fatalf("error %v expected to be nil", err)
	}
	if expected != value {
		t.Fatalf("value expected: %v, actual: %v", expected, value)
	}
}

func assertInt(t *testing.T, expected int, node *Node) {
	t.Helper()
	value, err := node.Int()
	if err != nil {
		t.Fatalf("error %v expected to be nil", err)
	}
	if expected != value {
		t.Fatalf("value expected: %v, actual: %v", expected, value)
	}
}

func format(jsonString string, indent int) (string, error) {
	if indent < 0 {
		return "", errors.New("negative indent")
	}
	var v any
	if err := json.Unmarshal([]byte(jsonString), &v); err != nil {
		return "", err
	}
	var result []byte
	var err error
	if indent == 0 {
		result, err = json.Marshal(v)
	} else {
		result, err = json.MarshalIndent(v, "", strings.Repeat(" ", indent))
	}
	if err != nil {
		return "", err
	}
	return string(result), nil
}
