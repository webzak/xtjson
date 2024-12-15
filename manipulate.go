package xtjson

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrInvalidNodeForOperation = errors.New("invalid node for operation")
	ErrInvalidIndex            = errors.New("invalid index")
	ErrInvalidKey              = errors.New("invalid key")
	ErrNoParent                = errors.New("node has no parent")
	ErrNodeHasParent           = errors.New("node has parent")
)

// NewArray creates new array node
func NewArray() *Node {
	return &Node{kind: Array}
}

// NewObject creates new object node
func NewObject() *Node {
	return &Node{kind: Object, keymap: make(map[string]int)}
}

// NewString creates new string node
func NewString(value string) *Node {
	return &Node{kind: String, value: value}
}

// NewBool creates new bool mode
func NewBool(value bool) *Node {
	return &Node{kind: Bool, value: value}
}

// NewNull creates new null node
func NewNull() *Node {
	return &Node{kind: Null}
}

// NewNumber creates new number node
func NewNumber(value float64) *Node {
	return &Node{kind: Number, value: value}
}

// NewInt creates new number mode, it internally converts it to default json float64 type
func NewInt(value int) *Node {
	return &Node{kind: Number, value: float64(value)}
}

// Append adds node to receiver children, error returned when receiver is not array node
func (n *Node) Append(node *Node) error {
	if n.kind != Array {
		return fmt.Errorf("%w %s", ErrInvalidNodeForOperation, "append")
	}
	if node.parent != nil {
		return fmt.Errorf("%w %s", ErrNodeHasParent, "attemt to append node linked to another parent")
	}
	node.parent = n
	node.idx = len(n.children)
	n.children = append(n.children, node)
	return nil
}

// AppendString adds string node to receiver children, error returned when receiver is not array node
func (n *Node) AppendString(value string) error {
	node := NewString(value)
	return n.Append(node)
}

// AppendBool adds boolean node to receiver children, error returned when receiver is not array node
func (n *Node) AppendBool(value bool) error {
	node := NewBool(value)
	return n.Append(node)
}

// AppendNumber adds numeric node to receiver children, error returned when receiver is not array node
func (n *Node) AppendNumber(value float64) error {
	node := NewNumber(value)
	return n.Append(node)
}

// AppendInt adds numeric integer node to receiver children, error returned when receiver is not array node
func (n *Node) AppendInt(value int) error {
	node := NewInt(value)
	return n.Append(node)
}

// AppendNull adds null node to receiver children, error returned when receiver is not array node
func (n *Node) AppendNull() error {
	node := NewNull()
	return n.Append(node)
}

// Set sets node associated with key as a receiver property, error is returned if receiver is nots object
func (n *Node) Set(key string, node *Node) error {
	if n.kind != Object {
		return fmt.Errorf("%w %s", ErrInvalidNodeForOperation, "set")
	}
	if node.parent != nil {
		return fmt.Errorf("%w %s", ErrNodeHasParent, "attemt to set property node linked to another parent")
	}
	node.key = key
	node.parent = n
	idx, ok := n.keymap[key]
	if ok {
		node.idx = idx
		if err := n.replaceIdx(idx, node); err != nil {
			panic(err)
		}
		return nil
	}
	node.idx = len(n.children)
	n.children = append(n.children, node)
	n.keymap[key] = node.idx
	return nil
}

// SetString sets string property, error returned when receiver is not array node
func (n *Node) SetString(key string, value string) error {
	node := NewString(value)
	return n.Set(key, node)
}

// SetBool sets bool property, error returned when receiver is not array node
func (n *Node) SetBool(key string, value bool) error {
	node := NewBool(value)
	return n.Set(key, node)
}

// SetNumber sets numeric property, error returned when receiver is not array node
func (n *Node) SetNumber(key string, value float64) error {
	node := NewNumber(value)
	return n.Set(key, node)
}

// SetInt sets numeric integer property, error returned when receiver is not array node
func (n *Node) SetInt(key string, value int) error {
	node := NewInt(value)
	return n.Set(key, node)
}

// SetNull sets null property, error returned when receiver is not array node
func (n *Node) SetNull(key string) error {
	node := NewNull()
	return n.Set(key, node)
}

