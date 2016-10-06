package external

// Podcast is a normalised version of either an itunes podcast
// or an internal API podcast.
type Podcast struct {
	Artwork60  string `json:"artwork60"`
	Artwork100 string `json:"artwork100"`
	Artwork600 string `json:"artwork600"`
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	FeedURL    string `json:"feed"`
	ItunesID   int    `json:"itunes_id"`

	ID int `json:"id"`
}

// SearchResponse is the json output for a search request.
type SearchResponse struct {
	Count   int        `json:"count"`
	Results []*Podcast `json:"results"`
}
