package itunes

type itunesSearchResults struct {
	ResultCount int                   `json:"resultCount"`
	Podcasts    []*itunesSearchResult `json:"results"`
}

type itunesSearchResult struct {
	CollectionID int `json:"collectionId"`
	TrackID      int `json:"trackId"`

	ArtistName             string `json:"artistName"`
	ArtworkURL600          string `json:"artworkUrl600"`
	ArtworkURL100          string `json:"artworkUrl100"`
	ArtworkURL60           string `json:"artworkUrl60"`
	ArtworkURL30           string `json:"artworkUrl30"`
	CollectionCensoredName string `json:"collectionCensoredName"`
	CollectionName         string `json:"collectionName"`
	Explicitness           string `json:"explicitness"`
	FeedURL                string `json:"feedUrl"`
	Kind                   string `json:"kind"`
	PreviewURL             string `json:"previewUrl"`
	TrackName              string `json:"trackName"`
	TrackCensoredName      string `json:"trackCensoredName"`
	TrackTimeMillis        int    `json:"trackTimeMillis"`
	WrapperType            string `json:"wrapperType"`
}
