package runtime_api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/*
three ways to login:
- User login with Username and Password, It’s important to say that
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
- Using Authentication providers. Here’s a example request.
{
    "msg": "method",
    "method": "login",
    "id":"42",
    "params": [
        {
            "oauth": {
                "credentialToken":"credential-token",
                "credentialSecret":"credential-secret"
            }
        }
    ]
}
- Using an authentication token
If you have a saved user authentication you may use the provided auth-token to automatically log the user in.

{
    "msg": "method",
    "method": "login",
    "id": "42",
    "params":[
        { "resume": "auth-token" }
    ]
}
As the token expires, you have to call the login method again in order to obtain a new token with a new expiration date.
NB: You don’t have to wait until the token is expired before asking for a new token.
reference :https://rocket.chat/docs/developer-guides/realtime-api/method-calls/login/
*/
type loginInfo struct {
	Msg         string        `json:"msg"`
	Method      string        `json:"method"`
	LoginParams []LoginParams `json:"params"`
	Id          string        `json:"id"`
}

//there are three ways to login
type LoginParams struct {
	User     *User     `json:"user,omitempty"` // for using username and password to login
	Password *Password `json:"password,omitempty"`
	OAuth    *OAuth    `json:"oauth,omitempty"`  //for Using Authentication providers to login
	Resume   string    `json:"resume,omitempty"` //for using an authentication token to login
}
type User struct {
	Username string `json:"username,omitempty"`
}
type Password struct {
	Digest string `json:"digest,omitempty"`
	Algo   string `json:"algorithm, omitempty"`
}
type OAuth struct {
	CredentialToken  string `josn:"credentialToken,omitempty"`
	CredentialSecret string `json:"credentialSecret,omitempty"`
}

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
				User: &User{
					Username: username,
				},
				Password: &Password{
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

//TODO test
func (wc *WebSocketClient) LoginWithOAuth(credentialToken string, credentialSecret string) error {
	loginMsg := loginInfo{
		Msg:    "method",
		Method: "login",
		LoginParams: []LoginParams{
			{
				OAuth: &OAuth{
					CredentialToken:  credentialToken,
					CredentialSecret: credentialSecret,
				},
			},
		},
		Id: uuid.New().String(),
	}
	loginJson, _ := json.Marshal(loginMsg)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: loginJson,
	}
	return nil
}

//TODO test
func (wc *WebSocketClient) LoginWithAuthToken(authToken string) error {
	loginMsg := loginInfo{
		Msg:    "method",
		Method: "login",
		LoginParams: []LoginParams{
			{
				Resume: authToken,
			},
		},
		Id: uuid.New().String(),
	}
	loginJson, _ := json.Marshal(loginMsg)
	wc.Request <- Request{
		mt:  websocket.TextMessage,
		msg: loginJson,
	}
	return nil
}
