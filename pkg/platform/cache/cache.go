package cache

import (
	"errors"
	"strconv"
	"time"

	"github.com/bnkamalesh/notes/pkg/platform/cache/redis"
)

// ErrNotFound is the error encountered when key not found in the cache
var ErrNotFound = errors.New("Key does not exist")

// ErrCacheNoHost is the error encountered when the provided Redis/Hosts are not available
var ErrCacheNoHost = errors.New("No valid host address(es) provided")

// ErrCacheNoHn is the error encountered when there is no valid handler initialized and
// still calling the cache methods
var ErrCacheNoHn = errors.New("No cache handler initialized")

// Service interface should have all the required methods
// This interface is implemented to remove isCluster check for every call
type Service interface {
	Set(key string, value interface{}, expiry time.Duration) error
	Get(key string, result interface{}) error
	// HSet(string, string, interface{}, time.Duration, bool) error
	// HGet(string, string, interface{}) (error)
	// Delete(...string) error
	// HDelete(string, ...string) error
	Ping() error
}

type Config struct {
	Hosts        []string
	Name         string
	Password     string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Handler struct {
	client Service
}

func (h *Handler) Set(key string, value interface{}, expiry time.Duration) error {
	return h.client.Set(key, value, expiry)
}

func (h *Handler) Get(key string, result interface{}) error {
	return h.client.Get(key, result)
}

func (h *Handler) Ping() error {
	return h.client.Ping()
}

func New(c Config) (*Handler, error) {
	h := &Handler{}
	db, _ := strconv.Atoi(c.Name)
	rh, err := redis.New(redis.Config{
		Hosts:        c.Hosts,
		DB:           db,
		Password:     c.Password,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.DialTimeout,
	})
	if err != nil {
		return nil, err
	}
	err = rh.Ping()
	if err != nil {
		return nil, err
	}

	h.client = rh
	return h, nil
}
