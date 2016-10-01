package itunes

import (
	"condenser/api/apierrors"
	"condenser/api/external"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
)

const itunesBaseURL = "https://itunes.apple.com"

// construct the itunes request url
func constructRequestURL(term, limit string) string {
	u := fmt.Sprintf("%s/search?entity=podcast&term=%s&limit=%s",
		itunesBaseURL,
		url.QueryEscape(term),
		limit,
	)
	return u
}

// Search hits the itunes API for data and returns a normalised
// object.
func Search(term, limit string) (*external.SearchResponse, error) {
	var err error
	if term == "" {
		return nil, apierrors.GenericValidation.WithDetails(
			"itunes/itunes.go|Search: no term provided")
	}

	if limit == "" {
		return nil, apierrors.GenericValidation.WithDetails(
			"itunes/itunes.go|Search: no limit provided")
	}

	itunesHTTPClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	u := constructRequestURL(term, limit)
	logrus.WithField("url", u).Info("Sending request to iTunes API.")

	resp, err := itunesHTTPClient.Get(u)
	if err != nil {
		return nil, apierrors.Generic.WithDetails(
			fmt.Sprintf("itunes/itunes.go|Search: %s", err.Error()))
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, apierrors.Generic.WithDetails(
			fmt.Sprintf("itunes/itunes.go|Search: %s", err.Error()))
	}

	if resp.StatusCode != 200 {
		s := string(buf[:])
		logrus.Error(s)
		return nil, apierrors.ITunesNon200.WithDetails(
			fmt.Sprintf("itunes/itunes.go|Search: %s", s))
	}

	var results itunesSearchResults
	err = json.Unmarshal(buf, &results)
	if err != nil {
		return nil, apierrors.Generic.WithDetails(
			fmt.Sprintf("itunes/itunes.go|Search: %s", err.Error()))
	}

	normalisedResults := make([]*external.Podcast, len(results.Podcasts))
	for i, result := range results.Podcasts {
		normalisedResults[i] = &external.Podcast{
			Owner:       result.ArtistName,
			PodcastName: result.CollectionName,
			FeedURL:     result.FeedURL,
			Artwork: &external.Artwork{
				Image60URL:  result.ArtworkURL60,
				Image100URL: result.ArtworkURL100,
				Image600URL: result.ArtworkURL600,
			},
		}
	}

	return &external.SearchResponse{
		Count:   len(normalisedResults),
		Results: normalisedResults,
	}, nil
}
