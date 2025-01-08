package xtjson

import (
	"errors"
	"testing"
)

func TestSearch(t *testing.T) {
	root, err := ParseString(`[{"a":1,"b":2},{"a":2,"b":3},{"a":3,"b":4},{"a":4,"b":5}]`)
	assertNil(t, err)

	ns, err := root.Search(NodeMatcherFunc(func(n *Node) bool {
		if !n.IsObject() {
			return false
		}
		nv, err := n.Key("b").Int()
		if err != nil || nv < 4 {
			return false
		}
		return true
	}), nil)
	assertNil(t, err)
	assertEqual(t, `[{"a":3,"b":4},{"a":4,"b":5}]`, ns.ToArray().Stringify())
}

func TestSearchKey(t *testing.T) {
	root, err := ParseString(`{"ka":"va", "kb":{"kkb1":"kkv1", "ka": 25}, "kc":123}`)
	assertNil(t, err)
	ns, err := root.SearchKey("ka", nil)
	assertNil(t, err)
	assertEqual(t, 2, len(ns))
	assertEqual(t, "va", ns[0].value)
	assertEqual(t, 25.0, ns[1].value)
	assertEqual(t, `["va",25]`, ns.ToArray().Stringify())
}

func TestSearchWithMatcherFunc(t *testing.T) {
	root, err := ParseString(`{"ka":"va", "kb":{"kkb1":"kkv1", "ka": 25}, "kc":123}`)
	assertNil(t, err)

	ns, err := root.Search(NodeMatcherFunc(func(n *Node) bool {
		return n.SelfKey() == "ka"
	}), nil)
	assertNil(t, err)
	assertEqual(t, 2, len(ns))
	assertEqual(t, "va", ns[0].value)
	assertEqual(t, 25.0, ns[1].value)
	assertEqual(t, `["va",25]`, ns.ToArray().Stringify())
}

func TestQueryPipe(t *testing.T) {
	root, err := ParseString(`{"ka":"va", "kb":{"kkb1":"kkv1", "ka": 25}, "kc":[1,2,3],"kd":{"kkd":{"kkb1":222}}}`)
	assertNil(t, err)
	ns, err := root.QueryPipe("$.kb", "${...}")
	assertNil(t, err)
	assertEqual(t, 2, len(ns))
	assertEqual(t, `["kkv1",25]`, ns.ToArray().Stringify())

	ns, err = root.QueryPipe("$.kc", "$[...]")
	assertNil(t, err)
	assertEqual(t, 3, len(ns))
	assertEqual(t, `[1,2,3]`, ns.ToArray().Stringify())

	ns, err = root.QueryPipe("$", "$...kkb1")
	assertNil(t, err)
	assertEqual(t, 2, len(ns))
	assertEqual(t, `["kkv1",222]`, ns.ToArray().Stringify())
}

func TestParseQuery(t *testing.T) {
	_, err := parseQuery("")
	if !errors.Is(err, ErrBadQuery) {
		t.Fatal("error is expected ErrBadQuery")
	}
	_, err = parseQuery("abc")
	if !errors.Is(err, ErrBadQuery) {
		t.Fatal("error is expected ErrBadQuery")
	}

	steps, err := parseQuery("$.foo")
	assertNil(t, err)
	assertEqual(t, []string{"$.foo"}, steps)

	steps, err = parseQuery("$.foo...boo")
	assertNil(t, err)
	assertEqual(t, []string{"$.foo", "$...boo"}, steps)

	steps, err = parseQuery("$...boo[35].x")
	assertNil(t, err)
	assertEqual(t, []string{"$...boo", "$[35].x"}, steps)
	steps, err = parseQuery("$[12]...boo.aaa{...}.x[10]")
	assertNil(t, err)
	assertEqual(t, []string{"$[12]", "$...boo", "$.aaa", "${...}", "$.x[10]"}, steps)
}

func TestQuery(t *testing.T) {
	root, err := ParseString(`{"ka":"va", "kb":{"kkb1":"kkv1", "ka": 25}, "kc":[1,2,3],"kd":{"kkd":{"kkb1":222}}}`)
	assertNil(t, err)
	ns, err := root.Query("$.kb{...}")
	assertNil(t, err)
	assertEqual(t, 2, len(ns))
	assertEqual(t, `["kkv1",25]`, ns.ToArray().Stringify())

	ns, err = root.Query("$.kc[...]")
	assertNil(t, err)
	assertEqual(t, 3, len(ns))
	assertEqual(t, `[1,2,3]`, ns.ToArray().Stringify())

	ns, err = root.Query("$...kkb1")
	assertNil(t, err)
	assertEqual(t, 2, len(ns))
	assertEqual(t, `["kkv1",222]`, ns.ToArray().Stringify())
}
