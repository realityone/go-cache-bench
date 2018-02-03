package cache_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/bluele/gcache"
	"github.com/golang/groupcache"
	"github.com/muesli/cache2go"
)

const (
	LRU_SIZE = 1024
)

var (
	keyGen func() string

	gcOnce    sync.Once
	groupOnce sync.Once
	c2goOnce  sync.Once

	gc   gcache.Cache
	sg   *groupcache.Group
	c2go *cache2go.CacheTable
)

func init() {
	method := os.Getenv("KEY_METHOD")
	if method == "" {
		method = "random"
	}
	fmt.Printf("Testing with `%s` policy\n", method)

	keyMethod := randomKey
	if method == "scheduled" {
		keyMethod = scheduledKey
	}
	init, gen := keyMethod()
	init(LRU_SIZE)
	keyGen = gen
}

func BenchmarkgcacheLRU(b *testing.B) {
	gcOnce.Do(func() {
		gc = gcache.New(LRU_SIZE).
			LRU().
			LoaderFunc(loader).
			Build()
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keyGen()
			value := vValue(key)
			value1, err := gc.Get(key)
			if err != nil {
				panic(err)
			}
			if value.(string) != value1.(string) {
				panic("value is not expected")
			}
		}
	})
}

func BenchmarkgroupcacheStringGroup(b *testing.B) {
	groupOnce.Do(func() {
		sg = groupcache.NewGroup("string-group", LRU_SIZE, groupcache.GetterFunc(func(_ groupcache.Context, key string, dest groupcache.Sink) error {
			return dest.SetString(vValue(key).(string))
		}))
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var s string
			key := keyGen()
			value := vValue(key)
			if err := sg.Get(context.TODO(), key, groupcache.StringSink(&s)); err != nil {
				panic(err)
			}
			if value.(string) != string(s) {
				panic("value is not expected")
			}
		}
	})
}

func Benchmarkcache2goExpire(b *testing.B) {
	c2goOnce.Do(func() {
		c2go = cache2go.Cache("cache2go-bench")
		c2go.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
			return cache2go.NewCacheItem(key, 15*time.Millisecond, vValue(key))
		})
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keyGen()
			value := vValue(key)
			value1, err := c2go.Value(key)
			if err != nil {
				panic(err)
			}
			if value.(string) != value1.Data().(string) {
				panic("value is not expected")
			}
		}
	})
}
