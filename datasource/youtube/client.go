package youtube

import (
	"fmt"

	"github.com/axatol/anywhere/config"
)

const (
	APIKeyName = "youtube"
)

type Client struct {
	BaseURL string
	key     string
}

func NewClient() (*Client, error) {
	key := config.Values.Datasource.APIKeys.For(APIKeyName)
	if key == nil {
		return nil, fmt.Errorf("key not available for source: %s", APIKeyName)
	}

	return &Client{
		key:     *key,
		BaseURL: "https://www.googleapis.com/youtube/v3",
	}, nil
}
