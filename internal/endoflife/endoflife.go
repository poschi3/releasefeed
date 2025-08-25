package endoflife

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const baseUrl = "https://endoflife.date/api/"

type BoolOrString struct {
	StringValue string
	BoolValue   bool
	IsString    bool
}

func (bos BoolOrString) AsString() string {
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

type CustomDate struct {
	time.Time
}

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		layout := "2006-01-02" // YYYY-MM-DD format
		// Parse the time using the custom layout
		t, err := time.Parse(layout, str)
		if err != nil {
			return err
		}
		cd.Time = t
		return nil
	}

	return fmt.Errorf("attribute could not be paresd as Time")
}

type Cycle struct {
	Cycle             string       // The release cycle which this release is part of. e.g. 1.39
	ReleaseDate       string       // Release date for the first release in this cycle.
	Eol               BoolOrString // End-of-Life date for this release cycle.
	Latest            string       // Latest release in this cycle.
	LatestReleaseDate CustomDate   // Release date for the latest release in this cycle.
	Link              string       // Link to changelog for the latest release in this cycle, or null if unavailable.
	// Lts          bool // Whether this release cycle has long-term-support (LTS), or the date it entered LTS status.
	// Support      string // Whether this release cycle has active support.
	// Discontinued string // Whether this device version is no longer in production.
}

type Product []Cycle

func GetProduct(product string) Product {
	escapedProduct := url.PathEscape(product)
	resp, err := http.Get(baseUrl + escapedProduct + ".json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Response status:", resp.Status)

	myProduct := Product{}
	err = json.NewDecoder(resp.Body).Decode(&myProduct)
	if err != nil {
		panic(err)
	}
	// myProduct.print(os.Stdout)
	return myProduct
}

func GetCycle(product string, cycle string) Cycle {
	escapedProduct := url.PathEscape(product)
	escapedCycle := url.PathEscape(cycle)
	resp, err := http.Get(baseUrl + escapedProduct + "/" + escapedCycle + ".json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Response status:", resp.Status)

	myCycle := Cycle{}
	err = json.NewDecoder(resp.Body).Decode(&myCycle)
	if err != nil {
		panic(err) // TODO
	}

	if myCycle.Cycle == "" {
		myCycle.Cycle = cycle
	}
	// myCycle.print(os.Stdout)
	// log.Println(myCycle.LatestReleaseDate.Time)
	return myCycle
}
