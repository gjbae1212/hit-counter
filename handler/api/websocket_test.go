package api_handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebSocketMessage_GetMessage(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		input string
		want  string
	}{"step1": {input: "hi", want: "hi"}}

	for _, v := range tests {
		wsm := &WebSocketMessage{Payload: []byte(v.input)}
		assert.Equal(wsm.Payload, wsm.GetMessage())
	}
}
