package main

import (
	//"fmt"
	//"single-cache/cache"
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
	//C := makecache(5)

	server.RunServer()
}
