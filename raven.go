package raven

//
// Raven defines a message delivery object
//
type Raven struct {
	// A message that raven carries.
	message Message
	// Message Destination
	destination Destination
	//
	farm *Farm
}

func (this *Raven) HandMessage(m Message) *Raven {
	this.message = m
	return this
}

func (this *Raven) SetDestination(d Destination) *Raven {
	this.destination = d
	return this
}

func (this *Raven) Fly() error {
	//validate
	if this.message == "" {
		return ErrNoMessage
	}
	if !this.destination.IsValid() {
		return ErrInvalidDestination
	}
	return this.farm.Manager.Send(string(this.message), this.destination.Name)
	//return nil
}
