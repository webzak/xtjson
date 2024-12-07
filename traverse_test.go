package xtjson

import (
	"testing"
)

func TestPath(t *testing.T) {
	node, err := Parse(`[20, ["v1", {"k1": "ov1", "k2": {"kk1": "value"}}]]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, undef, node.Path(""))
	assertEqual(t, node, node.Path("$"))
	assertEqual(t, float64(20), node.Path("$[0]").value)
	assertEqual(t, "v1", node.Path("$[1][0]").value)
	assertEqual(t, "ov1", node.Path("$[1][1].k1").value)
	assertEqual(t, "value", node.Path("$[1][1].k2.kk1").value)

	node, err = Parse(`{"k1": ["v1", "v2", {"kk1": "vv1"}], "k2": "foo"}`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, "foo", node.Path("$.k2").value)
	assertEqual(t, Array, node.Path("$.k1").kind)
	assertEqual(t, "v1", node.Path("$.k1[0]").value)
	assertEqual(t, Object, node.Path("$.k1[2]").kind)
	assertEqual(t, "vv1", node.Path("$.k1[2].kk1").value)
}

func TestSelfPath(t *testing.T) {
	node, err := Parse(`[20, ["v1", {"k1": "ov1", "k2": {"kk1": true}}]]`)
	assertNil(t, err)
	assertNotNil(t, node)
	if err != nil {
		t.Fatal("parse error must be nil")
	}
	assertEqual(t, "$[1][1].k2.kk1", node.Idx(1).Idx(1).Key("k2").Key("kk1").SelfPath())
	node = nil
	assertEqual(t, "", node.SelfPath())
	node = undef
	assertEqual(t, "", node.SelfPath())
}
