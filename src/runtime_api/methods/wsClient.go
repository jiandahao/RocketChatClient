package runtime_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
	//"time"
)

type WebSocketClient struct {
	// Username
	Username string
	// Password
	Password string
	// websocket connection
	Conn *websocket.Conn
	// ping ticker for heart beat
	PingTicker *time.Ticker
	// response from websocket server
	Response chan Response
	//
	Request chan Request
	//
	wg sync.WaitGroup
	// wait to close the connection
	Wait2Close chan string

	// subscription list
	SubscriptionList []Subscription
	// Ping message handler
	PingMessageHandler  func()
	PongMessageHandler  func()
	TextMessageHandler  func(mt int, msgStr string, msgObj map[string]interface{})
	CloseMessageHandler func()
}

// subscription list, you can add all the subscript actions into
// this list,  Run() function will execute all the actions after login
type Subscription struct {
	Handler func([]string) error
	Args    []string
}

type Response struct {
	Mt     int                    // message type
	MsgStr string                 // message json string
	MsgObj map[string]interface{} // message map object, per call defined data struct
}

type Request struct {
	mt  int    // message type
	msg []byte // request message
}

func NewWebSocketClient(wsUrl string) (*WebSocketClient, error) {
	wsClient := new(WebSocketClient)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return nil, err
	}
	wsClient.wg.Add(2)
	wsClient.Conn = conn
	wsClient.PingTicker = nil
	wsClient.Response = make(chan Response)
	wsClient.Request = make(chan Request)
	wsClient.Wait2Close = make(chan string)
	return wsClient, nil
}

func (wc *WebSocketClient) Close() {
	fmt.Println("====2=====")
	wc.Conn.Close()
	close(wc.Response)
	close(wc.Request)
	close(wc.Wait2Close)
}

func (wc *WebSocketClient) Run() error {
	if wc == nil {
		return errors.New("invalid websocket connector ")
	}

	go wc.ReadMessage()
	go wc.WriteMessage()

	/*Before requesting any method / subscription
	you have to send a connect message firstly*/
	if err := wc.Connect(); err != nil {
		log.Println(err)
		return errors.New("Error occurs when create a connection ")
	}

	if err := wc.Login(wc.Username, wc.Password); err != nil {
		log.Println(err)
		return errors.New("Error occurs when login ")
	}

	for _, sub := range wc.SubscriptionList {
		_ = sub.Handler(sub.Args)
	}

	// waiting to close the websocket client
	<-wc.Wait2Close
	wc.wg.Wait() //
	// do closing
	fmt.Println("going to close")
	wc.Close()
	return nil
}

func (wc *WebSocketClient) Done() error {
	return nil
}

func (wc *WebSocketClient) ReadMessage() {
	defer func() {
		fmt.Println("======3=====")
		wc.Conn.Close()
	}()
	for {
		mt, p, err := wc.Conn.ReadMessage()
		var msg map[string]interface{}
		_ = json.Unmarshal(p, &msg)

		if err != nil {
			log.Println(err)
			//TODO send a close message
			wc.wg.Done()
			wc.Wait2Close <- "close"
			return
		}

		if mt == websocket.CloseMessage {
			wc.CloseMessageHandler()
			wc.wg.Done()
			wc.Wait2Close <- "close"
			return
		}

		if mt == websocket.PongMessage || msg["msg"] == "pong" {
			fmt.Println("receive a pong message from server")
			wc.PongMessageHandler()
			wc.PingMessageHandler()
			continue
		}

		if mt == websocket.PingMessage || msg["msg"] == "ping" {
			fmt.Println("receive a ping message from server")
			wc.PingMessageHandler()
			//wc.PongMessageHandler()
			continue
		}

		if mt == websocket.TextMessage {
			// create a new goroutine to handler text message
			fmt.Println("receive message: ", string(p))
			go wc.TextMessageHandler(mt, string(p), msg)
		}

	}
}

func (wc *WebSocketClient) WriteMessage() {
	defer func() {
		fmt.Println("=====1======")
		wc.Conn.Close()
	}()
	for {
		select {
		case req := <-wc.Request:
			fmt.Println("send request: ", string(req.msg))
			if err := wc.Conn.WriteMessage(req.mt, req.msg); err != nil {
				//TODO send a close message
				log.Println(err)
				wc.wg.Done()
				return
			}
		case <-wc.Wait2Close:
			fmt.Println("Write routine going to close")
			wc.wg.Done()
		}
	}
}
