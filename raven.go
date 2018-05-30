package raven

import "fmt"

const FARM_TYPE_REDISCLUSTER = "redis-cluster"

//
// Entry point to this library.
// mtype: Farm magaer type.
// config: Farm manager config.
//
func InitializeFarm(mtype string, config interface{}, inlogger Logger) (*Farm, error) {
	f := new(Farm)

	//assign logger
	f.logger = new(DummyLogger)
	if inlogger != nil {
		f.logger = inlogger
	}

	//assign adapter
	switch mtype {
	case FARM_TYPE_REDISCLUSTER:
		conf := config.(RedisClusterConfig)
		redis := InitializeRedisCluster(conf)
		f.manager = redis
		return f, nil

	default:
		return nil, fmt.Errorf("Not a Valid Raven Manager supplied")
	}
}
