package xtjson

import "strconv"

type state int

const (
	enter state = iota
	object
	array
)

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
		if parent.kind == Array {
			ret = "[" + strconv.Itoa(node.idx) + "]" + ret
		} else if parent.kind == Object {
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
	if n == nil || n.kind != Object {
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
	if n == nil || n.kind != Object && n.kind != Array {
		return 0
	}
	return len(n.children)
}

// Walk goes over tree and the pass is done when it returns undefined node
func (n *Node) Walk() *Node {
	if n == nil {
		return undef
	}
	if (n.kind == Array || n.kind == Object) && len(n.children) > 0 {
		return n.children[0]
	}
	node := n
	for {
		parent := node.parent
		if parent == nil {
			return undef
		}
		if len(parent.children)-1 > node.idx {
			return parent.children[node.idx+1]
		}
		node = parent
	}
}
