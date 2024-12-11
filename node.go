package xtjson

import "errors"

var (
	ErrNodeDoesNotExist  = errors.New("node does not exist")
	ErrValueIsNotString  = errors.New("value is not string")
	ErrValueIsNotBool    = errors.New("value is not boolen")
	ErrValueIsNotNumber  = errors.New("value is not number")
	ErrValueIsNotInteger = errors.New("value is not integer")
)

var (
	undef = &Node{kind: Undefined}
)

// Type contains node type
type Type int

const (
	Null Type = iota
	Bool
	Number
	String
	Array
	Object
	Undefined
)

// Node structure represents the element of parsed json tree
type Node struct {
	kind     Type
	parent   *Node
	value    any
	idx      int
	key      string
	keymap   map[string]int
	children []*Node
}

// Idx returns the child by index or node of type Undefined
// it is safe method for chain access
func (n *Node) Idx(index int) *Node {
	if n == nil || n.kind != Array || index < 0 || index >= len(n.children) {
		return undef
	}
	return n.children[index]
}

// Key returns the child by index or nil
// it is safe method for chain access
func (n *Node) Key(key string) *Node {
	if n == nil {
		return undef
	}
	idx, ok := n.keymap[key]
	if !ok {
		return undef
	}
	return n.children[idx]
}

// RawValue returns raw node value, for Array or Object type it return nil
func (n *Node) RawValue() any {
	if n == nil {
		return nil
	}
	return n.value
}

// Exists indicate if node traversed exists in a tree
func (n *Node) Exists() bool {
	return n != nil && n.kind != Undefined
}

// IsParent returns true if node is Array or Object
func (n *Node) IsParent() bool {
	if n == nil {
		return false
	}
	return n.kind == Array || n.kind == Object
}

// IsScalar returns true when node type is not Array or Object or Undefined
func (n *Node) IsScalar() bool {
	if !n.Exists() {
		return false
	}
	return !(n.kind == Array || n.kind == Object)
}

// Type returns node type
func (n *Node) Type() Type {
	if n == nil {
		return Undefined
	}
	return n.kind
}

// SelfIdx returns the index of current node if it is a member of array
func (n *Node) SelfIdx() int {
	if n == nil || n.parent == nil || n.parent.kind != Array {
		return 0
	}
	return n.idx
}

// SelfKey returns the key of current node if it is a member of object
func (n *Node) SelfKey() string {
	if n == nil || n.parent == nil || n.parent.kind != Object {
		return ""
	}
	return n.key
}

// IsNull return true if node contains null json value
func (n *Node) IsNull() bool {
	return n.Exists() && n.kind == Null
}

// IsString return true if node contains string value
func (n *Node) IsString() bool {
	return n.Exists() && n.kind == String
}

// String return string value or error if value is not of string type
func (n *Node) String() (string, error) {
	if !n.Exists() {
		return "", ErrNodeDoesNotExist
	}
	if n.kind != String {
		return "", ErrValueIsNotString
	}
	return n.value.(string), nil
}

// IsBool return true if node contains boolean value
func (n *Node) IsBool() bool {
	return n.Exists() && n.kind == Bool
}

// Bool boolean node value or error if value is not of bool type
func (n *Node) Bool() (bool, error) {
	if !n.Exists() {
		return false, ErrNodeDoesNotExist
	}
	if n.kind != Bool {
		return false, ErrValueIsNotBool
	}
	return n.value.(bool), nil
}

// IsNumber return true if node contains numeric value
func (n *Node) IsNumber() bool {
	return n.Exists() && n.kind == Number
}

// Nuumber returns numeric node value or error if value is not numeric
func (n *Node) Number() (float64, error) {
	if !n.Exists() {
		return 0, ErrNodeDoesNotExist
	}
	if n.kind != Number {
		return 0, ErrValueIsNotNumber
	}
	return n.value.(float64), nil
}

// IsInt return true if node contains numeric value, that can be converted to integer without loss
func (n *Node) IsInt() bool {
	if !n.Exists() || n.kind != Number {
		return false
	}
	v := n.value.(float64)
	return v == float64(int(v))
}

// Nuumber returns numeric node value or error if value is not numeric
func (n *Node) Int() (int, error) {
	if !n.Exists() {
		return 0, ErrNodeDoesNotExist
	}
	if n.kind != Number {
		return 0, ErrValueIsNotInteger
	}
	fv := n.value.(float64)
	v := int(fv)
	if fv != float64(v) {
		return 0, ErrValueIsNotInteger
	}
	return int(v), nil
}
