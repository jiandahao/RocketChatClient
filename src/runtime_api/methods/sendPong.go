package runtime_api

import (
	"github.com/gorilla/websocket"
)

type PongMessage struct {
	Msg string `json:"msg"`
}

func (wc *WebSocketClient) SendPong() error {
	// for interacting with rocket chat server,
	// you must response with pong (message type: websocket.TextMessage)
	// it will not work if the message type is websocket.PongMessage
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: []byte("{\"msg\":\"pong\"}"),
	}
	return nil
}
