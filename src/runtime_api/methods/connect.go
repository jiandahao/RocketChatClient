package runtime_api

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type ConnectInfo struct {
	Msg     string   `json:"msg"`
	Version string   `json:"version"`
	Support []string `json:"support"`
}

/*
According to official runtime api reference,
Before requesting any method / subscription
you have to send a connect message:
{
    "msg": "connect",
    "version": "1",
    "support": ["1"]
}
*/
func (wc *WebSocketClient) Connect() error {
	connInfo := ConnectInfo{
		Msg:     "connect",
		Version: "1",
		Support: []string{
			"1", "pre2", "pre1",
		},
	}
	connInfoJson, _ := json.Marshal(connInfo)

	// write request to the channel to trigger the request handler
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: connInfoJson,
	}
	return nil
}
