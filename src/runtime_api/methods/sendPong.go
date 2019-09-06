package runtime_api

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type PongMessage struct {
	Msg string `json:"msg"`
}

func (wc *WebSocketClient) SendPong() error {
	msg := PongMessage{
		Msg: "pong",
	}
	res, _ := json.Marshal(msg)
	wc.Request <- Request{
		mt:  websocket.PongMessage,
		msg: res,
	}
	return nil
}
