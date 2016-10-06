package process

import "errors"

// Track represents an individual podcast track.
type Track struct {
	Duration string `json:"duration"`
	Name     string `json:"name"`
	MediaURL string `json:"media_url"`
	Summary  string `json:"summary"`
	Posted   string `json:"posted"`
}

// Validate ensures the minimum required attributes are present.
func (t *Track) Validate() error {
	if t.Name == "" {
		return errors.New("No name set.")
	}

	if t.MediaURL == "" {
		return errors.New("No media url set.")
	}

	return nil
}

type trackRequest struct {
	Duration string `json:"duration"`
	Name     string `json:"name"`
	MediaURL string `json:"media_url"`
	Summary  string `json:"summary"`
	Posted   string `json:"feed_posted"`

	PodcastPK int `json:"podcast"`
}

type podcastSuccessfulSaveResponse struct {
	Artwork60  string `json:"artwork60"`
	Artwork100 string `json:"artwork100"`
	Artwork600 string `json:"artwork600"`
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	FeedURL    string `json:"feed"`
	ItunesID   int    `json:"itunes_id"`

	ID      int    `json:"id"`
	APILink string `json:"api_link"`
}

type trackSuccessfulSaveResponse struct {
	Duration string `json:"duration"`
	Name     string `json:"name"`
	MediaURL string `json:"media_url"`
	Summary  string `json:"summary"`
	Posted   string `json:"feed_posted"`

	ID      int    `json:"id"`
	APILink string `json:"api_link"`
}
