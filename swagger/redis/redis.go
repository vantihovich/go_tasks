package redis

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	conn redis.Conn
}

func New(connType string, address string) (Redis, error) {
	conn, err := redis.Dial(connType, address)
	if err != nil {
		log.WithError(err).Fatal("Failed to establish connection with Redis")
		return Redis{
			conn: nil,
		}, err
	}

	return Redis{
		conn: conn,
	}, nil
}

func (r Redis) Get(login string) (int, bool, error) {
	cached, err := redis.Int(r.conn.Do("Get", login))
	if err != nil {
		if err == redis.ErrNil {
			return 0, false, nil
		}
		log.WithError(err).Error("Unable to get cached data")
		return 0, false, err
	}
	return cached, true, nil
}

func (r Redis) Set(login string, attempts int, ttl int) (err error) {
	_, err = r.conn.Do("SetEX", login, ttl, attempts)
	if err != nil {
		log.WithError(err).Error("Unable to set cache data")
		return err
	}
	return nil
}
