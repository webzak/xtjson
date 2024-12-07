package xtjson

import (
	"testing"
)

func TestIdx(t *testing.T) {
	node, err := Parse(`["aa", "bb", "cc"]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, node.children[0], node.Idx(0))
	assertEqual(t, undef, node.Idx(20))
	assertEqual(t, undef, node.Idx(-20))
	node = nil
	assertEqual(t, undef, node.Idx(0))
	node = undef
	assertEqual(t, undef, node.Idx(0))
}

func TestKey(t *testing.T) {
	node, err := Parse(`{"ka":"va", "kb":"vb", "kc":"vc"}`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, node.children[1], node.Key("kb"))
	assertEqual(t, undef, node.Key("foo"))
	node = nil
	assertEqual(t, undef, node.Key("foo"))
	node = undef
	assertEqual(t, undef, node.Key("foo"))
}

func TestRawValue(t *testing.T) {
	node, err := Parse(`"foo"`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, "foo", node.RawValue())
	node = nil
	assertEqual(t, nil, node.RawValue())
	node = undef
	assertEqual(t, nil, node.RawValue())
}

func TestExists(t *testing.T) {
	node, err := Parse(`{"key":"foo"}`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, true, node.Key("key").Exists())
	assertEqual(t, false, node.Key("wrong").Exists())
	assertEqual(t, false, node.Idx(0).Exists())
	node = nil
	assertEqual(t, false, node.Key("key").Exists())
	node = undef
	assertEqual(t, false, node.Key("key").Exists())
}

func TestIsParent(t *testing.T) {
	node, err := Parse(`[123]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, true, node.IsParent())
	assertEqual(t, false, node.Idx(0).IsParent())
	node = nil
	assertEqual(t, false, node.IsParent())
	node = undef
	assertEqual(t, false, node.IsParent())
}

func TestIsScalar(t *testing.T) {
	node, err := Parse(`[123]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, false, node.IsScalar())
	assertEqual(t, true, node.Idx(0).IsScalar())
	node = nil
	assertEqual(t, false, node.IsScalar())
	node = undef
	assertEqual(t, false, node.IsScalar())
}

func TestType(t *testing.T) {
	node, err := Parse(`[123]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, Array, node.Type())
	assertEqual(t, Number, node.Idx(0).Type())
	node = nil
	assertEqual(t, Undefined, node.Type())
	node = undef
	assertEqual(t, Undefined, node.Type())
}

func TestSelfIdx(t *testing.T) {
	node, err := Parse(`[1,2,3]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, 0, node.SelfIdx())
	assertEqual(t, 0, node.Idx(0).SelfIdx())
	assertEqual(t, 1, node.Idx(1).SelfIdx())
	assertEqual(t, 2, node.Idx(2).SelfIdx())
	node = nil
	assertEqual(t, 0, node.SelfIdx())
	node = undef
	assertEqual(t, 0, node.SelfIdx())
}

func TestSelfKey(t *testing.T) {
	node, err := Parse(`{"foo": 123}`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, "", node.SelfKey())
	assertEqual(t, "foo", node.Key("foo").SelfKey())
	node = nil
	assertEqual(t, "", node.SelfKey())
	node = undef
	assertEqual(t, "", node.SelfKey())
}

func TestIsNull(t *testing.T) {
	node, err := Parse(`[123, null]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, false, node.Idx(0).IsNull())
	assertEqual(t, true, node.Idx(1).IsNull())
	node = nil
	assertEqual(t, false, node.IsNull())
	node = undef
	assertEqual(t, false, node.IsNull())
}

