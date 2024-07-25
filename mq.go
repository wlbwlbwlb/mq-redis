package mq

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
)

var opt Options

func Init(pool *redis.Pool, opts ...OptionFunc) (e error) {
	opt = Options{
		pool: pool,
	}
	for _, fn := range opts {
		fn(&opt)
	}
	go func() {
		for {
			psc := redis.PubSubConn{
				Conn: opt.pool.Get(),
			}
			stop := func() {
				//取消订阅
				psc.Unsubscribe()
				//关闭连接
				psc.Close()
			}

			for channel, _ := range handlers {
				if e = psc.Subscribe(channel); e != nil {
					break //太多了，不写日志
				}
			}
			if e != nil {
				stop()
				continue //订阅了一部分
			}

			for {
				switch recv := psc.Receive().(type) {
				case redis.Message:
					if opt.logger != nil {
						opt.logger.Printf("channel=%s, pattern=%s, data=%s\n", recv.Channel, recv.Pattern, recv.Data)
					}
					if fn, ok := getHandler(recv.Channel); ok {
						fn(recv)
					}
				case redis.Subscription:
					if opt.logger != nil {
						opt.logger.Printf("kind=%s, channel=%s, count=%d\n", recv.Kind, recv.Channel, recv.Count)
					}
				case error:
					if opt.logger != nil {
						opt.logger.Printf("%+v\n", recv)
					}
					e = recv
				}
				if e != nil {
					break //出错了
				}
			}
			stop()
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

func getHandler(channel string) (fn HandlerFunc, ok bool) {
	fn, ok = handlers[channel]
	return
}

var handlers = make(map[string]HandlerFunc)

type HandlerFunc func(message redis.Message) error
