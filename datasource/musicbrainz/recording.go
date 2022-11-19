package musicbrainz

import (
	"github.com/axatol/anywhere/util"
)

// AKA a track
type Recording struct {
	ID           string `json:"id"`
	Score        int    `json:"score"`
	Title        string `json:"title"`
	Length       int64  `json:"length"`
	ArtistCredit []struct {
		Artist struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			SortName string `json:"sort-name"`
		} `json:"artist"`
	} `json:"artist-credit"`
	Releases []struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		StatusID     string `json:"status-id"`
		Status       string `json:"status"`
		ReleaseGroup struct {
			ID          string `json:"id"`
			PrimaryType string `json:"primary-type"`
		} `json:"release-group"`
		Date       string `json:"date"`
		Country    string `json:"country"`
		TrackCount int    `json:"track-count"`
		Media      []struct {
			Position int    `json:"position"`
			Format   string `json:"format"`
			Track    []struct {
				ID     string `json:"id"`
				Number string `json:"number"`
				Title  string `json:"title"`
				Length int    `json:"length"`
			} `json:"track"`
			TrackCount  int `json:"track-count"`
			TrackOffset int `json:"track-offset"`
		} `json:"media"`
		ArtistCredit []struct {
			Artist struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				SortName string `json:"sort-name"`
			} `json:"artist"`
		} `json:"artist-credit,omitempty"`
	} `json:"releases"`
	Tags []struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
	} `json:"tags"`
}

type RecordingList struct {
	Pagination
	Recordings []Recording `json:"recordings"`
}

func (c *Client) LookupRecording(query string) ([]Recording, error) {
	response, err := util.Request[RecordingList](baseURL+"/recording", util.RequestConfig{
		Headers: util.Headers{
			"UserAgent": c.UserAgent(),
		},
		Parameters: util.QueryParameters{
			"fmt":   "json",
			"limit": "5",
			"query": query,
		},
	})

	if err != nil {
		return nil, err
	}

	return response.Recordings, nil
}
