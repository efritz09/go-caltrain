package caltrain

import (
	"sync"
	"time"

	"github.com/benbjohnson/clock"
)

const (
	defaultCacheTimeout = 5 * time.Minute
)

type cache interface {
	set(key string, body []byte)
	get(key string) ([]byte, time.Time, bool)
	clearCache()
}

type caltrainCache struct {
	cache   map[string]cacheData // map of endpoint to body data and timestamp
	timeout time.Duration
	lock    sync.RWMutex

	clock clock.Clock // time package for unit testing
}

type cacheData struct {
	body       []byte    // response body
	entryTime  time.Time // time that 'body' was stored
	expiration int64     // expiration time
}

func newCache(expire time.Duration) *caltrainCache {
	return &caltrainCache{
		cache:   make(map[string]cacheData),
		timeout: expire,
		clock:   clock.New(),
	}
}

// set calculates the key's expiration and sets the cacheData for that key
func (c *caltrainCache) set(key string, body []byte) {
	exp := c.clock.Now().Add(c.timeout).UnixNano()
	c.lock.Lock()
	c.cache[key] = cacheData{
		body:       body,
		entryTime:  time.Now(),
		expiration: exp,
	}
	c.lock.Unlock()
}

// get will query the cache for an endpoint. if the endpoint exists, it will
// check if the cache has expired. If not, it returnstrue. if it has expired,
// it will return false. It always returns the body and entryTime for use in
// case the API limit has been reached
func (c *caltrainCache) get(key string) ([]byte, time.Time, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	data, ok := c.cache[key]
	var t time.Time
	if !ok {
		return nil, t, false
	}

	if c.clock.Now().UnixNano() > data.expiration {
		return data.body, data.entryTime, false
	}
	return data.body, data.entryTime, true
}

// clearCache clears the cache by creating a new cache map
func (c *caltrainCache) clearCache() {
	c.lock.Lock()
	c.cache = make(map[string]cacheData)
	c.lock.Unlock()
}

type mockCache struct {
	SetFunc func(string, []byte)
	GetFunc func(string) ([]byte, time.Time, bool)
}

func (c *mockCache) set(key string, body []byte) {
	if c.SetFunc != nil {
		c.SetFunc(key, body)
	}
}

func (c *mockCache) get(key string) ([]byte, time.Time, bool) {
	if c.GetFunc != nil {
		return c.GetFunc(key)
	}
	return nil, time.Now(), false
}

func (C *mockCache) clearCache() {}
