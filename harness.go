package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type sequence struct {
	Items []int
	Times []time.Time
}

func buildSeq() sequence {
	base := time.Now()
	seq := sequence{
		Items: []int{37, 63, 24, 37, 75, 51, 24, 95, 26, 37},
		Times: []time.Time{
			base.Add(1 * time.Second), base.Add(2 * time.Second), base.Add(3 * time.Second),
			base.Add(5 * time.Second), base.Add(6 * time.Second), base.Add(10 * time.Second),
			base.Add(11 * time.Second), base.Add(12 * time.Second), base.Add(14 * time.Second),
			base.Add(15 * time.Second),
		},
	}

	return seq
}

func parseCacheResponse(body []byte) string {
	// unmarshal the response

	var res map[string]any
	if err := json.Unmarshal(body, &res); err != nil {
		panic(err)
	}

	if res["outcome"] == "hit" {
		return "hit"
	} else {
		return "miss"
	}
}

func reqB(Burl string, seq sequence) {
	fmt.Println("================= TO BACKING STORE DIRECT =================")

	// to B
	baseep := Burl + "/fetch/"

	for i := 0; i < len(seq.Items); i++ {
		item := seq.Items[i]
		start := time.Now()
		// fmt.Println("==========")
		// fmt.Println(start)
		// fmt.Println(seq.Times[i])

		if seq.Times[i].Before(start) {
			continue
		}

		s := strconv.Itoa(item)
		endp := baseep + s

		resp, err := http.Get(endp)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Println("Item", item, "request to B failed:", err)
		}

		if resp.StatusCode != 200 {
			fmt.Println("Item", item, "got status", resp.StatusCode, "from B")
		}

		latency := elapsed.Nanoseconds()

		fmt.Printf("Item %d took %d ns to be fetched from B.\n", item, latency)

	}
}

func reqC(Curl string, seq sequence) {
	fmt.Println("================= TO CACHE FIRST =================")

	// to C
	baseep := Curl + "/request/"

	for i := 0; i < len(seq.Items); i++ {
		item := seq.Items[i]
		start := time.Now()
		// fmt.Println("==========")
		// fmt.Println(start)
		// fmt.Println(seq.Times[i])

		if seq.Times[i].Before(start) {
			continue
		}

		s := strconv.Itoa(item)
		endp := baseep + s

		resp, err := http.Get(endp)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Println("Item", item, "request to C failed:", err)
		}

		if resp.StatusCode != 200 {
			fmt.Println("Item", item, "got status", resp.StatusCode, "from C")
		}

		//We Read the response body on the line below.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		stat := parseCacheResponse(body)
		latency := elapsed.Nanoseconds()

		if stat == "hit" {
			fmt.Printf("Item %d took %d ns to be fetched from C (HIT).\n", item, latency)
		} else {
			fmt.Printf("Item %d took %d ns to be fetched from C (MISS -- went to B).\n", item, latency)
		}
	}
}

func harnessMain(Burl string, Curl string) {
	if Burl == "" {
		Burl = "http://localhost:4000"
	}

	if Curl == "" {
		Curl = "http://localhost:4040"
	}

	seq := buildSeq()
	reqB(Burl, seq)
	seq = buildSeq()
	reqC(Curl, seq)
}
