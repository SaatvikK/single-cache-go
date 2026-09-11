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
	fmt.Println(r)
	fmt.Fprint(w, b.data)
}

func (b *BackingStore) fetchItem(w http.ResponseWriter, r *http.Request) {
	var s string = mux.Vars(r)["item"]
	flag := false

	item, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.size; i++ {
		fmt.Println(("for loopin"))
		if item == b.data[i] {
			fmt.Println("Item found")
			flag = true
			break
		}
	}

	if flag == true {
		fmt.Fprint(w, item)
	} else {
		fmt.Fprint(w, "not found")
	}
}

func RunServer() {
	var port string = "4000"

	fmt.Println("server started")

	data := []int{37, 63, 51, 95, 26, 24, 75}
	B := BackingStore{
		data: data,
		size: len(data),
	}

	fmt.Printf("store size %d: %v\n", B.size, B.data)

	r := mux.NewRouter()
	r.HandleFunc("/b-dump", B.dump).Methods(http.MethodGet)
	r.HandleFunc("/fetch/{item}", B.fetchItem).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

}
