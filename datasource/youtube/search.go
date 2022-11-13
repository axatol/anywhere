package youtube

import (
	"fmt"

	"github.com/axatol/anywhere/util"
)

type YoutubeSearchList struct {
	Kind  string `json:"kind"`
	Items []struct {
		Kind string `json:"kind"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
}

func (c *Client) Search(query string, limit ...int) (*YoutubeSearchList, error) {
	parameters := util.QueryParameters{
		"key":        c.key,
		"part":       "snippet",
		"order":      "relevance",
		"safeSearch": "none",
		"type":       "video",
		"maxResults": "1",
		"q":          query,
	}

	if len(limit) == 1 {
		parameters["maxResults"] = fmt.Sprint(limit[0])
	}

	result, err := util.Request[YoutubeSearchList](c.BaseURL+"/search", util.RequestConfig{Parameters: parameters})
	if err != nil {
		return nil, err
	}

	return result, nil
}
