package runtime_api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/*
{
    "msg": "sub",
    "id": "unique-id",
    "name": "stream-room-messages",
    "params":[
        "room-id",
        false
    ]
}
*/

type StreamRoomMessage struct {
	Msg    string        `json:"msg"`
	Id     string        `json:"id"`
	Name   string        `json:"name"`
	Params []interface{} `json:"params"`
}

func (wc *WebSocketClient) SubscribeStreamRoomMessage(args []string) error {
	roomId := args[0]
	params := []interface{}{
		roomId,
		false,
	}
	msg := StreamRoomMessage{
		Msg:    "sub",
		Id:     uuid.New().String(),
		Name:   "stream-room-messages",
		Params: params,
	}

	subRoomMegJson, _ := json.Marshal(msg)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: subRoomMegJson,
	}
	return nil
}
