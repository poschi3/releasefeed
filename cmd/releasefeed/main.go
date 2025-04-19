package main

import (
	"fmt"
	"log"
	"net/http"
	"poschi3/releasefeed/internal/endoflife"
	"poschi3/releasefeed/internal/feed"
)

func main() {
	http.HandleFunc("/{product}", handleProduct)
	http.HandleFunc("/{product}/{cycle}", handleCycle)
	http.ListenAndServe(":8090", nil)
}

func handleProduct(w http.ResponseWriter, req *http.Request) {
	// TODO sanitation
	log.Println(req.URL)
	productName := req.PathValue("product")
	product := endoflife.GetProduct(productName)

	w.Header().Add("content-type", "application/atom+xml; charset=UTF-8")

	feed := feed.FeedProduct(productName, product)
	fmt.Fprintln(w, feed)
}

func handleCycle(w http.ResponseWriter, req *http.Request) {
	// TODO sanitation
	productName := req.PathValue("product")
	cycleName := req.PathValue("cycle")
	cycle := endoflife.GetCycle(productName, cycleName)
	w.Header().Add("content-type", "application/atom+xml; charset=UTF-8")
	//cycle.print(w)
	feed := feed.FeedCycle(productName, cycle)
	fmt.Fprintln(w, feed)
}
