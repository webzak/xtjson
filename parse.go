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

func parse(dec *json.Decoder) (*Node, error) {

	var parent *Node
	var node *Node
	var key string
	var keySet bool
	var setParent bool
	var level int

	for {
		token, err := dec.Token()
		if err != nil {
			return nil, errors.Join(ErrInvalidJson, err)
		}

		switch v := token.(type) {
		case json.Delim:
			switch v {
			case '{':
				node = &Node{value: make(keymap)}
				setParent = true
				level++
			case '[':
				node = &Node{value: Array}
				setParent = true
				level++
			case '}', ']':
				if parent != nil {
					parent, node = parent.parent, parent
				}
				level--
				if level < 0 {
					return nil, ErrInvalidJson
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

		if node != nil && node.parent == nil {
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
		if level <= 0 {
			break
		}
	}
	if node != nil && node.parent != nil {
		return nil, ErrInvalidJson
	}
	return node, nil
}

// Parse converts the bytes stream to tree and returns the top node of the tree
func Parse(stream io.Reader) (*Node, error) {
	decoder := json.NewDecoder(stream)
	node, err := parse(decoder)
	if err != nil {
		return nil, err
	}
	_, err = decoder.Token()
	if err != io.EOF {
		return nil, ErrInvalidJson
	}
	return node, nil
}

// ParseString converts json string to tree and returns the top node of the tree
func ParseString(s string) (*Node, error) {
	stream := strings.NewReader(s)
	return Parse(stream)
}

// ParseBytes converts json bytes to tree and returns the top node of the tree
func ParseBytes(b []byte) (*Node, error) {
	stream := bytes.NewReader(b)
	return Parse(stream)
}

// ParseFile converts json file contents to tree and returns the top node of the tree
func ParseFile(name string) (*Node, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}
