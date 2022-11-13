package youtube

import (
	"fmt"

	"github.com/axatol/anywhere/util"
)

type YoutubePlaylistItemList struct {
	Kind  string `json:"kind"`
	Items []struct {
		Kind    string `json:"kind"`
		ID      string `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			ChannelTitle string `json:"channelTitle"`
			ResourceID   struct {
				Kind    string `json:"kind"`
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
		} `json:"snippet"`
	} `json:"items"`
	PageInfo struct {
		TotalResults int `json:"totalResults"`
	} `json:"pageInfo"`
}

func (c *Client) List(playlistID string, limit ...int) (*YoutubePlaylistItemList, error) {
	parameters := util.QueryParameters{
		"key":        c.key,
		"part":       "snippet",
		"maxResults": "25",
		"playlistId": playlistID,
	}

	if len(limit) == 1 {
		parameters["maxResults"] = fmt.Sprint(limit[0])
	}

	result, err := util.Request[YoutubePlaylistItemList](c.BaseURL+"/playlistItems", util.RequestConfig{Parameters: parameters})
	if err != nil {
		return nil, err
	}

	return result, nil
}
