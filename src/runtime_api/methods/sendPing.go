package runtime_api

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type PingMessage struct {
	Msg string `json:"msg"`
}

func (wc *WebSocketClient) SendPing() error {
	msg := PingMessage{
		Msg: "ping",
	}
	res, _ := json.Marshal(msg)
	//fmt.Println("send a ping to server: ",string(res))

	wc.Request <- Request{
		mt:  websocket.PingMessage,
		msg: res,
	}
	return nil
}
