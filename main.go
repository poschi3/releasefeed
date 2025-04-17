package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/{product}", handleProduct)
	http.HandleFunc("/{product}/{cycle}", handleCycle)
	http.ListenAndServe(":8090", nil)
}

func handleProduct(w http.ResponseWriter, req *http.Request) {
	productName := req.PathValue("product")
	product := getProduct(productName)
	product.print(w)
}

func handleCycle(w http.ResponseWriter, req *http.Request) {
	productName := req.PathValue("product")
	cycleName := req.PathValue("cycle")
	cycle := getCycle(productName, cycleName)
	cycle.print(w)
}
