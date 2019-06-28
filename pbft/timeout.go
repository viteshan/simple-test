package pbft

import (
	"time"

	"github.com/hashicorp/golang-lru"
	"github.com/vitelabs/go-vite/log15"
)

// Blacklist define a data set interface with a timeout
type Blacklist interface {
	Add(key interface{})
	AddAddTimeout(key interface{}, duration time.Duration)
	Exists(key interface{}) bool
	Remove(key interface{})
}

// NewBlacklist returns timeout blacklist
func NewBlacklist() (Blacklist, error) {
	cache, err := lru.New(10 * 10000)
	if err != nil {
		return nil, err
	}
	return &blacklist{cache: cache, defaultTimeout: 0, log: log15.New("module", "pool/blacklist")}, nil
}

type timeout struct {
	timeoutT *time.Time
}

func (tt *timeout) reset(duration time.Duration) *timeout {
	if duration <= 0 {
		tt.timeoutT = nil
	} else {
		t := time.Now().Add(duration)
		tt.timeoutT = &t
	}

	return tt
}
func (tt *timeout) isTimeout() bool {
	if tt.timeoutT == nil {
		return false
	}
	return tt.timeoutT.Before(time.Now())
}

type blacklist struct {
	cache          *lru.Cache
	defaultTimeout time.Duration
	log            log15.Logger
}

func (bl *blacklist) Add(key interface{}) {
	bl.AddAddTimeout(key, bl.defaultTimeout)
}

func (bl *blacklist) AddAddTimeout(key interface{}, duration time.Duration) {
	value, ok := bl.cache.Get(key)
	if ok {
		value.(*timeout).reset(duration)
	} else {
		bl.cache.Add(key, (&timeout{}).reset(duration))
	}
}

func (bl *blacklist) Exists(key interface{}) bool {
	value, ok := bl.cache.Get(key)
	if ok {
		return !value.(*timeout).isTimeout()
	}
	return false
}

func (bl *blacklist) Remove(key interface{}) {
	bl.cache.Remove(key)
}
