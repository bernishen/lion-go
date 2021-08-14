package redisservice

import (
	"github.com/bernishen/lion-go/utils/exception"
	"time"

	"github.com/garyburd/redigo/redis"
)

var servers chan *redis.Conn

func init() {
	servers = make(chan *redis.Conn, 1)
	//option := redis.DialPassword("password")
	r, err := redis.Dial("tcp", "redis.dev.mine:6379", redis.DialConnectTimeout(time.Second*3))
	if err != nil {
		msg := "A error occourred connected redis server.[" + err.Error() + "]"
		exception.NewException(exception.Error, 1001, msg)
		return
	}
	servers <- &r
}

// Set : Setting a key|value into the redis.
func Set(key string, value string, expire time.Duration) (bool, *exception.Exception) {
	select {
	case server := <-servers:
		defer func() { servers <- server }()
		_, err := (*server).Do("SET", key, value)
		_, err1 := (*server).Do("EXPIRE", key, expire.Seconds())
		if err != nil {
			return false, exception.NewException(exception.Error, 1002, "had an exception occurred saving cache.["+err.Error()+"]")
		}
		if err1 != nil {
			return true, exception.NewException(exception.Warning, 1003, "Seted expire time error.["+err1.Error()+"]")
		}
		return true, nil
	case <-time.After(time.Second * 1):
		return false, exception.NewException(exception.Warning, 1001, "time out.")
	}
}

// Get : Getting a value  by the key.
func Get(key string) (string, *exception.Exception) {
	select {
	case server := <-servers:
		defer func() { servers <- server }()
		value, err := redis.String((*server).Do("GET", key))
		if err != nil {
			return "", exception.NewException(exception.Error, 1002, "had an exception occurred saving cache.["+err.Error()+"]")
		}

		return value, nil
	case <-time.After(time.Second * 1):
		return "", exception.NewException(exception.Warning, 1001, "time out.")
	}
}

// RefreshExpire : Set and update the key expire.
func RefreshExpire(key string, expire time.Duration) (bool, *exception.Exception) {
	select {
	case server := <-servers:
		defer func() { servers <- server }()
		_, err := (*server).Do("EXPIRE", key, expire.Seconds())
		if err != nil {
			return false, exception.NewException(exception.Error, 1002, err.Error())
		}
		return true, nil
	case <-time.After(time.Second * 5):
		return false, exception.NewException(exception.Warning, 1001, "time out.")
	}
}

// Exits :
func Exits(key string) (bool, *exception.Exception) {
	select {
	case server := <-servers:
		defer func() { servers <- server }()
		exist, err := (*server).Do("EXISTS", key)
		if err != nil {
			return false, exception.NewException(exception.Error, 1002, err.Error())
		}
		ret, err := redis.Bool(exist, err)
		if err != nil {
			return false, exception.NewException(exception.Error, 1002, err.Error())
		}
		return ret, nil
	case <-time.After(time.Second * 5):
		return false, exception.NewException(exception.Warning, 1001, "time out.")
	}
}
