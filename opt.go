package mq

import (
	"io"
	"log"

	"github.com/gomodule/redigo/redis"
)

type Options struct {
	pool   *redis.Pool
	logger *log.Logger
}

type OptionFunc func(*Options)

func Writer(writer io.Writer) OptionFunc {
	return func(opts *Options) {
		opts.logger = log.New(writer, "", log.Flags())
	}
}
