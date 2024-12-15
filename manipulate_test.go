package xtjson

import (
	"errors"
	"testing"
)

func TestNewArray(t *testing.T) {
	node := NewArray()
	assertEqual(t, Array, node.kind)
}

func TestNewObject(t *testing.T) {
	node := NewObject()
	assertEqual(t, Object, node.kind)
}

func TestNewString(t *testing.T) {
	node := NewString("foo")
	assertEqual(t, String, node.kind)
	assertEqual(t, "foo", node.value)
}

func TestNewBool(t *testing.T) {
	node := NewBool(true)
	assertEqual(t, Bool, node.kind)
	assertEqual(t, true, node.value)
}

func TestNewNull(t *testing.T) {
	node := NewNull()
	assertEqual(t, Null, node.kind)
}

func TestNewNumber(t *testing.T) {
	value := 1.2
	node := NewNumber(value)
	assertEqual(t, Number, node.kind)
	assertEqual(t, value, node.value)
}

func TestNewInt(t *testing.T) {
	value := 1
	node := NewInt(value)
	assertEqual(t, Number, node.kind)
	assertEqual(t, float64(value), node.value)
}

func TestAppend(t *testing.T) {
	root, err := Parse(`[1,2]`)
	assertParsed(t, root, err)
	node := NewString("s1")
	err = root.Append(node)
	assertNil(t, err)
	assertEqual(t, 2, node.idx)
	assertEqual(t, root, node.parent)
	assertEqual(t, node, root.children[2])

	err = root.Append(node)
	if !errors.Is(err, ErrNodeHasParent) {
		t.Fatal("expected error ErrNodeHasParent")
	}
	root = NewString("foo")
	node = NewString("boo")
	err = root.Append(node)
	if !errors.Is(err, ErrInvalidNodeForOperation) {
		t.Fatal("expected error ErrInvalidNodeForOperation")
	}
}

func TestAppendString(t *testing.T) {
	root, err := Parse(`["a","b"]`)
	assertParsed(t, root, err)
	err = root.AppendString("foo")
	assertNil(t, err)
	assertEqual(t, `["a","b","foo"]`, root.Stringify())
}

func TestAppendBool(t *testing.T) {
	root, err := Parse(`["a","b"]`)
	assertParsed(t, root, err)
	err = root.AppendBool(true)
	assertNil(t, err)
	assertEqual(t, `["a","b",true]`, root.Stringify())
}

func TestAppendNumber(t *testing.T) {
	root, err := Parse(`["a","b"]`)
	assertParsed(t, root, err)
	err = root.AppendNumber(2.111)
	assertNil(t, err)
	assertEqual(t, `["a","b",2.111]`, root.Stringify())
}

func TestAppendInt(t *testing.T) {
	root, err := Parse(`["a","b"]`)
	assertParsed(t, root, err)
	err = root.AppendInt(22)
	assertNil(t, err)
	assertEqual(t, `["a","b",22]`, root.Stringify())
}

func TestAppendNull(t *testing.T) {
	root, err := Parse(`["a","b"]`)
	assertParsed(t, root, err)
	err = root.AppendNull()
	assertNil(t, err)
	assertEqual(t, `["a","b",null]`, root.Stringify())
}

func TestSet(t *testing.T) {
	root, err := Parse(`{"a":1}`)
	assertParsed(t, root, err)
	node := NewString("foo")
	err = root.Set("k", node)
	assertNil(t, err)
	assertEqual(t, 1, node.idx)
	assertEqual(t, root, node.parent)
	assertEqual(t, node, root.children[1])
	assertEqual(t, 1, root.keymap["k"])
	assertEqual(t, `{"a":1,"k":"foo"}`, root.Stringify())

	node2 := NewString("aa")
	err = root.Set("k", node2)
	assertNil(t, err)
	assertEqual(t, 1, node2.idx)
	assertEqual(t, root, node2.parent)
	assertEqual(t, node2, root.children[1])
	assertEqual(t, 1, root.keymap["k"])
	assertEqual(t, `{"a":1,"k":"aa"}`, root.Stringify())

	err = root.Set("kk", node)
	if !errors.Is(err, ErrNodeHasParent) {
		t.Fatal("expected error ErrNodeHasParent")
	}

	root = NewString("foo")
	node = NewString("boo")
	err = root.Set("k", node)
	if !errors.Is(err, ErrInvalidNodeForOperation) {
		t.Fatal("expected error ErrInvalidNodeForOperation")
	}

	root, err = Parse(`{"a":0,"b":1,"c":2,"d":3,"e":4}`)
	assertParsed(t, root, err)
	node = NewString("foo")
	err = root.Set("c", node)
	assertNil(t, err)
	assertEqual(t, `{"a":0,"b":1,"c":"foo","d":3,"e":4}`, root.Stringify())
	assertEqual(t, root, node.parent)
	assertEqual(t, 2, node.idx)
}

