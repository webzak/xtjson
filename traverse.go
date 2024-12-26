package xtjson

import (
	"errors"
	"strconv"
)

type state int

const (
	enter state = iota
	object
	array
)

var (
	ErrNilNode  = errors.New("nil node")
	ErrWalkSkip = errors.New("skip can be done only right after WalkEnter state")
)

// Type contains node type
type WalkState int

const (
	WalkPass WalkState = iota
	WalkEnter
	WalkExit
	WalkDone
)

// Walker provides an api for traversing tree nodes
type Walker struct {
	root      *Node
	deepLimit int
	next      *Node
	nextState WalkState
}

// NewWalker creates Walker instance
// root is a start and end point
// deepLimit 0 means unlimited
func NewWalker(root *Node, deepLimit int) (*Walker, error) {
	if root == nil {
		return nil, ErrNilNode
	}
	if deepLimit > 0 {
		deepLimit += root.Level()
	} else {
		deepLimit = 0
	}
	ret := Walker{
		root:      root,
		deepLimit: deepLimit,
		next:      root,
		nextState: WalkEnter,
	}
	if root.IsScalar() {
		ret.nextState = WalkPass
	}
	return &ret, nil
}

// Next returns next node in walk
func (w *Walker) Next() (*Node, WalkState) {
	if w == nil || w.nextState == WalkDone {
		return nil, WalkDone
	}
	node := w.next
	state := w.nextState

	if state == WalkEnter {
		if len(node.children) == 0 || w.deepLimit > 0 && node.Level() >= w.deepLimit {
			w.next = node
			w.nextState = WalkExit
			return node, state
		}
		w.next = node.children[0]
		w.nextState = WalkPass
		if w.next.IsParent() {
			w.nextState = WalkEnter
		}
		return node, state
	}

	if state == WalkPass || state == WalkExit {
		parent := node.parent
		if parent == nil || node == w.root {
			w.nextState = WalkDone
			return node, state
		}
		if len(parent.children)-1 > node.idx {
			w.next = parent.children[node.idx+1]
			w.nextState = WalkPass
			if w.next.IsParent() {
				w.nextState = WalkEnter
			}
			return node, state
		}
		w.next = parent
		w.nextState = WalkExit
	}
	return node, state
}

// Skip method can be called to bypass going deep to current node children
// This method will have effect if called right after Walker returned WalkEnter state,
// otherwise it returns the error. Even after error the walker is able to continue.
func (w *Walker) Skip() error {
	if w.next == nil || w.next.idx != 0 || w.nextState == WalkExit || w.nextState == WalkDone {
		return ErrWalkSkip
	}
	parent := w.next.upper()
	if parent == nil {
		w.next = nil
		w.nextState = WalkDone
		return nil
	}
	w.next = parent
	w.nextState = WalkExit
	return nil
}

// Path returns the node in the tree referenced by json path
func (n *Node) Path(path string) *Node {
	if n == nil || path == "" {
		return undef
	}
	if path[0] != '$' {
		return undef
	}
	if path == "$" {
		return n
	}
	node := n
	var mode state
	var token []rune
	path = path[1:]

	for _, v := range path {
		if mode == enter {
			switch v {
			case '.':
				mode = object
			case '[':
				mode = array
			default:
				return undef
			}
			token = nil
			continue
		}
		if v != '.' && v != '[' && v != ']' {
			token = append(token, v)
			continue
		}
		if mode == object {
			if v == ']' {
				return undef
			}
			node = node.Key(string(token))
			if node == undef {
				return undef
			}
			token = nil
			if v == '[' {
				mode = array
			}
			continue
		}
		if mode == array {
			if !(v == ']' && len(token) > 0) {
				return undef
			}
			idx, err := strconv.Atoi(string(token))
			if err != nil || idx < 0 {
				return nil
			}
			node = node.Idx(idx)
			if node == undef {
				return undef
			}
			token = nil
			mode = enter
		}
	}
	if len(token) > 0 {
		if mode == array || mode == enter {
			return undef
		}

		if mode == object {
			node = node.Key(string(token))
		}
	}
	return node
}

// SelfPath returs json path of current node
func (n *Node) SelfPath() string {
	if n == nil || n == undef {
		return ""
	}
	ret := ""
	node := n
	for {
		parent := node.parent
		if parent == nil {
			ret = "$" + ret
			break
		}
		if parent.IsArray() {
			ret = "[" + strconv.Itoa(node.idx) + "]" + ret
		} else if parent.IsObject() {
			ret = "." + node.key + ret
		} else {
			panic("parent is neither array nor object")
		}
		node = parent
	}
	return ret
}

// Parent returns the parent of node
func (n *Node) Parent() *Node {
	if n == nil || n.parent == nil {
		return undef
	}
	return n.parent
}

// IsAncestorOf check if receiver is ancestor of node
func (n *Node) IsAncestorOf(node *Node) bool {
	if n == nil || node == nil {
		return false
	}
	for {
		if node.parent == nil {
			break
		}
		if node.parent == n {
			return true
		}
		node = node.parent
	}
	return false
}

// Children returns node children
func (n *Node) Children() []*Node {
	if n == nil {
		return nil
	}
	ret := make([]*Node, len(n.children))
	copy(ret, n.children)
	return ret
}

// ChildrenKeys returns children keys if node is object
func (n *Node) ChildrenKeys() []string {
	if !n.IsObject() {
		return nil
	}
	ret := make([]string, 0, len(n.children))
	for _, child := range n.children {
		ret = append(ret, child.key)
	}
	return ret
}

// ChildrenLength returns the length of node children
func (n *Node) ChildrenLength() int {
	if !n.IsParent() {
		return 0
	}
	return len(n.children)
}

func (n *Node) copy(parent *Node) *Node {
	if n == nil {
		return nil
	}
	node := &Node{
		parent: parent,
		value:  n.value,
		idx:    n.idx,
		key:    n.key,
	}
	if node.IsParent() {
		node.children = make([]*Node, len(n.children))
		for i, c := range n.children {
			node.children[i] = c.copy(node)
		}
	}
	if node.IsObject() {
		kmap := make(keymap)
		for key, idx := range n.value.(keymap) {
			kmap[key] = idx
		}
		node.value = kmap
	}
	return node
}

// Copy creates a tree copy starting from receiver which becomes the root
func (n *Node) Copy() *Node {
	node := n.copy(nil)
	if node == nil {
		return undef
	}
	return node
}
