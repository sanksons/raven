package raven

import (
	"fmt"
)

func newRavenReceiver(id string, source Source, reliable bool, ordering bool) (*RavenReceiver, error) {
	rr := new(RavenReceiver)
	err := rr.setId(id)
	if err != nil {
		return nil, err
	}
	rr.setSource(source)
	if reliable {
		rr.defineDeadQ().defineProcessingQ()
	}
	if ordering {
		rr.ordering = true
	}
	return rr, nil
}

//
// Message collector
//
type RavenReceiver struct {
	id          string
	source      Source
	isReliable  bool
	processingQ Q
	deadQ       Q
	ordering    bool
	farm        *Farm
}

func (this *RavenReceiver) setSource(s Source) *RavenReceiver {
	this.source = s
	return this
}

func (this *RavenReceiver) setId(id string) error {
	if id == "" {
		return fmt.Errorf("You need to define a unique ID for your consumer.")
	}
	return nil
}

func (this *RavenReceiver) defineProcessingQ() *RavenReceiver {

	qname := fmt.Sprintf("%s_processing_%s", this.source.GetName(), this.id)
	this.processingQ = CreateQ(qname, this.source.GetBucket())
	return this
}

func (this *RavenReceiver) defineDeadQ() *RavenReceiver {

	qname := fmt.Sprintf("%s_dead", this.source.GetName())
	this.deadQ = CreateQ(qname, this.source.GetBucket())
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
		msg, err := this.farm.Manager.Receive(this.source, this.processingQ)
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
			this.farm.Manager.MarkProcessed(msg, this.processingQ)
		} else {
			//store in DeadQ
			//@todo: do I need to check for error here.??
			this.farm.Manager.MarkFailed(msg, this.deadQ, this.processingQ)
		}

	}
}
