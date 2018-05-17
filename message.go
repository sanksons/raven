package raven

import (
	"encoding/json"
	"fmt"
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
	//Done we definitely need a counter, else multiserver will fail.
	//Counter int
	mtime time.Time
}

func (this Message) String() string {
	str, _ := json.Marshal(this)
	return string(str)
}

func (this Message) toJson() string {
	str, err := json.Marshal(this)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(str)
}

func (this *Message) isEmpty() bool {
	if this.Data == "" {
		return true
	}
	return false
}
