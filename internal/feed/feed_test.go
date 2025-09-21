package feed

import (
	"encoding/json"
	"poschi3/releasefeed/internal/endoflife"
	"poschi3/releasefeed/test/testdata"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// eol	"2025-11-30"
// latest	"1.39.12"
// latestReleaseDate	"2025-04-10"
// releaseDate	"2022-11-30"
// lts	true

func TestFeedCycle(t *testing.T) {
	testCycle := endoflife.Cycle{
		Cycle: "1.39",
		ReleaseDate: endoflife.CustomDate{
			Time: time.Date(2025, 4, 10, 0, 0, 0, 0, time.UTC),
		},
		Eol: endoflife.BoolOrString{
			StringValue: "2025-11-30",
		},
		Latest: "1.39.12",
	}
	feed, err := FeedCycle("localhost", "mediawiki", testCycle)
	if err != nil {
		t.Error(err)
	}
	t.Log(feed)

}

func TestNixos(t *testing.T) {
	feed, err := convertToProductFeed("nixos")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "25.05", feed.Entries[0].Title)
	assert.Equal(t, "https://endoflife.date/nixos", feed.Entries[0].Link.Href)
	assert.Equal(t, "tag:localhost,2025-05-23:nixos:25.05", feed.Entries[0].Id)
	assert.Equal(t, FeedTimeformat{
		Time: time.Date(2025, 5, 23, 0, 0, 0, 0, time.UTC),
	}, feed.Entries[0].Updated)
	assert.Equal(t, "nixos 25.05 updated to 25.05 (2025-05-23). Support until 2025-12-31", feed.Entries[0].Summary)
}

func TestMediawiki(t *testing.T) {
	feed, err := convertToProductFeed("mediawiki")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "1.43.3", feed.Entries[1].Title)
	assert.Equal(t, "https://endoflife.date/mediawiki", feed.Entries[1].Link.Href)
	assert.Equal(t, "tag:localhost,2025-07-01:mediawiki:1.43.3", feed.Entries[1].Id)
	assert.Equal(t, FeedTimeformat{
		Time: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
	}, feed.Entries[1].Updated)
	assert.Equal(t, "mediawiki 1.43 updated to 1.43.3 (2025-07-01). Support until 2027-12-31", feed.Entries[1].Summary)
}

func TestJiraSoftware(t *testing.T) {
	feed, err := convertToProductFeed("jira-software")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "11.0.1", feed.Entries[0].Title)
	assert.Equal(t, "https://endoflife.date/jira-software", feed.Entries[0].Link.Href)
	assert.Equal(t, "tag:localhost,2025-09-04:jira-software:11.0.1", feed.Entries[0].Id)
	assert.Equal(t, FeedTimeformat{
		Time: time.Date(2025, 9, 4, 0, 0, 0, 0, time.UTC),
	}, feed.Entries[0].Updated)
	assert.Equal(t, "jira-software 11.0 updated to 11.0.1 (2025-09-04). Support until 2027-08-13", feed.Entries[0].Summary)
}

func convertToProductFeed(productName string) (Feed, error) {
	fileContent, err := testdata.GetFileContent("endoflive", productName+".json")
	if err != nil {
		return Feed{}, err
	}
	var product endoflife.Product
	err = json.Unmarshal(fileContent, &product)
	if err != nil {
		return Feed{}, err
	}

	feed := createProductFeed(productName, "localhost", product)
	return feed, nil
}
