package raven

import (
	"encoding/json"
	"time"
)

type Message struct {
	Id      string
	Data    string
	Counter int
	Time    time.Duration
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
