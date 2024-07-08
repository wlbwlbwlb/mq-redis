package mq

import (
	"github.com/gomodule/redigo/redis"
)

type Options struct {
	pool *redis.Pool
}

type OptionFunc func(*Options)
