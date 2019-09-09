package main

import (
	"fmt"
	"log"
	parser "runtime_api/Parser"
	rc "runtime_api/methods"
	"time"
	//"time"
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

	// register a text message handler
	wsClient.TextMessageHandler = func(msgType int, msgStr string, msgObj map[string]interface{}) {
		m, ok := msgObj["msg"]
		if !ok {
			log.Println(msgStr)
			return
		}
		//if met, ok := msgObj["method"]; ok {
		//	log.Println(met)
		//}

		if m == "updated" {
			log.Println("receive updated")
			return
		}

		if m == "changed" {
			log.Println("receive changed")
			//var changedMsg parser.ChangedMessage
			//_ = json.Unmarshal([]byte(msgStr), &changedMsg)
			//log.Println(changedMsg)
			//return
			if v, ok := msgObj["collection"]; ok {
				if v == "stream-room-messages" {
					parser.ParseChangesForSRM(msgStr)
				} else if v == "stream-notify-room" {
					parser.ParseChangesForSNR(msgStr)
				}
			} else {
				fmt.Println("Unrecognizable changed message")
			}

		}

		if m == "result" {
			log.Println("receive result")
			return
		}
	}

	// register a ping message handler, usually you should response with a pong
	wsClient.PingMessageHandler = func() {
		if err := wsClient.SendPong(); err != nil {
			log.Println(err)
		}
	}

	// register a pong message handler, usually you should response with a ping
	wsClient.PongMessageHandler = func() {
		if err := wsClient.SendPing(); err != nil {
			log.Println(err)
		}
	}

	// register close message handler
	wsClient.CloseMessageHandler = func() {
		fmt.Println("Handle close message from server")
	}

	// it is the main entry that will be excuted right after connecting to the server successfully
	wsClient.MainEntry = func() {
		// always login first
		if err := wsClient.Login("jiandahao", "xdh5695565"); err != nil {
			log.Println(err)
		}
		//	// subscribe events that you interested
		_ = wsClient.SubscribeStreamRoomMessage("GENERAL")
		_ = wsClient.SubscribeStreamNotifyRoom("GENERAL", "typing")
		_ = wsClient.SubscribeStreamNotifyUser("qd9TRc82mkGGy5m5P", "message")
		//	_ = wsClient.GetUserRoles()
		_ = wsClient.SendMessage("GENERAL", "it is a test message from ."+" at "+time.Now().Format("2006-01-02 15:04:05"))
		_ = wsClient.SendMessage("GENERAL", "it is a test message 2 from ."+" at "+time.Now().Format("2006-01-02 15:04:05"))
		//	time.Sleep(35*time.Second)
		//	wsClient.SendPing()
	}

	log.Fatal(wsClient.Run())
}
