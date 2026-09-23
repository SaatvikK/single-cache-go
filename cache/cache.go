package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const Empty = 0

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type Cache struct {
	Items    []int
	Times    []time.Time
	Capacity int
	Server   string
}

func (C *Cache) eval(item int) bool {
	items := C.Items
	for i := 0; i < len(items); i++ {
		if item == items[i] {
			return true // HIT
		}
	}

	return false // MISS
}

func findEmpty(arr []int) (int, bool) { // finds empty indices
	emptyidx := Empty
	for i := 0; i < len(arr); i++ {
		if arr[i] == Empty {
			emptyidx = i
			return emptyidx, true
		}
	}

	return 0, false
}

func (C *Cache) evictItem(item int) (int, bool) {
	items := C.Items
	for i := 0; i < len(C.Items); i++ {
		if items[i] == item {
			items[i] = Empty
			return i, true
		}
	}

	return 0, false
}

func (C *Cache) admitItem(item int, idx int) bool {
	if idx >= len(C.Items) || idx < 0 {
		return false
	}
	C.Items[idx] = item
	C.Times[idx] = time.Now()
	return true
}

func (C *Cache) update(admit int, evict int) (bool, error) {
	idx, outEmpty, admitOut := 0, false, false
	if len(C.Items) < C.Capacity {
		idx, outEmpty = findEmpty(C.Items)
	}

	if outEmpty == false { // no free slots, must evict first
		evictOut := false
		idx, evictOut = C.evictItem(evict)
		if evictOut == false {
			return false, fmt.Errorf("Failed to evict an item.")
		}
	}

	admitOut = C.admitItem(admit, idx)
	if admitOut == false {
		return false, fmt.Errorf("Failed to admit item.")
	}

	// ## added by claude because i wasn't logging which slot the admitted item went into
	log.Printf("cache: admitted item %d at index %d (free slot: %t)", admit, idx, outEmpty)
	return true, nil
}

func (C *Cache) fetch(item int) (int, error) {
	endp := C.Server + fmt.Sprintf("/fetch/%d", item)
	// ## added by claude because i was printing the bare endpoint with fmt.Println
	log.Printf("cache: fetching item %d from %s", item, endp)
	resp, err := http.Get(endp)
	if err != nil {
		log.Fatalln(err)
	}

	// ## added by claude because i was printing "hiiiiiiiiii" here
	log.Printf("cache: backing store replied %s for item %d", resp.Status, item)

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//Convert the body to type string
	sb := string(body)

	if sb == "not found" {
		// ## added by claude because i wasn't logging a backing-store miss
		log.Printf("cache: backing store does not have item %d", item)
		return Empty, fmt.Errorf("Could not request item from cache.")
	}

	fetched, err := strconv.Atoi(sb)
	if err != nil {
		log.Fatalln(err)
	}

	return fetched, nil
}

func (C *Cache) cacheDump(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// ## added by claude because i was printing this with fmt so it never reached the log file
	log.Println(r)
	fmt.Fprint(w, C.Items)
	fmt.Fprint(w, C.Times)
}

func (C *Cache) fetchUpdateHandler(item int) error {
	fetched, err := C.fetch(item)
	if err != nil {
		return err
	}

	victim := LRUpol(C.Capacity, C.Items, C.Times)
	vicItem := C.Items[victim]

	out, err := C.update(fetched, vicItem)
	if err != nil {
		return err
	}

	log.Printf("Updated cache: admitted %d, evicted %d, code %t", fetched, vicItem, out)
	return nil
}

func (C *Cache) makeRequest(w http.ResponseWriter, r *http.Request) {
	var s string = mux.Vars(r)["item"]
	flag := false
	foundItem, foundTime := 0, time.Now()
	foundIdx := 0

	item, err := strconv.Atoi(s)
	if err != nil {
		// ## added by claude because i was panicking without saying what the bad input was
		log.Printf("cache: bad item %q in %s: %v", s, r.URL.Path, err)
		panic(err)
	}

	// ## added by claude because i was printing "for loopin" every iteration and "Item found" with no item or index
	log.Printf("cache: request for item %d", item)
	for i := 0; i < C.Capacity; i++ {
		if item == C.Items[i] {
			flag = true
			foundItem = C.Items[i]
			foundTime = C.Times[i]
			C.Times[i] = time.Now()
			break
		}
	}

	type cacheres struct {
		Item    int    `json:"item"`
		Time    int    `json:"time"`
		Err     string `json:"err"`
		Outcome string `json:"outcome"`
	}

	res := cacheres{
		Item:    foundItem,
		Time:    int(foundTime.Unix()),
		Err:     "null",
		Outcome: "hit",
	}

	if flag == false {
		res.Outcome = "miss"
		// ## added by claude because i was printing the miss with fmt and dropping fetchUpdateHandler's error
		log.Printf("cache: MISS item %d, fetching from backing store", item)
		if err := C.fetchUpdateHandler(item); err != nil {
			log.Printf("cache: fetch/update for item %d failed: %v", item, err)
		}

	} else {
		log.Printf("cache: HIT item %d at index %d, last requested %s", foundItem, foundIdx, foundTime.Format(time.StampMilli))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
	}

}

func (C *Cache) CacheServer() {
	var port string = "4040"

	// ## added by claude because i was printing this with fmt so it never reached the log file
	log.Printf("cache server started on :%s", port)

	r := mux.NewRouter()
	r.HandleFunc("/c-dump", C.cacheDump).Methods(http.MethodGet)
	r.HandleFunc("/request/{item}", C.makeRequest).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
