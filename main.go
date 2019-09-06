package main

import (
	"encoding/json"
	"fmt"
	"log"
	parser "runtime_api/Parser"
	rc "runtime_api/methods"
	"time"
)

func main() {
	//random := 9*rand.Intn(100)+100
	//sessId := fmt.Sprintf("jian%v",random)
	//wsUrl := fmt.Sprintf("ws://127.0.0.1:3000/sockjs/%v/%s/websocket",random,sessId)
	wsUrl := fmt.Sprintf("ws://127.0.0.1:3000/websocket")
	//wsUrl := "ws://127.0.0.1:3000/sockjs/824/4n2rq6yk/websocket"
	wsClient, err := rc.NewWebSocketClient(wsUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer wsClient.Close()

	wsClient.Username = "dahao"
	wsClient.Password = "xdh5695565"

	wsClient.TextMessageHandler = func(mt int, msgStr string, msgObj map[string]interface{}) {
		m, ok := msgObj["msg"]
		if !ok {
			log.Println(msgStr)
			return
		}
		if met, ok := msgObj["method"]; ok {
			log.Println(met)
		}

		if m == "updated" {
			log.Println("receive updated")
			fmt.Println(msgObj)
			return
		}

		if m == "changed" {
			log.Println("receive changed")
			var changedMsg parser.ChangedMessage
			_ = json.Unmarshal([]byte(msgStr), &changedMsg)
			log.Println(changedMsg)
			//return
		}

		if m == "result" {
			log.Println("receive result")
			return
		}
	}

	wsClient.PingMessageHandler = func() {
		if err := wsClient.SendPong(); err != nil {
			log.Println(err)
		}
	}

	wsClient.PongMessageHandler = func() {
		if err := wsClient.SendPing(); err != nil {
			log.Println(err)
		}
	}

	wsClient.SubscriptionList = []rc.Subscription{
		{
			wsClient.SubscribeStreamRoomMessage,
			[]string{
				"GENERAL",
			},
		},
		{
			wsClient.SubscribeStreamNotifyUser,
			[]string{
				"qd9TRc82mkGGy5m5P",
				"message",
			},
		},
	}
	wsClient.CloseMessageHandler = func() {
		fmt.Println("receive a close message from server")
	}

	go func() {
		//username := "dahao"
		//_ = wsClient.Login(username,"xdh5695565")
		_ = wsClient.GetUserRoles()
		for i := 0; i < 2; i++ {
			time.Sleep(2 * time.Second)
			_ = wsClient.SendMessage("GENERAL", "it is a test message from ."+" at "+time.Now().Format("2006-01-02 15:04:05"))
		}
	}()

	log.Fatal(wsClient.Run())
}
