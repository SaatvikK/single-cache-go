package cache

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
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

	return true, nil
}

func (C *Cache) fetch(item int) (int, error) {
	endp := C.Server + fmt.Sprintf("/fetch/%d", item)
	fmt.Println(endp)
	resp, err := http.Get(endp)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("hiiiiiiiiii")

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//Convert the body to type string
	sb := string(body)

	if sb == "not found" {
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
	fmt.Println(r)
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
	foundItem, foundTime := 0, 0

	item, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	for i := 0; i < C.Capacity; i++ {
		fmt.Println(("for loopin"))
		if item == C.Items[i] {
			fmt.Println("Item found")
			flag = true
			foundItem = C.Items[i]
			foundTime = int(C.Times[i].Unix())
			C.Times[i] = time.Now()
			break
		}
	}

	type cacheres struct {
		item    int
		time    int
		err     string
		outcome string
	}

	res := cacheres{
		item:    foundItem,
		time:    foundTime,
		err:     "null",
		outcome: "success",
	}

	if flag == false {
		res.outcome = "not found"
		fmt.Println("Item not found, sending fetch request now.")
		C.fetchUpdateHandler(item)

	}

	fmt.Fprint(w, res)
}

func (C *Cache) CacheServer() {
	var port string = "4040"

	fmt.Println("cache server started")

	r := mux.NewRouter()
	r.HandleFunc("/c-dump", C.cacheDump).Methods(http.MethodGet)
	r.HandleFunc("/request/{item}", C.makeRequest).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
