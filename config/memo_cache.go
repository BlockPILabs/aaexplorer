package config

import (
	"github.com/dgraph-io/ristretto"
)

func DefaultMemoCacheConfig() *ristretto.Config {
	return &ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	}
}
