package cache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLRUCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewLRUCache[string, int](10)
		require.Equal(t, 0, c.list.Len())
		_, ok := c.Get("a")
		require.Equal(t, false, ok)
		c.Put("b", 1)
		require.Equal(t, 1, c.list.Len())
	})

	t.Run("existing cache", func(t *testing.T) {
		c := NewLRUCache[string, int](3)
		c.Put("aaa", 100)
		c.Put("bbb", 200)
		c.Put("ccc", 300)

		_, ok := c.Get("ddd")
		require.False(t, ok)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		c.Put("ddd", 400)

		_, ok = c.Get("bbb")
		require.False(t, ok)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})
}
func TestCacheMultithreading(t *testing.T) {

	c := NewLRUCache[string, int](10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Put(strconv.Itoa(i), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(strconv.Itoa(rand.Intn(1_000_000)))
		}
	}()

	wg.Wait()
}
