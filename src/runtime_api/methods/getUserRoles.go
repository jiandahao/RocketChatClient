package runtime_api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/*
   {
       "msg": "method",
       "method": "getUserRoles",
       "id": "42",
       "params": []
   }
*/

type UserRole struct {
	Msg    string   `json:"msg"`
	Method string   `json:"method"`
	Id     string   `json:"id"`
	Params []string `json:"params"`
}

func (wc *WebSocketClient) GetUserRoles() error {
	uid := uuid.New().String()
	msg := UserRole{
		Msg:    "method",
		Method: "getUserRoles",
		Id:     uid,
		Params: []string{},
	}

	res, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: res,
	}
	return nil
}
