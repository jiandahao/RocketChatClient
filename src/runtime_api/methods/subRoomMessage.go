package runtime_api

import (
	"encoding/json"
	"errors"
	"fmt"
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
{
    "msg": "sub",
    "id": "unique-id",
    "name": "stream-notify-room",
    "params":[
        "room-id/event",
        false
    ]
}
*/

type SubscribeStreamMessage struct {
	Msg    string        `json:"msg"`
	Id     string        `json:"id"`
	Name   string        `json:"name"`
	Params []interface{} `json:"params"`
}

/*
Replace event from one in the list Events available:
- deleteMessage
- typing
So far, only these two events are supported
*/
func (wc *WebSocketClient) SubscribeStreamNotifyRoom(roomId string, event string) error {
	if event != "deleteMessage" && event != "typing" {
		return errors.New("Unsupported event type ")
	}
	params := []interface{}{
		fmt.Sprintf("%s/%s", roomId, event),
		false,
	}
	msg := SubscribeStreamMessage{
		Msg:    "sub",
		Id:     uuid.New().String(),
		Name:   "stream-notify-room",
		Params: params,
	}
	subMsgJson, _ := json.Marshal(msg)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: subMsgJson,
	}
	return nil
}

func (wc *WebSocketClient) SubscribeStreamRoomMessage(roomId string) error {
	params := []interface{}{
		roomId,
		false,
	}
	msg := SubscribeStreamMessage{
		Msg:    "sub",
		Id:     uuid.New().String(),
		Name:   "stream-room-messages",
		Params: params,
	}

	subRoomMsgJson, _ := json.Marshal(msg)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: subRoomMsgJson,
	}
	return nil
}
