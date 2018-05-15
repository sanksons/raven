package raven

func InitializeFarm() (*Farm, error) {

}

type Farm struct {
	Manager FarmManager
}

//
// This functions returns a raven that can be flied.
// Before flying a Raven do not forget to set the Destination
// and the message that raven needs to carry.
//
// ex: farm.GetRaven().SetMessage().SetDestination().Fly()
//
func (this *Farm) GetRaven() *Raven {
	r := new(Raven)
	r.helper = this.Manager
	return r
}

//
// This function returns a picker which can be used to pick messages sent via raven.
//
//
func GetPicker() *Picker {
	return new(Picker)
}
