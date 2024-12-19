package xtjson

import (
	"errors"
)

var (
	ErrInvalidRegexp = errors.New("invalid regexp")
)

// NodeMatcher interface to specify matcher for search operations
type NodeMatcher interface {
	Match(node *Node) bool
}

// NodeMatcherFunc provides a way to use functions as node matcher
type NodeMatcherFunc func(node *Node) bool

// Match the node by certain criteria defined by NodeMatcherFunc implementation
func (f NodeMatcherFunc) Match(node *Node) bool {
	return f(node)
}

// SearchOptions
type SearchOptions struct {
	DeepLimit  int
	SkipNested bool
}

// Search returns all nodes matched by provided matcher
func (n *Node) Search(matcher NodeMatcher, opt *SearchOptions) (Nodes, error) {
	ret := make(Nodes, 0)
	if n == nil {
		return ret, ErrNilNode
	}
	if opt == nil {
		opt = &SearchOptions{
			DeepLimit:  0,
			SkipNested: true,
		}
	}
	walker, err := NewWalker(n, opt.DeepLimit)
	if err != nil {
		return ret, err
	}
	for {
		node, state := walker.Next()
		if state == WalkDone {
			break
		}
		if state == WalkExit {
			continue
		}
		if matcher.Match(node) {
			ret = append(ret, node)
			if state == WalkEnter && opt.SkipNested {
				walker.Skip()
			}
		}
	}
	return ret, nil
}

type keyMatcher struct {
	key string
}

func (km *keyMatcher) Match(node *Node) bool {
	return node.key == km.key
}

// Search returns all nodes matched by provided matcher
func (n *Node) SearchKey(key string, opt *SearchOptions) (Nodes, error) {
	return n.Search(&keyMatcher{key}, opt)
}

// QueryPipe applies queries in the order of spedifying
// Each step can be path like $.foo, deeo keysearch like ...key,
// array [...] or object {...} cildren
func (n *Node) QueryPipe(steps ...string) (Nodes, error) {
	if n == nil {
		return nil, ErrNilNode
	}
	if len(steps) == 0 {
		return make(Nodes, 0), ErrBadQuery
	}
	var err error
	nodes := Nodes{n}
	for _, step := range steps {
		if len(nodes) == 0 {
			break
		}
		if len(step) == 0 || step[0] != '$' {
			return nodes, ErrBadQuery
		}
		if len(step) == 1 {
			continue
		}
		switch {
		case len(step) >= 4 && step[0:4] == "$...":
			key := step[4:]
			if len(key) == 0 {
				return nodes, ErrBadQuery
			}
			nodes, err = nodes.SearchKey(key, nil)
			if err != nil {
				return nodes, err
			}
		case step == "$[...]":
			children := make(Nodes, 0)
			for _, node := range nodes {
				if node.IsArray() {
					children = append(children, node.Children()...)
				}
			}
			nodes = children

		case step == "${...}":
			children := make(Nodes, 0)
			for _, node := range nodes {
				if node.IsObject() {
					children = append(children, node.Children()...)
				}
			}
			nodes = children
		default:
			nodes = nodes.Path(step)
		}
	}
	return nodes, nil
}

func parseQuery(path string) ([]string, error) {
	if len(path) == 0 || path[0] != '$' {
		return nil, ErrBadQuery
	}
	steps := []string{}
	step := ""
	buf := ""

	for _, v := range path {
		switch buf {
		case "":
			switch v {
			case '.', '[', '{':
				buf += string(v)
				if len(step) > 3 && step[0:3] == "..." {
					steps = append(steps, step)
					step = ""
				}
			default:
				step += string(v)
			}

		case ".", "[", "{":
			switch v {
			case '.':
				buf += string(v)
			case '}', ']':
				return nil, ErrBadQuery
			default:
				step += buf + string(v)
				buf = ""
			}

		case "..", "[.", "{.":
			switch v {
			case '.':
				buf += "."
			default:
				return nil, ErrBadQuery
			}

		case "...":
			switch v {
			case '.', '[', '{', ']', '}':
				return nil, ErrBadQuery
			default:
				steps = append(steps, step)
				step = buf + string(v)
				buf = ""
			}

		case "[..", "{..":
			switch v {
			case '.':
				buf += "."
			default:
				return nil, ErrBadQuery
			}

		case "[...":
			switch v {
			case ']':
				steps = append(steps, step, "[...]")
				step = ""
				buf = ""
			default:
				return nil, ErrBadQuery
			}

		case "{...":
			switch v {
			case '}':
				steps = append(steps, step, "{...}")
				step = ""
				buf = ""
			default:
				return nil, ErrBadQuery
			}
		default:
			return nil, ErrBadQuery
		}
	}
	if len(step) > 0 {
		steps = append(steps, step)
	}
	for i := range steps {
		step := steps[i]
		if len(steps) == 0 || step[0] != '$' {
			steps[i] = "$" + step
		}
	}
	return steps, nil
}

// Query extracts nodes using the combination of Path syntax with extensions
// deeo keysearch ...key, array [...] or object {...} cildren
// use QueryPipe instead when keys can contain dots or brackets
func (n *Node) Query(path string) (Nodes, error) {
	steps, err := parseQuery(path)
	if err != nil {
		return nil, err
	}
	return n.QueryPipe(steps...)
}
