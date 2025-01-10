package xtjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrDirWalk     = errors.New("dir walk error")
	ErrPathMatch   = errors.New("path match error")
	ErrIsNotArray  = errors.New("payload is not json array")
	ErrIsNotObject = errors.New("payload is not json object")
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

// ArrayReader reads one item per read
type ArrayReader struct {
	dec *json.Decoder
	cnt int
}

// NewArrayReader creates new array reader from stream
func NewArrayReader(stream io.Reader) (*ArrayReader, error) {
	reader := ArrayReader{
		dec: json.NewDecoder(stream),
	}
	token, err := reader.dec.Token()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidJson, err.Error())
	}
	vt, ok := token.(json.Delim)
	if !ok || vt != '[' {
		return nil, ErrIsNotArray
	}
	return &reader, nil
}

// Read next array value
func (r *ArrayReader) Read() (*Node, error) {
	if !r.dec.More() {
		return nil, io.EOF
	}
	node, err := parse(r.dec)
	if err != nil {
		return nil, err
	}
	node.idx = r.cnt
	r.cnt++
	return node, err
}

// ObjectReader reads one item per read
type ObjectReader struct {
	dec *json.Decoder
}

// NewObjectReader creates new object reader from stream
func NewObjectReader(stream io.Reader) (*ObjectReader, error) {
	reader := ObjectReader{
		dec: json.NewDecoder(stream),
	}
	token, err := reader.dec.Token()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidJson, err.Error())
	}
	vt, ok := token.(json.Delim)
	if !ok || vt != '{' {
		return nil, ErrIsNotObject
	}
	return &reader, nil
}

// Read next object value
func (r *ObjectReader) Read() (*Node, error) {
	if !r.dec.More() {
		return nil, io.EOF
	}
	token, err := r.dec.Token()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidJson, err.Error())
	}
	key, ok := token.(string)
	if !ok {
		return nil, ErrInvalidJson
	}
	node, err := parse(r.dec)
	if err != nil {
		return nil, err
	}
	node.key = key
	return node, err
}
