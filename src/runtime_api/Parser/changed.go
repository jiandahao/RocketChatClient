package Parser

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

/*
once a new message is coming, the received result will be like the following form:
(Note: if you haven't subscribed the "stream-room-messages" event, you will not received
the message notification. you could use the SubscribeStreamRoomMessage api to subscribe,
see runtime_api/methods/subRoomMessage.go)
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
*/
type (
	/*changed meesage for stream-room-messages*/
	ChangedMessageForSRM struct {
		Msg        string     `json:"msg"`
		Collection string     `json:"collection"`
		Id         string     `json:"id"`
		Fields     *SRMFields `json:"fields"`
	}

	// stream-room-message fileds
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
		Date int64 `json:"$date"`
	}

	/*changed message for stream-notify-room*/
	ChangedMessageForSNR struct {
		Msg        string     `json:"msg"`
		Collection string     `json:"collection"`
		Id         string     `json:"id"`
		Fields     *SNRFields `json:"fields"`
	}

	SNRFields struct {
		EventName string        `json:"eventName"`
		Args      []interface{} `json:"args"`
	}
)

/*Parse message for stream-room-message*/
func ParseChangesForSRM(msgJSONStr string) *ChangedMessageForSRM {
	var results ChangedMessageForSRM
	fmt.Println(msgJSONStr)
	err := json.Unmarshal([]byte(msgJSONStr), &results)
	if err != nil {
		log.Println("invalid json string, chould not parse the stream-room-message type message")
		return nil
	}
	fmt.Printf("You have a new message\n"+
		"from:%s\n"+
		"at:%v\n"+
		"message:%s", results.GetSenderName(),
		results.GetSendTime(),
		results.GetMessage())
	return &results
}

/*Parse message for stream-notify-room*/
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
		return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
	}
	return ""
}
