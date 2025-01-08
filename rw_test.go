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
	name, node, err := reader.NamedRead()
	assertNil(t, err)
	assertEqual(t, "fixtures/dir/one.json", name)
	assertEqual(t, "[1,2,3]", node.Stringify())
	name, node, err = reader.NamedRead()
	assertNil(t, err)
	assertEqual(t, "fixtures/dir/three.json", name)
	assertEqual(t, `{"value":3}`, node.Stringify())
	name, node, err = reader.NamedRead()
	assertNil(t, err)
	assertEqual(t, "fixtures/dir/two.json", name)
	assertEqual(t, `{"two":true}`, node.Stringify())

	reader, err = NewDirReader("./fixtures/dir", "*.json")
	assertNil(t, err)
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, "[1,2,3]", node.Stringify())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, `{"value":3}`, node.Stringify())
	node, err = reader.Read()
	assertNil(t, err)
	assertEqual(t, `{"two":true}`, node.Stringify())
}
