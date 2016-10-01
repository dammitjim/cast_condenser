package process

import "errors"

type Track struct {
	Duration string `json:"duration"`
	Name     string `json:"name"`
	MediaURL string `json:"media_url"`
	Summary  string `json:"summary"`
	Posted   int64  `json:"posted"`
}

func (t *Track) Validate() error {
	if t.Name == "" {
		return errors.New("No name set.")
	}

	if t.MediaURL == "" {
		return errors.New("No media url set.")
	}

	return nil
}
