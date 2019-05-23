package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetRoot(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(GetRoot())
}
