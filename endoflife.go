package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const baseUrl = "https://endoflife.date/api/"

type BoolOrString struct {
	StringValue string
	BoolValue   bool
	IsString    bool
}

func (bos BoolOrString) asString() string {
	if bos.IsString {
		return bos.StringValue
	} else {
		return strconv.FormatBool(bos.BoolValue)
	}
}

func (bos *BoolOrString) UnmarshalJSON(data []byte) error {
	// Try parsing as string
	if err := json.Unmarshal(data, &bos.StringValue); err == nil {
		bos.IsString = true
		return nil
	}

	// Try parsing as bool
	if err := json.Unmarshal(data, &bos.BoolValue); err == nil {
		bos.IsString = false
		return nil
	}

	return fmt.Errorf("attribute is neither string nor bool")
}

type Cycle struct {
	Cycle       string
	ReleaseDate string
	Eol         BoolOrString
	Latest      string
	Link        string
	// Lts          bool
	// Support      string
	// Discontinued string
}

func (c Cycle) print(w io.Writer) {
	fmt.Fprintf(w, "Version %s on %s (until %s)\n", c.Latest, c.ReleaseDate, c.Eol.asString())
}

type Product []Cycle

func (p Product) print(w io.Writer) {
	for _, c := range p {
		c.print(w)
	}
}

func getProduct(product string) Product {
	resp, err := http.Get(baseUrl + product + ".json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	myProduct := Product{}
	err = json.NewDecoder(resp.Body).Decode(&myProduct)
	if err != nil {
		panic(err)
	}
	// myProduct.print(os.Stdout)
	return myProduct
}

func getCycle(product string, cicle string) Cycle {
	resp, err := http.Get(baseUrl + product + "/" + cicle + ".json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	myCicle := Cycle{}
	err = json.NewDecoder(resp.Body).Decode(&myCicle)
	if err != nil {
		panic(err) // TODO
	}
	// myCicle.print(os.Stdout)
	return myCicle
}
