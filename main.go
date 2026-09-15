package main

import (
	//"fmt"
	"single-cache/cache"
	"single-cache/server"
	// ## added by claude because i needed time for the Times values below
	"time"
)

// func makecache(size int) *cache.Cache {
// 	C := cache.Cache{
// 		Items: make([]int, size),
// 		Times: make([]int, size),
// 	}

// 	return &C
// }

func main() {
	// make cache
	C := cache.Cache{
		// ## added by claude because i was starting the cache empty
		Items:    []int{37, 63, 51, 95, 26},
		Times:    []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0), time.Unix(4, 0), time.Unix(5, 0)},
		Capacity: 5,
		Server:   "http://localhost:4000",
	}

	go server.RunServer()

	C.CacheServer()
}
