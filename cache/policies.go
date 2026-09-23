package cache

import (
	"log"
	"time"
)

func LRUpol(cap int, items []int, times []time.Time) int {
	leastRecent := 0

	for i := 0; i < len(times); i++ {
		now := time.Now()
		currLeastElapsed := now.Sub(times[leastRecent])
		currElapsed := now.Sub(times[i])
		if currElapsed > currLeastElapsed {
			leastRecent = i
		}
	}
	// ## added by claude because i was printing raw time.Time values with %d and no newline
	// ## added by claude because these were fmt prints so they never reached the log file
	const layout = "2006-01-02 15:04:05.000"
	log.Printf("cache: LRU victim is item %d (index %d), last requested at %s", items[leastRecent], leastRecent, times[leastRecent].Format(layout))
	log.Println("cache: times at eviction:")
	for i := range times {
		log.Printf("  index %d: item %d at %s", i, items[i], times[i].Format(layout))
	}
	return leastRecent
}
