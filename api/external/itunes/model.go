package itunes

type itunesSearchResults struct {
	ResultCount int                   `json:"resultCount"`
	Results     []*itunesSearchResult `json:"results"`
}

type itunesSearchResult struct {
	CollectionID int `json:"collectionId"`
	TrackID      int `json:"trackId"`

	ArtistName             string `json:"artistName"`
	ArtworkURL100          string `json:"artworkUrl100"`
	ArtworkURL60           string `json:"artworkUrl60"`
	CollectionCensoredName string `json:"collectionCensoredName"`
	CollectionName         string `json:"collectionName"`
	Explicitness           string `json:"explicitness"`
	Kind                   string `json:"kind"`
	PreviewURL             string `json:"previewUrl"`
	TrackName              string `json:"trackName"`
	TrackCensoredName      string `json:"trackCensoredName"`
	TrackTimeMillis        int    `json:"trackTimeMillis"`
	WrapperType            string `json:"wrapperType"`
}
