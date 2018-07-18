package raven

import (
	"strings"
)

// createQ based on the supplied name and bucket.
func createQ(name string, bucket string) Q {
	return Q{name: name, bucket: bucket}
}

//
// A type for each Q.
//
type Q struct {
	//Name of the Queue
	name string

	//Bucket to which queue belongs.
	bucket string
}

//Check if Q is empty.
func (this *Q) IsEmpty() bool {
	if this.name == "" {
		return true
	}
	return false
}

func (this *Q) GetName() string {
	if this.name == "" {
		return ""
	}
	return strings.ToLower(this.name) + "-{" + strings.ToLower(this.bucket) + "}"
}

func (this *Q) GetRawName() string {
	return this.name
}

func (this *Q) GetBucket() string {
	return strings.ToLower(this.bucket)
}

//
// Exposed method for creation of new Source.
//
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

//
// Exposed methos for creation of new Destination.
//
func CreateDestination(name string, bucket string) Destination {

	d := Destination{
		createQ(name, bucket),
	}
	return d
}

//
// Destination specifies the Queue name to which the message needs to be sent.
//
type Destination struct {
	Q
}
