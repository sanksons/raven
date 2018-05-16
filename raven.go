package raven

import "fmt"

const FARM_TYPE_REDISCLUSTER = "redis-cluster"

//
// Entry point to this library.
// mtype: Farm magaer type.
// config: Farm manager config.
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
