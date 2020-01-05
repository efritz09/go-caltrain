package caltrain

import (
	"sync"
	"time"

	"github.com/benbjohnson/clock"
)

const (
	DefaultCacheTimeout = 5 * time.Minute
)

type Cache interface {
	set(key string, body []byte)
	get(key string) ([]byte, bool)
	clearCache()
}

type CaltrainCache struct {
	cache   map[string]cacheData // map of endpoint to body data and timestamp
	timeout time.Duration
	lock    sync.RWMutex

	clock clock.Clock // time package for unit testing
}

type cacheData struct {
	body       []byte
	expiration int64
}

func NewCache(expire time.Duration) *CaltrainCache {
	return &CaltrainCache{
		cache:   make(map[string]cacheData),
		timeout: expire,
		clock:   clock.New(),
	}
}

// set calculates the key's expiration and sets the cacheData for that key
func (c *CaltrainCache) set(key string, body []byte) {
	exp := c.clock.Now().Add(c.timeout).UnixNano()
	c.lock.Lock()
	c.cache[key] = cacheData{
		body:       body,
		expiration: exp,
	}
	c.lock.Unlock()
}

// get will query the cache for an endpoint. if the endpoint exists, it will
// check if the cache has expired. If not, it returns the value and true. if it
// has expired, it will return false
func (c *CaltrainCache) get(key string) ([]byte, bool) {
	c.lock.RLock()
	data, ok := c.cache[key]
	if !ok {
		c.lock.RUnlock()
		return nil, false
	}

	if c.clock.Now().UnixNano() > data.expiration {
		c.lock.RUnlock()
		return nil, false
	}
	c.lock.RUnlock()
	return data.body, true
}

// clearCache clears the cache by creating a new cache map
func (c *CaltrainCache) clearCache() {
	c.lock.Lock()
	c.cache = make(map[string]cacheData)
	c.lock.Unlock()
}

type MockCache struct {
	SetFunc func(string, []byte)
	GetFunc func(string) ([]byte, bool)
}

func (c *MockCache) set(key string, body []byte) {
	if c.SetFunc != nil {
		c.SetFunc(key, body)
	}
}

func (c *MockCache) get(key string) ([]byte, bool) {
	if c.GetFunc != nil {
		return c.GetFunc(key)
	}
	return nil, false
}

func (C *MockCache) clearCache() {}
