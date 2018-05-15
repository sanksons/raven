package raven

import (
	"fmt"
)

type MessageCollector struct {
	destination Destination
	farm        *Farm
}

func (this *MessageCollector) SetDestination(d Destination) *MessageCollector {
	this.destination = d
	return this
}

func (this *MessageCollector) Start(f func(string)) error {

	for {
		//this blocks
		msg, err := this.farm.Manager.Receive(this.destination.Name)
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
