package runtime_api

import (
	"encoding/json"
	"fmt"
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

var messageId int64 = 1
var sendMessageId int64 = 1

type messageParams struct {
	MsgId  string `json:"_id"` // The message id
	RoomId string `json:"rid"` // the room id for where to send this message
	Msg    string `json:"msg"` // message body (the text of the message itself)
}

func (wc *WebSocketClient) SendMessage(roomId string, message string) error {
	msgParams := []messageParams{
		{
			MsgId:  fmt.Sprintf("sendmsgid:%v", sendMessageId), //uuid.New().String(),
			RoomId: roomId,
			Msg:    message,
		},
	}

	msg := Message{
		Msg:    "method",
		Method: "sendMessage",
		Id:     fmt.Sprintf("msgid:%v", messageId), //uuid.New().String(),
		Params: msgParams,
	}

	messageId = messageId + 1
	sendMessageId = sendMessageId + 1
	res, _ := json.Marshal(msg)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: res,
	}
	return nil
}
