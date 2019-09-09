package runtime_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type WebSocketClient struct {
	// websocket connection
	Conn *websocket.Conn
	// ping ticker for heart beat
	PingTicker *time.Ticker
	// channel for detecting response from websocket server
	Response chan Response
	// channel for detecting request to websocket server
	Request chan Request

	wg sync.WaitGroup

	// wait to close the connection
	shouldClose chan string

	// subscription list, it is optional
	//	SubscriptionList []Subscription

	/*message callback handlers, will be trigger
	  when a specified message is received from server*/
	PingMessageHandler  func() // Ping, there is a default implementation if user haven't implemented
	PongMessageHandler  func() // Pong
	CloseMessageHandler func() // Close
	// Text message
	TextMessageHandler func(msgType int, msgStr string, msgObj map[string]interface{})
	MainEntry          func() // main task after connecting websocket socket, user specified
}

/*subscription list, you could add the subscriptions into
this list. once you have set this list , Run() function
will execute all the subscriptions after login*/
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
	wsClient.PingTicker = time.NewTicker(time.Second * 10)
	wsClient.Response = make(chan Response)
	wsClient.Request = make(chan Request)
	wsClient.shouldClose = make(chan string)
	return wsClient, nil
}

func (wc *WebSocketClient) Close() {
	log.Println("websocket client is going to close")
	wc.Conn.Close()
	close(wc.Response)
	close(wc.Request)
	close(wc.shouldClose)
	wc.PingTicker.Stop()
}

func (wc *WebSocketClient) Run() error {
	if wc == nil {
		return errors.New("invalid websocket connector ")
	}

	if wc.PingMessageHandler == nil {
		// if ping message callback is nil, set a default one
		wc.PingMessageHandler = func() {
			if err := wc.SendPong(); err != nil {
				log.Println(err)
			}
		}
	}

	if wc.PongMessageHandler == nil {
		// if pong message callback is nil, set a default one
		wc.PongMessageHandler = func() {
			if err := wc.SendPing(); err != nil {
				log.Println(err)
			}
		}
	}

	go wc.ReadMessage()
	go wc.WriteMessage()

	/*Before requesting any method / subscription
	you have to send a connect message firstly*/
	if err := wc.Connect(); err != nil {
		log.Println(err)
		return errors.New("Error occurs when create a connection ")
	}
	//for _, sub := range wc.SubscriptionList {
	//	_ = sub.Handler(sub.Args)
	//}
	if wc.MainEntry != nil {
		go wc.MainEntry()
	}

	// waiting to close the websocket client
	wc.wg.Wait() //
	// do closing
	fmt.Println("Main run routine going to close")
	wc.Close()
	return nil
}

func (wc *WebSocketClient) Done() error {
	return nil
}

func (wc *WebSocketClient) ReadMessage() {
	defer func() {
		log.Println("======ReadMessage Function exit=====")
		wc.Conn.Close()
	}()
	for {
		mt, p, err := wc.Conn.ReadMessage()
		if err != nil {
			log.Println("error occurs: ", mt, string(p), err)
			wc.wg.Done()
			wc.shouldClose <- "close" // error occurs, tell write goroutine to exit
			break
		}
		var msg map[string]interface{}
		_ = json.Unmarshal(p, &msg)

		// Note: actually, the message type could not be PongMessage / PingMessage/Close Message
		// the PongMessage/PingMessage/CloseMessage will be handle by function that
		// binded with wc.Conn.SetPongHandler/wc.Conn.SetPingHandler/wc.Conn.SetCloseHandler
		if mt == websocket.PongMessage || msg["msg"] == "pong" {
			log.Println("receive a pong message from server")
			wc.PongMessageHandler()
			continue
		}

		if mt == websocket.PingMessage || msg["msg"] == "ping" {
			log.Println("receive a ping message from server")
			wc.PingMessageHandler()
			continue
		}

		if mt == websocket.TextMessage {
			// create a new goroutine to handler text message
			log.Println("receive message: ", string(p))
			go wc.TextMessageHandler(mt, string(p), msg)
		}

	}
}

func (wc *WebSocketClient) WriteMessage() {
	defer func() {
		fmt.Println("======WriteMessage Function exit=====")
		wc.Conn.Close()
	}()
	for {
		select {
		case req := <-wc.Request:
			log.Println("send request: ", string(req.msg))
			if err := wc.Conn.WriteMessage(req.mt, req.msg); err != nil {
				log.Println(err)
				wc.wg.Done()
				return
			}
		case <-wc.shouldClose:
			log.Println("Write routine is going to close")
			wc.wg.Done()
			return
		}

	}
}