func TestSetString(t *testing.T) {
	root, err := Parse(`{"a":"b"}`)
	assertParsed(t, root, err)
	err = root.SetString("c", "d")
	assertNil(t, err)
	assertEqual(t, `{"a":"b","c":"d"}`, root.Stringify())
}

func TestSetBool(t *testing.T) {
	root, err := Parse(`{"a":"b"}`)
	assertParsed(t, root, err)
	err = root.SetBool("c", false)
	assertNil(t, err)
	assertEqual(t, `{"a":"b","c":false}`, root.Stringify())
}

func TestSetNumber(t *testing.T) {
	root, err := Parse(`{"a":"b"}`)
	assertParsed(t, root, err)
	err = root.SetNumber("c", 23.23)
	assertNil(t, err)
	assertEqual(t, `{"a":"b","c":23.23}`, root.Stringify())
}

func TestSetInt(t *testing.T) {
	root, err := Parse(`{"a":"b"}`)
	assertParsed(t, root, err)
	err = root.SetInt("c", 23)
	assertNil(t, err)
	assertEqual(t, `{"a":"b","c":23}`, root.Stringify())
}

func TestSetNull(t *testing.T) {
	root, err := Parse(`{"a":"b"}`)
	assertParsed(t, root, err)
	err = root.SetNull("c")
	assertNil(t, err)
	assertEqual(t, `{"a":"b","c":null}`, root.Stringify())
}

func TestRemoveIdx(t *testing.T) {
	root, err := Parse(`["a","b","c","d","e"]`)
	assertParsed(t, root, err)

	err = root.RemoveIdx(4)
	assertNil(t, err)
	assertEqual(t, `["a","b","c","d"]`, root.Stringify())

	err = root.RemoveIdx(0)
	assertNil(t, err)
	assertEqual(t, `["b","c","d"]`, root.Stringify())
	cnode := root.Idx(1)
	assertEqual(t, 1, cnode.idx)

	err = root.RemoveIdx(1)
	assertNil(t, err)
	assertEqual(t, `["b","d"]`, root.Stringify())

	err = root.RemoveIdx(0)
	assertNil(t, err)
	assertEqual(t, `["d"]`, root.Stringify())

	err = root.RemoveIdx(0)
	assertNil(t, err)
	assertEqual(t, `[]`, root.Stringify())

	err = root.RemoveIdx(0)
	if !errors.Is(err, ErrInvalidIndex) {
		t.Fatal("expected error ErrInvalidIndex")
	}
	node := NewString("boo")
	err = node.RemoveIdx(0)
	if !errors.Is(err, ErrInvalidNodeForOperation) {
		t.Fatal("expected error ErrInvalidNodeForOperation")
	}
}

