package xtjson

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// ParseStream converts the bytes stream to tree and returns the top node of the tree
func ParseStream(stream io.Reader) (*Node, error) {
	dec := json.NewDecoder(stream)

	var parent *Node
	var node *Node
	var key string
	var keySet bool
	var setParent bool

	for {
		token, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch v := token.(type) {
		case json.Delim:
			switch v {
			case '{':
				node = &Node{kind: Object}
				setParent = true
			case '[':
				node = &Node{kind: Array}
				setParent = true
			case '}', ']':
				if parent != nil {
					parent, node = parent.parent, parent
				}
			}
		case string:
			if parent != nil && parent.kind == Object && !keySet {
				key = v
				keySet = true
				continue
			}
			node = &Node{kind: String, value: v}
		case bool:
			node = &Node{kind: Bool, value: v}
		case float64:
			node = &Node{kind: Number, value: v}
		case nil:
			node = &Node{kind: Null, value: v}
		}

		if node.parent == nil {
			node.parent = parent
			if parent != nil {
				node.idx = len(parent.children)
				parent.children = append(parent.children, node)

				if parent.kind == Object {
					parent.keys = append(parent.keys, key)
					keySet = false
				}
			}
		}
		if setParent {
			parent = node
			setParent = false
		}
	}
	return node, nil
}

// Parse converts json string to tree and returns the top node of the tree
func Parse(s string) (*Node, error) {
	stream := strings.NewReader(s)
	return ParseStream(stream)
}

// DecodeStream converts json bytes to tree and returns the top node of the tree
func ParseBytes(b []byte) (*Node, error) {
	stream := bytes.NewReader(b)
	return ParseStream(stream)
}
