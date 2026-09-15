package cache

const TauVal = 6

func LCEpol(cap int, times []int) int {
	leastRecent := 0

	for i := 0; i < len(times); i++ {
		if times[i] > leastRecent {
			leastRecent = i
		}
	}

	return leastRecent
}
