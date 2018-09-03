package raven

import "fmt"

//
// Initiate a Raven Receiver.
//
func newRavenReceiver(id string, source Source) (*RavenReceiver, error) {
	rr := new(RavenReceiver)
	//Define source and Id for receiver.
	rr.setSource(source).setId("")

	// Generate Message Receivers, for each message box
	msgreceivers := make([]*MsgReceiver, len(source.MsgBoxes))
	for _, box := range source.MsgBoxes {
		m := &MsgReceiver{
			msgbox: box,
			parent: rr,
		}
		// Set Id for msgReceiver.
		m.setId(box.GetName())
		msgreceivers = append(msgreceivers, m)
	}

	return rr, nil
}

type RavenReceiver struct {

	//A Unique Id that distinguishes this Receiver from other receivers.
	id string

	// Source from which receiver will fetch Ravens.
	source Source

	// Receiving options.
	options struct {
		//Specifies if we want to use reliable Q or not
		//@todo: ordering is yet to be implemented.
		isReliable, ordering bool
	}

	//All the child receivers.
	msgReceivers []*MsgReceiver

	// Access to Raven farm and underlying adapters.
	farm *Farm
}

func (this *RavenReceiver) setSource(s Source) *RavenReceiver {
	this.source = s
	return this
}

func (this *RavenReceiver) setId(id string) *RavenReceiver {
	// Make sure we are setting source before Id.
	// Since a Source can have only one receiver at a time, it makes perfect sense to allot
	// source name as ID.
	this.id = this.source.GetName()
	return this
}

//
// Markall the allotted message receivers as reliable.
//
func (this *RavenReceiver) MarkReliable() *RavenReceiver {
	this.options.isReliable = true

	for _, msgReceiver := range this.msgReceivers {
		msgReceiver.MarkReliable()
	}
	return this
}

//@todo: implement all the necessary validations required for a receiver.
func (this *RavenReceiver) validate() error {

	//Check if Id, Source and farm are defined.
	// check if atleast one receiver is assigned.

	if this.id == "" {
		return fmt.Errorf("An Id needs to be assigned to Receiver. Make sure its unique within source")
	}
	if this.source.GetName() == "" {
		return fmt.Errorf("Receiver Source cannot be Empty")
	}
	if this.farm == nil {
		return fmt.Errorf("You need to define to which farm this receiver belongs.")
	}
	if len(this.msgReceivers) <= 0 {
		return fmt.Errorf("Atleast one msg Receiver needs to be assigned")
	}
	return nil
}

func (this *RavenReceiver) Start(f MessageHandler) error {
	// execute prestart hook of all receivers.
	// once all prestart hooks are successfull start receivers.
	for _, msgreceiver := range this.msgReceivers {
		if err := msgreceiver.preStart(); err != nil {
			return err
		}
	}

	//@todo: Start receivers.
	// Since the start functions of receivers block, we need to start
	// receivers as seperate goroutines.
	for _, msgreceiver := range this.msgReceivers {
		go msgreceiver.start(f)
	}

	return nil
}
