package raven

//
// Raven defines a message delivery object
//
type Raven struct {
	message     Message
	destination Destination
	helper      FarmManager
}

func (this *Raven) SetMessage(m Message) *Raven {
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
		return ErrNoDestination
	}
	return nil
}
