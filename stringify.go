package xtjson

import (
	"strconv"
	"strings"
)

// Format provide settings for json repesentation
type Format struct {
	Indent             int
	SpacesAfterColon   int
	SpacesAfterComma   int
	SpacesAfterBracket int
}

type stringifyState struct {
	nl           string
	indentSize   int
	indent       string
	afterColon   string
	afterComma   string
	afterBracket string
}

// Stringify returns json string representation of node tree
func (n *Node) Stringify(opts ...*Format) string {
	if n.IsScalar() {
		return stringifyScalar(n)
	}
	state := stringifyState{}
	if len(opts) == 0 || opts[0] == nil {
		return stringifyContainer(n, &state)
	}
	opt := opts[0]

	if opt.Indent > 0 {
		state.nl = "\n"
		state.indentSize = opt.Indent
		state.afterComma = ""
		state.afterColon = " "
	} else {
		if opt.SpacesAfterComma > 0 {
			state.afterComma = strings.Repeat(" ", opt.SpacesAfterComma)
		}
		if opt.SpacesAfterBracket > 0 {
			state.afterBracket = strings.Repeat(" ", opt.SpacesAfterBracket)
		}
	}
	return stringifyContainer(n, &state)
}

func stringifyScalar(n *Node) string {
	if n == nil || n == undef {
		return ""
	}
	var ret string
	switch n.Type() {
	case Null:
		ret = "null"
	case String:
		ret = strconv.Quote(n.value.(string))
	case Bool:
		v := n.value.(bool)
		if v {
			ret = "true"
			break
		}
		ret = "false"
	case Number:
		if n.IsInt() {
			v, err := n.Int()
			if err != nil {
				panic(err)
			}
			ret = strconv.Itoa(v)
			break
		}
		v, err := n.Number()
		if err != nil {
			panic(err)
		}
		ret = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return ret
}

func stringifyContainer(n *Node, s *stringifyState) string {
	if n == nil || n == undef {
		return ""
	}
	openBracket := "{"
	closeBracket := "}"
	if n.IsArray() {
		openBracket = "["
		closeBracket = "]"
	}
	ret := openBracket + s.afterBracket + s.nl
	s.indent = s.indent + strings.Repeat(" ", s.indentSize)
	last := len(n.children) - 1
	for idx, node := range n.children {
		ret += s.indent
		if n.IsObject() {
			ret += strconv.Quote(node.key) + ":" + s.afterColon
		}
		switch node.Type() {
		case Null, String, Bool, Number:
			ret += stringifyScalar(node)
		case Array, Object:
			ret += stringifyContainer(node, s)
		}
		if idx < last {
			ret += "," + s.afterComma
		}
		ret += s.nl
	}
	if len(s.indent) >= s.indentSize {
		s.indent = s.indent[0 : len(s.indent)-s.indentSize]
	}
	return ret + s.indent + closeBracket
}
