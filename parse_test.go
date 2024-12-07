package xtjson

import (
	"testing"
)

func TestParsSingle(t *testing.T) {
	node, err := Parse(`"foo"`)
	assertNil(t, err)
	assertEqual(t, String, node.kind)
	assertEqual(t, "foo", node.value)

	node, err = Parse(`true`)
	assertNil(t, err)
	assertEqual(t, Bool, node.kind)
	assertEqual(t, true, node.value)

	node, err = Parse(`100`)
	assertNil(t, err)
	assertEqual(t, Number, node.kind)
	assertEqual(t, 100.0, node.value)

	node, err = Parse(`100.1`)
	assertNil(t, err)
	assertEqual(t, Number, node.kind)
	assertEqual(t, 100.1, node.value)

	node, err = Parse(`null`)
	assertNil(t, err)
	assertEqual(t, Null, node.kind)
	assertEqual(t, nil, node.value)

	node, err = Parse(`{}`)
	assertNil(t, err)
	assertEqual(t, Object, node.kind)
	assertEqual(t, nil, node.value)

	node, err = Parse(`[]`)
	assertNil(t, err)
	assertEqual(t, Array, node.kind)
	assertEqual(t, nil, node.value)
}

func TestParsePlainObject(t *testing.T) {
	node, err := Parse(`{"key1":"value1", "key2": false, "key3":  20, "key4": null}`)
	assertNil(t, err)
	assertEqual(t, Object, node.kind)
	assertNil(t, node.value)
	assertEqual(t, []string{"key1", "key2", "key3", "key4"}, node.keys)
	assertEqual(t, 4, len(node.children))
	c0 := node.children[0]
	assertEqual(t, String, c0.kind)
	assertEqual(t, "value1", c0.value)
	assertEqual(t, node, c0.parent)
	assertEqual(t, 0, c0.idx)
	c1 := node.children[1]
	assertEqual(t, Bool, c1.kind)
	assertEqual(t, false, c1.value)
	assertEqual(t, node, c1.parent)
	assertEqual(t, 1, c1.idx)
	c2 := node.children[2]
	assertEqual(t, Number, c2.kind)
	assertEqual(t, float64(20), c2.value)
	assertEqual(t, node, c2.parent)
	assertEqual(t, 2, c2.idx)
	c3 := node.children[3]
	assertEqual(t, Null, c3.kind)
	assertEqual(t, nil, c3.value)
	assertEqual(t, node, c3.parent)
	assertEqual(t, 3, c3.idx)
}

func TestParsePlainArray(t *testing.T) {
	node, err := Parse(`["value1", false, 20, null]`)
	assertNil(t, err)
	assertEqual(t, Array, node.kind)
	assertNil(t, node.value)
	assertEqual(t, 0, len(node.keys))
	assertEqual(t, 4, len(node.children))
	c0 := node.children[0]
	assertEqual(t, String, c0.kind)
	assertEqual(t, "value1", c0.value)
	assertEqual(t, node, c0.parent)
	assertEqual(t, 0, c0.idx)
	c1 := node.children[1]
	assertEqual(t, Bool, c1.kind)
	assertEqual(t, false, c1.value)
	assertEqual(t, node, c1.parent)
	assertEqual(t, 1, c1.idx)
	c2 := node.children[2]
	assertEqual(t, Number, c2.kind)
	assertEqual(t, float64(20), c2.value)
	assertEqual(t, node, c2.parent)
	assertEqual(t, 2, c2.idx)
	c3 := node.children[3]
	assertEqual(t, Null, c3.kind)
	assertEqual(t, nil, c3.value)
	assertEqual(t, node, c3.parent)
	assertEqual(t, 3, c3.idx)
}

func TestParseNested(t *testing.T) {
	node, err := Parse(`[20, ["v1", "v2"], {"k1": "ov1", "k2": {"kk1": true}}]`)
	assertNil(t, err)
	assertEqual(t, Array, node.kind)
	assertNil(t, node.value)
	assertEqual(t, 0, len(node.keys))
	assertEqual(t, 3, len(node.children))
	c0 := node.children[0]
	assertEqual(t, Number, c0.kind)
	assertEqual(t, float64(20), c0.value)
	assertEqual(t, node, c0.parent)
	assertEqual(t, 0, c0.idx)
}
