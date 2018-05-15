package raven

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const BLOCK_FOR_DURATION = 10 * time.Second

//
// Configuration to Initialize redis cluster.
//
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

type RedisCluster struct {
	client *redis.ClusterClient
}

//
//  Implementation of Send() method exposed by raven manager.
//
func (this *RedisCluster) Send(message string, dest string) error {
	ret := this.client.LPush(dest, message)
	if ret.Err() != nil {
		return ret.Err()
	}
	return nil
}

func (this *RedisCluster) Receive(source string, tempQ string) (string, error) {

	if tempQ == "" {
		return this.receive(source)
	} else {
		return this.receiveReliable(source, tempQ)
	}
}

func (this *RedisCluster) receive(source string) (string, error) {
	ret := this.client.BRPop(BLOCK_FOR_DURATION, source)
	err := ret.Err()
	if err != nil && err == redis.Nil {
		//we got an error
		return "", ErrEmptyQueue
	}
	if err != nil {
		return "", err
	}
	sliceRes := ret.Val()
	fmt.Printf("%v\n", sliceRes)

	return "", nil
}

func (this *RedisCluster) receiveReliable(source, tempQ string) (string, error) {
	ret := this.client.BRPopLPush(source, tempQ, BLOCK_FOR_DURATION)
	err := ret.Err()
	if err != nil && err == redis.Nil {
		//we got an error
		return "", ErrEmptyQueue
	}
	if err != nil {
		return "", err
	}
	sliceRes := ret.Val()
	fmt.Printf("%v\n", sliceRes)

	return "", nil
}
