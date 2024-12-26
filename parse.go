package xtjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

var (
	ErrInvalidJson  = errors.New("invalid json")
	ErrDuplicateKey = errors.New("duplicate key")
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
			return nil, errors.Join(ErrInvalidJson, err)
		}

		switch v := token.(type) {
		case json.Delim:
			switch v {
			case '{':
				node = &Node{value: make(keymap)}
				setParent = true
			case '[':
				node = &Node{value: Array}
				setParent = true
			case '}', ']':
				if parent != nil {
					parent, node = parent.parent, parent
				}
			}
		case string:
			if parent != nil && parent.IsObject() && !keySet {
				key = v
				keySet = true
				continue
			}
			node = &Node{value: v}
		case bool, float64:
			node = &Node{value: v}
		case nil:
			node = &Node{value: Null}
		}

		if node.parent == nil {
			node.parent = parent
			if parent.IsObject() {
				node.key = key
				node.idx, err = parent.appendKey(key, node)
				if err != nil {
					return nil, err
				}
				keySet = false
			} else if parent.IsArray() {
				node.idx = parent.append(node)
			} else if parent != nil {
				return nil, ErrInvalidJson
			}
		}
		if setParent {
			parent = node
			setParent = false
		}
	}
	if node.parent != nil {
		return nil, ErrInvalidJson
	}
	return node, nil
}

// Parse converts json string to tree and returns the top node of the tree
func Parse(s string) (*Node, error) {
	stream := strings.NewReader(s)
	return ParseStream(stream)
}

// ParseBytes converts json bytes to tree and returns the top node of the tree
func ParseBytes(b []byte) (*Node, error) {
	stream := bytes.NewReader(b)
	return ParseStream(stream)
}

// ParseFile converts json file contents to tree and returns the top node of the tree
func ParseFile(name string) (*Node, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseStream(f)
}
