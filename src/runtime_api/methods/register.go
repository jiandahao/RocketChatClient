package runtime_api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/*
   "msg": "method",
   "method": "registerUser",
   "id":"42",
   "params": [{
       "email": "String",
       "pass": "String",
       "name": "String",
       "secretURL": "String" // Optional
   }]}`
*/
type RegisterInfo struct {
	Msg    string   `json:"msg"`
	Method string   `json:"method"`
	Id     string   `json:"id"`
	Params []Params `json:"params"`
}
type Params struct {
	Email     string `json:"email"`
	Pass      string `json:"pass"`
	Name      string `json:"name"`
	SecretUrl string `json:"secretURL"`
}

func (wc *WebSocketClient) Register(email string, password string, name string) error {
	uid := uuid.New().String()
	params := []Params{
		{
			Email: email,    //"jdachau@163.com",
			Pass:  password, //"123456",
			Name:  name,     //"dachau",
		},
	}
	registerInfo := RegisterInfo{
		Msg:    "method",
		Method: "registerUser",
		Id:     uid,
		Params: params,
	}
	registerJson, _ := json.Marshal(registerInfo)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: registerJson,
	}
	return nil
}
