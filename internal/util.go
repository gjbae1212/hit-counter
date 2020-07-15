package internal

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"runtime"

	"github.com/goware/urlx"
)

var (
	root string
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("No caller information")
	}
	root = filepath.Join(path.Dir(filename), "../")
}

// GetRoot returns root path.
func GetRoot() string {
	return root
}

// StringInSlice checks a string element is included to string array.
func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// ParseURL parses url.
func ParseURL(s string) (schema, host, port, path, query, fragment string, err error) {
	if s == "" {
		err = fmt.Errorf("[err] ParseURL %w", ErrorEmptyParams)
	}

	url, suberr := urlx.Parse(s)
	if suberr != nil {
		err = suberr
		return
	}

	schema = url.Scheme

	host, port, err = urlx.SplitHostPort(url)
	if err != nil {
		return
	}
	if schema == "http" && port == "" {
		port = "80"
	} else if schema == "https" && port == "" {
		port = "443"
	}

	path = url.Path
	query = url.RawQuery
	fragment = url.Fragment
	return
}
