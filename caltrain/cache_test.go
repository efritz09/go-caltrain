package caltrain

import (
	"bytes"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

var d1 = []byte{100, 49, 95, 101, 114, 105, 99}
var d2 = []byte{100, 50, 95, 105, 115}
var d3 = []byte{100, 51, 95, 103, 114, 101, 97, 116}

func TestCache(t *testing.T) {
	c := newCache(defaultCacheTimeout)

	v1, _, ok := c.get("a")
	if ok || v1 != nil {
		t.Errorf("Getting A found value that shouldn't exist: %s", v1)
	}

	v2, _, ok := c.get("b")
	if ok || v2 != nil {
		t.Errorf("Getting B found value that shouldn't exist: %s", v2)
	}

	v3, _, ok := c.get("c")
	if ok || v3 != nil {
		t.Errorf("Getting C found value that shouldn't exist: %s", v3)
	}

	c.set("a", d1)
	c.set("b", d2)
	c.set("c", d3)

	if v1, _, ok := c.get("a"); !ok {
		t.Error("'a' was not found")
	} else if !bytes.Equal(v1, d1) {
		t.Errorf("unexpected value for 'a': Expected %b, received %b", d1, v1)
	}

	if v2, _, ok := c.get("b"); !ok {
		t.Error("'b' was not found")
	} else if !bytes.Equal(v2, d2) {
		t.Errorf("unexpected value for 'b': Expected %b, received %b", d2, v2)
	}

	if v3, _, ok := c.get("c"); !ok {
		t.Error("'c' was not found")
	} else if !bytes.Equal(v3, d3) {
		t.Errorf("unexpected value for 'c': Expected %b, received %b", d3, v3)
	}
}

func TestExpiration(t *testing.T) {
	c := newCache(defaultCacheTimeout)
	mock := clock.NewMock()
	c.clock = mock

	c.set("a", d1)
	c.set("b", d2)

	// increment the time to less than the default
	mock.Add(defaultCacheTimeout - 3*time.Second)
	if v1, _, ok := c.get("a"); !ok {
		t.Error("'a' was not found")
	} else if !bytes.Equal(v1, d1) {
		t.Errorf("unexpected value for 'a': Expected %b, received %b", d1, v1)
	}
	if v2, _, ok := c.get("b"); !ok {
		t.Error("'b' was not found")
	} else if !bytes.Equal(v2, d2) {
		t.Errorf("unexpected value for 'b': Expected %b, received %b", d2, v2)
	}

	mock.Add(defaultCacheTimeout)

	if v1, _, ok := c.get("a"); ok {
		t.Error("'a' has not timed out")
	} else if !bytes.Equal(v1, d1) {
		t.Errorf("unexpected value for 'a': Expected %b, received %b", d1, v1)
	}
	if v2, _, ok := c.get("b"); ok {
		t.Error("'b' has not timed out")
	} else if !bytes.Equal(v2, d2) {
		t.Errorf("unexpected value for 'b': Expected %b, received %b", d2, v2)
	}
}

func TestReplacement(t *testing.T) {
	c := newCache(defaultCacheTimeout)
	mock := clock.NewMock()
	c.clock = mock

	c.set("a", d1)
	c.set("b", d2)

	// increment the time to less than the default
	mock.Add(defaultCacheTimeout - 3*time.Second)
	if v1, _, ok := c.get("a"); !ok {
		t.Error("'a' was not found")
	} else if !bytes.Equal(v1, d1) {
		t.Errorf("unexpected value for 'a': Expected %b, received %b", d1, v1)
	}
	if v2, _, ok := c.get("b"); !ok {
		t.Error("'b' was not found")
	} else if !bytes.Equal(v2, d2) {
		t.Errorf("unexpected value for 'b': Expected %b, received %b", d2, v2)
	}

	// now replace b with a new value
	c.set("b", d3)

	// expire 'a', but not 'b'
	mock.Add(defaultCacheTimeout)
	if v1, _, ok := c.get("a"); ok {
		t.Error("'a' has not timed out")
	} else if !bytes.Equal(v1, d1) {
		t.Errorf("unexpected value for 'a': Expected %b, received %b", d1, v1)
	}
	if v2, _, ok := c.get("b"); !ok {
		t.Error("'b' has timed out")
	} else if !bytes.Equal(v2, d3) {
		t.Errorf("unexpected value for 'b': Expected %b, received %b", d3, v2)
	}
}

func TestClearCache(t *testing.T) {
	c := newCache(defaultCacheTimeout)

	c.set("a", d1)
	c.set("b", d2)
	c.set("c", d2)

	if len(c.cache) != 3 {
		t.Fatalf("cache length is not 3: %d", len(c.cache))
	}

	c.clearCache()

	if len(c.cache) != 0 {
		t.Fatalf("cache is not empty! %v", c.cache)
	}
}
