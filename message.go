package raven

import (
	"encoding/json"
	"time"
)

const KEY_TYPE_PUB_SEQ = "publisher-seq"

func PrepareMessage(id string, mtype string, data string) Message {
	return Message{
		id:    id,
		data:  data,
		mtype: mtype,
	}
}

//
// The message that is sent and retrieved.
// @todo: need to check if we can avoid json encoding and decoding.
// @todo: check if we need counter or time is enough.
//
type Message struct {
	id    string
	mtype string
	data  string

	//Need to check if we need a counter here or time is sufficient.??
	//Counter int
	mtime time.Time
}

func (this Message) String() string {
	str, _ := json.Marshal(this)
	return string(str)
}

func (this Message) toJson() string {
	str, _ := json.Marshal(this)
	return string(str)
}

func (this *Message) isEmpty() bool {
	if this.data == "" {
		return true
	}
	return false
}
