package mq

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var opt Options

func Init(opts ...Option) (e error) {
	for _, f := range opts {
		f(&opt)
	}
	psc := redis.PubSubConn{
		Conn: opt.pool.Get(),
	}
	for channel, _ := range handlers {
		if e = psc.Subscribe(channel); e != nil {
			return
		}
	}
	go func() {
		defer func() {
			psc.Unsubscribe()
		}()
		for {
			switch msg := psc.Receive().(type) {
			case redis.Message:
				//fmt.Printf("channel=%s, pattern=%s, data=%s\n", msg.Channel, msg.Pattern, msg.Data)
				if handle, ok := getHandler(msg.Channel); ok {
					handle(msg)
				}
			case redis.Subscription:
				fmt.Printf("kind=%s, channel=%s, count=%d\n", msg.Kind, msg.Channel, msg.Count)
			case error:
				fmt.Printf("%+v\n", msg)
				//return
			}
		}
	}()
	return
}

func Pub(channel string, message json.RawMessage) (e error) {
	if nil == opt.pool {
		return errors.New("init first")
	}
	conn := opt.pool.Get()
	defer conn.Close()
	_, e = conn.Do("PUBLISH", channel, message)
	return
}

func Sub(channel string, handler HandlerFunc) (e error) {
	handlers[channel] = handler
	return
}

func getHandler(channel string) (f HandlerFunc, ok bool) {
	f, ok = handlers[channel]
	return
}

var handlers = make(map[string]HandlerFunc)

type HandlerFunc func(message redis.Message) error
