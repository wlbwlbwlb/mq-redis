package mq

import (
	"github.com/gomodule/redigo/redis"
)

type Options struct {
	pool *redis.Pool
	//conn redis.Conn
}

type Option func(*Options)

func Pool(pool *redis.Pool) Option {
	return func(opt *Options) {
		opt.pool = pool
	}
}
