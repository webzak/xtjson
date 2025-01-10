package xtjson

import (
	"io"
	"os"
	"testing"
)

func TestFindFiles(t *testing.T) {
	files, err := FindFiles("./fixtures", "*.json")
	assertNil(t, err)
	assertEqual(t, 6, len(files))

	files, err = FindFiles("./fixtures", "*ne.json")
	assertNil(t, err)
	assertEqual(t, 1, len(files))
}

func TestDirReader(t *testing.T) {
	reader, err := NewDirReader("./fixtures/dir", "*.json")
	assertNil(t, err)
	node, err := reader.Read()
	assertNil(t, err)
	assertEqual(t, "one.json", node.SelfKey())
	assertEqual(t, "[1,2,3]", node.Stringify())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, "three.json", node.SelfKey())
	assertEqual(t, `{"value":3}`, node.Stringify())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, "two.json", node.SelfKey())
	assertEqual(t, `{"two":true}`, node.Stringify())
}

func TestArrayReader(t *testing.T) {
	f, err := os.Open("./fixtures/array.json")
	assertNil(t, err)
	defer f.Close()
	reader, err := NewArrayReader(f)
	assertNil(t, err)
	node, err := reader.Read()
	assertNil(t, err)
	assertEqual(t, `{"one":"two"}`, node.Stringify())
	assertEqual(t, 0, node.SelfIdx())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, `[1,2,3]`, node.Stringify())
	assertEqual(t, 1, node.SelfIdx())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, `"str"`, node.Stringify())
	assertEqual(t, 2, node.SelfIdx())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, "null", node.Stringify())
	assertEqual(t, 3, node.SelfIdx())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, `{"x":[{"a":"b"},{"c":"d"}]}`, node.Stringify())
	assertEqual(t, 4, node.SelfIdx())
	node, err = reader.Read()
	assertNil(t, node)
	assertEqual(t, err, io.EOF)
}

func TestObjectReader(t *testing.T) {
	f, err := os.Open("./fixtures/object.json")
	assertNil(t, err)
	defer f.Close()
	reader, err := NewObjectReader(f)
	assertNil(t, err)
	node, err := reader.Read()
	assertNil(t, err)
	assertEqual(t, `{"one":"two"}`, node.Stringify())
	assertEqual(t, "first", node.SelfKey())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, `[1,2,3]`, node.Stringify())
	assertEqual(t, "second", node.SelfKey())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, `"str"`, node.Stringify())
	assertEqual(t, "third", node.SelfKey())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, "true", node.Stringify())
	assertEqual(t, "fourth", node.SelfKey())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, `{"x":[{"a":"b"},{"c":"d"}]}`, node.Stringify())
	assertEqual(t, "fifth", node.SelfKey())
	node, err = reader.Read()
	assertNil(t, node)
	assertEqual(t, err, io.EOF)
}
