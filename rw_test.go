package xtjson

import "testing"

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
