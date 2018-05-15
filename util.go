package raven

import (
	"fmt"
)

const FARM_TYPE_REDISCLUSTER = "redis-cluster"

//
// An outsider function, to initialize a farm.
//
func InitializeFarm(mtype string, config interface{}) (*Farm, error) {
	f := new(Farm)
	switch mtype {
	case FARM_TYPE_REDISCLUSTER:
		conf := config.(RedisClusterConfig)
		redis := InitializeRedisCluster(conf)
		f.Manager = redis
		return f, nil

	default:
		return nil, fmt.Errorf("Not a Valid Raven Manager supplied")
	}
}

type Farm struct {
	Manager RavenManager
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
	r.farm = this
	return r
}

//
// This function returns a picker which can be used to pick messages sent via raven.
// aka  Consumer Code
//
func (this *Farm) MessageCollector(d Destination) *MessageCollector {
	collector := new(MessageCollector)
	collector.SetDestination(d)
	collector.farm = this
	return collector
}
