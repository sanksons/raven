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

func (this *MessageCollector) Start(f func(string) error) error {

	//validate before start picking
	return fmt.Errorf("To be Impl")
}
