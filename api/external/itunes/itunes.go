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
const itunesHTTPClientTimeout = 10

var itunesHTTPClient *http.Client

// Setup initialises the package.
func Setup() error {
	itunesHTTPClient = &http.Client{
		Timeout: itunesHTTPClientTimeout * time.Second,
	}

	return nil
}

// construct the itunes request url
func constructRequestURL(term, limit string) string {
	u := fmt.Sprintf("%s/search?entity=podcast&term=%s&limit=%s",
		itunesBaseURL,
		url.QueryEscape(term),
		limit,
	)
	return u
}

// request itunes api for data
// error on stdlib err or non 200 code returned
func contactItunes(apiURL string) ([]byte, error) {
	resp, err := itunesHTTPClient.Get(apiURL)
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
		// TODO should probably make this not a thing
		s := string(buf[:])
		logrus.Error(s)
		return nil, apierrors.ITunesNon200.WithDetails(
			fmt.Sprintf("itunes/itunes.go|Search: %s", s))
	}

	return buf, nil
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

	apiURL := constructRequestURL(term, limit)
	logrus.WithField("url", apiURL).Info("Sending request to iTunes API.")

	// get our response body as a byte slice
	buf, err := contactItunes(apiURL)
	if err != nil {
		return nil, apierrors.Generic.WithDetails(
			fmt.Sprintf("itunes/itunes.go|Search: %s", err.Error()))
	}

	// unpack response buffer into struct
	var results itunesSearchResults
	err = json.Unmarshal(buf, &results)
	if err != nil {
		return nil, apierrors.Generic.WithDetails(
			fmt.Sprintf("itunes/itunes.go|Search: %s", err.Error()))
	}

	// extract normalised results
	normalisedResults := make([]*external.Podcast, len(results.Podcasts))
	for i, result := range results.Podcasts {
		normalisedResults[i] = &external.Podcast{
			ItunesID: result.CollectionID,
			Owner:    result.ArtistName,
			Name:     result.CollectionName,
			FeedURL:  result.FeedURL,
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