func TestRemoveKey(t *testing.T) {
	root, err := Parse(`{"a":0,"b":1,"c":2,"d":3,"e":4}`)
	assertParsed(t, root, err)

	err = root.RemoveKey("e")
	assertNil(t, err)
	assertEqual(t, `{"a":0,"b":1,"c":2,"d":3}`, root.Stringify())

	err = root.RemoveKey("a")
	assertNil(t, err)
	assertEqual(t, `{"b":1,"c":2,"d":3}`, root.Stringify())
	cnode := root.Key("c")
	assertEqual(t, 1, cnode.idx)

	err = root.RemoveKey("c")
	assertNil(t, err)
	assertEqual(t, `{"b":1,"d":3}`, root.Stringify())

	err = root.RemoveKey("b")
	assertNil(t, err)
	assertEqual(t, `{"d":3}`, root.Stringify())

	err = root.RemoveKey("d")
	assertNil(t, err)
	assertEqual(t, `{}`, root.Stringify())

	err = root.RemoveKey("foo")
	if !errors.Is(err, ErrInvalidKey) {
		t.Fatal("expected error ErrInvalidKey")
	}
	node := NewString("boo")
	err = node.RemoveKey("foo")
	if !errors.Is(err, ErrInvalidNodeForOperation) {
		t.Fatal("expected error ErrInvalidNodeForOperation")
	}
}

func TestRemove(t *testing.T) {
	root, err := Parse(`["a","b","c","d","e"]`)
	assertParsed(t, root, err)

	node := root.Idx(2)
	err = node.Remove()
	assertNil(t, err)
	assertNil(t, node.parent)
	assertEqual(t, `["a","b","d","e"]`, root.Stringify())
	err = node.Remove()
	if !errors.Is(err, ErrNoParent) {
		t.Fatal("expected error ErrNoParent")
	}

	root, err = Parse(`{"a":"b","c":"d","e":"f"}`)
	assertParsed(t, root, err)

	node = root.Key("c")
	err = node.Remove()
	assertNil(t, err)
	assertNil(t, node.parent)
	assertEqual(t, `{"a":"b","e":"f"}`, root.Stringify())
}

func TestReplaceIdx(t *testing.T) {
	root, err := Parse(`["a","b","c","d","e"]`)
	assertParsed(t, root, err)
	err = root.ReplaceIdx(100, NewNull())
	if !errors.Is(err, ErrInvalidIndex) {
		t.Fatal("expected error ErrInvalidIndex")
	}
	node := NewNull()
	err = root.ReplaceIdx(2, node)
	assertNil(t, err)
	assertEqual(t, root, node.parent)
	assertEqual(t, 2, node.idx)
	assertEqual(t, `["a","b",null,"d","e"]`, root.Stringify())
}

func TestReplace(t *testing.T) {
	root, err := Parse(`["a","b","c","d","e"]`)
	assertParsed(t, root, err)
	node := root.Idx(1)
	err = node.Replace(NewNull())
	assertNil(t, err)
	assertEqual(t, `["a",null,"c","d","e"]`, root.Stringify())

	root, err = Parse(`{"a":"b","c":"d","e":"f"}`)
	assertParsed(t, root, err)
	node = root.Key("e")
	err = node.Replace(NewNull())
	assertNil(t, err)
	assertEqual(t, `{"a":"b","c":"d","e":null}`, root.Stringify())
}

func TestReplaceByType(t *testing.T) {
	root, err := Parse(`["a","b","c","d","e"]`)
	assertParsed(t, root, err)
	err = root.Idx(0).ReplaceByString("aaa")
	assertNil(t, err)
	err = root.Idx(1).ReplaceByBool(true)
	assertNil(t, err)
	err = root.Idx(2).ReplaceByNumber(33.333)
	assertNil(t, err)
	err = root.Idx(3).ReplaceByInt(10)
	assertNil(t, err)
	err = root.Idx(4).ReplaceByNull()
	assertNil(t, err)
	assertEqual(t, `["aaa",true,33.333,10,null]`, root.Stringify())
}

func TestSortKeys(t *testing.T) {
	root, err := Parse(`{"c":"c","a":"a","b":"b"}`)
	assertParsed(t, root, err)
	err = root.SortKeys()
	assertNil(t, err)
	assertEqual(t, `{"a":"a","b":"b","c":"c"}`, root.Stringify())
}

func TestSortTreeKeys(t *testing.T) {
	root, err := Parse(`[1,2,3,{"c":"c","a":"a","b":"b"}]`)
	assertParsed(t, root, err)
	err = root.SortTreeKeys()
	assertNil(t, err)
	assertEqual(t, `[1,2,3,{"a":"a","b":"b","c":"c"}]`, root.Stringify())
}
