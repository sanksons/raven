package raven

//
// History tells us ravens are the reliable and quickest way for sending messages.
// Raven defines a message delivery object.
//
type Raven struct {
	// A message that raven carries.
	message Message

	// Message Destination
	destination Destination

	//To which farm the raven belongs.
	//This helps in identifying the Farm Manager of Raven.
	farm *Farm
}

//
// Hand over the Message to this raven.
//
func (this *Raven) HandMessage(m Message) *Raven {
	this.message = m
	return this
}

//
// Tell the Raven where to deliver message.
//
func (this *Raven) SetDestination(d Destination) *Raven {
	this.destination = d
	return this
}

//
// Send Message.
//
func (this *Raven) Fly() error {
	//Its a waste of raven if mesasge is empty.
	if this.message.IsEmpty() {
		return ErrNoMessage
	}
	// We dont want our raven to wander around world!!
	if !this.destination.IsEmpty() {
		return ErrInvalidDestination
	}
	return this.farm.Manager.Send(this.message, this.destination)
	//return nil
}
