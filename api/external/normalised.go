package external

// Artwork is the collection of artwork.
type Artwork struct {
	Image600URL string `json:"600x600"`
	Image100URL string `json:"100x100"`
	Image60URL  string `json:"60x60"`
}

// Podcast is a normalised version of either an itunes podcast
// or an internal API podcast.
type Podcast struct {
	Artwork  *Artwork `json:"artwork"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	FeedURL  string   `json:"-"`
	ItunesID int      `json:"-"`
}

// SearchResponse is the json output for a search request.
type SearchResponse struct {
	Count   int        `json:"count"`
	Results []*Podcast `json:"results"`
}
