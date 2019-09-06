package runtime_api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/*
{
    "msg": "method",
    "method": "sendMessage",
    "id": "42",
    "params": [
        {
            "_id": "message-id",
            "rid": "room-id",
            "msg": "Hello World!"
        }
    ]
}
*/

type Message struct {
	Msg    string          `json:"msg"`
	Method string          `json:"method"`
	Id     string          `json:"id"`
	Params []messageParams `json:"params"`
}

type messageParams struct {
	MsgId  string `json:"_id"` // The message id
	RoomId string `json:"rid"` // the room id for where to send this message
	Msg    string `json:"msg"` // message body (the text of the message itself)
}

func (wc *WebSocketClient) SendMessage(roomId string, message string) error {
	msgParams := []messageParams{
		{
			MsgId:  uuid.New().String(),
			RoomId: roomId,
			Msg:    message,
		},
	}

	msg := Message{
		Msg:    "method",
		Method: "sendMessage",
		Id:     uuid.New().String(),
		Params: msgParams,
	}

	res, _ := json.Marshal(msg)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: res,
	}
	return nil
}
