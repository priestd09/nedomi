package storage

import (
	"github.com/ironsmile/nedomi/cache"
	"github.com/ironsmile/nedomi/config"
	"github.com/ironsmile/nedomi/upstream"

	"github.com/ironsmile/nedomi/storage/disk"
)

type newStorageFunc func(cfg config.CacheZoneSection, cm cache.CacheManager,
	up upstream.Upstream) Storage

var storageTypes map[string]newStorageFunc = map[string]newStorageFunc{

	"disk": func(cfg config.CacheZoneSection, cm cache.CacheManager,
		up upstream.Upstream) Storage {
		return disk.New(cfg, cm, up)
	},
}
