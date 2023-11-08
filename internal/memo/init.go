package memo

import (
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/dgraph-io/ristretto"
	"sync"
	"time"
)

type cache struct {
	c      *ristretto.Cache
	config *ristretto.Config
	logger log.Logger
	lck    sync.Mutex
}

var instance = &cache{}

func (c *cache) init() (err error) {
	if c.c == nil {
		c.lck.Lock()
		defer c.lck.Unlock()
		if c.c == nil {
			instance.c, err = ristretto.NewCache(c.config)
		}
	}
	return
}
func Start(logger log.Logger, config *config.Config) (err error) {
	if config.MemoCache == nil {
		panic("memo cache config error")
	}

	instance.config = &ristretto.Config{
		NumCounters:        config.MemoCache.NumCounters,
		MaxCost:            config.MemoCache.MaxCost,
		BufferItems:        config.MemoCache.BufferItems,
		Metrics:            config.MemoCache.Metrics,
		IgnoreInternalCost: config.MemoCache.IgnoreInternalCost,
	}
	instance.logger = logger
	return instance.init()
}

func Get(key interface{}) (interface{}, bool) {
	instance.init()
	return instance.c.Get(key)
}

func Set(key, value interface{}, cost int64) bool {
	instance.init()
	return instance.c.Set(key, value, cost)
}

func SetWithTTL(key, value interface{}, cost int64, ttl time.Duration) bool {
	instance.init()
	return instance.c.SetWithTTL(key, value, cost, ttl)
}

func Del(key interface{}) {
	instance.init()
	instance.c.Del(key)
}
