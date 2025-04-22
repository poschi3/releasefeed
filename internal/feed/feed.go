package feed

import (
	"encoding/xml"
	"fmt"
	"log"
	"poschi3/releasefeed/internal/endoflife"
	"time"
)

// Template
// <?xml version="1.0" encoding="utf-8"?>
// <feed xmlns="http://www.w3.org/2005/Atom">
//   <author>
//     <name>Autor des Weblogs</name>
//   </author>
//   <title>Titel des Weblogs</title>
//   <id>urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6</id>
//   <updated>2003-12-14T10:20:09Z</updated>

//   <entry>
//     <title>Titel des Weblog-Eintrags</title>
//     <link href="http://example.org/2003/12/13/atom-beispiel"/>
//     <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
//     <updated>2003-12-13T18:30:02Z</updated>
//     <summary>Zusammenfassung des Weblog-Eintrags</summary>
//     <content>Volltext des Weblog-Eintrags</content>
//   </entry>
// </feed>

const baseUri = "https://endoflife.date/"

type Feed struct {
	XMLName xml.Name       `xml:"http://www.w3.org/2005/Atom feed"`
	Author  Author         `xml:"author"`
	Title   string         `xml:"title"`
	Id      string         `xml:"id"`
	Updated FeedTimeformat `xml:"updated"`
	Entries []Entry        `xml:"entry"`
}

type FeedTimeformat struct {
	time.Time
}

// MarshalXML implements the xml.Marshaler interface
func (ct FeedTimeformat) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	formattedTime := ct.Format(time.RFC3339)
	return e.EncodeElement(formattedTime, start)
}

type Author struct {
	Name string `xml:"name"`
}

type Entry struct {
	Title   string         `xml:"title"`
	Link    Link           `xml:"link"`
	Id      string         `xml:"id"`
	Updated FeedTimeformat `xml:"updated"`
	Summary string         `xml:"summary"`
	Content string         `xml:"content,omitempty"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

func FeedCycle(host string, productName string, cycle endoflife.Cycle) string {
	feed := Feed{
		Author:  Author{Name: "ReleaseFeed"},
		Title:   fmt.Sprintf("%s (%s)", productName, cycle.Cycle),
		Id:      getCycleTag(host, productName, cycle),
		Updated: FeedTimeformat{Time: cycle.LatestReleaseDate.Time},
		Entries: []Entry{
			createCycleEntry(host, productName, cycle),
		},
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		log.Printf("Error encoding XML: %v\n", err)
		return "" // TODO was machen
	}

	// TODO
	// <?xml version="1.0" encoding="utf-8"?>
	return string(output)
}

func FeedProduct(host string, productName string, product endoflife.Product) string {
	feed := Feed{
		Author: Author{Name: "ReleaseFeed"},
		Title:  productName,
		Id:     "https://" + host + "/" + productName,
	}

	var latestUpdate time.Time
	for _, cycle := range product {
		feedCycle := createCycleEntry(host, productName, cycle)
		feed.Entries = append(feed.Entries, feedCycle)
		if cycle.LatestReleaseDate.Time.After(latestUpdate) {
			latestUpdate = cycle.LatestReleaseDate.Time
		}
	}
	feed.Updated = FeedTimeformat{Time: latestUpdate}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		log.Printf("Error encoding XML: %v\n", err)
		return "" // TODO was machen
	}

	// TODO
	// <?xml version="1.0" encoding="utf-8"?>
	return string(output)
}

func createCycleEntry(host string, productName string, cycle endoflife.Cycle) Entry {
	summary := fmt.Sprintf(
		"%s %s updated to %s (%s). Support until %s",
		productName, cycle.Cycle, cycle.Latest, cycle.LatestReleaseDate.Time.Format("2006-01-02"), cycle.Eol.AsString(),
	)

	return Entry{
		Title:   cycle.Latest,
		Id:      getCycleTag(host, productName, cycle),
		Link:    Link{Href: baseUri + productName},
		Summary: summary,
		Updated: FeedTimeformat{Time: cycle.LatestReleaseDate.Time},
	}
}

func getCycleTag(host string, productName string, cycle endoflife.Cycle) string {
	return "tag:" + host + "," + cycle.LatestReleaseDate.Time.Format("2006-01-02") + ":" + productName + ":" + cycle.Latest
}
