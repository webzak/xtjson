package xtjson

import (
	"testing"
)

func TestPath(t *testing.T) {
	node, err := Parse(`[20, ["v1", {"k1": "ov1", "k2": {"kk1": "value"}}]]`)
	assertParsed(t, node, err)
	assertEqual(t, undef, node.Path(""))
	assertEqual(t, node, node.Path("$"))
	assertEqual(t, float64(20), node.Path("$[0]").value)
	assertEqual(t, "v1", node.Path("$[1][0]").value)
	assertEqual(t, "ov1", node.Path("$[1][1].k1").value)
	assertEqual(t, "value", node.Path("$[1][1].k2.kk1").value)

	node, err = Parse(`{"k1": ["v1", "v2", {"kk1": "vv1"}], "k2": "foo"}`)
	assertParsed(t, node, err)
	assertEqual(t, "foo", node.Path("$.k2").value)
	assertEqual(t, Array, node.Path("$.k1").kind)
	assertEqual(t, "v1", node.Path("$.k1[0]").value)
	assertEqual(t, Object, node.Path("$.k1[2]").kind)
	assertEqual(t, "vv1", node.Path("$.k1[2].kk1").value)
}

func TestSelfPath(t *testing.T) {
	node, err := Parse(`[20, ["v1", {"k1": "ov1", "k2": {"kk1": true}}]]`)
	assertParsed(t, node, err)
	assertEqual(t, "$[1][1].k2.kk1", node.Idx(1).Idx(1).Key("k2").Key("kk1").SelfPath())
	node = nil
	assertEqual(t, "", node.SelfPath())
	node = undef
	assertEqual(t, "", node.SelfPath())
}

func TestParent(t *testing.T) {
	node, err := Parse(`[20]`)
	assertParsed(t, node, err)
	assertEqual(t, undef, node.Parent())
	assertEqual(t, node, node.Idx(0).Parent())
	node = nil
	assertEqual(t, undef, node.Parent())
	node = undef
	assertEqual(t, undef, node.Parent())
}

func TestChildren(t *testing.T) {
	node, err := Parse(`[20, 21]`)
	assertParsed(t, node, err)
	assertEqual(t, 2, len(node.Children()))
	assertEqual(t, 0, len(node.Idx(0).Children()))
	node = nil
	assertEqual(t, 0, len(node.Children()))
	node = undef
	assertEqual(t, 0, len(node.Children()))
}

func TestChildrenKeys(t *testing.T) {
	node, err := Parse(`{"ka":"va", "kb":"vb", "kc":[1,2,3]}`)
	assertParsed(t, node, err)
	assertEqual(t, []string{"ka", "kb", "kc"}, node.ChildrenKeys())
	assertEqual(t, 0, len(node.Key("kc").ChildrenKeys()))
	node = nil
	assertEqual(t, 0, len(node.ChildrenKeys()))
	node = undef
	assertEqual(t, 0, len(node.ChildrenKeys()))
}

func TestChildrenLength(t *testing.T) {
	node, err := Parse(`{"ka":"va", "kb":"vb", "kc":[1,2,3]}`)
	assertParsed(t, node, err)
	assertEqual(t, 3, node.ChildrenLength())
	assertEqual(t, 0, node.Key("kb").ChildrenLength())
	assertEqual(t, 3, node.Key("kc").ChildrenLength())
	node = nil
	assertEqual(t, 0, node.ChildrenLength())
	node = undef
	assertEqual(t, 0, node.ChildrenLength())
}

