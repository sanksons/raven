package raven

import (
	"fmt"
	"strconv"
	"strings"
)

// createQ based on the supplied name and bucket.
func createQ(name string, bucket string) Q {
	return Q{name: name, bucket: bucket}
}

func createMsgBox(name string, bucket string) MsgBox {
	return MsgBox{name: name, boxId: bucket}
}

type MsgBox struct {
	//Name of the Queue
	name string

	//Bucket to which queue belongs.
	boxId string
}

func (this *MsgBox) GetName() string {
	if this.name == "" {
		return ""
	}
	return strings.ToLower(this.name) + "-{" + strings.ToLower(this.boxId) + "}"
}

func (this *MsgBox) GetRawName() string {
	return this.name
}

func (this *MsgBox) GetBoxId() string {
	return strings.ToLower(this.boxId)
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
// Exposed method for creation of new Destination.
//
func CreateDestination(name string, boxes int, shardlogic func(Message, int) (string, error)) Destination {

	if boxes < 1 {
		boxes = 1
	}
	var i = 1

	msgBoxes := make([]MsgBox, boxes)
	for i <= boxes {
		q := createMsgBox(name, strconv.Itoa(i))
		msgBoxes[i-1] = q
		i++
	}
	//incase no shardlogic is provided use default.
	if shardlogic == nil {
		shardlogic = DefaultShardlogic
	}

	d := Destination{
		Name:       name,
		shardLogic: shardlogic,
		MsgBoxes:   msgBoxes,
	}
	return d
}

//
// Destination specifies the Queue name to which the message needs to be sent.
//
type Destination struct {
	Name       string
	MsgBoxes   []MsgBox
	shardLogic func(Message, int) (string, error)
}

func (this *Destination) GetAllBoxes() ([]MsgBox, error) {
	return this.MsgBoxes, nil
}

func (this *Destination) GetBox4Msg(m Message) (*MsgBox, error) {
	boxId, err := this.shardLogic(m, len(this.MsgBoxes))
	if err != nil {
		return nil, err
	}
	for _, b := range this.MsgBoxes {
		if b.GetBoxId() == boxId {
			return &b, nil
		}
	}
	return nil, fmt.Errorf("Shard Logic Seems to be Incorrect, Msg Box with Id [%s] does not exists", boxId)
}

//@todo: check if its a valid destination.
// and a lot more
func (this *Destination) Validate() error {
	return nil
}
