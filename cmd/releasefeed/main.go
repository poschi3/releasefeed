package main

import (
	"fmt"
	"net/http"
	"poschi3/releasefeed/internal/endoflife"
	"poschi3/releasefeed/internal/feed"
	"strings"
)

func main() {
	http.HandleFunc("/{product}", handleProduct)
	http.HandleFunc("/{product}/{cycle}", handleCycle)
	http.ListenAndServe(":8090", nil)
}

func handleProduct(w http.ResponseWriter, req *http.Request) {
	productName := req.PathValue("product")
	product, err := endoflife.GetProduct(productName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	feed, err := feed.FeedProduct(strings.Split(req.Host, ":")[0], productName, product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeFeed(w, feed)
}

func handleCycle(w http.ResponseWriter, req *http.Request) {
	productName := req.PathValue("product")
	cycleName := req.PathValue("cycle")
	cycle, err := endoflife.GetCycle(productName, cycleName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	feed, err := feed.FeedCycle(strings.Split(req.Host, ":")[0], productName, cycle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeFeed(w, feed)
}

func writeFeed(w http.ResponseWriter, feed string) {
	w.Header().Add("content-type", "application/atom+xml; charset=UTF-8")
	fmt.Fprintln(w, feed)
}
