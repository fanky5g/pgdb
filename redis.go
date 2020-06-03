package database

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	MASTER_SLAVE = iota
	SLAVE_ONLY
	MASTER_ONLY
)

var (
	RedisDefaultMaxIdle     = 3
	RedisDefaultMaxActive   = 10
	RedisDefaultMaxLifetime = func() time.Duration {
		return time.Now().Sub(time.Now().Add(-time.Minute * 5))
	}
)

// RedisClient runs redis commands by running
// SET operations agains master pool and READ operations against slave pool
type RedisClient struct {
	Mode   int
	master *redis.Pool //private redis master connection pool
	slave  *redis.Pool //private redis slave connection pool
}

type RedisConnection struct {
	Master      string
	Slave       string
	MaxActive   int
	MaxLifetime time.Duration
	MaxIdle     int
}

func GetRedisClient(cfg RedisConnection) (*RedisClient, error) {
	client := &RedisClient{}
	if cfg.Master != "" {
		pool, err := newPool(cfg.Master, cfg.MaxActive, cfg.MaxLifetime, cfg.MaxIdle)
		if err != nil {
			return nil, err
		}
		client.master = pool
	}

	if cfg.Slave != "" {
		pool, err := newPool(cfg.Slave, cfg.MaxActive, cfg.MaxLifetime, cfg.MaxIdle)
		if err != nil {
			return nil, err
		}
		client.slave = pool
	}

	if client.master == nil && client.slave == nil {
		return nil, errors.New("Failed to create redis client, no master or slave connections found")
	}

	if client.master == nil {
		client.Mode = SLAVE_ONLY
	} else if client.slave == nil {
		client.Mode = MASTER_ONLY
	} else {
		client.Mode = MASTER_SLAVE
	}

	return client, nil
}

func newPool(addr string, maxActive int, maxLifetime time.Duration, maxIdle int) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxActive:       maxActive,
		MaxConnLifetime: maxLifetime,
		MaxIdle:         maxIdle,
		IdleTimeout:     240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	_, err := pool.Get().Do("PING")
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (c *RedisClient) GetConn(commandName string) redis.Conn {
	if c.Mode == SLAVE_ONLY {
		return c.slave.Get()
	} else if c.Mode == MASTER_ONLY {
		return c.master.Get()
	}

	if _, ok := readOnly[commandName]; ok {
		return c.slave.Get()
	}

	return c.master.Get()
}

var readOnly = map[string]bool{
	"hgetall":              true,
	"sdiff":                true,
	"zcard":                true,
	"exists":               true,
	"dbsize":               true,
	"zrank":                true,
	"smembers":             true,
	"psync":                true,
	"hvals":                true,
	"xrevrange":            true,
	"hkeys":                true,
	"strlen":               true,
	"hlen":                 true,
	"pttl":                 true,
	"mget":                 true,
	"zrevrangebyscore":     true,
	"sunion":               true,
	"touch":                true,
	"type":                 true,
	"lrange":               true,
	"hscan":                true,
	"xinfo":                true,
	"ttl":                  true,
	"get":                  true,
	"llen":                 true,
	"xread":                true,
	"geohash":              true,
	"object":               true,
	"sync":                 true,
	"zrange":               true,
	"bitpos":               true,
	"scard":                true,
	"zrevrangebylex":       true,
	"zscan":                true,
	"hexists":              true,
	"zcount":               true,
	"getbit":               true,
	"scan":                 true,
	"pfcount":              true,
	"psubscribe":           true,
	"geodist":              true,
	"sismember":            true,
	"xlen":                 true,
	"randomkey":            true,
	"georadiusbymember_ro": true,
	"sscan":                true,
	"hget":                 true,
	"zrangebyscore":        true,
	"zrevrank":             true,
	"xpending":             true,
	"hstrlen":              true,
	"srandmember":          true,
	"dump":                 true,
	"keys":                 true,
	"lolwut":               true,
	"sinter":               true,
	"bitcount":             true,
	"substr":               true,
	"memory":               true,
	"getrange":             true,
	"georadius_ro":         true,
	"xrange":               true,
	"zrangebylex":          true,
	"hmget":                true,
	"zscore":               true,
	"subscribe":            true,
	"unsubscribe":          true,
	"zlexcount":            true,
	"zrevrange":            true,
}