// RemoveIdx removes the child from array node
func (n *Node) RemoveIdx(idx int) error {
	if n.kind != Array {
		return fmt.Errorf("%w %s", ErrInvalidNodeForOperation, "remove index")
	}
	lc := len(n.children) - 1
	if idx > lc {
		return fmt.Errorf("%w %d", ErrInvalidIndex, idx)
	}
	for i := idx; i < lc; i++ {
		node := n.children[i+1]
		node.idx = i
		n.children[i] = n.children[i+1]
	}
	n.children = n.children[:lc]
	return nil
}

// RemoveKey removes the child from object node
func (n *Node) RemoveKey(key string) error {
	if n.kind != Object {
		return fmt.Errorf("%w %s", ErrInvalidNodeForOperation, "remove key")
	}
	idx, ok := n.keymap[key]
	if !ok {
		return fmt.Errorf("%w %s", ErrInvalidKey, key)
	}
	lc := len(n.children) - 1
	if idx > lc {
		panic(fmt.Errorf("%w %d", ErrInvalidIndex, idx))
	}
	for i := idx; i < lc; i++ {
		node := n.children[i+1]
		node.idx = i
		n.children[i] = n.children[i+1]
	}
	n.children = n.children[:lc]
	delete(n.keymap, key)
	for key, i := range n.keymap {
		if i > idx {
			n.keymap[key] = i - 1
		}
	}
	return nil
}

// Remove unlinks the node from parent
func (n *Node) Remove() error {
	parent := n.parent
	if parent == nil {
		return ErrNoParent
	}
	if parent.IsScalar() {
		panic("parent is scalar")
	}
	n.parent = nil
	if parent.kind == Array {
		return parent.RemoveIdx(n.idx)
	}
	return parent.RemoveKey(n.key)
}

func (n *Node) replaceIdx(idx int, node *Node) error {
	lc := len(n.children)
	if idx >= lc {
		return fmt.Errorf("%w %d", ErrInvalidIndex, idx)
	}
	node.idx = idx
	n.children[idx] = node
	return nil
}

// ReplaceIdx replaces the child of array node
func (n *Node) ReplaceIdx(idx int, node *Node) error {
	if n.kind != Array {
		return fmt.Errorf("%w %s", ErrInvalidNodeForOperation, "replace index")
	}
	node.parent = n
	return n.replaceIdx(idx, node)
}

// Replace replaces the node
func (n *Node) Replace(node *Node) error {
	parent := n.parent
	if parent == nil {
		return ErrNoParent
	}
	if parent.IsScalar() {
		panic("parent is scalar")
	}
	n.parent = nil
	if parent.kind == Array {
		return parent.ReplaceIdx(n.idx, node)
	}
	return parent.Set(n.key, node)
}

// ReplaceByString replaces receiver with string node
func (n *Node) ReplaceByString(value string) error {
	return n.Replace(NewString(value))
}

// ReplaceByBool replaces receiver with bool node
func (n *Node) ReplaceByBool(value bool) error {
	return n.Replace(NewBool(value))
}

// ReplaceByNumber replaces receiver with numeric node
func (n *Node) ReplaceByNumber(value float64) error {
	return n.Replace(NewNumber(value))
}

// ReplaceByInt replaces receiver with numeric int node
func (n *Node) ReplaceByInt(value int) error {
	return n.Replace(NewInt(value))
}

// ReplaceByNull replaces receiver with null node
func (n *Node) ReplaceByNull() error {
	return n.Replace(NewNull())
}

// SortKeys sorts keys alphabetically reordering children nodes
func (n *Node) SortKeys() error {
	if n.kind != Object {
		return ErrInvalidNodeForOperation
	}
	klen := len(n.keymap)
	if klen == 0 {
		return nil
	}

	keys := make([]string, 0, klen)
	for k := range n.keymap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	children := make([]*Node, klen)
	for i, key := range keys {
		idx := n.keymap[key]
		children[i] = n.children[idx]
		n.keymap[key] = i
	}
	n.children = children
	return nil
}

// SortTreeKeys sorts keys alphabetically in current and all children nodes
func (n *Node) SortTreeKeys() error {
	walker, err := NewWalker(n, 0)
	if err != nil {
		return err
	}
	for {
		node, state := walker.Next()
		if state == WalkDone {
			break
		}
		if node.kind != Object || state == WalkExit {
			continue
		}
		err = node.SortKeys()
		if err != nil {
			return err
		}
	}
	return nil
}
