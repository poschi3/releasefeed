package main

import (
	"encoding/xml"
	"fmt"
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

type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Author  Author   `xml:"author"`
	Title   string   `xml:"title"`
	Id      string   `xml:"id"`      // TODO
	Updated string   `xml:"updated"` // TODO
	Entries []Entry  `xml:"entry"`
}

type Author struct {
	Name string `xml:"name"`
}

type Entry struct {
	Title   string `xml:"title"`
	Link    Link   `xml:"link"`    // TODO
	Id      string `xml:"id"`      // TODO
	Updated string `xml:"updated"` // TODO
	Summary string `xml:"summary"`
	Content string `xml:"content"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

func feedCycle(cycle Cycle) string {
	feed := Feed{
		Author: Author{
			Name: "ReleaseFeed",
		},
		Title:   "Bla",
		Id:      "TODO",
		Updated: "xxx",
		Entries: []Entry{
			createCycle(cycle),
		},
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding XML: %v\n", err)
		return "" // TODO was machen
	}

	// TODO
	// <?xml version="1.0" encoding="utf-8"?>

	fmt.Println(string(output))
	return string(output)
}

func feedProduct(product Product) string {
	feed := Feed{
		Author: Author{
			Name: "ReleaseFeed",
		},
		Title:   "Bla",
		Id:      "TODO",
		Updated: "xxx",
	}

	for _, cycle := range product {
		feedCycle := createCycle(cycle)
		feed.Entries = append(feed.Entries, feedCycle)
		// builder.WriteString(c.format())
		// builder.WriteString("\n")
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding XML: %v\n", err)
		return "" // TODO was machen
	}

	// TODO
	// <?xml version="1.0" encoding="utf-8"?>

	fmt.Println(string(output))
	return string(output)
}

func createCycle(cycle Cycle) Entry {
	return Entry{
		Title:   cycle.Latest,
		Link:    Link{Href: "http://example.com"},
		Summary: cycle.format(),
	}
}
