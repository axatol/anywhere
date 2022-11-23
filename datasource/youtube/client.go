package youtube

import (
	"fmt"
	"net/url"
	"time"

	"github.com/axatol/anywhere/config"
)

const (
	APIKeyName = "youtube"
)

var (
	client       *Client
	youtubeHosts = []string{"youtu.be", "youtube.com", "music.youtube.com"}
)

type Client struct {
	BaseURL string
	key     string
}

func GetClient() (*Client, error) {
	if client != nil {
		return client, nil
	}

	key := config.Values.Datasource.APIKeys.For(APIKeyName)
	if key == nil {
		return nil, fmt.Errorf("key not available for source: %s", APIKeyName)
	}

	client = &Client{
		key:     *key,
		BaseURL: "https://www.googleapis.com/youtube/v3",
	}

	return client, nil
}

type QueryResult struct {
	VideoID      string    `json:"video_id"`
	VideoTitle   string    `json:"video_title"`
	ChannelID    string    `json:"channel_id"`
	ChannelTitle string    `json:"channel_title"`
	Description  string    `json:"description"`
	Thumbnail    Thumbnail `json:"thumbnail"`
	PublishedAt  time.Time `json:"published_at"`
}

func isYoutubeHost(input string) bool {
	for _, host := range youtubeHosts {
		if input == host {
			return true
		}
	}

	return false
}

func (c *Client) Query(input string, limit ...int) ([]QueryResult, error) {
	videoURL, err := url.Parse(input)
	if err == nil && isYoutubeHost(videoURL.Hostname()) {
		if playlistID := videoURL.Query().Get("list"); playlistID != "" {
			config.Log.Debugw("query by playlist id",
				"input", input,
				"playlist_id", playlistID,
			)

			return c.GetPlaylist(playlistID, limit...)
		}

		if videoID := videoURL.Query().Get("v"); videoID != "" {
			config.Log.Debugw("query by video id",
				"input", input,
				"video_id", videoID,
			)

			return c.GetVideo(videoID)
		}

		return nil, fmt.Errorf("url did not contain any ids parameters: %s", input)
	}

	config.Log.Debugw("query by fuzzy search",
		"input", input,
		"err", err,
	)

	return c.FuzzySearch(input, limit...)
}
