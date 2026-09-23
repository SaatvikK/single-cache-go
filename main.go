package main

import (
	//"fmt"
	"single-cache/cache"
	"single-cache/server"

	// ## added by claude because i needed time for the Times values below
	"log"
	"os"
	"time"
)

func setLog() {
	// ## added by claude because i was appending with O_APPEND, so old runs piled up in the file
	f, err := os.OpenFile("logfile.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	// ## added by claude because i was defer-closing f here, which shut the file as soon as setLog returned
	// f now stays open for the life of the process

	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("started")
}

func main() {

	setLog()

	// make cache
	C := cache.Cache{
		// ## added by claude because i was starting the cache empty
		Items:    []int{37, 63, 51, 95, 26},
		Times:    []time.Time{time.Unix(1, 0), time.Unix(2, 0), time.Unix(3, 0), time.Unix(4, 0), time.Unix(5, 0)},
		Capacity: 5,
		Server:   "http://localhost:4000",
	}
	// ## added by claude because i wasn't logging the starting cache state
	log.Printf("cache initialised: capacity %d, backing store %s, items %v", C.Capacity, C.Server, C.Items)

	go server.RunServer()
	go C.CacheServer()

	time.Sleep(1 * time.Second)

	harnessMain("http://localhost:4000", "http://localhost:4040")
}
