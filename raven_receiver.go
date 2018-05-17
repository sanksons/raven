package raven

import (
	"fmt"
)

func newRavenReceiver(id string, source Source) (*RavenReceiver, error) {
	rr := new(RavenReceiver)

	rr.setSource(source).setId(id)
	return rr, nil
}

//
// Message collector
//
type RavenReceiver struct {
	id     string
	source Source

	//Options define characteristics of a receiver.
	options struct {
		isReliable, ordering bool
	}

	//Q to store processing and dead messages.
	// used only when marked reliable.
	processingQ Q
	deadQ       Q

	// Farm to which reveiver belongs.
	farm *Farm
}

// Mark the Q as reliable.
func (this *RavenReceiver) MarkReliable() *RavenReceiver {
	this.options.isReliable = true
	this.defineProcessingQ().defineDeadQ()
	return this
}

// Mark the Q as ordered.
func (this *RavenReceiver) MarkOrdered() *RavenReceiver {
	this.options.ordering = true
	return this
}

func (this *RavenReceiver) setSource(s Source) *RavenReceiver {
	this.source = s
	return this
}

func (this *RavenReceiver) setId(id string) *RavenReceiver {
	this.id = id
	return this
}

func (this *RavenReceiver) defineProcessingQ() *RavenReceiver {

	qname := fmt.Sprintf("%s_processing_%s", this.source.GetName(), this.id)
	this.processingQ = createQ(qname, this.source.GetBucket())
	return this
}

func (this *RavenReceiver) defineDeadQ() *RavenReceiver {

	qname := fmt.Sprintf("%s_dead", this.source.GetName())
	this.deadQ = createQ(qname, this.source.GetBucket())
	return this
}

//@todo: implement all the necessary validations required for a receiver.
func (this *RavenReceiver) validate() error {
	return nil
}

func (this *RavenReceiver) Start(f func(string) error) error {

	if verr := this.validate(); verr != nil {
		return verr
	}
	// this blocks
	for {
		//this blocks
		msg, err := this.farm.manager.Receive(this.source, this.processingQ)
		if err != nil && err == ErrEmptyQueue {
			//Q is empty, Simple recheck.
			fmt.Println("Queue is empty recheck")
			continue
		}
		if err != nil {
			//add a wait here.
			//log error
			fmt.Println(err.Error())
			return err
		}

		execerr := f(msg.String()) //process message
		if execerr == nil {
			//free up message from processing Q
			//@todo: do I need to check for error here.??
			this.farm.manager.MarkProcessed(msg, this.processingQ)
		} else {
			//store in DeadQ
			//@todo: do I need to check for error here.??
			this.farm.manager.MarkFailed(msg, this.deadQ, this.processingQ)
		}

	}
}