func TestIsString(t *testing.T) {
	node, err := Parse(`["abc"]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, false, node.IsString())
	assertEqual(t, true, node.Idx(0).IsString())
	node = nil
	assertEqual(t, false, node.IsString())
	node = undef
	assertEqual(t, false, node.IsString())
}

func TestString(t *testing.T) {
	node, err := Parse(`["abc"]`)
	assertNil(t, err)
	assertNotNil(t, node)
	v, err := node.String()
	assertEqual(t, ErrValueIsNotString, err)
	assertEqual(t, "", v)
	v, err = node.Idx(0).String()
	assertNil(t, err)
	assertEqual(t, "abc", v)
	node = nil
	v, err = node.String()
	assertEqual(t, ErrValueDoesNotExist, err)
	assertEqual(t, "", v)
	node = undef
	v, err = node.String()
	assertEqual(t, ErrValueDoesNotExist, err)
	assertEqual(t, "", v)
}

func TestIsBool(t *testing.T) {
	node, err := Parse(`[true]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, false, node.IsBool())
	assertEqual(t, true, node.Idx(0).IsBool())
	node = nil
	assertEqual(t, false, node.IsBool())
	node = undef
	assertEqual(t, false, node.IsBool())
}

func TestBool(t *testing.T) {
	node, err := Parse(`[true]`)
	assertNil(t, err)
	assertNotNil(t, node)
	v, err := node.Bool()
	assertEqual(t, ErrValueIsNotBool, err)
	assertEqual(t, false, v)
	v, err = node.Idx(0).Bool()
	assertNil(t, err)
	assertEqual(t, true, v)
	node = nil
	v, err = node.Bool()
	assertEqual(t, ErrValueDoesNotExist, err)
	assertEqual(t, false, v)
	node = undef
	v, err = node.Bool()
	assertEqual(t, ErrValueDoesNotExist, err)
	assertEqual(t, false, v)
}

func TestIsNumber(t *testing.T) {
	node, err := Parse(`[12.4, 100]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, false, node.IsNumber())
	assertEqual(t, true, node.Idx(0).IsNumber())
	assertEqual(t, true, node.Idx(1).IsNumber())
	node = nil
	assertEqual(t, false, node.IsNumber())
	node = undef
	assertEqual(t, false, node.IsNumber())
}

func TestIsInt(t *testing.T) {
	node, err := Parse(`[12.4, 100, 100.0]`)
	assertNil(t, err)
	assertNotNil(t, node)
	assertEqual(t, false, node.IsInt())
	assertEqual(t, false, node.Idx(0).IsInt())
	assertEqual(t, true, node.Idx(1).IsInt())
	assertEqual(t, true, node.Idx(2).IsInt())
	node = nil
	assertEqual(t, false, node.IsInt())
	node = undef
	assertEqual(t, false, node.IsInt())
}

func TestNumber(t *testing.T) {
	node, err := Parse(`[12.4, 12]`)
	assertNil(t, err)
	assertNotNil(t, node)
	v, err := node.Number()
	assertEqual(t, ErrValueIsNotNumber, err)
	assertEqual(t, float64(0), v)
	v, err = node.Idx(0).Number()
	assertNil(t, err)
	assertEqual(t, 12.4, v)

	v, err = node.Idx(1).Number()
	assertNil(t, err)
	assertEqual(t, float64(12), v)

	node = nil
	v, err = node.Number()
	assertEqual(t, ErrValueDoesNotExist, err)
	assertEqual(t, float64(0), v)
	node = undef
	v, err = node.Number()
	assertEqual(t, ErrValueDoesNotExist, err)
	assertEqual(t, float64(0), v)
}

func TestInt(t *testing.T) {
	node, err := Parse(`[12.4, 12]`)
	assertNil(t, err)
	assertNotNil(t, node)
	v, err := node.Int()
	assertEqual(t, ErrValueIsNotInteger, err)
	assertEqual(t, 0, v)
	v, err = node.Idx(0).Int()
	assertEqual(t, ErrValueIsNotInteger, err)
	assertEqual(t, 0, v)

	v, err = node.Idx(1).Int()
	assertNil(t, err)
	assertEqual(t, 12, v)

	node = nil
	v, err = node.Int()
	assertEqual(t, ErrValueDoesNotExist, err)
	assertEqual(t, 0, v)
	node = undef
	v, err = node.Int()
	assertEqual(t, ErrValueDoesNotExist, err)
	assertEqual(t, 0, v)
}
