package musicbrainz

import (
	"github.com/axatol/anywhere/util"
)

type Artist struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Score    int      `json:"score"`
	Name     string   `json:"name"`
	SortName string   `json:"sort-name"`
	Country  string   `json:"country"`
	ISNIS    []string `json:"isnis"`
	Aliases  []struct {
		Name     string `json:"name"`
		SortName string `json:"sort-name"`
		Type     string `json:"type"`
	} `json:"alises"`
	Disambiguation string `json:"disambiguation,omitempty"`
	Tags           []struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
	} `json:"tags,omitempty"`
}

type ArtistList struct {
	Pagination
	Artists []Artist `json:"artists"`
}

func (c *Client) LookupArtist(query string) ([]Artist, error) {
	response, err := util.Request[ArtistList](baseURL+"/artist", util.RequestConfig{
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

	return response.Artists, nil
}
