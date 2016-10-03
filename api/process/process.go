package process

import (
	"condenser/api/external"
	"errors"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

// Run is the main function for processing a set of podcasts.
func Run(podcasts ...*external.Podcast) {
	// Iterate through podcasts in a single routine.
	//
	// Initially I thought making this concurrent would make sense, however
	// after the service has been running for a while we should only
	// need to process new podcasts rarely.
	//
	// It's more important that the new, individual podcasts are processed
	// quickly.
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
		// TODO do something with the tracks
	}

	logrus.Info("done processing new podcasts")
}

const itemConcurrencyLimit = 10

// Concurrently process feed items.
func extractTracks(items []*gofeed.Item) ([]*Track, error) {
	// Process x items at once.
	concurrentLimit := itemConcurrencyLimit
	if len(items) < concurrentLimit {
		concurrentLimit = len(items)
	}

	// Semaphore channel to block extra goroutines from spawning.
	sem := make(chan bool, concurrentLimit)

	// Data queue to be read from.
	queue := make(chan *Track, len(items))

	for _, item := range items {
		// Signify a new goroutine has started by writing to it.
		sem <- true

		go func(i *gofeed.Item) {
			// After the function has finished, read from the channel to unblock the next
			// goroutine.
			defer func() { <-sem }()

			logrus.Debug("Processing " + i.Title)
			track, err := processFeedItem(i)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"enclosures": i.Enclosures,
					"published":  i.Published,
					"itunes":     i.Extensions["itunes"],
				}).Error(err)
				return
			}

			// Write to our data queue.
			queue <- track
		}(item)
	}

	// Block until all our goroutines have finished.
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	tracks := make([]*Track, len(queue))
	count := 0

	// Close the queue otherwise the following loop will
	// read infinitely.
	close(queue)

	// Read from the data queue channel until empty.
	for track := range queue {
		tracks[count] = track
		count++
	}

	return tracks, nil
}

// Process a single feed item, normalise the data
func processFeedItem(item *gofeed.Item) (*Track, error) {
	var err error
	if len(item.Enclosures) != 1 {
		return nil, errors.New("Invalid enclosures found for " + item.Title)
	}

	if _, ok := item.Extensions["itunes"]; !ok {
		err = errors.New("No itunes data parsable from " + item.Title)
		return nil, err
	}

	itunesData := ext.NewITunesItemExtension(item.Extensions["itunes"])

	publishedParsed, err := attemptTimeParsing(item.Published)
	if err != nil {
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
		return nil, err
	}

	return track, nil
}

// templates for attemptTimeParsing
var timeTemplates = []string{
	time.RFC1123,
	time.RFC1123Z,
	"Mon, 2 Jan 2006 15:04 MST",
	"Mon, 2 Jan 2006 15:04 -0700",
	"2 January 2006 15:04",
	"2 January 2006 15:04 MST",
	"Mon, 2 January 2006 15:04 MST",
}

// Iterate through our templates and attempt to parse a time object out.
// Error if all templates have failed.
func attemptTimeParsing(timeString string) (parsed time.Time, err error) {
	success := false
	for _, template := range timeTemplates {
		parsed, err = time.Parse(template, timeString)
		if err != nil {
			logrus.Debug(err)
			continue
		}

		success = true
		break
	}

	if !success {
		err = errors.New("could not parse time " + timeString)
	}

	return
}

// Small wrapper function for retreiving a feed object
func getFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	return fp.ParseURL(url)
}
