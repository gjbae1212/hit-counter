package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRoot(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(GetRoot())
}
