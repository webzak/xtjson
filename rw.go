package xtjson

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrDirWalk   = errors.New("dir walk error")
	ErrPathMatch = errors.New("path match error")
)

// Reader interface
type Reader interface {
	Read() (*Node, error)
}

// FindFiles recursively searches for the files in specified root dir
func FindFiles(path, pattern string) ([]string, error) {

	var ret []string

	err := filepath.Walk(path, func(fpath string, stat os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w: %s", ErrDirWalk, fpath)
		}
		if stat.IsDir() {
			return nil
		}
		matched, err := filepath.Match(pattern, filepath.Base(fpath))
		if err != nil {
			return fmt.Errorf("%w: %s", ErrPathMatch, fpath)
		}
		if matched {
			ret = append(ret, fpath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDirWalk, path)
	}
	return ret, nil
}

// DirReader provides Reader and NamedReader interfaces to read and parse files in directory
type DirReader struct {
	basePath string
	files    []string
	next     int
}

// NewDirReader creates new directory reader
func NewDirReader(path string, pattern string) (*DirReader, error) {
	var dr DirReader
	var err error
	switch {
	case len(path) >= 2 && path[0:2] == "./":
		dr.basePath = path[2:] + "/"
	case len(path) == 1 && path[0] == '.':
		dr.basePath = ""
	default:
		dr.basePath = path
	}
	dr.files, err = FindFiles(path, pattern)
	if err != nil {
		return nil, err
	}
	return &dr, nil
}

// Read gets next pased tree
func (dr *DirReader) Read() (*Node, error) {
	if dr.next >= len(dr.files) {
		return nil, io.EOF
	}
	name := dr.files[dr.next]
	dr.next++
	node, err := ParseFile(name)
	if dr.basePath != "" {
		name, _ = strings.CutPrefix(name, dr.basePath)
	}
	node.key = name
	return node, err
}
