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
	undef = &Node{}
)

const maxDeep = 10000

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
	parent   *Node
	idx      int
	key      string
	value    any
	children []*Node
}

type keymap map[string]int

func (n *Node) Type() Type {
	if n == nil || n == undef {
		return Undefined
	}
	switch v := n.value.(type) {
	case Type:
		return v
	case keymap:
		return Object
	case string:
		return String
	case bool:
		return Bool
	case float64:
		return Number
	default:
		panic("node value type is not supported")
	}
}

func (n *Node) IsArray() bool {
	if n == nil {
		return false
	}
	v, ok := n.value.(Type)
	return ok && v == Array
}

func (n *Node) IsObject() bool {
	if n == nil {
		return false
	}
	switch n.value.(type) {
	case keymap:
		return true
	default:
		return false
	}
}

func (n *Node) IsParent() bool {
	if n == nil {
		return false
	}
	switch v := n.value.(type) {
	case keymap:
		return true
	case Type:
		return v == Array
	}
	return false
}

func (n *Node) IsScalar() bool {
	if n == nil {
		return false
	}
	switch v := n.value.(type) {
	case string, bool, float64:
		return true
	case Type:
		return v == Null
	}
	return false
}

func (n *Node) upper() *Node {
	if n == nil || n.parent == nil {
		return nil
	}
	pt := n.parent.Type()
	if pt != Object && pt != Array {
		panic("node parent is scalar!")
	}
	return n.parent
}

func (n *Node) append(node *Node) int {
	if !n.IsArray() {
		panic("attempt to append to wrong node")
	}
	n.children = append(n.children, node)
	return len(n.children) - 1
}

func (n *Node) appendKey(key string, node *Node) (int, error) {
	if !n.IsObject() {
		panic("attempt to append key to wrong node")
	}
	n.children = append(n.children, node)
	idx := len(n.children) - 1
	keymap := n.value.(keymap)
	if _, ok := keymap[key]; ok {
		return 0, errors.Join(ErrDuplicateKey, errors.New("key already exists: "+key))
	}
	keymap[key] = idx
	return idx, nil
}

// Idx returns the child by index or node of type Undefined
// it is safe method for chain access
func (n *Node) Idx(index int) *Node {
	if n == nil || !n.IsArray() || index < 0 || index >= len(n.children) {
		return undef
	}
	if index < 0 || index > len(n.children) {
		return undef
	}
	return n.children[index]
}

// Key returns the child by index or nil
// it is safe method for chain access
func (n *Node) Key(key string) *Node {
	if n == nil || !n.IsObject() {
		return undef
	}
	idx, ok := n.value.(keymap)[key]
	if !ok {
		return undef
	}
	if idx < 0 || idx > len(n.children) {
		return undef
	}
	return n.children[idx]
}

// Level returns the node deep level
func (n *Node) Level() int {
	ret := 0
	node := n
	for {
		if node == nil || node.parent == nil {
			break
		}
		node = node.parent
		ret++
		if ret > maxDeep {
			panic("very deep node or parent loop dependecy detected")
		}
	}
	return ret
}

// Exists indicate if node traversed exists in a tree
func (n *Node) Exists() bool {
	return n != nil && n != undef
}

// SelfIdx returns the index of current node if it is a member of array
func (n *Node) SelfIdx() int {
	if n == nil || n.parent == nil || !n.parent.IsArray() {
		return 0
	}
	return n.idx
}

// SelfKey returns the key of current node if it is a member of object
func (n *Node) SelfKey() string {
	if n == nil || n.parent == nil || !n.parent.IsObject() {
		return ""
	}
	return n.key
}

// IsNull return true if node contains null json value
func (n *Node) IsNull() bool {
	if n == nil {
		return false
	}
	v, ok := n.value.(Type)
	return ok && v == Null
}

// IsString return true if node contains string value
func (n *Node) IsString() bool {
	if n == nil {
		return false
	}
	_, ok := n.value.(string)
	return ok
}

// String return string value or error if value is not of string type
func (n *Node) String() (string, error) {
	if n == nil || n == undef {
		return "", ErrNodeDoesNotExist
	}
	v, ok := n.value.(string)
	if !ok {
		return "", ErrValueIsNotString
	}
	return v, nil
}

// IsBool return true if node contains boolean value
func (n *Node) IsBool() bool {
	if n == nil {
		return false
	}
	_, ok := n.value.(bool)
	return ok
}

// Bool boolean node value or error if value is not of bool type
func (n *Node) Bool() (bool, error) {
	if n == nil || n == undef {
		return false, ErrNodeDoesNotExist
	}
	v, ok := n.value.(bool)
	if !ok {
		return false, ErrValueIsNotBool
	}
	return v, nil
}

// IsNumber return true if node contains numeric value
func (n *Node) IsNumber() bool {
	if n == nil {
		return false
	}
	_, ok := n.value.(float64)
	return ok
}

// Nuumber returns numeric node value or error if value is not numeric
func (n *Node) Number() (float64, error) {
	if n == nil || n == undef {
		return 0, ErrNodeDoesNotExist
	}
	v, ok := n.value.(float64)
	if !ok {
		return 0, ErrValueIsNotNumber
	}
	return v, nil
}

// IsInt return true if node contains numeric value, that can be converted to integer without loss
func (n *Node) IsInt() bool {
	if n == nil {
		return false
	}
	v, ok := n.value.(float64)
	if !ok {
		return false
	}
	return v == float64(int(v))
}

// Int returns numeric node value if the one can be converted to integert without loss
func (n *Node) Int() (int, error) {
	if n == nil || n == undef {
		return 0, ErrNodeDoesNotExist
	}
	fv, ok := n.value.(float64)
	if !ok {
		return 0, ErrValueIsNotNumber
	}
	v := int(fv)
	if fv != float64(v) {
		return 0, ErrValueIsNotInteger
	}
	return v, nil
}
