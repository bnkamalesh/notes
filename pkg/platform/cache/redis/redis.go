package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

var (
	// ErrInvHosts is returned when the provided host(s) is/are invalid
	ErrInvHosts = errors.New("Invalid hosts provided")
)

// Config struct has all the configurations required for redis
type Config struct {
	Hosts           []string
	DB              int
	Password        string
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PoolSize           int
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

// Handler struct does all the cache operations
type Handler struct {
	ring  *redis.Ring
	codec *cache.Codec
}

// Set saves a new value in Redis with the given key, value and expiry
func (h *Handler) Set(key string, value interface{}, expiry time.Duration) error {
	return h.codec.Set(&cache.Item{
		Key:        key,
		Object:     value,
		Expiration: expiry,
	})
}

// Get loads the value of the given key, from Redis to result
func (h *Handler) Get(key string, result interface{}) error {
	return h.codec.Get(key, result)
}

// Ping pings the redis server
func (h *Handler) Ping() error {
	result := h.ring.Ping()
	fmt.Println("ping result:", result)
	return nil
}

// New returns a handler instance with all the required attributes initialized
func New(c Config) (*Handler, error) {
	if len(c.Hosts) == 0 {
		return nil, ErrInvHosts
	}
	hosts := make(map[string]string, len(c.Hosts))
	for i, h := range c.Hosts {
		hosts[fmt.Sprintf("server-%d", i)] = h
	}

	redisRing := redis.NewRing(
		&redis.RingOptions{
			Addrs:              hosts,
			DB:                 c.DB,
			Password:           c.Password,
			MaxRetries:         c.MaxRetries,
			MinRetryBackoff:    c.MinRetryBackoff,
			MaxRetryBackoff:    c.MaxRetryBackoff,
			DialTimeout:        c.DialTimeout,
			ReadTimeout:        c.ReadTimeout,
			WriteTimeout:       c.WriteTimeout,
			PoolSize:           c.PoolSize,
			PoolTimeout:        c.PoolTimeout,
			IdleCheckFrequency: c.IdleCheckFrequency,
			IdleTimeout:        c.IdleTimeout,
		},
	)

	h := &Handler{
		ring: redisRing,
		codec: &cache.Codec{
			Redis: redisRing,
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		},
	}
	return h, nil
}
