package xtjson

import (
	"testing"
)

func TestStrignifySingle(t *testing.T) {
	for _, json := range []string{`"foo"`, "true", "100", "100.1", "null", "{}", "[]"} {
		node, err := ParseString(json)
		assertParsed(t, node, err)
		assertEqual(t, json, node.Stringify())
	}
}

func TestStrignifyArray(t *testing.T) {
	json := `["aa", "bb", "cc", "dd", 1, 2, 3, 4, 5, false, true, null]`
	node, err := ParseString(json)
	assertParsed(t, node, err)
	assertEqual(t, json, node.Stringify(&Format{SpacesAfterColon: 1, SpacesAfterComma: 1}))
}

func TestStrignifyObject(t *testing.T) {
	json := `{"k1":"v1","k2":false,"k3":22.33}`
	node, err := ParseString(json)
	assertParsed(t, node, err)
	assertEqual(t, json, node.Stringify())
}

func TestStrignifyMultiLevel(t *testing.T) {
	json := `{"k1": "v1", "k2": [1, 2, 3], "k3": {"kk1": "vv1", "kk2": [true, false]}}`
	exp, err := format(json, 0)
	assertNil(t, err)
	node, err := ParseString(json)
	assertParsed(t, node, err)
	assertEqual(t, exp, node.Stringify())
}

func TestStrignifyPretty(t *testing.T) {
	json := `{"k1": "v1", "k2": [1, 2, 3], "k3": {"kk1": "vv1", "kk2": [true, false]}}`
	exp, err := format(json, 2)
	assertNil(t, err)
	node, err := ParseString(json)
	assertParsed(t, node, err)
	assertEqual(t, exp, node.Stringify(&Format{Indent: 2}))
}
