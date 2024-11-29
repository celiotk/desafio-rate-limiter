package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Função para conectar ao Redis
func NewRedisPool(redisAddress string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddress)
		},
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
	}
}

type RedisStorage struct {
	pool *redis.Pool
}

func NewRedisStorage(pool *redis.Pool) *RedisStorage {
	return &RedisStorage{pool: pool}
}

func (r *RedisStorage) Increment(key string) (int, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", key))
}

func (r *RedisStorage) SetExpiration(key string, ttl int) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, ttl)
	return err
}

func (r *RedisStorage) GetTTL(key string) (int, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("TTL", key))
}

func (r *RedisStorage) Exists(key string) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", key))
}

func (r *RedisStorage) Block(key string, duration int) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, duration, 1)
	return err
}

func (r *RedisStorage) SetTokenRateConfig(token string, limit int, interval int) error {
	conn := r.pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("config:token:%s", token)
	_, err := conn.Do("HMSET", key, "limit", limit, "interval", interval)
	return err
}

func (r *RedisStorage) GetTokenRateConfig(key string) (limit int, interval int, err error) {
	conn := r.pool.Get()
	defer conn.Close()

	key = fmt.Sprintf("config:token:%s", key)
	values, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		return 0, 0, err
	}

	var limitStr, intervalStr string
	for i := 0; i < len(values); i += 2 {
		switch string(values[i].([]byte)) {
		case "limit":
			limitStr = string(values[i+1].([]byte))
		case "interval":
			intervalStr = string(values[i+1].([]byte))
		}
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, err
	}

	interval, err = strconv.Atoi(intervalStr)
	if err != nil {
		return 0, 0, err
	}

	return limit, interval, nil
}
