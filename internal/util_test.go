package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRoot(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(GetRoot())
}

func TestStringInSlice(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		str  string
		list []string
		ok   bool
	}{
		"false": {str: "d", list: []string{"a", "b", "c"}, ok: false},
		"true":  {str: "b", list: []string{"a", "b", "c"}, ok: true},
	}

	for _, t := range tests {
		assert.Equal(StringInSlice(t.str, t.list), t.ok)
	}
}

func TestParseURL(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		input    string
		schema   string
		domain   string
		port     string
		path     string
		query    string
		fragment string
	}{
		"sample1": {
			input: "http://naver.com/aa/bb?cc=dd&ee=ff#fragment", schema: "http", domain: "naver.com", port: "80",
			path: "/aa/bb", query: "cc=dd&ee=ff", fragment: "fragment",
		},
		"sample2": {
			input: "cc.com:8080/aa/bb", schema: "http", domain: "cc.com", port: "8080", path: "/aa/bb", query: "", fragment: "",
		},
		"sample3": {
			input: "https://naver.com/aa/bb?cc=dd&ee=ff#fragment", schema: "https", domain: "naver.com", port: "443",
			path: "/aa/bb", query: "cc=dd&ee=ff", fragment: "fragment",
		},
	}

	for _, t := range tests {
		schema, domain, port, path, query, fragment, err := ParseURL(t.input)
		assert.NoError(err)
		assert.Equal(t.schema, schema)
		assert.Equal(t.domain, domain)
		assert.Equal(t.port, port)
		assert.Equal(t.path, path)
		assert.Equal(t.query, query)
		assert.Equal(t.fragment, fragment)
	}
}