func TestWalker(t *testing.T) {
	root, err := Parse(`{"ka":"va", "kb":{"kkb1":"kkv1", "kkb2":[1,2]}, "kkb3":[]}`)
	assertParsed(t, root, err)
	walker, err := NewWalker(root, 0)
	assertNil(t, err)

	node, state := walker.Next()
	assertEqual(t, root, node)
	assertEqual(t, WalkEnter, state)

	node, state = walker.Next()
	assertEqual(t, "ka", node.key)
	assertEqual(t, WalkPass, state)

	node, state = walker.Next()
	assertEqual(t, "kb", node.key)
	assertEqual(t, WalkEnter, state)

	node, state = walker.Next()
	assertEqual(t, "kkb1", node.key)
	assertEqual(t, WalkPass, state)

	node, state = walker.Next()
	assertEqual(t, "kkb2", node.key)
	assertEqual(t, WalkEnter, state)

	node, state = walker.Next()
	assertEqual(t, 1.0, node.value)
	assertEqual(t, WalkPass, state)

	node, state = walker.Next()
	assertEqual(t, 2.0, node.value)
	assertEqual(t, WalkPass, state)

	node, state = walker.Next()
	assertEqual(t, "kkb2", node.key)
	assertEqual(t, WalkExit, state)

	node, state = walker.Next()
	assertEqual(t, "kb", node.key)
	assertEqual(t, WalkExit, state)

	node, state = walker.Next()
	assertEqual(t, "kkb3", node.key)
	assertEqual(t, WalkEnter, state)

	node, state = walker.Next()
	assertEqual(t, "kkb3", node.key)
	assertEqual(t, WalkExit, state)

	node, state = walker.Next()
	assertEqual(t, root, node)
	assertEqual(t, WalkExit, state)

	node, state = walker.Next()
	assertNil(t, node)
	assertEqual(t, WalkDone, state)

	node, state = walker.Next()
	assertNil(t, node)
	assertEqual(t, WalkDone, state)
}

func TestWalkerWithDeepLevel(t *testing.T) {
	root, err := Parse(`{"ka":"va", "kb":{"kkb1":"kkv1", "kkb2":[1,2]}, "kkb3":[]}`)
	assertParsed(t, root, err)
	walker, err := NewWalker(root, 1)
	assertNil(t, err)

	node, state := walker.Next()
	assertEqual(t, root, node)
	assertEqual(t, WalkEnter, state)

	node, state = walker.Next()
	assertEqual(t, "ka", node.key)
	assertEqual(t, WalkPass, state)

	node, state = walker.Next()
	assertEqual(t, "kb", node.key)
	assertEqual(t, WalkEnter, state)

	node, state = walker.Next()
	assertEqual(t, "kb", node.key)
	assertEqual(t, WalkExit, state)

	node, state = walker.Next()
	assertEqual(t, "kkb3", node.key)
	assertEqual(t, WalkEnter, state)

	node, state = walker.Next()
	assertEqual(t, "kkb3", node.key)
	assertEqual(t, WalkExit, state)

	node, state = walker.Next()
	assertEqual(t, root, node)
	assertEqual(t, WalkExit, state)

	node, state = walker.Next()
	assertNil(t, node)
	assertEqual(t, WalkDone, state)

	node, state = walker.Next()
	assertNil(t, node)
	assertEqual(t, WalkDone, state)
}

func TestWalkerStartDeeper(t *testing.T) {
	root, err := Parse(`{"ka":"va", "kb":{"kkb1":"kkv1", "kkb2":[1,2]}, "kkb3":[]}`)
	assertParsed(t, root, err)
	root = root.Key("kb")
	walker, err := NewWalker(root, 2)
	assertNil(t, err)

	node, state := walker.Next()
	assertEqual(t, root, node)
	assertEqual(t, WalkEnter, state)

	node, state = walker.Next()
	assertEqual(t, "kkb1", node.key)
	assertEqual(t, WalkPass, state)

	node, state = walker.Next()
	assertEqual(t, "kkb2", node.key)
	assertEqual(t, WalkEnter, state)

	node, state = walker.Next()
	assertEqual(t, 1.0, node.value)
	assertEqual(t, WalkPass, state)

	node, state = walker.Next()
	assertEqual(t, 2.0, node.value)
	assertEqual(t, WalkPass, state)

	node, state = walker.Next()
	assertEqual(t, "kkb2", node.key)
	assertEqual(t, WalkExit, state)

	node, state = walker.Next()
	assertEqual(t, root, node)
	assertEqual(t, WalkExit, state)

	node, state = walker.Next()
	assertNil(t, node)
	assertEqual(t, WalkDone, state)

	node, state = walker.Next()
	assertNil(t, node)
	assertEqual(t, WalkDone, state)
}
