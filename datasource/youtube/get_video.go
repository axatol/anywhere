package youtube

import (
	"github.com/axatol/anywhere/config"
	"github.com/axatol/anywhere/util"
)

func (c *Client) GetVideo(id string) ([]QueryResult, error) {
	parameters := util.QueryParameters{
		"key":  c.key,
		"part": "snippet",
		"id":   id,
	}

	response, err := util.Request[VideoList](c.BaseURL+"/videos", util.RequestConfig{Parameters: parameters})
	if err != nil {
		return nil, err
	}

	config.Log.Debugw("youtube get",
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
