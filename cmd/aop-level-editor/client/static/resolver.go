package static

import (
	"encoding/base64"
	"errors"
	"fmt"
)

//go:generate web-compiler -root=../src -output=./mapping.go

// Resolve the path of a statically compiled file to its contents.
func Resolve(path string) ([]byte, error) {
	if path == "/" {
		// Default '/' to index.html
		return Resolve("/index.html")
	}
	val, ok := staticMapping[path]
	if !ok {
		return nil, errors.New("not found")
	}
	d, err := base64.StdEncoding.DecodeString(string(val))
	if err != nil {
		return nil, fmt.Errorf("base 64 decode: %s", err)
	}
	return []byte(d), nil
}
