package raven

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const DEFAULT_MSG_TYPE = "DEF"

//
// Prepare message based on the specified details.
//
func PrepareMessage(id string, mtype string, data string) Message {

	if mtype == "" {
		mtype = DEFAULT_MSG_TYPE
	}
	if id == "" {
		uid, _ := uuid.NewUUID()
		id = uid.String()
	}
	return Message{
		Id:   id,
		Data: data,
		Type: mtype,
	}
}

//
// The message that is sent and retrieved.
// @todo: need to check if we can avoid json encoding and decoding.
type Message struct {
	Id   string
	Type string
	Data string

	//Need to check if we need a counter here or time is sufficient.??
	//Done we definitely need a counter, else multiserver will fail.
	//Counter int
	mtime time.Time
}

// String representation of message.
func (this Message) String() string {
	str, _ := json.Marshal(this)
	return string(str)
}

func (this *Message) toJson() string {
	str, err := json.Marshal(this)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(str)
}

func (this *Message) fromJson(data string) error {
	err := json.Unmarshal([]byte(data), this)
	return err
}

//Check if its an empty message.
func (this *Message) isEmpty() bool {
	if this.Data == "" {
		return true
	}
	return false
}
