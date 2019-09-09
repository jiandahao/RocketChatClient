package runtime_api

/*
This is the user stream.
Events available:
	- message
	- otr (Off the Record Message)
	- webrtc
	- notification
	- rooms-changed
subscriptions-changed
{
    "msg": "sub",
    "id": "unique-id",
    "name": "stream-notify-user",
    "params":[
        "user-id/event",
        false
    ]
}
*/
import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type StreamUserMessage struct {
	Msg    string        `json:"msg"`
	Id     string        `json:"id"`
	Name   string        `json:"name"`
	Params []interface{} `json:"params"`
}

/*
Params:
	args[0]     userId
	args[1] 	eventName
*/
func (wc *WebSocketClient) SubscribeStreamNotifyUser(userId string, eventName string) error {
	params := []interface{}{
		fmt.Sprintf("%s/%s", userId, eventName),
		false,
	}
	msg := StreamUserMessage{
		Msg:    "sub",
		Id:     uuid.New().String(),
		Name:   "stream-notify-user",
		Params: params,
	}

	res, _ := json.Marshal(msg)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: res,
	}
	return nil
}
