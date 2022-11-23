package youtube

import "time"

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type SearchList struct {
	Kind          string   `json:"kind"`
	Etag          string   `json:"etag"`
	NextPageToken string   `json:"nextPageToken"`
	RegionCode    string   `json:"regionCode"`
	PageInfo      PageInfo `json:"pageInfo"`
	Items         []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default Thumbnail `json:"default"`
				Medium  Thumbnail `json:"medium"`
				High    Thumbnail `json:"high"`
			} `json:"thumbnails"`
			ChannelTitle         string    `json:"channelTitle"`
			LiveBroadcastContent string    `json:"liveBroadcastContent"`
			PublishTime          time.Time `json:"publishTime"`
		} `json:"snippet"`
	} `json:"items"`
}

type VideoList struct {
	Kind     string   `json:"kind"`
	Etag     string   `json:"etag"`
	PageInfo PageInfo `json:"pageInfo"`
	Items    []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default  Thumbnail `json:"default"`
				Medium   Thumbnail `json:"medium"`
				High     Thumbnail `json:"high"`
				Standard Thumbnail `json:"standard"`
				Maxres   Thumbnail `json:"maxres"`
			} `json:"thumbnails"`
			ChannelTitle         string   `json:"channelTitle"`
			Tags                 []string `json:"tags"`
			CategoryID           string   `json:"categoryId"`
			LiveBroadcastContent string   `json:"liveBroadcastContent"`
			Localized            struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"localized"`
			DefaultAudioLanguage string `json:"defaultAudioLanguage"`
		} `json:"snippet"`
	} `json:"items"`
}

type PlaylistItemList struct {
	Kind     string   `json:"kind"`
	Etag     string   `json:"etag"`
	PageInfo PageInfo `json:"pageInfo"`
	Items    []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default  Thumbnail `json:"default"`
				Medium   Thumbnail `json:"medium"`
				High     Thumbnail `json:"high"`
				Standard Thumbnail `json:"standard"`
				Maxres   Thumbnail `json:"maxres"`
			} `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
			PlaylistID   string `json:"playlistId"`
			Position     int    `json:"position"`
			ResourceID   struct {
				Kind    string `json:"kind"`
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
			VideoOwnerChannelTitle string `json:"videoOwnerChannelTitle"`
			VideoOwnerChannelID    string `json:"videoOwnerChannelId"`
		} `json:"snippet"`
	} `json:"items"`
}
