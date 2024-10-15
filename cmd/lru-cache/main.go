package main

import (
	"fmt"

	"github.com/RacovAlex/lru-cache/internal/cache"
)

type StringIntCache interface {
	Get(key string) (int, bool)
	Put(key string, value int)
}

func main() {
	var lru StringIntCache = cache.NewLRUCache[string, int](2)

	lru.Put("a", 1)
	lru.Put("b", 2)

	if value, found := lru.Get("a"); found {
		fmt.Println("Key 'a' found with value:", value)
	} else {
		fmt.Println("Key 'a' not found")
	}

	lru.Put("d", 4)

	if _, found := lru.Get("b"); !found {
		fmt.Println("Key 'b' was removed from cache")
	}
}
