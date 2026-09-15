package cache

import (
	"fmt"
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
	const layout = "2006-01-02 15:04:05.000"
	fmt.Printf("Cached item %d (index %d) is the LRU-item, last requested at %s\n", items[leastRecent], leastRecent, times[leastRecent].Format(layout))
	fmt.Println("Cache times:")
	for i := range times {
		fmt.Printf("  index %d: item %d at %s\n", i, items[i], times[i].Format(layout))
	}
	return leastRecent
}
