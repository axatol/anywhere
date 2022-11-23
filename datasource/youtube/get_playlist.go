package youtube

import (
	"fmt"

	"github.com/axatol/anywhere/config"
	"github.com/axatol/anywhere/util"
)

func (c *Client) GetPlaylist(playlistID string, limit ...int) ([]QueryResult, error) {
	maxResults := 25
	if len(limit) == 1 {
		maxResults = limit[0]
	}

	parameters := util.QueryParameters{
		"key":        c.key,
		"part":       "snippet",
		"maxResults": fmt.Sprint(maxResults),
		"playlistId": playlistID,
	}

	if len(limit) == 1 {
		parameters["maxResults"] = fmt.Sprint(limit[0])
	}

	response, err := util.Request[PlaylistItemList](c.BaseURL+"/playlistItems", util.RequestConfig{Parameters: parameters})
	if err != nil {
		return nil, err
	}

	config.Log.Debugw("youtube list",
		"kind", response.Kind,
		"count", len(response.Items),
	)

	result := make([]QueryResult, len(response.Items))

	for index, item := range response.Items {
		result[index] = QueryResult{
			VideoID:      item.ID,
			VideoTitle:   item.Snippet.Title,
			ChannelID:    item.Snippet.ChannelID,
			ChannelTitle: item.Snippet.ChannelTitle,
			Description:  item.Snippet.Description,
			Thumbnail:    item.Snippet.Thumbnails.Default,
			PublishedAt:  item.Snippet.PublishedAt,
		}
	}

	return result, nil

}
