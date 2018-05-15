package raven

//
// Destination specifies the Queue name to which the message needs to be sent.
// If the Destination is marked as reliable, an additional queue is maintained to
// store messages that have failed processing.
//
type Destination struct {
	Name       string
	IsReliable bool
}

func (this *Destination) IsValid() bool {
	destinations := GetValidDestinations()
	for _, di := range destinations {
		if di == this.Name {
			return true
		}
	}
	return false
}

//@todo: this needs to be picked from config.
func GetValidDestinations() []string {
	return []string{"London", "Asia"}
}
