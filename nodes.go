package xtjson

import (
	"errors"
)

var (
	ErrBadQuery = errors.New("bad query format")
)

// Nodes represents list of nodes
type Nodes []*Node

// Search runs the Search method on all nodes and combines the output
func (ns Nodes) Search(matcher NodeMatcher, opt *SearchOptions) (Nodes, error) {
	ret := make(Nodes, 0)
	for _, node := range ns {
		result, err := node.Search(matcher, opt)
		if err != nil {
			return ret, err
		}
		ret = append(ret, result...)
	}
	return ret, nil
}

// Search returns all nodes matched by provided matcher
func (ns Nodes) SearchKey(key string, opt *SearchOptions) (Nodes, error) {
	return ns.Search(&keyMatcher{key}, opt)
}

// Applying path to nodes
func (ns Nodes) Path(path string) Nodes {
	ret := make(Nodes, 0)
	for _, node := range ns {
		result := node.Path(path)
		if !result.Exists() {
			continue
		}
		ret = append(ret, result)
	}
	return ret
}

// StringifyValues returns json representation of values combined with top level array
func (ns Nodes) StringifyValues(opts ...*Format) string {
	root := &Node{
		kind:     Array,
		children: make([]*Node, len(ns)),
	}
	for i, node := range ns {
		if node == nil {
			continue
		}
		nc := node.copy(root)
		nc.idx = i
		nc.key = ""
		root.children[i] = nc
	}
	return root.Stringify(opts...)
}
