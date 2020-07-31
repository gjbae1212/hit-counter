package api_handler

type WebSocketMessage struct {
	Payload []byte
}

// GetMessage is a specific method implemented Message interface in websocket.
func (wsm *WebSocketMessage) GetMessage() []byte {
	return wsm.Payload
}
