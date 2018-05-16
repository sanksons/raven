package raven

import (
	"strings"
)

type Q struct {
	Name string
}

func (this Q) String() string {
	return strings.ToLower(this.Name)
}

func (this Q) IsEmpty() bool {
	if this.Name == "" {
		return true
	}
	return false
}

//
// Specifies the Queue Name from which messages needs to be retrieved.
//
type Source struct {
	Q
}

//
// Destination specifies the Queue name to which the message needs to be sent.
//
type Destination struct {
	Q
}
