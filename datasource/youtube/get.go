package youtube

import (
	"github.com/axatol/anywhere/util"
)

type YoutubeVideoList struct {
	Kind  string `json:"kind"`
	Items []struct {
		Kind    string `json:"kind"`
		ID      string `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
}

func (c *Client) Get(id string) (*YoutubeVideoList, error) {
	parameters := util.QueryParameters{
		"key":  c.key,
		"part": "snippet",
		"id":   id,
	}

	result, err := util.Request[YoutubeVideoList](c.BaseURL+"/videos", util.RequestConfig{Parameters: parameters})
	if err != nil {
		return nil, err
	}

	return result, nil
}
