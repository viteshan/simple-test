package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/golang-lru"
)

func TestLRU(t *testing.T) {
	cache, _ := lru.NewWithEvict(20, func(key interface{}, value interface{}) {
		fmt.Println(key)
	})

	for i := 0; i < 100; i++ {
		cache.Add(i, i)
	}
}
