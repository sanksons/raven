package raven

import (
	"fmt"
	"strings"
)

func createQ(name string, bucket string) Q {
	return Q{name: name, bucket: bucket}
}

type Q struct {
	name   string
	bucket string
}

func (this *Q) IsEmpty() bool {
	if this.name == "" {
		return true
	}
	return false
}

func (this *Q) GetName() string {
	return strings.ToLower(this.name) + "-{" + strings.ToLower(this.bucket) + "}"
}

func (this *Q) GetBucket() string {
	return strings.ToLower(this.bucket)
}

func CreateSource(name string, bucket string) Source {
	return Source{
		createQ(name, bucket),
	}
}

//
// Specifies the Queue Name from which messages needs to be retrieved.
//
type Source struct {
	Q
}

func CreateDestination(name string, bucket string) Destination {

	d := Destination{
		createQ(name, bucket),
	}
	fmt.Printf("Destination:%v \n", d)
	return d
}

//
// Destination specifies the Queue name to which the message needs to be sent.
//
type Destination struct {
	Q
}
