package feed

import (
	"poschi3/releasefeed/internal/endoflife"
	"testing"
)

// eol	"2025-11-30"
// latest	"1.39.12"
// latestReleaseDate	"2025-04-10"
// releaseDate	"2022-11-30"
// lts	true

func TestFeedCycle(t *testing.T) {
	t.Log("fssoo")

	testCycle := endoflife.Cycle{
		Cycle:       "1.39",
		ReleaseDate: "2025-04-10",
		Eol: endoflife.BoolOrString{
			StringValue: "2025-11-30",
		},
		Latest: "1.39.12",
	}
	feed, _ := FeedCycle("localhost", "mediawiki", testCycle)
	t.Log(feed)

}
