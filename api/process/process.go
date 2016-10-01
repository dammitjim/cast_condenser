package process

import (
	"condenser/api/external"
	"errors"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

var timeTemplates = []string{
	time.RFC1123,
	time.RFC1123Z,
	"Mon, 2 Jan 2006 15:04 MST",
	"Mon, 2 Jan 2006 15:04 -0700",
}

// Run is the main function for processing a set of podcasts.
func Run(podcasts ...*external.Podcast) {
	for _, podcast := range podcasts {
		logrus.WithFields(logrus.Fields{
			"owner": podcast.Owner,
			"feed":  podcast.FeedURL,
		}).Info("processing " + podcast.Name)

		feed, err := getFeed(podcast.FeedURL)
		if err != nil {
			logrus.Error(err)
			continue
		}

		tracks, err := extractTracks(feed.Items)
		if err != nil {
			logrus.Error(err)
			continue
		}

		logrus.WithField("len", len(tracks)).Info("processed " + podcast.Name)
	}
}

func extractTracks(items []*gofeed.Item) ([]*Track, error) {
	var err error
	var tracks []*Track
	for _, item := range items {
		if len(item.Enclosures) != 1 {
			err = errors.New("Multiple enclosures found for " + item.Title)
			logrus.WithField("enclosures", item.Enclosures).Warn(err)
			continue
		}

		if _, ok := item.Extensions["itunes"]; !ok {
			err = errors.New("No itunes data parsable from " + item.Title)
			return nil, err
		}

		itunesData := ext.NewITunesItemExtension(item.Extensions["itunes"])

		var publishedParsed time.Time
		success := false

		// TODO template checking order optimisation here?
		for _, template := range timeTemplates {
			err = nil
			publishedParsed, err = time.Parse(template, item.Published)
			if err != nil {
				logrus.Debug(err)
				continue
			}
			success = true
			break
		}

		if !success {
			err = errors.New("Could not parse publish time " + item.Published)
			return nil, err
		}

		track := &Track{
			Name:     item.Title,
			Duration: itunesData.Duration,
			MediaURL: item.Enclosures[0].URL,
			Summary:  itunesData.Summary,
			Posted:   publishedParsed.Unix(),
		}

		err = track.Validate()
		if err != nil {
			logrus.Error(err)
			continue
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func getFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	return fp.ParseURL(url)
}
