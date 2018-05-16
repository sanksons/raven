package raven

import (
	"strings"
)

func CreateQ(name string, bucket string) Q {
	return Q{name: name, bucket: bucket}
}

type Q struct {
	name   string
	bucket string
}

func (this Q) String() string {
	return strings.ToLower(this.name)
}

func (this Q) IsEmpty() bool {
	if this.name == "" {
		return true
	}
	return false
}

func (this *Q) GetName() string {
	return strings.ToLower(this.name)
}

func (this *Q) GetBucket() string {
	return strings.ToLower(this.bucket)
}

func CreateSource(name string, bucket string) Source {

	return Source{
		Q{name: name, bucket: bucket},
	}
}

//
// Specifies the Queue Name from which messages needs to be retrieved.
//
type Source struct {
	Q
}

func CreateDestination(name string, bucket string) Destination {

	return Destination{
		Q{name: name, bucket: bucket},
	}
}

//
// Destination specifies the Queue name to which the message needs to be sent.
//
type Destination struct {
	Q
}
