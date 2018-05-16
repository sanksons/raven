package raven

import (
	"encoding/json"
	"time"
)

const KEY_TYPE_PUB_SEQ = "publisher-seq"

func PrepareMessage(id string, mtype string, data string) Message {
	return Message{
		Id:   id,
		Data: data,
		Type: mtype,
	}
}

//
// The message that is sent and retrieved.
// @todo: need to check if we can avoid json encoding and decoding.
// @todo: check if we need counter or time is enough.
//
type Message struct {
	Id   string
	Type string
	Data string

	//Need to check if we need a counter here or time is sufficient.??
	//Counter int
	Time time.Time
}

//
// @todo: need to imnplement thi shit.
//
func (this *Message) GetKeyName(ktype string, prefix string, postfix string) string {
	switch ktype {
	case KEY_TYPE_PUB_SEQ:
		key := "counter_" + this.Type + "_" + this.Id
		return key
	}
	return ""
}

func (this Message) String() string {
	str, _ := json.Marshal(this)
	return string(str)
}

func (this *Message) IsEmpty() bool {
	if this.Data == "" {
		return true
	}
	return false
}
