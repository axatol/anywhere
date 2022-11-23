package youtube

import (
	"fmt"

	"github.com/axatol/anywhere/config"
	"github.com/axatol/anywhere/util"
)

func (c *Client) FuzzySearch(query string, limit ...int) ([]QueryResult, error) {
	maxResults := 5
	if len(limit) == 1 {
		maxResults = limit[0]
	}

	parameters := util.QueryParameters{
		"key":        c.key,
		"part":       "snippet",
		"order":      "relevance",
		"safeSearch": "none",
		"type":       "video",
		"maxResults": fmt.Sprint(maxResults),
		"q":          query,
	}

	if len(limit) == 1 {
		parameters["maxResults"] = fmt.Sprint(limit[0])
	}

	response, err := util.Request[SearchList](c.BaseURL+"/search", util.RequestConfig{Parameters: parameters})
	if err != nil {
		return nil, err
	}

	config.Log.Debugw("youtube search",
		"kind", response.Kind,
		"count", len(response.Items),
	)

	result := make([]QueryResult, len(response.Items))

	for index, item := range response.Items {
		result[index] = QueryResult{
			VideoID:      item.ID.VideoID,
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
