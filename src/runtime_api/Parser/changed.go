package Parser

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

/*
once a new message is coming, the received result will be like the following form:
(Note: if you haven't subscribed the "stream-room-messages" event, you will not received
the message notification. you could use the SubscribeStreamRoomMessage api to subscribe,
see runtime_api/methods/subRoomMessage.go)

for stream-room-messages:
{
	"msg": "changed",
	"collection": "stream-room-messages",
	"id": "id",
	"fields": {
		"eventName": "GENERAL",
		"args": [{
			"_id": "PinxLkpkLd34pLaQn",
			"rid": "GENERAL",
			"msg": "12121212121",
			"ts": {
				"$date": 1567996428549
			},
			"u": {
				"_id": "qd9TRc82mkGGy5m5P",
				"username": "jiandahao",
				"name": "jiandahao"
			},
			"_updatedAt": {
				"$date": 1567996428585
			},
			"mentions": [],
			"channels": []
		}]
	}
}

for stream-notify-room :
{
	"msg": "changed",
	"collection": "stream-notify-room",
	"id": "id",
	"fields": {
		"eventName": "GENERAL/typing",
		"args": ["jiandahao", false]
	}
}
*/
type (
	/*changed message for stream-room-messages*/
	ChangedMessageForSRM struct {
		Msg        string     `json:"msg"`
		Collection string     `json:"collection"`
		Id         string     `json:"id"`
		Fields     *SRMFields `json:"fields"`
	}

	// stream-room-message fields
	SRMFields struct {
		EventName string     `json:"eventName"`
		Args      []*SRMArgs `json:"args"`
	}

	// stream-room-message args
	SRMArgs struct {
		Id        string        `json:"_id"`
		Rid       string        `json:"rid"`
		Msg       string        `json:"msg"`
		Timestamp *Date         `json:"ts"`
		User      *User         `json:"u"`
		UpdateAt  *Date         `json:"_updateAt"`
		Mentions  []interface{} `json:"mentions"`
		Channels  []interface{} `json:"channels"`
	}
	User struct {
		Id       string `json:"_id"`
		Username string `json:"username"`
		Name     string `json:"name"`
	}

	Date struct {
		Date int64 `json:"$date"` //rocket.chat use timestamp in nanoseconds
	}

	/*changed message for stream-notify-room*/
	ChangedMessageForSNR struct {
		Msg        string     `json:"msg"`
		Collection string     `json:"collection"`
		Id         string     `json:"id"`
		Fields     *SNRFields `json:"fields"`
	}

	SNRFields struct {
		EventName string    `json:"eventName"`
		Args      []*SNRArg `json:"args"`
	}

	SNRArg struct {
		MsgId string `json:"_id"`
	}
)

/*Parse message for stream-room-message*/
func ParseChangesForSRM(msgJSONStr string) *ChangedMessageForSRM {
	var results ChangedMessageForSRM
	err := json.Unmarshal([]byte(msgJSONStr), &results)
	if err != nil {
		log.Println("invalid json string, chould not parse the stream-room-message type message")
		return nil
	}
	return &results
}

// GetMessage gets the message body
func (cm *ChangedMessageForSRM) GetMessage() string {
	if args := cm.GetArgs(); args != nil {
		return args.Msg
	}
	return ""
}

// GetSenderName gets the name of user who send this message
func (cm *ChangedMessageForSRM) GetSenderName() string {
	if u := cm.GetSender(); u != nil {
		return u.Username
	}
	return ""
}

func (cm *ChangedMessageForSRM) GetFields() *SRMFields {
	if cm == nil {
		log.Println("Error ChangedMessage results: empty")
		return nil
	}
	return cm.Fields
}

func (cm *ChangedMessageForSRM) GetCollection() string {
	if cm == nil {
		log.Println("Error ChangedMessage results: empty")
		return ""
	}
	return cm.Collection
}

func (cm *ChangedMessageForSRM) GetArgs() *SRMArgs {
	if fields := cm.GetFields(); fields != nil {
		return fields.Args[0]
	}
	return nil
}

// eventName field represents the room id
func (cm *ChangedMessageForSRM) GetEventName() string {
	if fields := cm.GetFields(); fields != nil {
		return fields.EventName
	}
	return ""
}

func (cm *ChangedMessageForSRM) GetRoomId() string {
	if args := cm.GetArgs(); args != nil {
		return args.Rid
	}
	return ""
}

func (cm *ChangedMessageForSRM) GetSender() *User {
	if args := cm.GetArgs(); args != nil {
		return args.User
	}
	return nil
}

func (cm *ChangedMessageForSRM) GetSendTime() string {
	if args := cm.GetArgs(); args != nil {
		ts := args.Timestamp.Date
		// TODO correct the time parsing
		return time.Unix(0, ts*1000000).Format("2006-01-02 15:04:05")
	}
	return ""
}

/*Parse message for stream-notify-room*/
/*{"msg":"changed","collection":"stream-notify-room","id":"id","fields":
{"eventName":"GENERAL/deleteMessage","args":[{"_id":"KfZZZfMzv76enWwWG"}]}}*/
func ParseChangesForSNR(msgJSONStr string) *ChangedMessageForSNR {
	var results ChangedMessageForSNR
	fmt.Println(msgJSONStr)
	err := json.Unmarshal([]byte(msgJSONStr), &results)
	if err != nil {
		log.Println("invalid json string, chould not parse the stream-notify-room type message: " + err.Error())
		return nil
	}
	return &results
}

func (cs *ChangedMessageForSNR) GetCollection() string {
	if cs == nil {
		log.Println("Error ChangedMessage results: empty")
		return ""
	}
	return cs.Collection
}

func (cs *ChangedMessageForSNR) GetFields() *SNRFields {
	if cs == nil {
		log.Println("Error ChangedMessage results: empty")
		return nil
	}
	return cs.Fields
}

// GetEventName returns the specified eventName, the roomId and the event(typing or deleteMessage)
func (cs *ChangedMessageForSNR) GetEventName() (eventName string, roomId string, event string) {
	if f := cs.GetFields(); f != nil {
		if len(f.EventName) > 0 && string(f.EventName[0]) != "/" && strings.Contains(f.EventName, "/") {
			event := strings.Split(f.EventName, "/")
			return f.EventName, event[0], event[1]
		}
	}
	return "", "", ""
}

func (cs *ChangedMessageForSNR) GetArgs() []*SNRArg {
	if f := cs.GetFields(); f != nil && len(f.Args) > 0 {
		return f.Args
	}
	return nil
}

func (cs *ChangedMessageForSNR) GetMessageId() string {
	if args := cs.GetArgs(); args != nil && len(args) > 0 {
		return args[0].MsgId
	}
	return ""
}
