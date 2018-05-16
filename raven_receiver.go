package raven

import (
	"fmt"
)

//
// User to collect Messages
//
type RavenReceiver struct {
	source      Source
	processingQ Q
	deadQ       Q
	farm        *Farm
}

func (this *RavenReceiver) SetSource(s Source) *RavenReceiver {
	this.source = s
	return this
}

func (this *RavenReceiver) SetProcessingQ(q Q) *RavenReceiver {
	this.processingQ = q
	return this
}

func (this *RavenReceiver) Start(f func(string) error) error {

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
