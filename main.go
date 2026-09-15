package main

import (
	//"fmt"
	"single-cache/cache"
	"single-cache/server"
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
		Items:    []int{37, 63, 51, 95, 26},
		Times:    []int{1, 2, 3, 4, 5},
		Capacity: 5,
		Server:   "http://localhost:4000",
	}

	go server.RunServer()

	C.CacheServer()
}
