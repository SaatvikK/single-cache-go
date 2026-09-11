package cache

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

const Empty = 0

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type Cache struct {
	Items    []int
	Times    []int
	Capacity int
	server   string
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

func (C *Cache) admitItem(item int, time int, idx int) bool {
	if idx >= len(C.Items) || idx < 0 {
		return false
	}
	C.Items[idx] = item
	C.Times[idx] = time
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

	admitOut = C.admitItem(admit, TauVal, idx)
	if admitOut == false {
		return false, fmt.Errorf("Failed to admit item.")
	}

	return true, nil
}

func (C *Cache) fetch(item int) (int, error) {
	endp := C.server + fmt.Sprintf("/fetch/%s", item)

	resp, err := http.Get(endp)
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//Convert the body to type string
	sb := string(body)

	if sb == "not found" {
		return Empty, fmt.Errorf("Could not fetch item from B.")
	}

	fetched, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return fetched, nil
}

func (C *Cache) cacheDump(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println(r)
	fmt.Fprint(w, C.Items)
	fmt.Fprint(w, C.Times)
}

func (C *Cache) cacheServer() {
	var port string = "4040"

	fmt.Println("cache server started")

	r := mux.NewRouter()
	r.HandleFunc("/c-dump", B.dump).Methods(http.MethodGet)
	r.HandleFunc("/request/{item}", B.fetchItem).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
