package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/vininta-solution/bid/model/ads"
	"github.com/vininta-solution/bid/model/placement"
)

var allAds map[int]ads.Ads

type pickRequest struct {
	MinBid   float64 `json:"minBid"`
	Category []int   `category:"category"`
}

func main() {
	allAds = make(map[int]ads.Ads)
	r := http.NewServeMux()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/add", addHandler)
	r.HandleFunc("/delete", deleteHandler)
	r.HandleFunc("/clear", clearHandler)
	r.HandleFunc("/list", listHandler)
	r.HandleFunc("/runtime-status", runtimeHandler)
	r.HandleFunc("/pick", pickHandler)
	r.HandleFunc("/init-random", randomHandler)

	s := &http.Server{
		Addr:           ":8081",
		Handler:        r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	response("No Response", w, r)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	var ad ads.Ads
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &ad)
	m, _ := json.Marshal(ad)
	response(string(m), w, r)

	allAds[ad.Id] = ad
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var ad ads.Ads
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &ad)
	m, _ := json.Marshal(ad)
	response(string(m), w, r)

	delete(allAds, ad.Id)
}

func clearHandler(w http.ResponseWriter, r *http.Request) {
	allAds = make(map[int]ads.Ads)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	allAds, _ := json.Marshal(allAds)
	response(string(allAds), w, r)
}

func pickHandler(w http.ResponseWriter, r *http.Request) {
	var winnerAds ads.Ads
	var count int
	var placement placement.Placement
	var adRequest pickRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &adRequest)
	placement.Category = adRequest.Category
	for _, ad := range allAds {
		if ad.IsMatch(placement) {
			if ad.Bid > winnerAds.Bid {
				winnerAds = ad
			}
		}
		count++
	}
	resp, _ := json.Marshal(winnerAds)
	//fmt.Println(count)
	response(string(resp), w, r)
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	var ad ads.Ads
	var i int
	for i = 0; i < 100; i++ {
		ad.Id = rand.Intn(4000000000)
		ad.Bid = rand.Float64() * 100
		ad.Category = make([]int, 2)
		ad.Category[0] = rand.Intn(3)
		ad.Category[1] = rand.Intn(3)
		allAds[ad.Id] = ad
	}
}

func runtimeHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Fprintf(w, "Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Fprintf(w, "\nTotalAlloc = %d byte, %v MiB", m.TotalAlloc, bToMb(m.TotalAlloc))
	fmt.Fprintf(w, "\nSys = %v MiB", bToMb(m.Sys))
	fmt.Fprintf(w, "\nNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func response(resp string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "%s", resp)
}
