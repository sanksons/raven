package raven

type Picker struct {
	destination Destination
}

func (this *Picker) SetDestination(d Destination) *Picker {
	this.destination = d
	return this
}

func StartPicking(f func()) error {

	//validate before start picking
	return nil
}
