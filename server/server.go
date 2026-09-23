package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type BackingStore struct {
	data []int
	size int
}

func (b *BackingStore) dump(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println(r)
	fmt.Fprint(w, b.data)
}

func (b *BackingStore) fetchItem(w http.ResponseWriter, r *http.Request) {
	var s string = mux.Vars(r)["item"]
	flag := false

	item, err := strconv.Atoi(s)
	if err != nil {
		// ## added by claude because i was panicking without saying what the bad input was
		log.Printf("backing store: bad item %q in %s: %v", s, r.URL.Path, err)
		panic(err)
	}

	// ## added by claude because i was printing "for loopin" every iteration and "Item found" with no item or index
	log.Printf("backing store: fetch request for item %d", item)
	for i := 0; i < b.size; i++ {
		if item == b.data[i] {
			log.Printf("backing store: item %d found at index %d", item, i)
			flag = true
			break
		}
	}

	if flag == true {
		fmt.Fprint(w, item)
	} else {
		// ## added by claude because i wasn't logging misses
		log.Printf("backing store: item %d not found", item)
		fmt.Fprint(w, "not found")
	}
}

func RunServer() {
	var port string = "4000"

	log.Println("server started")

	data := []int{37, 63, 51, 95, 26, 24, 75}
	B := BackingStore{
		data: data,
		size: len(data),
	}

	// ## added by claude because i was printing this with fmt so it never reached the log file
	log.Printf("backing store on :%s, size %d: %v", port, B.size, B.data)

	r := mux.NewRouter()
	r.HandleFunc("/b-dump", B.dump).Methods(http.MethodGet)
	r.HandleFunc("/fetch/{item}", B.fetchItem).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

}
