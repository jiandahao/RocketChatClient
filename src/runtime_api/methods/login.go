package runtime_api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type loginInfo struct {
	Msg         string        `json:"msg"`
	Method      string        `json:"method"`
	LoginParams []LoginParams `json:"params"`
	Id          string        `json:"id"`
}
type LoginParams struct {
	User     User     `json:"user"`
	Password Password `json:"password"`
}
type User struct {
	Username string `json:"username"`
}
type Password struct {
	Digest string `json:"digest"`
	Algo   string `json:"algorithm"`
}

/*
User login with Username and Password, It’s important to say that
we must not pass the user password as plain-text, applying a hashing
algorithm makes things better (sha-256). Make sure your digest is lower-case!
The request should form as:
{
    "msg": "method",
    "method": "login",
    "id":"42",
    "params":[
        {
            "user": { "username": "example-user" },
            "password": {
                "digest": "some-digest",
                "algorithm":"sha-256"
            }
        }
    ]
}
*/
func (wc *WebSocketClient) Login(username string, password string) error {
	hdigest := sha256.New()
	uid := uuid.New().String()
	hdigest.Write([]byte(password))
	h := hdigest.Sum(nil)
	loginMsg := loginInfo{
		Msg:    "method",
		Method: "login",
		LoginParams: []LoginParams{
			{
				User: User{
					Username: username,
				},
				Password: Password{
					Digest: fmt.Sprintf("%x", h),
					Algo:   "sha-256",
				},
			},
		},
		Id: uid,
	}
	loginJson, _ := json.Marshal(loginMsg)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: loginJson,
	}
	return nil
}
