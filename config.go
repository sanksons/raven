package raven

import "github.com/go-redis/redis"

var _ RavenManager = (*RedisCluster)(nil)

type RedisCluster struct {
	client *redis.ClusterClient
}

//
//  Implementation of Send() method exposed by raven manager.
//
func (this *RedisCluster) Send(message string) error {
	return nil
}

type RedisClusterConfig struct {
	Addrs    []string
	Password string
	PoolSize int
}

func InitializeRedisCluster(config RedisClusterConfig) *RedisCluster {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    config.Addrs,
		Password: config.Password,
		PoolSize: config.PoolSize,
	})
	redisCluster := new(RedisCluster)
	redisCluster.client = client
	return redisCluster
}
