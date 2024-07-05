package mq

import (
	"github.com/gomodule/redigo/redis"
)

type Options struct {
	pool *redis.Pool
	//conn redis.Conn
}

type OptionFunc func(*Options)

func Pool(pool *redis.Pool) OptionFunc {
	return func(opt *Options) {
		opt.pool = pool
	}
}
