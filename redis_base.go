package raven

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

//var _ RedisClient = (*RedisSimpleClient)(nil)
var _ RedisClient = (*RedisClusterClient)(nil)

type RedisClient interface {
	LPush(string, ...interface{}) *redis.IntCmd
	BRPop(time.Duration, ...string) *redis.StringSliceCmd
	BRPopLPush(string, string, time.Duration) *redis.StringCmd
	RPop(string) *redis.StringCmd
	RPopLPush(string, string) *redis.StringCmd
}

type RedisSimpleClient struct {
	*redis.Client
}

type RedisClusterClient struct {
	*redis.ClusterClient
}

//
// A Base client to be implemented by redis and redis cluster.
//
type redisbase struct {
	Client RedisClient
}

//
//  Implementation of Send() method exposed by raven manager.
//
func (this *redisbase) Send(message Message, dest Destination) error {

	ret := this.Client.LPush(dest.GetName(), message.toJson())
	if ret.Err() != nil {
		return ret.Err()
	}
	return nil
}

func (this *redisbase) Receive(source Source, procQ Q) (*Message, error) {

	var message string
	var err error
	if procQ.IsEmpty() {
		message, err = this.receive(source)
	} else {
		message, err = this.receiveReliable(source, procQ)
	}
	if err != nil {
		return nil, err
	}
	var m Message
	err = json.Unmarshal([]byte(message), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (this *redisbase) receive(source Source) (string, error) {
	ret := this.Client.BRPop(BLOCK_FOR_DURATION, source.GetName())
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

func (this *redisbase) receiveReliable(source Source, procQ Q) (string, error) {
	ret := this.Client.BRPopLPush(source.GetName(), procQ.GetName(), BLOCK_FOR_DURATION)
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

func (this *redisbase) MarkProcessed(m *Message, procQ Q) error {

	if procQ.IsEmpty() {
		return nil
	}
	return failSafeExec(func() error {
		ret := this.Client.RPop(procQ.GetName())
		err := ret.Err()
		if err != nil && err != redis.Nil {
			return err
		}
		return nil
	}, MAX_TRY_LIMIT)
}

func (this *redisbase) MarkFailed(m *Message, deadQ Q, processingQ Q) error {

	if m == nil || (deadQ.IsEmpty() && processingQ.IsEmpty()) {
		return nil //nothing to do
	}

	// Simply remove from processing Q
	if deadQ.IsEmpty() {
		return failSafeExec(func() error {
			ret := this.Client.RPop(processingQ.GetName())
			err := ret.Err()
			if err != nil && err != redis.Nil {
				return err
			}
			return nil
		}, MAX_TRY_LIMIT)
	}

	// Simply put in deadQ
	if processingQ.IsEmpty() {
		return failSafeExec(func() error {
			ret := this.Client.LPush(deadQ.GetName(), m.String())
			err := ret.Err()
			if err != nil && err != redis.Nil {
				return err
			}
			return nil
		}, MAX_TRY_LIMIT)
	}

	// else remove from processingQ and put in deadQ
	return failSafeExec(func() error {
		ret := this.Client.RPopLPush(processingQ.GetName(), deadQ.GetName())
		err := ret.Err()
		if err != nil && err != redis.Nil {
			return err
		}
		return nil
	}, MAX_TRY_LIMIT)
}
