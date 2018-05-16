package raven

//
// Ravens are not like street dogs, they belong to a farm.
// Each farm has a Raven manager, whose role is to contain implementation details of
// each raven.
//
type Farm struct {
	Manager RavenManager
}

//
// Pick a Raven from Farm.
//
// This functions returns a raven that can be used.
// Before flying a Raven do not forget to set the Destination
// and the message that raven needs to carry.
//
// ex: farm.GetRaven().HandMessage().SetDestination().Fly()
//
func (this *Farm) GetRaven() *Raven {
	r := new(Raven)
	r.farm = this
	return r
}

//
// This function returns a picker which can be used to pick messages sent via raven.
// aka  Consumer Code
//
func (this *Farm) GetRavenReceiver(s Source) *RavenReceiver {
	collector := new(RavenReceiver)
	collector.SetSource(s)
	collector.farm = this
	return collector
}