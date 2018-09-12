package childlock

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

var ERR_LOCK_BUSY error = errors.New("lock busy")
var ERR_NOT_MINE_LOCK error = errors.New("not my lock")

type RedisOptions struct {
	Addres     []string
	MaxRetries int
	PoolSize   int
}

func New(name string, expiry int, options RedisOptions) *Lock {

	opt := redis.UniversalOptions{
		Addrs:      options.Addres,
		MaxRetries: options.MaxRetries,
		PoolSize:   options.PoolSize,
	}
	return &Lock{
		name:        name,
		expire:      time.Duration(expiry) * time.Second,
		redisClient: redis.NewUniversalClient(&opt),
	}
}

type Lock struct {
	name        string
	value       string
	expire      time.Duration
	redisClient redis.UniversalClient
}

func (this *Lock) Acquire(val string) error {
	ok, err := this.redisClient.SetNX(this.name, val, this.expire).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ERR_LOCK_BUSY
	}
	this.value = val
	return nil
}

func (this *Lock) Refresh() error {
	res, err := this.redisClient.Get(this.name).Result()
	if err != nil {
		return err
	}
	if this.value != res {
		return ERR_NOT_MINE_LOCK
	}
	err = this.redisClient.Watch(func(tx *redis.Tx) error {
		ok, err := tx.Expire(this.name, this.expire).Result()
		if err != nil {
			return err
		}
		if !ok {
			return ERR_NOT_MINE_LOCK
		}
		return nil
	}, this.name)
	return nil
}

func (this *Lock) Release() error {
	res, err := this.redisClient.Get(this.name).Result()
	if err != nil {
		return err
	}
	if this.value != res {
		return ERR_NOT_MINE_LOCK
	}
	err = this.redisClient.Watch(func(tx *redis.Tx) error {
		res, err := tx.Del(this.name).Result()
		if err != nil {
			return err
		}
		if res == 0 {
			return ERR_NOT_MINE_LOCK
		}
		return nil
	}, this.name)

	return err
}
