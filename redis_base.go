package raven

import (
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

	fmt.Printf("Publishing to Q [%s]\n\n", dest.GetName())
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
		fmt.Printf("Receiving from Q [%s] Non Reliable.\n", source.GetName())
		message, err = this.receive(source)
	} else {
		fmt.Printf("Receiving from Q [%s] Reliable.\n", source.GetName())
		message, err = this.receiveReliable(source, procQ)
	}
	if err != nil {
		return nil, err
	}
	var m *Message = new(Message)
	err = m.fromJson(message)
	fmt.Printf("Message is: %+v\n", m)
	return m, nil
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
	fmt.Printf("%v\n", sliceRes)
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
	fmt.Printf("Marking message [%d] processed.\n", m.Id)
	return failSafeExec(func() error { //@todo: use ltrim instead of rpop.
		//to make sure no previous message remains.
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

	fmt.Printf("DeadQ is [%s], ProcessingQ is [%s]\n", deadQ.GetName(), processingQ.GetName())
	// Simply remove from processing Q
	if deadQ.IsEmpty() {
		fmt.Printf("Marking message dead [%d] . Simply remove from processing Q.\n", m.Id)
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
		fmt.Printf("Marking message dead [%s] . Simply put in dead Q.\n", m.Id)
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
	fmt.Printf("Marking message dead [%s] . Fully functional.\n", m.Id)
	return failSafeExec(func() error {
		ret := this.Client.RPopLPush(processingQ.GetName(), deadQ.GetName())
		err := ret.Err()
		if err != nil && err != redis.Nil {
			return err
		}
		return nil
	}, MAX_TRY_LIMIT)
}
