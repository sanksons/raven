package raven

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const BLOCK_FOR_DURATION = 10 * time.Second

const MAX_TRY_LIMIT = 3

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
func (this *RedisCluster) Send(message Message, dest Destination) error {

	ret := this.client.LPush(dest.String(), message.String())
	if ret.Err() != nil {
		return ret.Err()
	}
	return nil
}

func (this *RedisCluster) Receive(source Source, procQ Q) (*Message, error) {

	var message string
	var err error
	if procQ.IsEmpty() {
		message, err = this.receive(source)
	} else {
		message, err = this.receiveReliable(source, procQ)
	}
	var m Message
	err = json.Unmarshal([]byte(message), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (this *RedisCluster) receive(source Source) (string, error) {
	ret := this.client.BRPop(BLOCK_FOR_DURATION, source.String())
	err := ret.Err()
	if err != nil && err == redis.Nil {
		//we got an error
		return "", ErrEmptyQueue
	}
	if err != nil {
		return "", err
	}
	sliceRes := ret.Val()
	if len(sliceRes) == 2 { //check if its what we expected.
		return sliceRes[1], nil
	}
	return "", fmt.Errorf("An unexpected error occured while fetching message from Q: %s", source)
}

func (this *RedisCluster) receiveReliable(source Source, procQ Q) (string, error) {
	ret := this.client.BRPopLPush(source.String(), procQ.String(), BLOCK_FOR_DURATION)
	err := ret.Err()
	if err != nil && err == redis.Nil {
		//we got an error
		return "", ErrEmptyQueue
	}
	if err != nil {
		return "", err
	}
	sliceRes := ret.Val()
	return sliceRes, nil
}

func (this *RedisCluster) MarkProcessed(m *Message, procQ Q) error {

	if procQ.IsEmpty() {
		return nil
	}
	return FailSafeExec(func() error {
		ret := this.client.RPop(procQ.String())
		err := ret.Err()
		if err != nil && err != redis.Nil {
			return err
		}
		return nil
	}, MAX_TRY_LIMIT)
}

func (this *RedisCluster) MarkFailed(m *Message, deadQ Q, processingQ Q) error {

	if m == nil || (deadQ.IsEmpty() && processingQ.IsEmpty()) {
		return nil //nothing to do
	}

	// Simply remove from processing Q
	if deadQ.IsEmpty() {
		return FailSafeExec(func() error {
			ret := this.client.RPop(processingQ.String())
			err := ret.Err()
			if err != nil && err != redis.Nil {
				return err
			}
			return nil
		}, MAX_TRY_LIMIT)
	}

	// Simply put in deadQ
	if processingQ.IsEmpty() {
		return FailSafeExec(func() error {
			ret := this.client.LPush(deadQ.String(), m.String())
			err := ret.Err()
			if err != nil && err != redis.Nil {
				return err
			}
			return nil
		}, MAX_TRY_LIMIT)
	}

	// else remove from processingQ and put in deadQ
	return FailSafeExec(func() error {
		ret := this.client.RPopLPush(processingQ.String(), deadQ.String())
		err := ret.Err()
		if err != nil && err != redis.Nil {
			return err
		}
		return nil
	}, MAX_TRY_LIMIT)
}
