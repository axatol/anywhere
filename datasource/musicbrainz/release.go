package musicbrainz

import (
	"github.com/axatol/anywhere/util"
)

// AKA album
type Release struct {
	ID           string `json:"id"`
	Score        int    `json:"score"`
	Count        int    `json:"count"`
	Title        string `json:"title"`
	StatusID     string `json:"status-id"`
	Status       string `json:"status"`
	ArtistCredit []struct {
		Artist struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			SortName string `json:"sort-name"`
		} `json:"artist"`
	} `json:"artist-credit"`
	ReleaseGroup struct {
		ID          string `json:"id"`
		PrimaryType string `json:"primary-type"`
	} `json:"release-group"`
	Date       string `json:"date"`
	Country    string `json:"country"`
	TrackCount int    `json:"track-count"`
}

type ReleaseList struct {
	Pagination
	Releases []Release `json:"releases"`
}

func (c *Client) LookupRelease(query string) ([]Release, error) {
	response, err := util.Request[ReleaseList](baseURL+"/release", util.RequestConfig{
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

	return response.Releases, nil
}
