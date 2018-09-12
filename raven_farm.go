package raven

import (
	"github.com/newrelic/go-agent"
	"github.com/sanksons/raven/childlock"
)

//
// Ravens are not like street dogs, they belong to a farm.
// Each farm has a Raven manager, whose role is to contain implementation details of
// each raven.
//
type Farm struct {
	manager     RavenManager
	logger      Logger
	newrelicApp newrelic.Application
	lockOptions *childlock.RedisOptions
}

func (this *Farm) AttachNewRelicApp(app newrelic.Application) {
	this.newrelicApp = app
}

func (this *Farm) AttachLockOptions(options childlock.RedisOptions) {
	this.lockOptions = &options
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
// aka Consumer Code
//
func (this *Farm) GetRavenReceiver(id string, s Source) (*RavenReceiver, error) {

	receiver, err := newRavenReceiver(id, s)
	if err != nil {
		return nil, err
	}
	receiver.farm = this

	//Add lock details to receiver.
	if this.lockOptions != nil {
		receiver.lock = childlock.New(receiver.GetId(), CHILD_LOCK_TIMEOUT, *this.lockOptions)
	}
	return receiver, nil
}
