package xtjson

import (
	"errors"
	"testing"
)

func TestParseSingle(t *testing.T) {
	node, err := Parse(`"foo"`)
	assertParsed(t, node, err)
	assertString(t, "foo", node)

	node, err = Parse(`true`)
	assertParsed(t, node, err)
	assertBool(t, true, node)

	node, err = Parse(`100`)
	assertParsed(t, node, err)
	assertInt(t, 100, node)

	node, err = Parse(`100.1`)
	assertParsed(t, node, err)
	assertNumber(t, 100.1, node)

	node, err = Parse(`null`)
	assertParsed(t, node, err)
	assertEqual(t, true, node.IsNull())

	node, err = Parse(`{}`)
	assertParsed(t, node, err)
	assertEqual(t, true, node.IsObject())

	node, err = Parse(`[]`)
	assertParsed(t, node, err)
	assertEqual(t, true, node.IsArray())
}

func TestParsePlainObject(t *testing.T) {
	node, err := Parse(`{"key1":"value1", "key2": false, "key3":  20, "key4": null}`)
	assertParsed(t, node, err)
	assertEqual(t, Object, node.Type())
	assertEqual(t, 4, len(node.children))
	assertEqual(t, 0, node.value.(keymap)["key1"])
	assertEqual(t, 1, node.value.(keymap)["key2"])
	assertEqual(t, 2, node.value.(keymap)["key3"])
	assertEqual(t, 3, node.value.(keymap)["key4"])
	c0 := node.children[0]
	assertEqual(t, String, c0.Type())
	assertEqual(t, "value1", c0.value)
	assertEqual(t, node, c0.parent)
	assertEqual(t, 0, c0.idx)
	c1 := node.children[1]
	assertEqual(t, Bool, c1.Type())
	assertEqual(t, false, c1.value)
	assertEqual(t, node, c1.parent)
	assertEqual(t, 1, c1.idx)
	c2 := node.children[2]
	assertEqual(t, Number, c2.Type())
	assertEqual(t, float64(20), c2.value)
	assertEqual(t, node, c2.parent)
	assertEqual(t, 2, c2.idx)
	c3 := node.children[3]
	assertEqual(t, Null, c3.Type())
	assertEqual(t, node, c3.parent)
	assertEqual(t, 3, c3.idx)
}

func TestParsePlainArray(t *testing.T) {
	node, err := Parse(`["value1", false, 20, null]`)
	assertParsed(t, node, err)
	assertEqual(t, Array, node.Type())
	assertEqual(t, 4, len(node.children))
	c0 := node.children[0]
	assertEqual(t, String, c0.Type())
	assertEqual(t, "value1", c0.value)
	assertEqual(t, node, c0.parent)
	assertEqual(t, 0, c0.idx)
	c1 := node.children[1]
	assertEqual(t, Bool, c1.Type())
	assertEqual(t, false, c1.value)
	assertEqual(t, node, c1.parent)
	assertEqual(t, 1, c1.idx)
	c2 := node.children[2]
	assertEqual(t, Number, c2.Type())
	assertEqual(t, float64(20), c2.value)
	assertEqual(t, node, c2.parent)
	assertEqual(t, 2, c2.idx)
	c3 := node.children[3]
	assertEqual(t, Null, c3.Type())
	assertEqual(t, node, c3.parent)
	assertEqual(t, 3, c3.idx)
}

func TestParseNested(t *testing.T) {
	node, err := Parse(`[20, ["v1", "v2"], {"k1": "ov1", "k2": {"kk1": true}}]`)
	assertParsed(t, node, err)
	assertEqual(t, Array, node.Type())
	assertEqual(t, 3, len(node.children))
	c0 := node.children[0]
	assertEqual(t, Number, c0.Type())
	assertEqual(t, float64(20), c0.value)
	assertEqual(t, node, c0.parent)
	assertEqual(t, 0, c0.idx)
}

func TestIncorrectJson(t *testing.T) {
	inputs := []string{
		`{"k1": "v1", "k2": "v2",}`,
		`{"k1": "v1", "k2": {"kk1": "vv1"}`,
		`{"k1": "v1", "k2": {"kk1": "vv1"}}}`,
	}
	for _, s := range inputs {
		_, err := Parse(s)
		if !errors.Is(err, ErrInvalidJson) {
			t.Fatal("expected ErrInvalid json")
		}
	}
}

func TestDuplicateKey(t *testing.T) {

	_, err := Parse(`{"k1": "v1", "k1": "v2"}`)
	if !errors.Is(err, ErrDuplicateKey) {
		t.Fatal("expected ErrDuplicateKey error")
	}

}
