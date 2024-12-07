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
		idx := node.idx
		parent := node.parent
		if parent == nil {
			ret = "$" + ret
			break
		}
		if parent.kind == Array {
			ret = "[" + strconv.Itoa(idx) + "]" + ret
		} else if parent.kind == Object {
			ret = "." + parent.keys[idx] + ret
		} else {
			panic("parent is neither array nor object")
		}
		node = parent
	}
	return ret
}
