package raven

import (
	"fmt"
)

//
// User to collect Messages
//
type MessageCollector struct {
	destination Destination
	tempQ       string
	farm        *Farm
}

func (this *MessageCollector) SetDestination(d Destination) *MessageCollector {
	this.destination = d
	return this
}

func (this *MessageCollector) SetTempQ(d string) *MessageCollector {
	this.tempQ = d
	return this
}

func (this *MessageCollector) Start(f func(string)) error {

	for {
		//this blocks
		msg, err := this.farm.Manager.Receive(this.destination.Name, this.tempQ)
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

		f(msg) //process message

	}
}
