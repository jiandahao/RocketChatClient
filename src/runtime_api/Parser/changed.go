package Parser

/*
a["{\"msg\":\"changed\",\"collection\":\"stream-room-messages\",\"id\":\"id\",\"fields\":
{\"eventName\":\"GENERAL\",\"args\":[{\"_id\":\"W66XWX7DcEcGE7SJD\",\"rid\":\"GENERAL\",\"msg\":\"jiandahao\",\"ts\":
{\"$date\":1567489961718},\"u\":{\"_id\":\"qd9TRc82mkGGy5m5P\",\"username\":\"jiandahao\",
\"name\":\"jiandahao\"},\"_updatedAt\":{\"$date\":1567489961755},\"mentions\":[],\"channels\":[]}]}}"]
*/
type ChangedMessage struct {
	Msg        string `json:"msg"`
	Collection string `json:"collection"`
	Id         string `json:"id"`
	Fields     Fields `json:"fields"`
}

type Fields struct {
	EventName string `json:"eventName"`
	Args      []Args `json:"args"`
}

type Args struct {
	Id        string        `json:"_id"`
	Rid       string        `json:"rid"`
	Msg       string        `json:"msg"`
	Timestamp Date          `json:"ts"`
	User      User          `json:"u"`
	UpdateAt  Date          `json:"_updateAt"`
	Mentions  []interface{} `json:"mentions"`
	Channels  []interface{} `json:"channels"`
}
type User struct {
	Id       string `json:"_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type Date struct {
	Date int64 `json:"$date"`
}
