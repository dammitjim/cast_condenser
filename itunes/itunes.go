package itunes

import (
	"encoding/json"
	"errors"
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
func Search(term, limit string) error {
	var err error
	if term == "" {
		return errors.New("itunes/itunes.go|Search: no term provided")
	}

	if limit == "" {
		return errors.New("itunes/itunes.go|Search: no limit provided")
	}

	itunesHTTPClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	u := constructRequestURL(term, limit)
	fmt.Println(u)

	resp, err := itunesHTTPClient.Get(u)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		s := string(buf[:])
		logrus.Error(s)
		return errors.New("Non 200 response code received.")
	}

	var results itunesSearchResults
	err = json.Unmarshal(buf, &results)
	if err != nil {
		return err
	}

	b, _ := json.MarshalIndent(results, "", " ")
	logrus.Info(string(b))

	return nil
}
